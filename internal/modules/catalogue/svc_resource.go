package catalogue

type resourcesInterface interface {
}
type ResourceService struct {
	repo resourcesInterface
}

func NewResourceService(repo resourcesInterface) *ResourceService {
	return &ResourceService{
		repo: repo,
	}
}
