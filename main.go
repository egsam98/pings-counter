package main

import (
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"

	"addreality/controllers"
	"addreality/db"
	j "addreality/jobs"
)

func main() {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	initJobs(client)
	r := initRouter(client)
	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

func initJobs(client *db.PrismaClient) {
	var jobs = []j.Job{
		j.NewRobotWatcherService(client),
	}

	for _, job := range jobs {
		go job.Run()
	}
}

func initRouter(client *db.PrismaClient) *fasthttprouter.Router {
	r := fasthttprouter.New()
	r.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		log.Printf("%+v\n", i)
		ctx.SetContentType("application/json")
		ctx.Error("{\"error\": \"Internal error occured\"}", http.StatusInternalServerError)
	}

	controllers.InitLogController(r, client)
	return r
}
