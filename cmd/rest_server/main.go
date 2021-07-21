package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tnynlabs/wyrm/pkg/devices"
	"github.com/tnynlabs/wyrm/pkg/endpoints"
	"github.com/tnynlabs/wyrm/pkg/http/rest"
	"github.com/tnynlabs/wyrm/pkg/http/rest/middleware"
	"github.com/tnynlabs/wyrm/pkg/pipelines"
	"github.com/tnynlabs/wyrm/pkg/projects"
	"github.com/tnynlabs/wyrm/pkg/storage/postgres"
	"github.com/tnynlabs/wyrm/pkg/tunnels"
	"github.com/tnynlabs/wyrm/pkg/users"
)

func main() {
	if devFlag := os.Getenv("ENV_FILE"); devFlag == "1" {
		// Load environment variables from .env file
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file (error: %v)", err)
		}
	}

	db, err := postgres.GetFromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	userRepo := postgres.CreateUserRepository(db)
	userService := users.CreateService(userRepo)
	userHandler := rest.CreateUserHandler(userService)

	projectRepo := postgres.CreateProjectRepository(db)
	projectService := projects.CreateService(projectRepo)
	projectHandler := rest.CreateProjectHandler(projectService, userService)

	deviceRepo := postgres.CreateDeviceRepository(db)
	deviceService := devices.CreateDeviceService(deviceRepo)
	deviceHandler := rest.CreateDeviceHandler(deviceService)

	endpointRepo := postgres.CreateEndpointRepository(db)
	endpointService := endpoints.CreateEndpointService(endpointRepo)
	endpointHandler := rest.CreateEndpointHandler(endpointService)

	pipelineRepo := postgres.CreatePipelineRepository(db)
	pipelineService := pipelines.CreateService(pipelineRepo)
	pipelineHandler := rest.CreatePipelineHandler(pipelineService, projectService)

	tunnel_host := os.Getenv("TUNNEL_HOST")
	tunnel_port := os.Getenv("TUNNEL_PORT")
	tunnelService := tunnels.CreateHttpGrpcService(tunnel_host + ":" + tunnel_port)
	grpcHandler := rest.CreateGrpcHandler(tunnelService)

	r := chi.NewRouter()

	if devFlag := os.Getenv("WYRM_DEV"); devFlag == "1" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterWithPwd)
		r.Post("/login", userHandler.LoginWithEmailPwd)
		r.Post("/logout", userHandler.Logout)

		r.Route("/users/{userID}", func(r chi.Router) {
			r.Use(middleware.Auth(userService))
			r.Get("/", userHandler.Get)
			r.Patch("/", userHandler.Update)
			r.Delete("/", userHandler.Delete)

			r.Post("/projects", projectHandler.Create)
			r.Get("/projects", projectHandler.GetAllowed)
		})
		r.Route("/projects/{projectID}", func(r chi.Router) {
			r.Use(middleware.Auth(userService))
			r.Get("/", projectHandler.Get)
			r.Patch("/", projectHandler.Update)
			r.Delete("/", projectHandler.Delete)

			r.Post("/devices", deviceHandler.Create)
			r.Get("/devices", deviceHandler.GetByProjectID)

			r.Post("/pipelines", pipelineHandler.Create)
			r.Get("/pipelines", pipelineHandler.GetByProjectID)
		})
		r.Route("/pipelines/{pipelineID}", func(r chi.Router) {
			r.Use(middleware.Auth(userService))
			r.Get("/", pipelineHandler.Get)
			r.Patch("/", pipelineHandler.Update)
			r.Delete("/", pipelineHandler.Delete)

		})
		r.Route("/devices/{deviceID}", func(r chi.Router) {
			r.Get("/", deviceHandler.Get)
			r.Patch("/", deviceHandler.Update)
			r.Delete("/", deviceHandler.Delete)

			r.Post("/endpoints", endpointHandler.Create)
			r.Get("/endpoints", endpointHandler.GetbyDeviceID)

			r.Get("/invoke/{pattern}", grpcHandler.InvokeDevice)
		})

		r.Route("/endpoints/{endpoint_id}", func(r chi.Router) {
			r.Get("/", endpointHandler.Get)
			r.Patch("/", endpointHandler.Update)
			r.Delete("/", endpointHandler.Delete)
		})
	})

	log.Println("Server running...")
	http.ListenAndServe(":8080", r)
}
