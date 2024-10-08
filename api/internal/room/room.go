package room

import (
	"errors"
	"fmt"

	"github.com/Caiollaz/Planning-Poker/backend/internal/vote"
	"github.com/google/uuid"

	"github.com/Caiollaz/Planning-Poker/backend/internal/user"
)

type GameState string

const (
	InProgress    GameState = "IN_PROGRESS"
	CardsRevealed GameState = "CARDS_REVEALED"

	maxUsers int = 12
)

type Room struct {
	id        string      `json:"id"`
	name      string      `json:"name"`
	users     []user.User `json:"users"`
	admin     *user.User  `json:"admin"`
	votes     vote.Votes  `json:"votes"`
	gameState GameState   `json:"gameState"`
}

func NewRoom(name string) *Room {
	return &Room{
		id:        uuid.New().String(),
		name:      name,
		users:     []user.User{},
		admin:     nil,
		votes:     vote.NewVotes(),
		gameState: InProgress,
	}
}

func NewRoomWithAdmin(name string, admin *user.User) *Room {
	room := NewRoom(name)
	room.SetAdmin(admin)
	room.AddUser(admin)
	return room
}

func (r *Room) GetId() string {
	return r.id
}

func (r *Room) GetName() string {
	return r.name
}

func (r *Room) AddUser(user *user.User) {
	r.users = append(r.users, *user)
}

func (r *Room) RemoveUser(userId string) {
	for i, user := range r.users {
		if user.GetId() == userId {
			r.users = append(r.users[:i], r.users[i+1:]...)
			break
		}
	}

	if r.admin != nil && r.admin.GetId() == userId {
		r.admin = nil
	}

	r.votes.RemoveVote(userId)
}

func (r *Room) IsFull() bool {
	return len(r.users) >= maxUsers
}

func (r *Room) GetUserWithId(userId string) (*user.User, error) {
	for _, u := range r.users {
		if u.GetId() == userId {
			return &u, nil
		}
	}
	return nil, errors.New("não foi possível encontrar o usuário com o ID fornecido")
}

func (r *Room) GetAdmin() *user.User {
	return r.admin
}

func (r *Room) SetAdmin(user *user.User) {
	r.admin = user
}

func (r *Room) GetVotes() *vote.Votes {
	return &r.votes
}

func (r *Room) GetGameState() GameState {
	return r.gameState
}

func (r *Room) SetGameState(state GameState) {
	r.gameState = state
}

func (r *Room) ResetVotes() {
	r.votes = vote.NewVotes()
}

func (r *Room) String() string {
	return fmt.Sprintln("Id:", r.id, "Name:", r.name, "Users:", r.users, "Admin:", r.admin, "Votes:", r.votes, "State:", r.gameState)
}
