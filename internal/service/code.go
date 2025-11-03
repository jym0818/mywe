package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jym0818/mywe/internal/repository"
	"github.com/jym0818/mywe/internal/service/sms"
)

const codeTplId = "2367159"

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = repository.ErrUnknownForCode
)

type codeService struct {
	svc  sms.Service
	repo repository.CodeRepository
}

func (c *codeService) Send(ctx context.Context, biz, phone string) error {
	code := c.generateCode()
	//存储到redis中
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//发送验证码
	err = c.svc.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {

		//记录日志就可以了
	}
	return nil
}

func (c *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return c.repo.Verify(ctx, biz, phone, inputCode)
}

func (c *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}
func NewCodeService(svc sms.Service, repo repository.CodeRepository) CodeService {
	return &codeService{
		svc:  svc,
		repo: repo,
	}
}
