package grpcclient

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"my-app/module/product/query"
	"my-app/proto/category"
)

type catGRPCClient struct {
	client category.CategoryClient
}

func NewCatGRPCClient(client category.CategoryClient) *catGRPCClient {
	return &catGRPCClient{client: client}
}

func (c *catGRPCClient) FindWithIds(ctx context.Context, ids []uuid.UUID) ([]query.CategoryDTO, error) {
	catIds := make([]string, len(ids))

	for i := range catIds {
		catIds[i] = ids[i].String()
	}

	resp, err := c.client.GetCategoriesByIds(ctx, &category.GetCatIdsRequest{Ids: catIds})

	if err != nil {
		return nil, errors.New("cannot get categories")
	}

	result := make([]query.CategoryDTO, len(resp.Data))

	for i, dto := range resp.Data {
		result[i] = query.CategoryDTO{
			Id:    uuid.MustParse(dto.Id),
			Title: dto.Title,
		}
	}

	return result, nil
}
