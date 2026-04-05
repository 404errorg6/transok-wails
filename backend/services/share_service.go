package services

import (
	"context"
	"sync"
)

type ShareService struct {
	ctx context.Context
}

var share *ShareService
var shareOnce sync.Once

func Share() *ShareService {
	if share == nil {
		shareOnce.Do(func() {
			share = &ShareService{}
		})
	}
	return share
}

func (s *ShareService) Start(ctx context.Context) {
	s.ctx = ctx
}

/* Set sharing captcha */
func (s *ShareService) SetCaptcha(captcha string) {
	store := Storage()
	store.Set("captcha", captcha)
}

/* Get sharing captcha */
func (s *ShareService) GetCaptcha() string {
	store := Storage()
	captcha, exists := store.Get("captcha")
	if !exists {
		return ""
	}
	return captcha.(string)
}
