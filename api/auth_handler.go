package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/adalbertjnr/hotel-project/db"
	"github.com/adalbertjnr/hotel-project/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStorer
}

func NewAuthHandler(userStore db.UserStorer) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(genericResp{
		Type:    "error",
		Message: "invalid credentials",
	})
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	tokenByte := "5as4c56as4d654as569C8AS908"
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenByte))
	if err != nil {
		fmt.Println("failed to sign with secret", err)
	}
	return tokenString
}
