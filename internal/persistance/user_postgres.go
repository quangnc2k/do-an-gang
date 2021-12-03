package persistance

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"time"
)

type UserRepositorySQL struct {
	connection *pgxpool.Pool
}

func NewUserSQLRepo(ctx context.Context, conn *pgxpool.Pool) (UserRepository, error) {
	if conn == nil {
		return nil, errors.New("invalid sql connection")
	}
	return &UserRepositorySQL{connection: conn}, nil
}

func (ds *UserRepositorySQL) List(ctx context.Context) (users []model.User, err error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users`
	rows, err := ds.connection.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println(err)
			return
		}

		users = append(users, user)
	}

	return
}

func (ds *UserRepositorySQL) FindOneByEmail(ctx context.Context, email string) (*model.User, error) {
	var id, password string
	var scopes []string
	var createdAt, updatedAt time.Time

	query := "SELECT id, email, password, scopes, created_at, updated_at FROM users WHERE email = $1"
	err := ds.connection.QueryRow(ctx, query, email).Scan(&id, &email, &password, &scopes, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (ds *UserRepositorySQL) FindOneByID(ctx context.Context, id string) (*model.User, error) {
	var _id, email, password string
	var scopes []string
	var createdAt, updatedAt time.Time

	query := "SELECT id, email, password, scopes, created_at, updated_at FROM users WHERE id = $1"
	err := ds.connection.QueryRow(ctx, query, id).Scan(&_id, &email, &password, &scopes, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        _id,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (ds *UserRepositorySQL) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := ds.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (ds *UserRepositorySQL) Create(ctx context.Context, id, email, pw string, scopes []string) (_id string, err error) {
	re := regexp.MustCompile("^(([^<>()[\\]\\\\.,;:\\s@\"]+(\\.[^<>()[\\]\\\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$")
	if !re.MatchString(email) {
		return "", ErrInvalidEmail
	}

	user, err := ds.FindOneByEmail(ctx, email)
	if err != pgx.ErrNoRows && err != nil {
		return "", ErrCannotFindMail
	}
	if user != nil {
		return "", ErrMailAlreadyExist
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	row := ds.connection.QueryRow(
		ctx,
		"INSERT INTO users (id, email, password, scopes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		id,
		email,
		hashedPw,
		scopes,
		time.Now(),
		time.Now())

	err = row.Scan(&_id)
	if err != nil {
		return "", err
	}

	return _id, nil
}

func (ds *UserRepositorySQL) UpdateValidate(ctx context.Context, id, oldPassword string) (err error) {
	user, err := ds.FindOneByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if oldPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
		if err != nil {
			return ErrInvalidCredentials
		}
	}
	return
}

func (ds *UserRepositorySQL) Update(ctx context.Context, id, password string, scopes []string) (_id string, err error) {
	_, err = ds.FindOneByID(ctx, id)
	if err != nil {
		return "", ErrUserNotFound
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	err = ds.connection.QueryRow(
		ctx,
		"UPDATE users SET password = $1, scopes=$2, updated_at=$3 WHERE id = $4 RETURNING id;",
		hashedPw,
		scopes,
		time.Now(),
		id).Scan(&_id)
	if err != nil {
		return "", err
	}

	return _id, nil
}

func (ds *UserRepositorySQL) Delete(ctx context.Context, id string) (rows int, err error) {
	_, err = ds.FindOneByID(ctx, id)
	if err != nil {
		return 0, ErrUserNotFound
	}

	query := "DELETE FROM users WHERE id = $1;"
	tag, err := ds.connection.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return -1, err
	}

	return int(tag.RowsAffected()), nil
}
