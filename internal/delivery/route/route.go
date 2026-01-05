package route

import (
	"net/http"

	"github.com/NiflheimDevs/dyslexics-clock/wire"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RouteInit(app *wire.App) http.Handler {
	mux:= chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	mux.Use(app.Middlewares.PanicWall.Recovery)

	mux.Route("/alarm", func(r chi.Router) {
		r.Use(app.Middlewares.Auth.AuthRequired)
		r.Get("/", app.Handlers.AlarmHandler.GetAlarms)
		r.Post("/", app.Handlers.AlarmHandler.CreateAlarm)
		r.Delete("/{id}", app.Handlers.AlarmHandler.DeleteAlarm)
		r.Patch("/{id}", app.Handlers.AlarmHandler.UpdateAlarm)
	})

	mux.Post("/login", app.Handlers.DeviceHandler.Login)

	mux.Route("/device", func(r chi.Router) {
		r.Use(app.Middlewares.Auth.AuthRequired)
		r.Get("/", app.Handlers.DeviceHandler.GetDevice)
		r.Get("/color", app.Handlers.DeviceHandler.GetColor)
		r.Get("/brightness", app.Handlers.DeviceHandler.GetBrightness)
		r.Patch("/color", app.Handlers.DeviceHandler.UpdateColor)
		r.Patch("/volume", app.Handlers.DeviceHandler.UpdateVolume)
		r.Patch("/brightness", app.Handlers.DeviceHandler.UpdateBrightness)
		r.Post("/bezanbekob", app.Handlers.DeviceHandler.Ring)
		r.Post("/snooze", app.Handlers.DeviceHandler.Snooze)
		r.Post("/silent", app.Handlers.DeviceHandler.Silent)
	})

	return mux
}
