package room

import (
	"fmt"
	"net/http"

	"github.com/Caiollaz/Planning-Poker/backend/internal/auth"
	"github.com/Caiollaz/Planning-Poker/backend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Controller interface {
	CreateRoom(ctx *gin.Context)
	JoinRoom(ctx *gin.Context)
	GetRoom(ctx *gin.Context)
}

type controller struct {
	service Service
	logger  log.Logger
}

func NewRoomController(service Service, logger log.Logger) Controller {
	return &controller{service, logger}
}

func (c *controller) CreateRoom(ctx *gin.Context) {
	var roomCreation RoomCreation
	if err := ctx.BindJSON(&roomCreation); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Corpo de solicitação inválido")
		return
	}

	err := validate.Struct(roomCreation)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "O campo 'roomName' ou 'username' não foi fornecido ou não é válido.")
		return
	}
	c.logger.Info(fmt.Sprintln("criando sala com nome:", roomCreation.RoomName, "e nome de usuário do adm:", roomCreation.Username))

	room, err := c.service.CreateRoom(roomCreation.RoomName, roomCreation.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, err)
		return
	}
	c.logger.Info(fmt.Sprintln("Sala criada:", room))

	t, err := auth.GenerateToken(room.GetAdmin(), room.GetId(), true)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusCreated,
		CreatedRoomWithUser{
			&auth.AuthUser{
				Id:        room.GetAdmin().GetId(),
				RoomId:    room.GetId(),
				Name:      room.GetAdmin().GetName(),
				Token:     t.Token,
				ExpiresAt: t.ExpiresAt},
			room})
}

func (c *controller) GetRoom(ctx *gin.Context) {
	roomId := ctx.Param("id")
	c.logger.Info(fmt.Sprintln("Buscando sala com id:", roomId))

	if roomId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "'roomId' deve ser fornecido")
		return
	}

	userId, _ := auth.GetUserId(ctx)
	fmt.Println("userId é", userId)

	room, err := c.service.GetRoomWithVotesBasedOnGameState(roomId, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	_, err = room.GetUserWithId(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
	}

	c.logger.Info(fmt.Sprintln("Espaço encontrado:", room))
	ctx.AbortWithStatusJSON(http.StatusOK, &room)
}

func (c *controller) JoinRoom(ctx *gin.Context) {
	var joinRoom JoinRoom
	if err := ctx.BindJSON(&joinRoom); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Corpo de solicitação inválido")
		return
	}

	err := validate.Struct(joinRoom)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "O campo 'roomId' ou 'username' não foi fornecido")
		return
	}

	c.logger.Info(fmt.Sprintln("Entrando na sala com id:", joinRoom.RoomId, "e nome de usuário:", joinRoom.Username))
	room, user, err := c.service.JoinRoom(joinRoom.RoomId, joinRoom.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, err)
		return
	}
	c.logger.Info(fmt.Sprintln("Usuario:", user, "entrou na sala:", room))

	t, err := auth.GenerateToken(user, room.GetId(), false)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK,
		CreatedRoomWithUser{
			&auth.AuthUser{
				Id:        user.GetId(),
				RoomId:    room.GetId(),
				Name:      user.GetName(),
				Token:     t.Token,
				ExpiresAt: t.ExpiresAt},
			room})
}

func init() {
	validate = validator.New()
}
