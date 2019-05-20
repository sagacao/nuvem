package gate

import (
	"github.com/gin-gonic/gin"
)

type ApiContext struct {
	Quit  chan struct{}
	Ctext *gin.Context
}

type sapi struct {
	apiHub map[string]*ApiContext //gin.Context
}

func newSAPI() *sapi {
	return &sapi{
		apiHub: make(map[string]*ApiContext),
	}
}

func (s *sapi) push(unionid string, c *ApiContext) {
	s.apiHub[unionid] = c
}

func (s *sapi) pop(unionid string) *ApiContext {
	return s.apiHub[unionid]
}
