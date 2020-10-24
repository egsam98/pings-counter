package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

type JSON map[string]interface{}

// Ответ сервера в формате JSON
func respondJSON(ctx *fasthttp.RequestCtx, body JSON) {
	data, err := json.Marshal(body)
	if err != nil {
		respondErrorJSON(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Success("application/json", data)
}

// Ответ сервера ошибкой в формате JSON {"error": "ТЕКСТ ОШИБКИ"}
func respondErrorJSON(ctx *fasthttp.RequestCtx, status int, err error) {
	msg := ""
	if status >= 500 {
		log.Printf("%+v\n", err)
		msg = "{\"error\": \"Internal error occured\"}"
	} else {
		msg = "{\"error\": \"" + strings.ReplaceAll(err.Error(), "\"", "\\\"") + "\"}"
	}
	ctx.Error(msg, status)
	ctx.SetContentType("application/json")
}
