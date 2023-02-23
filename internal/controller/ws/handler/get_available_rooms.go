package handler

import (
	"grovo/internal/controller/ws"
	"log"

	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Res struct {
	*BaseResponse
	Data *[]RoomList `json:"data"`
}

type RoomList struct {
	RoomName string `json:"roomName"`
	RoomId   string `json:"roomId"`
}

func GetAvailableRooms(c *fiber.Ctx, h *ws.Hub) error {

	rooms := make([]RoomList, 0)
	log.Println("rooms", rooms)
	for _, room := range h.Rooms {
		rooms = append(rooms, RoomList{
			RoomName: room.RoomName,
			RoomId:   room.RoomId,
		})
	}

	res := Res{
		BaseResponse: &BaseResponse{
			Success: true,
			Code:    200,
			Message: "success get rooms",
		},
		Data: &rooms,
	}
	log.Println("res", res)
	return c.JSON(res)
}
