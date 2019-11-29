package xhttp

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/thecodingmachine/gotenberg/internal/pkg/conf"
)

// New returns a custom echo.Echo.
func New(config conf.Config) *echo.Echo {
	srv := echo.New()
	srv.HideBanner = true
	srv.HidePort = true
	namespace := config.HttpPathNamespace()
	fmt.Println(namespace)
	srv.Use(contextMiddleware(config))
	srv.Use(loggerMiddleware())
	srv.Use(cleanupMiddleware())
	srv.Use(errorMiddleware())
	srv.GET(namespace+pingEndpoint, pingHandler)
	srv.POST(namespace+mergeEndpoint, mergeHandler)
	if config.DisableGoogleChrome() && config.DisableUnoconv() {
		return srv
	}
	g := srv.Group(convertGroupEndpoint)
	if !config.DisableGoogleChrome() {
		g.POST(namespace+htmlEndpoint, htmlHandler)
		g.POST(namespace+urlEndpoint, urlHandler)
		g.POST(namespace+markdownEndpoint, markdownHandler)
	}
	if !config.DisableUnoconv() {
		g.POST(namespace+officeEndpoint, officeHandler)
	}
	return srv
}
