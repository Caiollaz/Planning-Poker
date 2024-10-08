package room

import (
	"errors"
	"fmt"

	db "github.com/Caiollaz/Planning-Poker/backend/pkg/db"
	"github.com/Caiollaz/Planning-Poker/backend/pkg/log"
)

type Repository interface {
	CreateRoom(room *Room) error
	GetRoom(roomId string) (*Room, error)
	SetRoom(room *Room) error
}

type repository struct {
	db     db.DBContext
	logger log.Logger
}

func NewRoomRepository(db db.DBContext, logger log.Logger) Repository {
	return &repository{db, logger}
}

func (r *repository) CreateRoom(newRoom *Room) error {
	if err := r.db.Get(newRoom.id, &Room{}); err == nil {
		r.logger.Error(fmt.Sprintln("Sala com", newRoom.id, "já existe"))
		return errors.New(fmt.Sprintln("Sala com", newRoom.id, "já existe"))
	}

	if err := r.db.Set(newRoom.id, newRoom); err != nil {
		r.logger.Error(fmt.Sprintln("Erro ao criar uma sala com id", newRoom.id))
		return errors.New(fmt.Sprintln("Erro ao criar uma sala com id", newRoom.id))
	}

	return nil
}

func (r *repository) GetRoom(roomId string) (*Room, error) {
	var room *Room
	if err := r.db.Get(roomId, &room); err == nil {
		return room, nil
	}
	return nil, errors.New(fmt.Sprintln("Sala com id", roomId, "não existe"))
}

func (r *repository) SetRoom(room *Room) error {
	if err := r.db.Set(room.id, room); err != nil {
		r.logger.Error(fmt.Sprintln("Erro ao atualizar a sala com id", room.id))
		return err
	}

	return nil
}
