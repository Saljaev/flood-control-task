package usecase

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	flood_control "task/internal/flood-control"
	"task/internal/models"
	"time"
)

type UserRateUseCase struct {
	*redis.Client
}

var _ flood_control.UserRate = (*UserRateUseCase)(nil)

func NewUserRateUseCase(client *redis.Client) *UserRateUseCase {
	return &UserRateUseCase{client}
}

func (ur *UserRateUseCase) Add(ctx context.Context, u models.UserRate) (int64, error) {
	const op = "UserRateRepo - Add"

	result, err := ur.ZAdd(ctx, string(u.ID), redis.Z{1, u.Request}).Result()
	if err != nil {
		return 0, fmt.Errorf("%s - ur.ZAdd: %w", op, err)
	}

	return result, nil
}

func (ur *UserRateUseCase) Increment(ctx context.Context, ID int64) (int64, error) {
	const op = "UserRateRepo - Increment"

	result, err := ur.ZIncrBy(ctx, string(ID), 1, "1").Result()
	if err != nil {
		return 0, fmt.Errorf("%s - ur.Incryby: %w", op, err)
	}

	return int64(result), nil
}

func (ur *UserRateUseCase) SetExpires(ctx context.Context, ID int64, duration time.Duration) error {
	const op = "UserRateRepo - Increment"

	err := ur.Expire(ctx, string(ID), duration).Err()
	if err != nil {
		return fmt.Errorf("%s - ur.Expire: %w", op, err)
	}

	return nil
}

func (ur *UserRateUseCase) CheckExists(ctx context.Context, ID int64) (int64, error) {
	const op = "UserRateRepo - Exists"

	result, err := ur.Exists(ctx, string(ID)).Result()
	if err != nil {
		return 0, fmt.Errorf("%s - ur.Exists: %w", op, err)
	}

	return result, nil
}
