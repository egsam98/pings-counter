package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"

	"addreality/db"
	"addreality/services"
)

type LogController struct {
	logService  *services.LogService
	userService *services.UserService
}

func InitLogController(r *fasthttprouter.Router, client *db.PrismaClient) {
	lc := &LogController{
		logService:  services.NewLogService(client),
		userService: services.NewUserService(client),
	}

	r.POST("/", lc.post)
	r.GET("/count", lc.get)
}

func (lc *LogController) post(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	userId, err := strconv.Atoi(string(ctx.QueryArgs().Peek("user_id")))
	if err != nil {
		respondErrorJSON(ctx, http.StatusNotFound, err)
		return
	}

	go lc.log(ctx, userId)
}

func (lc *LogController) get(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	count, err := lc.userService.CountRobots(ctx)
	if err != nil {
		respondErrorJSON(ctx, http.StatusInternalServerError, err)
		return
	}

	respondJSON(ctx, JSON{
		"count": count,
	})
}

func (lc *LogController) log(ctx *fasthttp.RequestCtx, userId int) {
	if err := lc.logService.Log(ctx, userId); err != nil {
		log.Printf("%+v\n", err)
	}
}
