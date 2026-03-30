package server

import (
	"mkit/example/ginapp/internal/controller/technology"

	ctrl "mkit/example/ginapp/internal/controller"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type controller interface {
	Register(r gin.IRouter)
}
type Server struct {
	dep      *ctrl.DependencyContainer
	errGroup errgroup.Group
}

func (s *Server) RegisterRouter(router *gin.Engine) {
	var (
		api         = router.Group("/api/v1")
		d           = s.dep
		private     = api.Group("")
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

	// public routes
	// if d.Cfg.Environment != enum.EnvironmentProduction {
	// 	apidocs.NewController(d).Register(router, api)
	// }
	// auth.NewController(d).Register(api)

	// private routes
	// private.Use(middleware.Auth(d.UserService.AuthenticateUser))
	for prefix, c := range controllers {
		c.Register(private.Group(prefix))
	}
}

func (s *Server) Init() error {
	// begin background jobs

	return nil
}

func (s *Server) Name() string {
	return "gin-app"
}

func (s *Server) Close() error {
	// close dependencies

	return nil
}
