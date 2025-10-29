package ioc

import "github.com/gin-gonic/gin"

func InitWeb() *gin.Engine {
	s := gin.New()
	return s
}
