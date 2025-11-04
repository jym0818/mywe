package failover

import (
	"context"
	"sync/atomic"

	"github.com/jym0818/mywe/internal/service/sms"
)

type FailoverService struct {
	svcs []sms.Service
	idx  *atomic.Int64
}

func (f *FailoverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := f.idx.Add(1)
	length := int64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return nil
		default:
			//记录日志和监控
		}
	}
	return errors.New("failover service not available")

}

func NewFailoverService(svcs []sms.Service) sms.Service {
	return &FailoverService{
		svcs: svcs,
		idx:  &atomic.Int64{},
	}
}
