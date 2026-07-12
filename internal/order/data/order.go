package data

import (
	"context"

	"github.com/go-rio/rio"
	"github.com/samber/do/v2"
	"github.com/samber/oops"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
)

type orderRepo struct {
	db *rio.DB
}

func NewOrderRepo(i do.Injector) (biz.OrderRepo, error) {
	return &orderRepo{
		db: do.MustInvoke[*rio.DB](i),
	}, nil
}

func (r *orderRepo) List(ctx context.Context, page, limit uint) ([]*biz.Order, int64, error) {
	total, err := rio.From[biz.Order]().Count(ctx, r.db)
	if err != nil {
		return nil, 0, oops.In("order").Wrapf(err, "count orders")
	}

	list, err := rio.From[biz.Order]().
		Offset(int((page-1)*limit)).
		Limit(int(limit)).
		OrderBy("id").
		All(ctx, r.db)
	if err != nil {
		return nil, 0, oops.In("order").Wrapf(err, "list orders")
	}

	orders := make([]*biz.Order, len(list))
	for i := range list {
		orders[i] = &list[i]
	}

	return orders, total, nil
}

func (r *orderRepo) Get(ctx context.Context, id uint) (*biz.Order, error) {
	order, err := rio.Find[biz.Order](ctx, r.db, id)
	if err != nil {
		return nil, oops.In("order").Wrapf(err, "get order %d", id)
	}

	return order, nil
}

func (r *orderRepo) Create(ctx context.Context, order *biz.Order) error {
	if err := rio.Insert(ctx, r.db, order); err != nil {
		return oops.In("order").Wrapf(err, "create order")
	}

	return nil
}

func (r *orderRepo) Delete(ctx context.Context, id uint) error {
	n, err := rio.From[biz.Order]().Where("id = ?", id).DeleteAll(ctx, r.db)
	if err != nil {
		return oops.In("order").Wrapf(err, "delete order %d", id)
	}
	if n == 0 {
		return oops.In("order").Wrap(rio.ErrNotFound)
	}

	return nil
}
