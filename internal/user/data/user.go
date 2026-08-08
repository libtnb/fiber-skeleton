package data

import (
	"context"
	"errors"

	"github.com/go-rio/rio"
	"github.com/samber/oops"

	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

type userRepo struct {
	db *rio.DB
}

var (
	usersQuery        = rio.From[biz.User]().Must()
	orderedUsersQuery = rio.From[biz.User]().OrderBy("id").Must()
	userByNameQuery   = rio.From[biz.User]().Where("name = ?").Must()
	userByIDQuery     = rio.From[biz.User]().Where("id = ?").Must()
)

func NewUserRepo(db *rio.DB) biz.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) List(ctx context.Context, page, limit int) ([]*biz.User, int64, error) {
	if page < 1 { // guard the callers that skip HTTP validation
		page = 1
	}
	total, err := usersQuery.Count(ctx, r.db)
	if err != nil {
		return nil, 0, oops.In("user").Wrapf(err, "count users")
	}

	list, err := orderedUsersQuery.
		Offset((page-1)*limit).
		Limit(limit).
		All(ctx, r.db)
	if err != nil {
		return nil, 0, oops.In("user").Wrapf(err, "list users")
	}

	users := make([]*biz.User, len(list))
	for i := range list {
		users[i] = &list[i]
	}

	return users, total, nil
}

func (r *userRepo) Get(ctx context.Context, id uint) (*biz.User, error) {
	user, err := rio.Find[biz.User](ctx, r.db, id)
	if err != nil {
		return nil, oops.In("user").Wrapf(err, "get user %d", id)
	}

	return user, nil
}

func (r *userRepo) ExistsName(ctx context.Context, name string) (bool, error) {
	exists, err := userByNameQuery.Exists(ctx, r.db, name)
	if err != nil {
		return false, oops.In("user").Wrapf(err, "check name %q", name)
	}

	return exists, nil
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) error {
	if err := rio.Insert(ctx, r.db, user); err != nil {
		if errors.Is(err, rio.ErrDuplicateKey) {
			return biz.ErrNameTaken
		}
		return oops.In("user").Wrapf(err, "create user")
	}

	return nil
}

// Update is read-modify-write: keeps CreatedAt, reports missing rows.
func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	existing, err := rio.Find[biz.User](ctx, r.db, user.ID)
	if err != nil {
		return nil, oops.In("user").Wrapf(err, "get user %d", user.ID)
	}

	existing.Name = user.Name
	if err := rio.Update(ctx, r.db, existing); err != nil {
		return nil, oops.In("user").Wrapf(err, "update user %d", user.ID)
	}

	return existing, nil
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	n, err := userByIDQuery.DeleteAll(ctx, r.db, id)
	if err != nil {
		return oops.In("user").Wrapf(err, "delete user %d", id)
	}
	if n == 0 {
		return oops.In("user").Wrap(rio.ErrNotFound)
	}

	return nil
}
