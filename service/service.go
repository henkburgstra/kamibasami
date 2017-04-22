package service

import (
	"github.com/henkburgstra/kamibasami/node"
)

type Service struct {
	nodeRepo node.INodeRepo
}

func (s *Service) NodeRepo() node.INodeRepo {
	return s.nodeRepo
}

func (s *Service) SetNodeRepo(r node.INodeRepo) {
	s.nodeRepo = r
}

func NewService() *Service {
	s := new(Service)
	return s
}
