package account

import "context"

type Service interface {
	CreateUser(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, id string) (string, error)
}
