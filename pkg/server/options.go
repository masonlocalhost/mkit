package server

import (
	"log/slog"
	"mkit/pkg/config"
	"mkit/pkg/cron"
	"mkit/pkg/tracing"
	"net/http"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"gorm.io/gorm"
)

type Dependency func(o *Dependencies)
type Dependencies struct {
	Name                 string
	AppConfig            *config.App
	Postgres             *gorm.DB
	Redis                *redis.Client
	Logger               *slog.Logger
	GRPCServer           *grpc.Server
	HealthServer         *health.Server
	ChiRouter            chi.Router
	Tracing              *tracing.Service
	RabbitMQ             *rabbitmq.Conn
	CronManager          *cron.Service
	ConnectMux           *http.ServeMux
	ConnectHandlerOptions []connect.HandlerOption
}

func NewDependencies() *Dependencies {
	return &Dependencies{}
}

func Postgres(g *gorm.DB) Dependency {
	return func(o *Dependencies) {
		o.Postgres = g
	}
}

func AppConfig(cf *config.App) Dependency {
	return func(o *Dependencies) {
		o.AppConfig = cf
	}
}

func Redis(r *redis.Client) Dependency {
	return func(o *Dependencies) {
		o.Redis = r
	}
}

func Logger(l *slog.Logger) Dependency {
	return func(o *Dependencies) {
		o.Logger = l
	}
}

func GRPCServer(gs *grpc.Server) Dependency {
	return func(o *Dependencies) {
		o.GRPCServer = gs
	}
}

func ChiRouter(r chi.Router) Dependency {
	return func(o *Dependencies) {
		o.ChiRouter = r
	}
}

func Tracing(t *tracing.Service) Dependency {
	return func(o *Dependencies) {
		o.Tracing = t
	}
}

func HealthServer(h *health.Server) Dependency {
	return func(o *Dependencies) {
		o.HealthServer = h
	}
}

func RabbitMQ(r *rabbitmq.Conn) Dependency {
	return func(o *Dependencies) {
		o.RabbitMQ = r
	}
}

func CronManager(c *cron.Service) Dependency {
	return func(o *Dependencies) {
		o.CronManager = c
	}
}

func ConnectServer(mux *http.ServeMux, opts []connect.HandlerOption) Dependency {
	return func(o *Dependencies) {
		o.ConnectMux = mux
		o.ConnectHandlerOptions = opts
	}
}
