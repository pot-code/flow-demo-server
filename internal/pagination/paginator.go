package pagination

import (
	"context"
	"fmt"
)

type EntQuery[Q any, V any] interface {
	Offset(int) Q
	Limit(int) Q
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) (V, error)
}

func EntPaginator[Q EntQuery[Q, V], V any](
	ctx context.Context,
	query Q,
	pagination *Pagination,
	dataType V,
) (V, int, error) {
	count, err := query.Count(ctx)
	if err != nil {
		return dataType, 0, fmt.Errorf("query count: %w", err)
	}

	data, err := query.
		Offset((pagination.Page - 1) * pagination.PageSize).
		Limit(pagination.PageSize).
		All(ctx)
	if err != nil {
		return dataType, 0, fmt.Errorf("query pagination data: %w", err)
	}

	return data, count, nil
}
