package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/scheerer/arcade-screen-colors/arcade"
	"github.com/scheerer/arcade-screen-colors/lights"
	"go.uber.org/zap"
)

type LightsController struct {
	screenColorConfig arcade.ScreenColorConfig
	lightService      lights.LightService
	mu                sync.Mutex
	currentMode       lightMode
	currentModeCancel context.CancelFunc
}

func NewLightsController(screenColorConfig arcade.ScreenColorConfig, lightService lights.LightService) *LightsController {
	lc := &LightsController{
		screenColorConfig: screenColorConfig,
		lightService:      lightService,
		currentMode:       LightModeNormal,
	}

	lc.SetLightMode(LightModeNormal)

	return lc
}

var errUnsupportedLightMode = errors.New("unsupported light mode")

type lightMode string

const LightModeArcadeScreen lightMode = "arcade-screen"
const LightModeNormal lightMode = "normal"
const LightModeOff lightMode = "off"
const LightModeRocksmith lightMode = "rocksmith"

func (lc *LightsController) LightsIndex(c echo.Context) error {
	viewData := map[string]interface{}{
		"title": "Light Control Panel",
	}

	err := c.Render(http.StatusOK, "index", viewData)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (lc *LightsController) LightsEdit(c echo.Context) error {
	mode := lightMode(strings.TrimSpace(strings.ToLower(c.QueryParam("mode"))))
	logger.With(zap.String("mode", string(mode))).Info("Setting light mode")
	err := lc.SetLightMode(mode)
	if err != nil {
		responseCode := http.StatusInternalServerError
		if errors.Is(err, errUnsupportedLightMode) {
			responseCode = http.StatusBadRequest
		} else {
			logger.Error(err)

		}
		return c.JSON(responseCode, err.Error())
	}

	return c.JSON(http.StatusOK, "OK")
}

func (lc *LightsController) SetLightMode(mode lightMode) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if lc.currentModeCancel != nil {
		lc.currentModeCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	lc.currentModeCancel = cancel
	lc.currentMode = mode

	if mode == LightModeArcadeScreen {
		go arcade.RunScreenColors(ctx, lc.screenColorConfig, lc.lightService)
	} else if mode == LightModeNormal {
		time.Sleep(500 * time.Millisecond) // allow time for arcade.RunScreenColors to stop
		lc.lightService.SetColorWithDuration(ctx, lights.Color{Red: 255, Green: 255, Blue: 255}, 250*time.Millisecond)
	} else {
		return errUnsupportedLightMode
	}

	return nil
}
