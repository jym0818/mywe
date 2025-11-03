package ioc

import (
	"github.com/jym0818/mywe/internal/service/sms"
	"github.com/jym0818/mywe/internal/service/sms/memory"
)

func InitSMS() sms.Service {
	return memory.NewService()
}
