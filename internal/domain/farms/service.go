package farms

import (
	"context"

	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/nmluci/da-farm-be/internal/core/httpres"
	"github.com/rs/zerolog"
)

type FarmService interface {
	GetAll(context.Context, *FarmRequestQuery) (*ListFarmResponse, error)
	// 	GetOne(context.Context, *FarmQuery) (*FarmResponse, error)
	// 	Create(context.Context, *FarmQuery, *FarmPayload) error
	// 	Update(context.Context, *FarmQuery, *FarmPayload) error
	// 	Delete(context.Context, *FarmQuery) error
}

type farmService struct {
	repo FarmRepository
}

func NewService(repo FarmRepository) FarmService {
	return &farmService{repo: repo}
}

func (svc *farmService) GetAll(ctx context.Context, params *FarmRequestQuery) (res *ListFarmResponse, err error) {
	logger := zerolog.Ctx(ctx)

	repoParams := &farmQuery{
		Keyword: params.Keyword,
		Limit:   params.Limit,
		Page:    params.Page,
	}

	if params.Limit >= 100 {
		repoParams.Limit = 100
	}

	res = &ListFarmResponse{
		Farms: []*FarmResponse{},
		Meta: httpres.ListPagination{
			Limit:     params.Limit,
			Page:      0,
			TotalPage: 0,
		},
	}

	count, err := svc.repo.Count(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if count == 0 {
		return nil, errs.ErrNotFound
	}
	res.Meta.TotalPage = count / params.Limit

	farms, err := svc.repo.GetAll(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, farm := range farms {
		res.Farms = append(res.Farms, &FarmResponse{
			ID:   farm.ID,
			Name: farm.Name,
		})
	}

	return
}
