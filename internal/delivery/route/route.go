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
		r.Put("/{id}", app.Handlers.AlarmHandler.UpdateAlarm)
	})

	mux.Post("/login", app.Handlers.DeviceHandler.Login)

	mux.Route("/device", func(r chi.Router) {
		r.Use(app.Middlewares.Auth.AuthRequired)
		r.Get("/color", app.Handlers.DeviceHandler.GetColor)
		r.Put("/color", app.Handlers.DeviceHandler.UpdateColor)
	})

	return mux
}
