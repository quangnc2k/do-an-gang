package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	FindOneByEmail(ctx context.Context, email string) (*model.User, error)
	FindOneByID(ctx context.Context, id string) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, error)
	Create(ctx context.Context, id, email, pw string, scopes []string) (string, error)
	Update(ctx context.Context, id, password string) (string, error)
	UpdateValidate(ctx context.Context, id, oldPassword string) error
	Delete(ctx context.Context, id string) (int, error)
}
