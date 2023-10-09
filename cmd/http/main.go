package main

import (
	"errors"
	"fmt"
	"gobit-demo/app/flow"
	"gobit-demo/config"
	"gobit-demo/infra/api"
	"gobit-demo/infra/cache"
	"gobit-demo/infra/db"
	"gobit-demo/infra/event"
	"gobit-demo/infra/logging"
	"gobit-demo/infra/mq"
	"gobit-demo/infra/orm"
	"gobit-demo/infra/validate"
	"gobit-demo/middlewares"
	"gobit-demo/services/audit"
	"gobit-demo/services/auth"
	"gobit-demo/services/auth/rbac"
	"gobit-demo/services/auth/session"
	"net/http"

	"github.com/labstack/echo/v4"
	v "github.com/pot-code/gobit/pkg/validate"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	va := validate.New()
	cfg := config.LoadConfig()
	logging.Init(cfg.Logging.Level)

	log.Debug().Any("config", cfg).Msg("config")

	rc := cache.NewRedisCache(cfg.Cache.Address)
	dc := db.NewMysqlDB(cfg.Database.GetDSN())
	gd := orm.NewGormDB(dc, log.Logger)
	kb := mq.NewKafkaPublisher(cfg.MessageQueue.GetBrokerList(), log.Logger)
	eb := event.NewKafkaEventBus(kb)
	ts := auth.NewTokenService(cfg.Token.Secret, cfg.Token.Key)
	sm := session.NewSessionManager(rc, cfg.Session.Exp)
	as := audit.NewService(gd)
	r := rbac.NewRBAC(gd)

	e := api.NewAppEngine()
	e.SetErrorHandler(func(err error, c echo.Context) {
		if errors.Is(err, rbac.ErrUnAuthorized) {
			api.JsonNoPermission(c, "无权限")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.NoContent(http.StatusNotFound)
			return
		}

		switch e := err.(type) {
		case v.ValidationError:
			api.JsonBusinessError(c, e[0].Error())
		case *v.ValidationResult:
			api.JsonBadRequest(c, e.Error())
		case *api.BindError:
			log.Debug().Err(err).Msg("bind error")
			api.JsonBadRequest(c, "数据解析失败，请检查输入")
		case *echo.HTTPError:
			api.Json(c, e.Code, map[string]any{
				"code": e.Code,
				"msg":  e.Message,
			})
		default:
			log.Err(err).Msg("")
			api.JsonServerError(c, e.Error())
		}
	})
	e.Use(middlewares.LoggingMiddleware)

	e.AddRouteGroup("/auth", auth.NewRoute(auth.NewService(gd, eb, sm, r, ts), ts, sm, va))
	e.AddRouteGroup("/flow", flow.NewRoute(flow.NewService(gd, r, eb, as), va),
		middlewares.AuthMiddleware(ts, sm, cfg.Session.RefreshExp))

	if err := e.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)); err != http.ErrServerClosed {
		log.Err(err).Msg("error starting server")
	}
}
