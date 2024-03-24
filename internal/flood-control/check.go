package flood_control

import (
	"context"
	"fmt"
	"task/internal/models"
)

func (h FloodHandler) Check(ctx context.Context, userID int64) (bool, error) {
	const op = "FloodHandler - Check"

	exists, err := h.userRate.CheckExists(ctx, userID)
	if err != nil {
		return true, fmt.Errorf("%s - h.userRate.Exists: %w", op, err)
	}

	user := models.UserRate{
		ID:      userID,
		Request: 1,
	}

	if exists == 0 {
		_, err = h.userRate.Add(ctx, user)
		if err != nil {
			return true, fmt.Errorf("%s - h.userRate.Add: %w", op, err)
		}

		err = h.userRate.SetExpires(ctx, userID, h.interval)
		if err != nil {
			return true, fmt.Errorf("%s - h.userRate.SetExpires: %w", op, err)
		}
		return true, nil
	}
	result, err := h.userRate.Increment(ctx, userID)
	if err != nil {
		return true, fmt.Errorf("%s - h.userRate.Increment: %w", op, err)
	}

	return !(result > h.requestCount), nil
}
