package booking

type svcInterface interface {
}

type Handler struct {
	svc svcInterface
}

func NewHandler(bookingSvc svcInterface) *Handler {
	return &Handler{
		svc: bookingSvc,
	}
}
