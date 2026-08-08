package data

import (
	"context"
	"errors"

	"github.com/go-rio/rio"

	orderbiz "github.com/libtnb/fiber-skeleton/internal/order/biz"
	userbiz "github.com/libtnb/fiber-skeleton/internal/user/biz"
)

// users adapts the user module's public usecase to the Users port; replace
// this file with an RPC client to split user into its own service.
type users struct {
	uc *userbiz.UserUsecase
}

func NewUsers(uc *userbiz.UserUsecase) orderbiz.Users {
	return &users{uc: uc}
}

func (u *users) Exists(ctx context.Context, id uint) (bool, error) {
	if _, err := u.uc.Get(ctx, id); err != nil {
		if errors.Is(err, rio.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
