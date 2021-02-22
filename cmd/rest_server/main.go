package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tnynlabs/wyrm/pkg/devices"
	"github.com/tnynlabs/wyrm/pkg/http/rest"
	"github.com/tnynlabs/wyrm/pkg/http/rest/middleware"
	"github.com/tnynlabs/wyrm/pkg/projects"
	"github.com/tnynlabs/wyrm/pkg/storage/postgres"
	"github.com/tnynlabs/wyrm/pkg/users"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file (error: %v)", err)
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

	r := chi.NewRouter()

	if devFlag := os.Getenv("WYRM_DEV"); devFlag == "1" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterWithPwd)
		r.Post("/login", userHandler.LoginWithEmailPwd)

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
		})

		r.Route("/devices/{deviceID}", func(r chi.Router) {
			r.Get("/", deviceHandler.Get)
			r.Patch("/", deviceHandler.Update)
			r.Delete("/", deviceHandler.Delete)
		})
	})

	log.Println("Server running...")
	http.ListenAndServe(":8080", r)
}
