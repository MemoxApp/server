package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"context"
	"fmt"
	"time_speak_server/graph/generated"
	"time_speak_server/src/service/subscribe"
	"time_speak_server/src/service/user"
)

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*user.User, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	getUser, err := r.userSvc.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &getUser, nil
}

// Used is the resolver for the used field.
func (r *userResolver) Used(ctx context.Context, obj *user.User) (int, error) {
	panic(fmt.Errorf("not implemented: Used - used"))
}

// Subscribe is the resolver for the subscribe field.
func (r *userResolver) Subscribe(ctx context.Context, obj *user.User) (*subscribe.Subscribe, error) {
	panic(fmt.Errorf("not implemented: Subscribe - subscribe"))
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
