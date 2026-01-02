package service

import (
	"context"
	"math"
	"time"

	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	intrv1 "github.com/jym0818/mywe/api/proto/gen/intr/v1"
	"github.com/jym0818/mywe/internal/domain"
	"github.com/jym0818/mywe/internal/repository"
)

type RankingService interface {
	TopN(ctx context.Context) error
}
type BatchRankingService struct {
	batchSize int
	n         int
	scoreFunc func(t time.Time, likeCnt int64) float64
	artSvc    ArticleService
	intrSvc   intrv1.InteractiveServiceClient
	repo      repository.RankingRepository
}

func (svc *BatchRankingService) TopN(ctx context.Context) error {
	now := time.Now()
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	topN := queue.NewConcurrentPriorityQueue[Score](svc.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score < dst.score {
			return -1
		} else {
			return 0
		}
	})

	for {
		arts, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return err
		}
		ids := slice.Map[domain.Article, int64](arts, func(idx int, src domain.Article) int64 {
			return src.Id
		})

		intrs, err := svc.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{Biz: "article", Ids: ids})
		if err != nil {
			return err
		}

		for _, art := range arts {
			score := svc.scoreFunc(now, intrs.Intrs[art.Id].LikeCnt)

			er := topN.Enqueue(Score{art, score})
			if er == queue.ErrOutOfCapacity {
				val, _ := topN.Dequeue()
				if val.score < score {
					_ = topN.Enqueue(Score{art, score})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}

		if len(arts) < svc.batchSize || now.Sub(arts[len(arts)-1].Utime).Hours() > 7*24 {
			break
		}

		offset += len(arts)

	}

	res := make([]domain.Article, svc.n)
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			break
		}
		res[i] = val.art
	}
	return svc.repo.ReplaceTopN(ctx, res)

}

func NewBatchRankingService(artSvc ArticleService, intrSvc intrv1.InteractiveServiceClient, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		n:         100,
		batchSize: 100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			sec := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(float64(sec+2), 1.5)
		},
		repo: repo,
	}
}
