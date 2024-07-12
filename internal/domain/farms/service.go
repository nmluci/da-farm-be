package farms

import (
	"context"

	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/nmluci/da-farm-be/internal/core/httpres"
	"github.com/rs/zerolog"
)

// FarmService contains public API available to be interacted with
type FarmService interface {
	GetAll(context.Context, *FarmRequestQuery) (*ListFarmResponse, error)
	GetOne(context.Context, *FarmRequestQuery) (*FarmResponse, error)
	Create(context.Context, *FarmPayload) error
	Update(context.Context, *FarmPayload) error
	Delete(context.Context, *FarmRequestQuery) error
}

type farmService struct {
	repo FarmRepository
}

// NewService return an instance of FarmService containing available usecases
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

	if params.Limit >= 100 || params.Limit <= 0 {
		repoParams.Limit = 100
	}

	if params.Page <= 0 {
		repoParams.Page = 1
	}

	res = &ListFarmResponse{
		Farms: []*FarmResponse{},
		Meta: httpres.ListPagination{
			Limit:     repoParams.Limit,
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
	res.Meta.TotalPage = count / repoParams.Limit

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

func (svc *farmService) GetOne(ctx context.Context, params *FarmRequestQuery) (res *FarmResponse, err error) {
	logger := zerolog.Ctx(ctx)

	repoParams := &farmQuery{ID: params.ID}

	farm, err := svc.repo.GetOne(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if farm == nil {
		return nil, errs.ErrNotFound
	}

	res = &FarmResponse{
		ID:   farm.ID,
		Name: farm.Name,
	}

	return
}

func (svc *farmService) Create(ctx context.Context, payload *FarmPayload) (err error) {
	logger := zerolog.Ctx(ctx)

	data := &FarmType{
		Name: payload.Name,
	}

	err = svc.repo.Store(ctx, data)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (svc *farmService) Update(ctx context.Context, payload *FarmPayload) (err error) {
	logger := zerolog.Ctx(ctx)

	data := &FarmType{
		ID:   payload.ID,
		Name: payload.Name,
	}

	err = svc.repo.Upsert(ctx, data)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (svc *farmService) Delete(ctx context.Context, params *FarmRequestQuery) (err error) {
	logger := zerolog.Ctx(ctx)

	repoParams := &farmQuery{
		ID: params.ID,
	}

	err = svc.repo.Delete(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
