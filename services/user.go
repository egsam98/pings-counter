package services

import (
	"context"

	"github.com/pkg/errors"

	"github.com/egsam98/pings-counter/db"
)

type UserService struct {
	client *db.PrismaClient
}

func NewUserService(client *db.PrismaClient) *UserService {
	return &UserService{client: client}
}

// Подсчет кол-ва роботов среди пользователей
func (us *UserService) CountRobots(ctx context.Context) (int, error) {
	var result []struct{ Count int }
	err := us.client.QueryRaw(`select count(*) as count from User where isRobot = true`).Exec(ctx, &result)
	if err != nil && err != db.ErrNotFound {
		return 0, errors.WithStack(err)
	}
	return result[0].Count, nil
}
