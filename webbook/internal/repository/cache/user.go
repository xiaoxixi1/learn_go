package cache

import (
	"context"
	"project_go/webbook/internal/domain"
)

type UserCache interface {
	Get(cxt context.Context, userid int64) (domain.User, error)
	Set(cxt context.Context, us domain.User) error
}
