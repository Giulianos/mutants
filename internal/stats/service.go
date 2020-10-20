package stats

type Service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return Service{
		repo: repository,
	}
}

func (s Service) NotifyVerification(verification DNAVerification) error {
	return s.repo.Persist(verification)
}

func (s Service) CountVerifications() (mutants, humans int64, ratio float64, err error) {
	mutants, err = s.repo.CountByResult(true)
	if err != nil {
		return 0, 0, 0, err
	}
	noMutants, err := s.repo.CountByResult(false)
	if err != nil {
		return 0, 0, 0, err
	}

	humans = mutants + noMutants
	// avoid zero division
	if humans == 0 {
		ratio = 0
	} else {
		ratio = float64(mutants)/float64(humans)
	}

	return
}