package server

import (
	"mkit/example/ginapp/internal/controller/technology"

	ctrl "mkit/example/ginapp/internal/controller"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

type controller interface {
	Register(r chi.Router)
}
type Server struct {
	dep      *ctrl.DependencyContainer
	errGroup errgroup.Group
}

func (s *Server) RegisterRouter(router chi.Router) {
	var (
		d           = s.dep
		controllers = map[string]controller{
			// "/profile":           profile.NewController(d),
			// "/graph":             graph.NewController(d),
			// "/entity":            entity.NewController(d),
			// "/collection":        collection.NewController(d),
			"/technology": technology.NewController(d),
			// "/asset":             assets.NewController(d),
			// "/user":              user.NewController(d),
			// "/issue":             issue.NewController(d),
			// "/notify":            notification.NewController(d),
			// "/discovery-history": dishistory.NewController(d),
		}
	)

	router.Route("/api/v1", func(r chi.Router) {
		// public routes
		// if d.Cfg.Environment != enum.EnvironmentProduction {
		// 	apidocs.NewController(d).Register(router, r)
		// }
		// auth.NewController(d).Register(r)

		// private routes
		// r.Use(middleware.Auth(d.UserService.AuthenticateUser))
		for prefix, c := range controllers {
			r.Route(prefix, c.Register)
		}
	})
}

func (s *Server) Init() error {
	// begin background jobs

	return nil
}

func (s *Server) Name() string {
	return "chi-app"
}

func (s *Server) Close() error {
	// close dependencies

	return nil
}
