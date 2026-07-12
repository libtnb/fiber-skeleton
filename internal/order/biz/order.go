// Package biz holds the order module's business logic. Its cross-module needs
// are expressed as interfaces it defines here (Users); the wiring provides the
// implementation, so this package never imports another module.
package biz

import (
	"context"
	"errors"
	"time"

	"github.com/samber/oops"

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
	List(ctx context.Context, page, limit uint) ([]*Order, int64, error)
	Get(ctx context.Context, id uint) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uint) error
}

// Users is the slice of the user module that orders need. The order module
// defines it here; the data layer implements it over the user module's public
// usecase, so order never touches the user tables. Swap that implementation
// for an RPC client to split user into its own service — this interface and
// OrderUsecase stay unchanged.
type Users interface {
	Exists(ctx context.Context, id uint) (bool, error)
}

// OrderPlaced is published after an order is created; subscribers react to it
// (confirmation email, analytics, ...) without the order module knowing them.
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

func (uc *OrderUsecase) List(ctx context.Context, page, limit uint) ([]*Order, int64, error) {
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
		return nil, oops.In("order").Code("order.user_not_found").Public("user not found").Wrap(ErrUserNotFound)
	}

	order := &Order{UserID: userID, Amount: amount}
	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	_ = uc.bus.Publish(ctx, OrderPlaced{OrderID: order.ID, UserID: userID, Amount: amount})

	return order, nil
}

func (uc *OrderUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
