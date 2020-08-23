package account

import "context"

// Service is exposed to the external traffic
type Service interface {
	CreateUser(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, id string) (string, error)
}
