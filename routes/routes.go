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

		r.Route("/organization", func(r chi.Router) {
			r.Use(middleware.Authentication)
			r.Post("/", controllers.CreateOrganization)
			r.Put("/", controllers.UpdateOrganizationDetails)
			r.Post("/otp              ", controllers.OTPForDeleteOrganization)

			r.Route("/members", func(r chi.Router) {

				r.Delete("/", controllers.RemoveMemberFromOrganization)
				r.Delete("/{organization}/leave", controllers.LeaveOrganization)

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
			r.Use(middleware.Authentication)
			r.Get("/organizations", controllers.FetchAllOrganizationDetailsOfCurrentUser)
			r.Get("/{organization}/organization", controllers.FetchOrganizationDetailsOfCurrentUser)
		})

		r.Route("/internal", func(r chi.Router) {
			r.Get("/jwt", controllers.GetJWT)
			r.Route("/members", func(r chi.Router) {
				r.Post("/organizations", controllers.FetchOragnizationListOfUsers)	
			})
			r.Delete("/organization/{organization}", controllers.DeleteOrganization)
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
