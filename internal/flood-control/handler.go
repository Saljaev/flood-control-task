package flood_control

import (
	"context"
	"task/internal/models"
	"time"
)

type (
	UserRate interface {
		Add(ctx context.Context, u models.UserRate) (int64, error)
		Increment(ctx context.Context, ID int64) (int64, error)
		SetExpires(ctx context.Context, ID int64, duration time.Duration) error
		CheckExists(ctx context.Context, ID int64) (int64, error)
	}

	// FloodHandler интерфейс, который нужно реализовать.
	// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
	FloodControl interface {
		// Check возвращает false если достигнут лимит максимально разрешенного
		// кол-ва запросов согласно заданным правилам флуд контроля.
		Check(ctx context.Context, userID int64) (bool, error)
	}

	Handler struct {
		FloodControl
	}

	FloodHandler struct {
		userRate     UserRate
		requestCount int64
		interval     time.Duration
	}
)

func NewHandler(userRate UserRate, requestCount int64, interval time.Duration) *Handler {
	return &Handler{FloodHandler{
		userRate:     userRate,
		requestCount: requestCount,
		interval:     interval,
	}}
}
