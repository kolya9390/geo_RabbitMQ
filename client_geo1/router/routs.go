package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	handler_auth "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/auth"
	handler_geo "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/geo"
	handler_notific "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/notific"
	handler_user "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/user"
	reverproxy "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/infrastructure/reverProxy"
	swaggerui "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/infrastructure/swaggerUI"
	auth_token_midw "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/middleware/auth"
	auth_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/auth"
	geo_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/geo"
	service_notifications "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/notifications"
	user_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/user"
)

func NewApiRouter( /*RPC CLIENT*/ ) http.Handler {
	r := chi.NewRouter()

	proxy := reverproxy.NewReverseProxy("hugo", "1313")

	r.Use(middleware.Logger)
	r.Use(proxy.ReverseProxy)

	//SwaggerUI
	r.Get("/swagger", swaggerui.SwaggerUI)

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./client_app/public"))).ServeHTTP(w, r)
	})

	// handlers

	geo_handler := handler_geo.NewGeoHandler(geo_client.NewGeoClient())
	user_handler := handler_user.NewUserHandler(user_client.NewUserClient())
	auth_handler := handler_auth.NewAuthHandler(auth_client.NewUserClient())
	notific_handler := handler_notific.NewNotificHandler(service_notifications.NewNotificClient())



		// API
			r.Route("/api",func(r chi.Router) {
				r.Post("/register", auth_handler.Registeretion)

				r.Post("/login", auth_handler.Login)

				r.Get("/sms/send",notific_handler.GetSMS)
				r.Get("/email/send",notific_handler.GetEMail)

			// Group Adress
			r.Route("/address", func(r chi.Router) {

				//r.Use(auth_token_midw.TokenAuthMiddleware)
				r.Post("/search", geo_handler.SearchAPI)

				r.Post("/geocode", geo_handler.GeocodeAPI)

			})

			r.Route("/user", func(r chi.Router) {

				r.Use(auth_token_midw.TokenAuthMiddleware)
				r.Get("/getuser", user_handler.GetUser)

				r.Get("/getusers", user_handler.GetUsers)

			})

		})

	return r
}
