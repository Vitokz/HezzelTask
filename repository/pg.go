package repository

import (
	"HezzelTask/config"
	"HezzelTask/models"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"time"
)

const (
	queryTimeout = time.Second * 10
)

type Pg struct {
	Db *pg.DB
}

func Connect(url string) (*pg.DB, error) {
	pgOpts, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	conn := pg.Connect(pgOpts)

	// Test connection
	var ping int
	_, err = conn.QueryOne(pg.Scan(&ping), "SELECT 1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the db")
	}

	return conn, nil
}

func (p *Pg) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, user).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Pg) DeleteUser(ctx context.Context, email string) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, &models.User{}).Where("email = ?",email).Delete()
	if err != nil {
		return  err
	}

	return nil
}

func (p *Pg) UserList(ctx context.Context) (*[]models.User, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	users := make([]models.User,0)
	err := p.Db.ModelContext(timeout, &users).Select()

	if err != nil {
		return  nil,err
	}

	return &users, nil
}

func RunPgMigrations(cfg *config.Config) error {

	if cfg.Pg.PgMigrationPath == "" {
		return nil
	}

	if cfg.Pg.PgUrl == "" {
		return errors.New("No cfg.PgURL provided")
	}

	m, err := migrate.New(
		cfg.Pg.PgMigrationPath,
		cfg.Pg.PgUrl,
	)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
