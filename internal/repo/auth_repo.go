package repo

import (
	"gitlab.com/sample_projects/technonext-assessment/internal/model"
	"gitlab.com/sample_projects/technonext-assessment/internal/pg"
)

type AuthRepo struct{ db *pg.DB }

func NewAuthRepo(db *pg.DB) *AuthRepo { return &AuthRepo{db: db} }

func (r *AuthRepo) FindUserByUsername(username, pass string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.Model(u).Where("username = ?", username).Select(); err != nil {
		return nil, err
	}
	isAuthenticated, err := CheckPasswordHash(pass, u.Password)
	if err != nil {
		return nil, err
	}
	if !isAuthenticated {
		return nil, err
	}
	return u, nil
}
