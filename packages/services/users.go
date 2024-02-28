package services

import (
	"context"

	"github.com/geril2207/gochat/packages/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersService struct {
	pool *pgxpool.Pool
}

func ProvideUsersService(pool *pgxpool.Pool) UsersService {
	return UsersService{
		pool: pool,
	}
}

func (s *UsersService) GetUserById(id int) (models.User, error) {
	var user models.User

	err := s.pool.
		QueryRow(context.Background(), "select id,login,password from users where id=$1", id).
		Scan(&user.Id, &user.Login, &user.Password)

	return user, err
}

func (s *UsersService) GetUserByLogin(login string) (models.User, error) {
	var user models.User

	err := s.pool.
		QueryRow(context.Background(), "select id,login,password from users where login=$1", login).
		Scan(&user.Id, &user.Login, &user.Password)

	return user, err
}

func (s *UsersService) IsUserExistsInDb(login string) (bool, error) {
	var result int
	err := s.pool.QueryRow(context.Background(), "select 1 from users where login=$1", login).
		Scan(&result)

	return result == 1, err
}

func (s *UsersService) InsertUser(login, password string) (models.User, error) {
	var user models.User
	err := s.pool.QueryRow(context.Background(),
		"insert into users(login,password) values($1,$2) RETURNING id,login,password", login, password).
		Scan(&user.Id, &user.Login, &user.Password)

	return user, err
}
