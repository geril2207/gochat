package auth

import (
	"os"
	"time"

	"github.com/geril2207/gochat/apps/server/utils"
	"github.com/geril2207/gochat/packages/db/users"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Pool *pgxpool.Pool
}

const (
	WrongCredentialsMessage = "Wrong credentials"
	HashPasswordCost        = bcrypt.DefaultCost
)

func Login(ctx *fiber.Ctx) error {
	var body *LoginBody
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	validate := utils.NewValidator()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.FormatValidationErrors(err))
	}

	user, err := users.GetUserByLogin(body.Login)
	if err != nil || user.Id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": WrongCredentialsMessage,
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": WrongCredentialsMessage,
		})
	}

	tokens, err := GenerateTokens(user)

	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Value:    tokens.Refresh,
		HTTPOnly: true,
		Secure:   true,
		Name:     "refresh",
	})
	return ctx.JSON(fiber.Map{
		"user":   user,
		"tokens": tokens,
	})
}

func Register(ctx *fiber.Ctx) error {
	var body *LoginBody
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	validate := utils.NewValidator()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.FormatValidationErrors(err))
	}
	isExists, _ := users.IsUserExistsInDb(body.Login)
	if isExists {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User with this login already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), HashPasswordCost)
	if err != nil {
		return err
	}

	user, err := users.InsertUser(body.Login, string(hashedPassword))
	if err != nil {
		return err
	}

	tokens, err := GenerateTokens(user)

	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"user":   user,
		"tokens": tokens,
	})
}

type AuthTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func GenerateTokens(user users.User) (AuthTokens, error) {
	key := os.Getenv("JWT_KEY")

	if key == "" {
		panic("JWT_KEY not found")
	}

	accessPayload := jwt.MapClaims{"id": user.Id, "iat": time.Now().Add(time.Minute * 30).Unix()}
	refreshPayload := jwt.MapClaims{"id": user.Id, "iat": time.Now().Add(time.Hour * 24 * 30).Unix()}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)

	accessToken, err := jwtAccessToken.SignedString([]byte(key))
	refreshToken, err := jwtRefreshToken.SignedString([]byte(key))

	return AuthTokens{Access: accessToken, Refresh: refreshToken}, err
}
