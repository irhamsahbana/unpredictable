package route

import (
	// integlocalstorage "codebase-app/internal/integration/localstorage"
	m "codebase-app/internal/middleware"
	appLogHandler "codebase-app/internal/module/app_log/handler"
	member "codebase-app/internal/module/member/handler"
	product "codebase-app/internal/module/product/handler"

	"codebase-app/pkg/response"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	// storage := integlocalstorage.NewLocalStorageIntegration()

	app.Get("/storage/private/:filename", m.ValidateSignedURL, storageFile)

	appLogHandler.NewAppLogHandler().Register(app.Group("/logs"))
	member.NewMemberHandler().Register(app.Group("/members"))
	product.NewProductHandler().Register(app.Group("/products"))

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Debug().
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found."))
	})
}

func storageFile(c *fiber.Ctx) error {
	var (
		fileName = c.Params("filename")
		filePath = filepath.Join("storage", "private", fileName)
	)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Err(err).Any("url", filePath).Msg("handler::storageFile - File not found")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("File not found"))
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Any("url", filePath).Msg("handler::storageFile - Failed to read file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Send(fileBytes)
}
