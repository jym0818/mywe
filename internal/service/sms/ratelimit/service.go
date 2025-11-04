package ratelimit

import (
	"context"
	"fmt"

	"github.com/jym0818/mywe/internal/service/sms"
	"github.com/jym0818/mywe/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}
func (s RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//你在这里加上代码  新特性
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return errLimited
	}
	if limited {
		return errLimited
	}

	err = s.svc.Send(ctx, tpl, args, numbers...)
	//你也可以在这里代码  新特性
	return err
}
