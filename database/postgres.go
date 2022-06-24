package database

import (
	"agchavez/go/rest-ws/models"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (p *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(
		ctx, "INSERT INTO app_user(first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
		user.FirstName, user.LastName, user.Email, user.Password,
	)

	return err

}

// Get all users from database
func (p *PostgresRepository) GetUsers(ctx context.Context, params models.ParamsQuery) ([]models.User, error) {
	fmt.Println("GetUsers", params)
	rows, err := p.db.QueryContext(ctx, "SELECT id, first_name, last_name, email, created_at FROM app_user ORDER BY id DESC LIMIT $1 OFFSET $2", params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Error closing rows: ", err)
		}
	}()

	var users []models.User
	for rows.Next() {
		user, err := p.parseRowToUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

// Get user by id from database
func (p *PostgresRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	row, err := p.db.QueryContext(ctx, "SELECT id, first_name, last_name, email, password, created_at FROM app_user WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := row.Close(); err != nil {
			log.Println("Error closing row: ", err)
		}
	}()

	var user models.User
	for row.Next() {
		err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil

}

// Get user by email from database
func (p *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := p.db.QueryContext(ctx, "SELECT id, first_name, last_name, email, password, created_at FROM app_user WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := row.Close(); err != nil {
			log.Println("Error closing row: ", err)
		}
	}()
	var user models.User
	for row.Next() {
		err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// Funtion parse row to user
func (p *PostgresRepository) parseRowToUser(row *sql.Rows) (*models.User, error) {
	var user models.User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Close database connection
func (p *PostgresRepository) Close() error {
	return p.db.Close()
}
