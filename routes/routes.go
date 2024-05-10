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

		r.Route("/auth/microsoft", func(r chi.Router) {
			r.Get("/", controllers.MicrosoftAuth)
			r.Get("/tokens", controllers.GetMicrosoftTokens)
		})
		
		r.Route("/organization", func(r chi.Router) {
			r.Use(middleware.Authentication)
			r.Post("/", controllers.CreateOrganization)
			r.Put("/", controllers.UpdateOrganization)
			r.Post("/delete/otp", controllers.OTPForDeleteOrganization)

			r.Route("/members", func(r chi.Router) {

				r.Delete("/", controllers.RemoveMemberFromOrganization)
				r.Put("/transfer-ownership", controllers.TransferOwnership)

				r.Route("/role", func(r chi.Router) {
					r.Put("/", controllers.UpdateMemberRole)
				})

			})

			r.Route("/invitation", func(r chi.Router) {
				r.Post("/", controllers.InvitationToOrganization)
			})

			r.Route("/{organization-id}/member", func(r chi.Router) {
				r.Delete("/leave", controllers.LeaveOrganization)
			})

		})

		r.Route("/user", func(r chi.Router) {
			r.Use(middleware.Authentication)
			r.Get("/organizations", controllers.FetchAllOrganizationDetailsOfCurrentUser)
			r.Get("/organization/{organization-id}", controllers.FetchOrganizationDetailsOfCurrentUser)

			r.Route("/invitations", func(r chi.Router) {
				r.Get("/", controllers.TrackAllInvitations)
				r.Post("/", controllers.RespondToInvitation)
			})

		})

		r.Route("/internal", func(r chi.Router) {

			r.Get("/jwt", controllers.GetJWTForOragnizationService)
			r.Route("/users", func(r chi.Router) {

				r.Route("/organizations", func(r chi.Router) {
					r.Post("/", controllers.FetchOragnizationListOfUsers)
				})

			})

			r.Delete("/organization/{organization-id}", controllers.DeleteOrganization)

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
