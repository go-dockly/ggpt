package speech

import (
	"github.com/go-dockly/ggpt/pkg/entity"
	"github.com/go-dockly/ggpt/pkg/gpt"
)

type IService interface {
	Listen(session *entity.Session, hotwordFile, detectorFile string) error
	Play(fileName string) (err error)
}

type Service struct {
	ai gpt.IClient
}

func New(ai gpt.IClient) IService {
	return &Service{ai: ai}
}
