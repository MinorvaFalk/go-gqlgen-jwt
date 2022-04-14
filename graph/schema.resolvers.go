package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/MinorvaFalk/go-gqlgen-jwt/graph/generated"
	"github.com/MinorvaFalk/go-gqlgen-jwt/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	todo := &model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", rand.Int()),
		UserID: input.UserID, // fix this line
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.ResponseString, error) {
	if username == "jon" && password == "test" {
		token, err := r.Middleware.SignKey(username)
		if err != nil {
			return nil, &gqlerror.Error{
				Path:    graphql.GetPath(ctx),
				Message: "Internal Server Error",
				Extensions: map[string]interface{}{
					"code": http.StatusInternalServerError,
				},
			}
		}

		return &model.ResponseString{String: token}, nil
	}

	return nil, errors.New("invalid user")
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	// gc, err := r.GinContextFromContext(ctx)

	return r.todos, nil
}

func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
