package stats

type Repository interface {
	Persist(verification DNAVerification) error
	CountByResult(value bool) (int64, error)
}
