package service

import (
	"context"

	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewinteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{repo: repo}
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func (svc *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.IncrLike(ctx, biz, bizId, uid)
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.DecrLike(ctx, biz, bizId, uid)
}

func (svc *interactiveService) Collect(ctx context.Context, biz string, bizId, cid, uid int64) error {
	return svc.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	res, err := svc.repo.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		res.Liked, er = svc.repo.Liked(ctx, biz, bizId, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		res.Collected, er = svc.repo.Collected(ctx, biz, bizId, uid)
		return er
	})
	return res, eg.Wait()
}
