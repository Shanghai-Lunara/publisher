package scheduler

import (
	"context"
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/Shanghai-Lunara/publisher/pkg/conf"
	"github.com/Shanghai-Lunara/publisher/pkg/dao"
	"github.com/Shanghai-Lunara/publisher/pkg/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	connections *connections
	login       *Login
	httpServer  *http.Server
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewServer(c *conf.Config, rbacPath string) *Server {
	zaplogger.Sugar().Info(1111111)
	ctx, cancel := context.WithCancel(context.Background())
	zaplogger.Sugar().Info(22222)
	_ = dao.New(&c.Mysql)
	s := &Server{
		connections: NewConnections(context.Background(), c),
		login:       NewLogin(rbacPath),
		ctx:         ctx,
		cancel:      cancel,
	}
	zaplogger.Sugar().Info(33333)
	router := gin.New()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sessionStore", store))
	router.Use(cors.Default())
	router.GET(types.HttpHandlerLogin, s.login.LoginHandler)
	router.GET(types.HttpHandlerLogout, s.login.LogoutHandler)
	router.GET(types.WebsocketHandlerDashboard, s.dashboard)
	router.GET(types.WebsocketHandlerRunner, s.runner)
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", c.PublisherService.ListenPort),
		Handler: router,
	}
	s.httpServer = server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				zaplogger.Sugar().Info("Server closed under request")
			} else {
				zaplogger.Sugar().Error("Server closed unexpected err:", err)
			}
		}
	}()
	return s
}

func (s *Server) dashboard(c *gin.Context) {
	zaplogger.Sugar().Infow("dashboard print token", "value", c.Request.Header.Get("Token"))
	s.connections.handlerDashboard(c.Writer, c.Request)
}

func (s *Server) runner(c *gin.Context) {
	s.connections.handlerRunner(c.Writer, c.Request)
}

func (s *Server) Shutdown() {
	s.cancel()
}
