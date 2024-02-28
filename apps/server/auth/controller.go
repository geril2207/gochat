package auth

import (
	"net/http"
	"time"

	"github.com/geril2207/gochat/packages/config"
	"github.com/geril2207/gochat/packages/models"
	"github.com/geril2207/gochat/packages/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	pool         *pgxpool.Pool
	usersService services.UsersService
	config       config.EnvConfig
}

const (
	WrongCredentialsMessage = "Wrong credentials"
	HashPasswordCost        = bcrypt.DefaultCost
)

var WrongCredentialsResponse = map[string]interface{}{
	"message": WrongCredentialsMessage,
}

func ProvideAuthController(pool *pgxpool.Pool, usersService services.UsersService, config config.EnvConfig) *AuthController {
	return &AuthController{
		pool:         pool,
		usersService: usersService,
		config:       config,
	}
}

func (c *AuthController) Login(ctx echo.Context) error {
	var body *LoginBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	if err := ctx.Validate(body); err != nil {
		return err
	}

	user, err := c.usersService.GetUserByLogin(body.Login)
	if err != nil || user.Id == 0 {
		return ctx.JSON(http.StatusBadRequest, WrongCredentialsResponse)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, WrongCredentialsResponse)
	}

	tokens, err := c.GenerateTokens(user)

	if err != nil {
		return err
	}

	refreshCookie := &http.Cookie{
		Value:    tokens.Refresh,
		HttpOnly: true,
		Secure:   true,
		Name:     "refresh",
		Path:     "/",
	}
	ctx.SetCookie(refreshCookie)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (c *AuthController) Register(ctx echo.Context) error {
	var body *LoginBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	if err := ctx.Validate(body); err != nil {
		return err
	}
	isExists, _ := c.usersService.IsUserExistsInDb(body.Login)
	if isExists {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "User with this login already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), HashPasswordCost)
	if err != nil {
		return err
	}

	user, err := c.usersService.InsertUser(body.Login, string(hashedPassword))
	if err != nil {
		return err
	}

	tokens, err := c.GenerateTokens(user)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (c *AuthController) Refresh(ctx echo.Context) error {
	userId := ctx.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims).Id

	return ctx.JSON(http.StatusOK, userId)
}

type AuthTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type JwtCustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func (c *AuthController) GenerateTokens(user models.User) (AuthTokens, error) {
	key := c.config.JwtKey

	if key == "" {
		panic("JwtKey not found")
	}

	accessPayload := &JwtCustomClaims{user.Id, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30))}}
	refreshPayload := &JwtCustomClaims{user.Id, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30))}}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)

	accessToken, err := jwtAccessToken.SignedString([]byte(key))
	refreshToken, err := jwtRefreshToken.SignedString([]byte(key))

	return AuthTokens{Access: accessToken, Refresh: refreshToken}, err
}
