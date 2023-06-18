package app

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/tastycrayon/go-chat/app/ws"
	"github.com/tastycrayon/go-chat/app/ws/handler"
)

func LoadCookieTokenToHeader(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("pb_auth")
			if err != nil {
				return next(c)
			}

			decodedValue, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return next(c)
			}
			var t struct {
				Token string `json:"token"`
			}
			err = json.Unmarshal([]byte(decodedValue), &t)
			if err != nil {
				return next(c)
			}

			c.Request().Header.Set("Authorization", t.Token)
			return next(c)
		}
	}
}
func CustomAuthMiddleware(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if authRecord == nil {
				fmt.Println("got hit 2")
				return apis.NewForbiddenError("Only registered user can access this endpoint", nil)
			}
			return next(c)
		}
	}
}

func InitRoutes(pb *pocketbase.PocketBase, e *core.ServeEvent, hub *ws.Hub) {

	var group *echo.Group = e.Router.Group("/rooms")

	group.Use(CustomAuthMiddleware(pb)) // requires auth
	group.Use(apis.ActivityLogger(pb))  // print logs

	group.POST("", handler.CreateRoom(hub))
	e.Router.GET("/rooms/:roomSlug",
		handler.JoinRoom(hub),
		LoadCookieTokenToHeader(pb),
		apis.LoadAuthContext(pb), CustomAuthMiddleware(pb),
	)
	group.GET("", handler.GetAvailableRooms(hub))

	e.Router.GET("/participants/:roomSlug",
		handler.GetParticipantByRoom(hub),
		CustomAuthMiddleware(pb),
		apis.ActivityLogger(pb),
	)
}
