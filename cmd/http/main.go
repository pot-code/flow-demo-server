package main

import (
	"fmt"
	"gobit-demo/config"
	"gobit-demo/features/auth"
	"gobit-demo/features/flow"
	"gobit-demo/features/perm"
	"gobit-demo/features/user"
	"gobit-demo/internal/api"
	"gobit-demo/internal/cache"
	"gobit-demo/internal/db"
	"gobit-demo/internal/event"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/mq"
	"gobit-demo/internal/token"
	"gobit-demo/internal/validate"
	"net/http"

	"github.com/labstack/echo/v4"
	v "github.com/pot-code/gobit/pkg/validate"
	"github.com/rs/zerolog/log"
)

func main() {
	validate.Init()
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)

	log.Debug().Any("config", cfg).Msg("config")

	rc := cache.NewRedisCache(cfg.Cache.Address)
	dc := db.NewDB(cfg.Database.String())
	gc := db.NewGormClient(dc, log.Logger)

	pub := mq.NewKafkaPublisher(cfg.MessageQueue.BrokerList(), log.Logger)
	eb := event.NewKafkaEventBus(pub)

	js := auth.NewJwtTokenService(
		token.NewJwtIssuer(cfg.Token.Secret),
		auth.NewRedisTokenBlacklist(rc, cfg.Token.Exp),
		cfg.Token.Exp,
	)
	en := perm.NewCasbinEnforcer(gc)
	ps := perm.NewService(en)

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		switch e := err.(type) {
		case v.ValidationError:
			api.JsonBusinessError(c, e[0].Error())
		case *v.ValidationResult:
			api.JsonBadRequest(c, e.Error())
		case *api.BindError:
			api.JsonBadRequest(c, e.Error())
		case *perm.NoPermissionError:
			api.JsonUnauthorized(c, "权限不足")
		case *echo.HTTPError:
			api.Json(c, e.Code, map[string]any{
				"code": e.Code,
				"msg":  e.Message,
			})
		default:
			log.Err(err).Msg("")
			api.JsonServerError(c, e.Error())
		}
	}
	e.Use(api.LoggingMiddleware)

	api.NewRouteGroup(e, "/auth", auth.NewRoute(auth.NewService(gc, eb, auth.NewBcryptPasswordHash()), js))
	api.NewRouteGroup(e, "/flow", api.RouteFn(func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(js))
		flow.NewRoute(flow.NewService(gc), ps).Append(g)
	}))
	api.NewRouteGroup(e, "/user", api.RouteFn(func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(js))
		user.NewRoute(user.NewService(gc)).Append(g)
	}))

	if err := e.Start(fmt.Sprintf(":%d", cfg.HttpPort)); err != http.ErrServerClosed {
		log.Err(err).Msg("error starting server")
	}
}
