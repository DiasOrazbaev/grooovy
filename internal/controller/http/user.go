package http

import (
	"errors"
	"grovo/internal/common"
	"grovo/internal/common/fiber/util"
	"grovo/internal/common/jwt"
	"grovo/internal/controller/ws"
	"grovo/internal/controller/ws/handler"
	"grovo/internal/usecase"
	"grovo/pkg/mw"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserController struct {
	usecase *usecase.UserUsecase
	log     *zap.Logger
}

func NewUserController(usecase *usecase.UserUsecase, logger *zap.Logger) *UserController {
	return &UserController{usecase: usecase, log: logger}
}

func (c *UserController) Routes(app *fiber.App) {
	hub := ws.NewHub()
	go hub.Run()

	app.Post("/login", c.Login)
	app.Post("/register", c.Register)

	app.Post("/ws", mw.JWTAuth, func(c *fiber.Ctx) error {
		return handler.CreateRoom(c, hub)
	})

	app.Get("/ws", func(c *fiber.Ctx) error {
		return handler.GetAvailableRooms(c, hub)
	})

	app.Get("/ws/rooms/:roomId", func(c *fiber.Ctx) error {
		return handler.GetClientInRoom(c, hub)
	})

	app.Get("/ws/:roomId", handler.JoinRoom(hub))
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&loginRequest); err != nil {
		return util.Send(ctx, fiber.StatusBadRequest, "Invalid request")
	}

	user, err := c.usecase.Login(ctx.Context(), loginRequest.Username)

	if errors.Is(err, common.ErrInvalidCredentials) {
		return util.Send(ctx, fiber.StatusUnauthorized, "Invalid credentials")
	} else if err != nil {
		c.log.Error("Error", zap.Error(err))
		return util.Send(ctx, fiber.StatusInternalServerError, "Error")
	}

	if user.Password != loginRequest.Password {
		return util.Send(ctx, fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := jwt.GenerateToken(jwt.Claims{
		Username: loginRequest.Username,
		Email:    user.Email,
		Id:       user.Id,
	})
	if err != nil {
		c.log.Error("Error", zap.Error(err))
		return util.Send(ctx, fiber.StatusInternalServerError, "Error")
	}

	return util.Send(ctx, fiber.StatusOK, fiber.Map{
		"token": token,
	})
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var registerRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&registerRequest); err != nil {
		return util.Send(ctx, fiber.StatusBadRequest, "Invalid request")
	}

	u, err := c.usecase.Register(ctx.Context(), registerRequest.Username, registerRequest.Password)

	if errors.Is(err, common.ErrUsernameTaken) {
		return util.Send(ctx, fiber.StatusConflict, "Username taken")
	} else if err != nil {
		c.log.Error("Error", zap.Error(err))
		return util.Send(ctx, fiber.StatusInternalServerError, "Error")
	}

	token, err := jwt.GenerateToken(jwt.Claims{
		Username: u.Username,
		Email:    u.Email,
		Id:       u.Id,
	})
	if err != nil {
		c.log.Error("Error", zap.Error(err))
		return util.Send(ctx, fiber.StatusInternalServerError, "Error")
	}

	return util.Send(ctx, fiber.StatusCreated, fiber.Map{"token": token})
}
