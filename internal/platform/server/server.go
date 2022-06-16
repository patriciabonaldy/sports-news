package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/patriciabonaldy/sports-news/internal/config"
	"github.com/patriciabonaldy/sports-news/internal/platform/server/handler"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine
	handler  handler.ArticleHandler

	shutdownTimeout time.Duration
}

func New(ctx context.Context, config *config.Config, handler handler.ArticleHandler) (context.Context, Server) {
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		handler:  handler,

		shutdownTimeout: time.Duration(config.ShutdownTimeout) + time.Second,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

// Middleware is a gin.HandlerFunc that set CORS
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func (s *Server) registerRoutes() {
	s.engine.Use(Middleware())
	s.engine.GET("/health", handler.CheckHandler())
	articles := s.engine.Group("/articles")
	{
		articles.GET("", s.handler.GetArticles())
		articles.GET("/:id", s.handler.GetArticles())
	}
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
