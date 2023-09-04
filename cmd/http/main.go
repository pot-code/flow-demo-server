package main

import (
	"fmt"
	"gobit-demo/audit"
	"gobit-demo/auth"
	"gobit-demo/config"
	"gobit-demo/features/flow"
	"gobit-demo/features/user"
	"gobit-demo/internal/api"
	"gobit-demo/internal/cache"
	"gobit-demo/internal/db"
	"gobit-demo/internal/event"
	"gobit-demo/internal/logging"
	"gobit-demo/internal/mq"
	"gobit-demo/internal/orm"
	"gobit-demo/internal/uuid"
	"gobit-demo/internal/validate"
	"net/http"

	"github.com/labstack/echo/v4"
	v "github.com/pot-code/gobit/pkg/validate"
	"github.com/rs/zerolog/log"
)

func main() {
	va := validate.New()
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)
	uuid.InitSonyflake(cfg.NodeID)

	log.Debug().Any("config", cfg).Msg("config")

	rc := cache.NewRedisCache(cfg.Cache.Address)
	dc := db.NewMysqlDB(cfg.Database.GetDSN())
	gd := orm.NewGormDB(dc, log.Logger)
	kb := mq.NewKafkaPublisher(cfg.MessageQueue.GetBrokerList(), log.Logger)
	eb := event.NewKafkaEventBus(kb)
	ts := auth.NewJwtTokenService(cfg.Token.Secret, cfg.Token.Key)
	sm := auth.NewRedisSessionManager(rc, cfg.Session.Exp)
	as := audit.NewService(gd, sm)
	rb := auth.NewRBAC(gd, sm)

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		switch e := err.(type) {
		case v.ValidationError:
			api.JsonBusinessError(c, e[0].Error())
		case *v.ValidationResult:
			api.JsonBadRequest(c, e.Error())
		case *api.BindError:
			log.Debug().Err(err).Msg("bind error")
			api.JsonBadRequest(c, "数据解析失败，请检查输入")
		case *auth.UnAuthorizedError:
			api.JsonNoPermission(c, "无权限")
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

	api.NewRouteGroup(e, "/auth", auth.NewRoute(auth.NewService(gd, eb), ts, sm, va))
	api.NewRouteGroup(e, "/flow", api.RouteFn(func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(ts, sm, cfg.Session.RefreshExp))
		flow.NewRoute(flow.NewService(gd, sm, eb, as), rb, va).Append(g)
	}))
	api.NewRouteGroup(e, "/user", api.RouteFn(func(g *echo.Group) {
		g.Use(auth.AuthMiddleware(ts, sm, cfg.Session.RefreshExp))
		user.NewRoute(user.NewService(gd), rb).Append(g)
	}))

	if err := e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)); err != http.ErrServerClosed {
		log.Err(err).Msg("error starting server")
	}
}
