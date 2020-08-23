package account

import "context"

// User model
type Users1 struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Repository connects to database
type Repository interface {
	CreateUser(ctx context.Context, user Users1) error
	GetUser(ctx context.Context, id string) (string, error)
}
