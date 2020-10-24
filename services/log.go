package services

import (
	"context"

	"github.com/pkg/errors"

	"github.com/egsam98/pings-counter/db"
	"github.com/egsam98/pings-counter/utils/gqlerrors"
)

type LogService struct {
	client *db.PrismaClient
}

func NewLogService(client *db.PrismaClient) *LogService {
	return &LogService{client: client}
}

// Логирование пинга со стороны пользователя
func (s *LogService) Log(ctx context.Context, userId int) error {
	_, err := s.client.User.CreateOne(db.User.ID.Set(userId)).Exec(ctx)
	if err != nil && !gqlerrors.Is(err, gqlerrors.UniqueError) {
		return errors.WithStack(err)
	}

	_, err = s.client.Log.CreateOne(
		db.Log.User.Link(db.User.ID.Equals(userId)),
	).Exec(ctx)
	return errors.WithStack(err)
}
