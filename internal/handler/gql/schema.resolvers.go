package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/kargotech/go-testapp/gen/graph/generated"
)

func (r *mutationResolver) Health(ctx context.Context) (string, error) {
	return "Healthy", nil
}

func (r *queryResolver) Health(ctx context.Context) (string, error) {
	return "Healthy", nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
