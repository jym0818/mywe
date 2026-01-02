package grpc

import (
	"context"

	intrv1 "github.com/jym0818/mywe/api/proto/gen/intr/v1"
	"github.com/jym0818/mywe/interactive/domain"
	"github.com/jym0818/mywe/interactive/service"
	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}

func (s *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := s.svc.IncrReadCnt(ctx, request.Biz, request.BizId)
	return &intrv1.IncrReadCntResponse{}, err
}

func (s *InteractiveServiceServer) Register(server *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(server, s)
}

func (s *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := s.svc.Like(ctx, request.Biz, request.BizId, request.Uid)
	return &intrv1.LikeResponse{}, err
}
func (s *InteractiveServiceServer) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := s.svc.CancelLike(ctx, request.Biz, request.BizId, request.Uid)
	return &intrv1.CancelLikeResponse{}, err
}
func (s *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := s.svc.Collect(ctx, request.Biz, request.BizId, request.Cid, request.Uid)
	return &intrv1.CollectResponse{}, err
}
func (s *InteractiveServiceServer) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	res, err := s.svc.Get(ctx, request.Biz, request.BizId, request.Uid)
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: s.toDTO(res),
	}, nil
}
func (s *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	if len(request.Ids) == 0 {
		return &intrv1.GetByIdsResponse{}, nil
	}
	data, err := s.svc.GetByIds(ctx, request.Biz, request.Ids)
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*intrv1.Interactive, len(data))
	for k, v := range data {
		res[k] = s.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: res,
	}, nil
}
func (s *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
