package handler

import (
	"grovo/internal/controller/ws"
	"log"

	"github.com/gofiber/fiber/v2"
)

type MyRoom struct {
	RoomName string
	RoomId   string
}

func CreateRoom(c *fiber.Ctx, h *ws.Hub) error {

	room := new(MyRoom)

	if err := c.BodyParser(room); err != nil {
		panic(err)
	}

	h.Rooms[room.RoomId] = &ws.Room{
		RoomId:   room.RoomId,
		RoomName: room.RoomName,
		Clients:  make(map[string]*ws.Client),
	}

	log.Println("room", room)
	log.Println(len(h.Rooms))

	return c.JSON(room)

}
