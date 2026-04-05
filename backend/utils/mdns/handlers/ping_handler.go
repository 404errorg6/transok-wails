package handlers

import (
	"fmt"
	"transok/backend/consts"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) GetType() string {
	return "PING"
}

func (h *PingHandler) Handle(payload consts.DiscoverPayload) {
	// Logic for handling messages
	fmt.Println("PingHandler: ", "PONG")
}
