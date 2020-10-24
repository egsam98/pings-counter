package jobs

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/egsam98/pings-counter/db"
)

const (
	RobotPingsCount = 100         // Кол-во пингов со стороны пользователя для определения его как робот
	Latency         = time.Minute // Пауза между вызовами поиска роботов
)

var _ Job = (*RobotWatcherService)(nil)

// Сервис-наблюдатель за базой пользователей, выявляющих среди них роботов
type RobotWatcherService struct {
	client *db.PrismaClient
}

func NewRobotWatcherService(client *db.PrismaClient) *RobotWatcherService {
	return &RobotWatcherService{client: client}
}

func (ws *RobotWatcherService) Run() {
	ctx := context.TODO()
	for {
		log.Println("RobotWatcherService's watching...")

		start := time.Now()
		time.Sleep(Latency)
		end := time.Now()

		//Поиск всех логов за последнее время Latency
		logs, err := ws.client.Log.FindMany(
			db.Log.User.Where(
				db.User.Not(db.User.IsRobot.Equals(true)),
			),
			db.Log.CreatedAt.AfterEquals(start),
			db.Log.Not(
				db.Log.CreatedAt.AfterEquals(end),
			),
		).Exec(ctx)

		if err != nil {
			if err == db.ErrNotFound {
				continue
			}

			log.Printf("RobotWatcherService's been stopped with failure: %+v\n", errors.WithStack(err))
			break
		}

		for userID, count := range ws.countUserPings(logs) {
			if count < RobotPingsCount {
				continue
			}

			// Определить пользователя как робот
			if _, err := ws.client.User.FindOne(
				db.User.ID.Equals(userID),
			).Update(
				db.User.IsRobot.Set(true),
			).Exec(ctx); err != nil {
				log.Printf("%+v\n", errors.WithStack(err))
			}
		}
	}
}

// Подсчет кол-ва пингов по каждому ID пользователя
func (ws *RobotWatcherService) countUserPings(logs []db.LogModel) map[int]int {
	counts := map[int]int{}
	for _, logModel := range logs {
		counts[logModel.UserID] += 1
	}
	return counts
}
