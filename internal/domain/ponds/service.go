package ponds

import (
	"context"

	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/nmluci/da-farm-be/internal/core/httpres"
	"github.com/rs/zerolog"
)

// PondService contains public API available to be interacted with
type PondService interface {
	GetAll(context.Context, *PondRequestQuery) (*ListPondResponse, error)
	GetOne(context.Context, *PondRequestQuery) (*PondResponse, error)
	Create(context.Context, *PondPayload) error
	Update(context.Context, *PondPayload) error
	Delete(context.Context, *PondRequestQuery) error
}

type pondService struct {
	repo PondRepository
}

// NewService return an instance of PondService containing available usecases
func NewService(repo PondRepository) PondService {
	return &pondService{repo: repo}
}

func (svc *pondService) GetAll(ctx context.Context, params *PondRequestQuery) (res *ListPondResponse, err error) {
	logger := zerolog.Ctx(ctx)

	repoParams := &pondQuery{
		FarmID:  params.FarmID,
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

	res = &ListPondResponse{
		Ponds: []*PondResponse{},
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

	ponds, err := svc.repo.GetAll(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	for _, pond := range ponds {
		res.Ponds = append(res.Ponds, &PondResponse{
			ID:       pond.ID,
			FarmID:   pond.FarmID,
			FarmName: pond.FarmName,
			Name:     pond.Name,
		})
	}

	return
}

func (svc *pondService) GetOne(ctx context.Context, params *PondRequestQuery) (res *PondResponse, err error) {
	logger := zerolog.Ctx(ctx)

	repoParams := &pondQuery{ID: params.ID}

	pond, err := svc.repo.GetOne(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if pond == nil {
		return nil, errs.ErrNotFound
	}

	res = &PondResponse{
		ID:       pond.ID,
		FarmID:   pond.FarmID,
		FarmName: pond.FarmName,
		Name:     pond.Name,
	}

	return
}

func (svc *pondService) Create(ctx context.Context, payload *PondPayload) (err error) {
	logger := zerolog.Ctx(ctx)

	data := &PondType{
		FarmID: payload.FarmID,
		Name:   payload.Name,
	}

	err = svc.repo.Store(ctx, data)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (svc *pondService) Update(ctx context.Context, payload *PondPayload) (err error) {
	logger := zerolog.Ctx(ctx)

	data := &PondType{
		ID:     payload.ID,
		FarmID: payload.FarmID,
		Name:   payload.Name,
	}

	err = svc.repo.Upsert(ctx, data)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}

func (svc *pondService) Delete(ctx context.Context, params *PondRequestQuery) (err error) {
	logger := zerolog.Ctx(ctx)

	repoParams := &pondQuery{
		ID: params.ID,
	}

	err = svc.repo.Delete(ctx, repoParams)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	return
}
