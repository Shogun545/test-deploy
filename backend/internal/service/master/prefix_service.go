package master

import (
	"backend/internal/app/dto"
	"backend/internal/app/repository"
)

type PrefixService struct {
	repo repository.PrefixRepository
}

func NewPrefixService(repo repository.PrefixRepository) *PrefixService {
	return &PrefixService{repo: repo}
}

func (s *PrefixService) GetAllPrefixes() ([]dto.PrefixResponse, error) {
	prefixes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	res := make([]dto.PrefixResponse, 0, len(prefixes))
	for _, p := range prefixes {
		res = append(res, dto.PrefixResponse{
			ID:     p.ID,
			Prefix: p.Prefix,
		})
	}

	return res, nil
}
