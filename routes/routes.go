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

			r.Route("/member", func(r chi.Router) {
				r.Post("/invite", controllers.InvitationToOrganization)
			})
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
