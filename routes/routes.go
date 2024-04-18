package routes

import (
	"net/http"
	"organization/controller"
	"organization/middleware"

	"github.com/go-chi/chi"
)

func InitializeRouter(controllers *controller.UserController) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.HandleCORS)
		r.Use(middleware.Authentication)
		r.Route("/organization", func(r chi.Router) {

			r.Post("/", controllers.CreateOrganization)
			r.Put("/", controllers.UpdateOrganizationDetails)

			r.Route("/members", func(r chi.Router) {

				r.Delete("/", controllers.RemoveMemberFromOrganization)
				r.Delete("/leave", controllers.LeaveOrganization)

				r.Route("/role", func(r chi.Router) {
					r.Put("/", controllers.UpdateMemberRole)
					r.Put("/owner", controllers.TransferOwnership)
				})

				r.Route("/invitation", func(r chi.Router) {
					r.Post("/", controllers.InvitationToOrganization)
					r.Get("/", controllers.TrackAllInvitations)
					r.Delete("/", controllers.RespondToInvitation)
				})

			})

		})

		r.Route("/member", func(r chi.Router) {
			r.Get("/organizations",controllers.FetchAllOrganizationDetailsOfCurrentUser)
		})

		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(405)
			w.Write([]byte("wrong method"))
		})
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("route does not exist"))
		})

	})

	return router

}
