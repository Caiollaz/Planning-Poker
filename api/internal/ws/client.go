package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Caiollaz/Planning-Poker/backend/internal/vote"

	"github.com/Caiollaz/Planning-Poker/backend/internal/auth"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var (
	openingBracket = []byte{'['}
	closingBracket = []byte{']'}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	id       string
	game     *Game
	conn     *websocket.Conn
	send     chan []byte
	wsServer *WSServer
}

func newClient(conn *websocket.Conn, wsServer *WSServer, id string, room *Game) *Client {
	return &Client{
		id:       id,
		game:     room,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// O WSServer fechou o canal.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(openingBracket)
			w.Write(message)

			// Se houver várias mensagens na fila, anexe-as juntas
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte(","))
				w.Write(<-c.send)
			}

			w.Write(closingBracket)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.setAsAFK()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Inicia um loop de leitura sem fim, aguardando mensagens do cliente
	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("erro de fechamento inesperado: %v", err)
			}
			break
		}

		c.handleNewMessage(jsonMessage)
	}
}

func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Erro na mensagem JSON unmarshal %s", err)
		return
	}

	log.Println("Message:", message)
	switch message.Action {
	case VoteSubmittedAction:
		c.handleVoteSubmittedMessage(message)
	case RevealCardsAction:
		c.handleRevealCardsMessage(message)
	case StartNewVotingAction:
		c.handleStartNewVotingMessage(message)
	case KickAction:
		c.handleKickMessage(message)
	}
}

func (c *Client) setAsAFK() {
	c.sendClientIsAFKMessage()
	c.game.setAFK <- c
	go c.game.disconnectAFKWithTimeout(c)
	close(c.send)
	c.conn.Close()
}

func ServeWS(wsServer *WSServer, w http.ResponseWriter, r *http.Request, claims *auth.UserClaims) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := newClient(conn, wsServer, claims.UserId, nil)
	client.joinGame(claims.RoomId, claims.UserId)

	go client.readPump()
	go client.writePump()
}

func (c *Client) joinGame(gameId, userId string) {
	game := c.wsServer.findGameById(gameId)
	if game == nil {
		game = c.wsServer.createGame(gameId)
	}

	user, err := c.wsServer.findUserById(gameId, userId)
	if err != nil {
		return
	}

	if c.game == nil {
		c.game = game
		game.register <- c

		message := Message{Action: RoomJoinedAction, ClientId: user.GetId(), Message: user.GetName()}
		c.game.broadcast <- &message
	}
}

func (c *Client) handleVoteSubmittedMessage(message Message) {
	value, err := strconv.Atoi(message.Message)
	if err != nil || !vote.IsValidValue(value) {
		return
	}

	v := vote.NewVote(message.ClientId, value)
	c.wsServer.saveVote(c.game.id, v)

	if vote.IsValueAccountable(value) {
		message.Message = strconv.Itoa(vote.Private)
	}

	c.game.broadcast <- &message
}

func (c *Client) handleRevealCardsMessage(message Message) {
	votes, err := c.wsServer.revealCards(c.game.id)
	if err != nil {
		return
	}

	votesJson, err := json.Marshal(votes)
	if err != nil {
		return
	}

	revealedVotes := Message{Action: CardsRevealedAction, ClientId: message.ClientId, Message: string(votesJson)}
	c.game.broadcast <- &revealedVotes
}

func (c *Client) handleKickMessage(message Message) {
	room, err := c.wsServer.roomService.GetRoom(c.game.id)
	if err != nil {
		return
	}

	if room.GetAdmin() != nil && room.GetAdmin().GetId() != c.id {
		return
	}

	err = c.wsServer.removeUser(c.game.id, message.Message)
	if err != nil {
		return
	}

	clientKickedMessage := Message{Action: DisconnectedAction, ClientId: message.Message}
	c.game.broadcast <- &clientKickedMessage
}

func (c *Client) handleStartNewVotingMessage(message Message) {
	err := c.wsServer.resetVotingSession(c.game.id)
	if err != nil {
		return
	}

	message.Action = NewVotingStartedAction
	c.game.broadcast <- &message
}

func (c *Client) sendClientIsAFKMessage() {
	game := c.wsServer.findGameById(c.game.id)
	if game == nil {
		return
	}

	user, err := c.wsServer.findUserById(game.id, c.id)
	if err != nil {
		return
	}

	message := Message{Action: IsAFKAction, ClientId: user.GetId()}
	c.game.broadcast <- &message
}
