package handler

import (
	"fmt"
	"grovo/internal/controller/ws"

	"github.com/gofiber/fiber/v2"
)

type ClientList struct {
	*BaseResponse
	Data []ClientInRoom `json:"data"`
}

type ClientInRoom struct {
	ClientId string `json:"clientId"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
}

func GetClientInRoom(c *fiber.Ctx, h *ws.Hub) error {

	var clients []ClientInRoom
	roomId := c.Params("roomId")
	fmt.Println(roomId)

	if _, isExist := h.Rooms[roomId]; !isExist {
		res := ClientList{
			BaseResponse: &BaseResponse{
				Success: true,
				Code:    200,
				Message: "no client",
			},
			Data: make([]ClientInRoom, 0),
		}
		return c.JSON(res)
	}

	if len(h.Rooms[roomId].Clients) == 0 {
		res := ClientList{
			BaseResponse: &BaseResponse{
				Success: true,
				Code:    200,
				Message: "no client",
			},
			Data: make([]ClientInRoom, 0),
		}
		return c.JSON(res)
	}

	for _, client := range h.Rooms[roomId].Clients {
		clients = append(clients, ClientInRoom{
			ClientId: client.ClientId,
			Username: client.Username,
			RoomId:   client.RoomId,
		})
	}

	res := ClientList{
		BaseResponse: &BaseResponse{
			Success: true,
			Code:    200,
			Message: "success get clients this room",
		},
		Data: clients,
	}

	return c.JSON(res)
}
