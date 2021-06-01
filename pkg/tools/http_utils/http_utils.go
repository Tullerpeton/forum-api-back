package http_utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

func SetJSONResponse(ctx *fasthttp.RequestCtx, body interface{}, statusCode int) {
	ctx.SetContentType("application/json")

	result, err := json.Marshal(body)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		if _, err := ctx.Write([]byte("{\"error\": \"can't marshal body\"}")); err != nil {
			log.Fatal(err)
		}
		return
	}
	ctx.SetStatusCode(statusCode)
	if _, err := ctx.Write(result); err != nil {
		log.Fatal(err)
	}
}
