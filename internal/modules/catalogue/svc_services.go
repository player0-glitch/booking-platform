package catalogue

type servicesRepoInterface interface {
}
type ServicesService struct {
	repo servicesRepoInterface
}

func NewServicesService(servicesRepo servicesRepoInterface) *ServicesService {
	return &ServicesService{
		repo: servicesRepo,
	}
}
