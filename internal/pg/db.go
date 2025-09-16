package pg

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"gitlab.com/sample_projects/technonext-assessment/internal/model"
)

type DB struct {
	*pg.DB
}

func Connect() (*DB, error) {
	addr := getenv("PG_ADDR", "localhost:5432")
	user := getenv("PG_USER", "root")
	pass := getenv("PG_PASS", "123456")
	name := getenv("PG_DB", "ordersdb")

	opt := &pg.Options{
		Addr:     addr,
		User:     user,
		Password: pass,
		Database: name,
	}

	db := pg.Connect(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pg ping: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	//seeding given user credentials
	if err := seedUser(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*model.User)(nil),
		(*model.Order)(nil),
	}
	for _, m := range models {
		err := db.Model(m).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func seedUser(db *pg.DB) error {
	const username = "01901901901@mailinator.com"
	const password = "321dsaf"
	bytesPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u := &model.User{Username: username, Password: string(bytesPass)}
	_, err = db.Model(u).Where("username = ?", username).SelectOrInsert()
	return err
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
