package service

import (
	"github.com/blevesearch/bleve"
	"github.com/henkburgstra/kamibasami/node"
)

type Service struct {
	index    bleve.Index
	nodeRepo node.INodeRepo
}

func (s *Service) NodeRepo() node.INodeRepo {
	return s.nodeRepo
}

func (s *Service) SetNodeRepo(r node.INodeRepo) {
	s.nodeRepo = r
}

func (s *Service) Index() bleve.Index {
	return s.index
}

func (s *Service) SetIndex(index bleve.Index) {
	s.index = index
}

func NewService() *Service {
	s := new(Service)
	return s
}
