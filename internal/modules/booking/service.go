package booking

import "booking-platform/internal/database"

type repoInterface interface {
}
type Service struct {
	repo repoInterface
	trx  database.TransactorInterface
}

func NewService(bookingRepo repoInterface, transactor database.TransactorInterface) *Service {
	return &Service{
		repo: bookingRepo,
		trx:  transactor,
	}
}
