package analyzer

import "github.com/Giulianos/mutants/internal/stats"

type EventPublisher interface {
	PublishVerification(verification stats.DNAVerification)
}
