package graphql

import (
	"context"

	"github.com/google/uuid"

	userapp "github.com/soliton-go/application/internal/application/user"
	"github.com/soliton-go/application/internal/domain/user"
)

// Resolver holds the dependencies for GraphQL resolvers.
type Resolver struct {
	repo          user.UserRepository
	createHandler *userapp.CreateUserHandler
	getHandler    *userapp.GetUserHandler
	listHandler   *userapp.ListUsersHandler
}

// NewResolver creates a new Resolver.
func NewResolver(
	repo user.UserRepository,
	createHandler *userapp.CreateUserHandler,
	getHandler *userapp.GetUserHandler,
	listHandler *userapp.ListUsersHandler,
) *Resolver {
	return &Resolver{
		repo:          repo,
		createHandler: createHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// User represents the GraphQL User type.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateUserInput represents the input for creating a user.
type CreateUserInput struct {
	Name string `json:"name"`
}

// QueryResolver implements the Query resolvers.
type QueryResolver struct {
	*Resolver
}

// MutationResolver implements the Mutation resolvers.
type MutationResolver struct {
	*Resolver
}

// Query returns the query resolver.
func (r *Resolver) Query() *QueryResolver {
	return &QueryResolver{r}
}

// Mutation returns the mutation resolver.
func (r *Resolver) Mutation() *MutationResolver {
	return &MutationResolver{r}
}

// User resolver for Query.user(id: ID!)
func (r *QueryResolver) User(ctx context.Context, id string) (*User, error) {
	query := userapp.GetUserQuery{ID: id}
	u, err := r.getHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:   string(u.ID),
		Name: u.Name,
	}, nil
}

// Users resolver for Query.users
func (r *QueryResolver) Users(ctx context.Context) ([]*User, error) {
	query := userapp.ListUsersQuery{}
	users, err := r.listHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	var result []*User
	for _, u := range users {
		result = append(result, &User{
			ID:   string(u.ID),
			Name: u.Name,
		})
	}
	return result, nil
}

// CreateUser resolver for Mutation.createUser(input: CreateUserInput!)
func (r *MutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	id := uuid.New().String()

	cmd := userapp.CreateUserCommand{
		ID:   id,
		Name: input.Name,
	}

	u, err := r.createHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:   string(u.ID),
		Name: u.Name,
	}, nil
}
