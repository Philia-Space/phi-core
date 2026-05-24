package domain

import (
	"context"
)

// Repository is the base interface for all aggregate repositories.
type Repository[T any] interface {
	FindByID(ctx context.Context, id string) (*T, error)
	Save(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}

// ListOptions provides pagination and filtering options.
type ListOptions struct {
	Page     int
	PageSize int
	SortBy   string
	SortDesc bool
}

// ListResult wraps paginated results.
type ListResult[T any] struct {
	Items      []T
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
