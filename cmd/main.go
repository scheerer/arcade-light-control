package main

import (
	"context"
	"flag"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/scheerer/arcade-screen-colors/arcade"
	"github.com/scheerer/arcade-screen-colors/lights/lifx"
	"github.com/scheerer/light-control/internal/app/controllers"
	"github.com/scheerer/light-control/internal/app/initializers"
	"github.com/scheerer/light-control/internal/app/logging"
	"github.com/scheerer/light-control/web"
	"go.uber.org/zap"
)

var logger = logging.New("main")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func init() {
	initializers.LoadEnvVariables()
	initializers.LoadTemplates()
}

func main() {
	defer logger.Sync()
	logger.Info("Starting light control application")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	e := echo.New()
	e.Renderer = initializers.TemplateRenderer
	e.Logger.SetLevel(log.INFO)

	var screenColorConfig arcade.ScreenColorConfig
	err := env.Parse(&screenColorConfig)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Failed to parse screen color config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	lightServiceCtx, shutdownLightService := context.WithCancel(ctx)
	lightService := lifx.NewLifxFromScreenColorConfig(screenColorConfig)
	go lightService.Start(lightServiceCtx)

	time.Sleep(2 * time.Second)

	lightsController := controllers.NewLightsController(screenColorConfig, lightService)

	staticContent, _ := fs.Sub(web.StaticContent, "static")
	e.StaticFS("/", staticContent)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "UP")
	})
	e.GET("/lights", lightsController.LightsIndex)
	e.PUT("/lights", lightsController.LightsEdit)

	// Start server
	var httpWaitGroup sync.WaitGroup
	httpWaitGroup.Add(1)
	go func() {
		defer httpWaitGroup.Done()
		logger.Info("HTTP server starting")
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			logger.With(zap.Error(err)).Fatal("error running http server")
		}
		logger.Info("HTTP server stopped")
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	logger.Info("Shutting down application")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	lightsController.SetLightMode(controllers.LightModeNormal)
	shutdownLightService()

	httpWaitGroup.Wait()
	logger.With(zap.Error(err)).Info("application stopped")
}
