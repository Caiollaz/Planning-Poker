package main

import (
	"net/http"
	"os"

	"github.com/Caiollaz/Planning-Poker/backend/internal/auth"
	"github.com/Caiollaz/Planning-Poker/backend/internal/room"
	"github.com/Caiollaz/Planning-Poker/backend/internal/ws"
	"github.com/Caiollaz/Planning-Poker/backend/pkg/config"
	db "github.com/Caiollaz/Planning-Poker/backend/pkg/db"
	"github.com/Caiollaz/Planning-Poker/backend/pkg/log"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

var (
	logger         = log.New()
	database       = db.SetupDatabaseConnection(config.GetDbClient())
	roomRepository = room.NewRoomRepository(database, logger)
	roomService    = room.NewRoomService(roomRepository, logger)
	roomController = room.NewRoomController(roomService, logger)
	wsServer       = ws.NewWSServer(roomService)
)

func main() {
	defer database.Close()

	r := gin.Default()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("CLIENT_URL")},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
		},
		AllowedHeaders: []string{"*"},
	})
	r.Use(c)

	roomRoutes := r.Group("api/room")
	{
		roomRoutes.POST("/", roomController.CreateRoom)
		roomRoutes.POST("/:id", roomController.JoinRoom)
		roomRoutes.GET("/:id", auth.IsUserAuthorizedInRoom, roomController.GetRoom)

		roomRoutes.GET("/ws/:token", func(ctx *gin.Context) {
			token := ctx.Param("token")
			userClaims, err := auth.GetUserClaimsFromToken(token)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid token")
				return
			}

			ws.ServeWS(wsServer, ctx.Writer, ctx.Request, userClaims)
		})
	}

	r.Run()
}
