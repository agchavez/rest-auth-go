package repository

import (
	"agchavez/go/rest-ws/models"
	"context"
)

// UserRepository interface
type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUsers(ctx context.Context, params models.ParamsQuery) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	//UpdateUser(ctx context.Context, user *models.User) error
	Close() error
}

var implementer UserRepository

func SetRepository(repository UserRepository) {
	implementer = repository
}

// Inser new user in database
func InsertUser(ctx context.Context, user *models.User) error {
	return implementer.InsertUser(ctx, user)
}

// Get all users from database
func GetUsers(ctx context.Context, params models.ParamsQuery) ([]models.User, error) {
	return implementer.GetUsers(ctx, models.ParamsQuery{})
}

// Get user by id
func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return implementer.GetUserByID(ctx, id)
}

// Get user by email
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementer.GetUserByEmail(ctx, email)
}
