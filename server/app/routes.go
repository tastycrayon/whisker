package app

import (
	"encoding/json"
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
				return apis.NewForbiddenError("Only registered user can access this endpoint", nil)
			}
			return next(c)
		}
	}
}

func InitRoutes(pb *pocketbase.PocketBase, e *core.ServeEvent, hub *ws.Hub) {

	var roomGroup *echo.Group = e.Router.Group("/rooms")

	roomGroup.Use(CustomAuthMiddleware(pb)) // requires auth
	roomGroup.Use(apis.ActivityLogger(pb))  // print logs

	roomGroup.POST("", handler.CreateRoom(hub))
	e.Router.GET("/rooms/:roomSlug",
		handler.JoinRoom(hub),
		LoadCookieTokenToHeader(pb),
		apis.LoadAuthContext(pb), CustomAuthMiddleware(pb),
	)
	roomGroup.GET("", handler.GetAvailableRooms(hub))

	// participants
	var participantGroup *echo.Group = e.Router.Group("/participants")
	participantGroup.GET("",
		handler.GetParticipantsByRoom(hub),
		CustomAuthMiddleware(pb),
		apis.ActivityLogger(pb),
	)
	participantGroup.GET("/:participantId",
		handler.GetParticipantById(pb),
		CustomAuthMiddleware(pb),
		apis.ActivityLogger(pb),
	)
}
