package catalogue

import "net/http"

type HandlerParams struct {
}
type Handler struct {
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{}
}

// POST /admin/orgs
// platform-control
func (h *Handler) CreateOrg(w http.ResponseWriter, r *http.Request) {

}

// GET /orgs/{org_id}
// platform-control
func (h *Handler) GetOrganisationById(w http.ResponseWriter, r *http.Request) {

}

// GET /orgs/{org_id}/services
// PATCH /orgs/{org_id}/services/{svc_id}
// org admin and configured platform wide and members
func (h *Handler) GetOrganisationServices(w http.ResponseWriter, r *http.Request) {

}

//POST / PATCH /orgs/{orgs_id}/services/{svc_1}
// org admin and configured platform wide and members
