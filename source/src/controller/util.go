package controller

import (
	"data-management/src/model"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetMetadataResponse(ctx *gin.Context, startTime time.Time, resp *model.Response) {
	code := http.StatusOK
	resp.Metadata.Status = true
	if resp.Metadata.Message == "" {
		resp.Metadata.Message = "OK"
	}

	if resp.Metadata.Message != "OK" {
		code = http.StatusBadRequest
		resp.Metadata.Status = false
	}

	resp.Metadata.TimeExecution = time.Since(startTime).String()
	m, _ := json.Marshal(resp)
	log.Println(string(m))
	log.Println(code, resp)
	ctx.JSON(code, resp)
}
