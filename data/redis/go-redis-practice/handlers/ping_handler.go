package handlers

import (
	"go-redis-practice/utils"
	"net/http"
)

type PingHandler struct {
	PingUtils *utils.PingUtils
}

func NewPingHandler() *PingHandler {
	pingUtils := utils.NewPingUtils()

	return &PingHandler{
		PingUtils: pingUtils,
	}
}

func (p *PingHandler) DoPing(w http.ResponseWriter, r *http.Request) {
	status := p.PingUtils.Ping()
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(status.String()))
}
