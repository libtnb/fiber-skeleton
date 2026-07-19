// Package biz holds the order module's business logic; cross-module needs are
// interfaces defined here (Users), so it never imports another module.
package biz

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/libtnb/fiber-skeleton/internal/pkg/apperr"
	"github.com/libtnb/fiber-skeleton/internal/pkg/event"
)

type Order struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	Amount    int64      `json:"amount"` // minor units, e.g. cents
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `rio:",softdelete" json:"-"`
}

var ErrUserNotFound = errors.New("user not found")

// OrderRepo is the persistence boundary; a missing row is rio.ErrNotFound.
type OrderRepo interface {
	List(ctx context.Context, page, limit int) ([]*Order, int64, error)
	Get(ctx context.Context, id uint) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uint) error
}

// Users is the slice of the user module that orders need; swap its adapter
// for an RPC client and the module splits into a service unchanged.
type Users interface {
	Exists(ctx context.Context, id uint) (bool, error)
}

// OrderPlaced is published after an order is created.
type OrderPlaced struct {
	OrderID uint
	UserID  uint
	Amount  int64
}

func (OrderPlaced) Name() string { return "order.placed" }

// OrderUsecase holds transport-independent order business logic.
type OrderUsecase struct {
	repo  OrderRepo
	users Users
	bus   event.Bus
}

func NewOrderUsecase(repo OrderRepo, users Users, bus event.Bus) *OrderUsecase {
	return &OrderUsecase{
		repo:  repo,
		users: users,
		bus:   bus,
	}
}

func (uc *OrderUsecase) List(ctx context.Context, page, limit int) ([]*Order, int64, error) {
	return uc.repo.List(ctx, page, limit)
}

func (uc *OrderUsecase) Get(ctx context.Context, id uint) (*Order, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *OrderUsecase) Place(ctx context.Context, userID uint, amount int64) (*Order, error) {
	exists, err := uc.users.Exists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperr.Unprocessable("order.user_not_found", "user not found").In("order").Wrap(ErrUserNotFound)
	}

	order := &Order{UserID: userID, Amount: amount}
	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// the order is committed — log and move on; a durable broker needs an outbox
	if err := uc.bus.Publish(ctx, OrderPlaced{OrderID: order.ID, UserID: userID, Amount: amount}); err != nil {
		slog.ErrorContext(ctx, "publish order.placed failed",
			slog.Uint64("order_id", uint64(order.ID)),
			slog.Any("err", err),
		)
	}

	return order, nil
}

func (uc *OrderUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
