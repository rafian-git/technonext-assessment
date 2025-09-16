package repo

import (
	"gitlab.com/sample_projects/technonext-assessment/internal/model"
	"gitlab.com/sample_projects/technonext-assessment/internal/pg"
)

type OrderRepo struct{ db *pg.DB }

func NewOrderRepo(db *pg.DB) *OrderRepo { return &OrderRepo{db: db} }

func (r *OrderRepo) Insert(o *model.Order) error {
	_, err := r.db.Model(o).Insert()
	return err
}

func (r *OrderRepo) FindByConsignmentID(id string) (*model.Order, error) {
	o := &model.Order{}
	if err := r.db.Model(o).Where("consignment_id = ?", id).Select(); err != nil {
		return nil, err
	}
	return o, nil
}

func (r *OrderRepo) Update(o *model.Order) error {
	_, err := r.db.Model(o).WherePK().Update()
	return err
}

func (r *OrderRepo) List(limit, offset int) ([]model.Order, int, error) {
	var items []model.Order
	if err := r.db.Model(&items).Order("created_at DESC").Limit(limit).Offset(offset).Select(); err != nil {
		return nil, 0, err
	}
	total, err := r.db.Model((*model.Order)(nil)).Count()
	return items, total, err
}
