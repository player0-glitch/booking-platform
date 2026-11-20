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

/*
snippet on how to get organisationID from the request
r.Route("/organisations/{organisationID}", func(r chi.Router) {
    r.With(
        authorisationMiddleware.RequirePermission(
            authorisation.PermissionBookingCreate,
            OrganisationIDFromURL,
        ),
    ).Post("/bookings", bookingHandler.Create)
})

func OrganisationIDFromURL(r *http.Request) (*uuid.UUID, error) {
    id := chi.URLParam(r, "organisationID")

    if id == "" {
        return nil, errors.New("organisation ID is required")
    }

    organisationID, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid organisation ID")
    }

    return &organisationID, nil
}
*/
