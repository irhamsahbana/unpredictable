package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/internal/module/product/repository"
	"codebase-app/internal/module/product/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type productHandler struct {
	service ports.ProductService
}

func NewProductHandler() *productHandler {
	var (
		repo    = repository.NewProductRepository()
		service = service.NewProductService(repo)
		handler = new(productHandler)
	)
	handler.service = service

	return handler
}

func (h *productHandler) Register(router fiber.Router) {
	router.Post("/import", h.importProducts)
	router.Post("/import-transactions", h.importProductTransactions)
	router.Post("/import-grammage", h.importProductsGrammage)

	router.Post("/transactions", h.createProductTransaction)

	router.Get("/transactions", h.getProductTransactions)
	router.Get("/data", h.getProducts)
	router.Get("/grammages", h.getProductGrammages)

}

func (h *productHandler) importProducts(c *fiber.Ctx) error {
	var (
		req = new(entity.ImportProductsReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	file, err := c.FormFile("file")
	if err != nil {
		if err != fasthttp.ErrMissingFile {
			log.Warn().Err(err).Msg("handler::importProducts - Missing file")
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}
	}
	req.File = file

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::importProducts - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.ImportProducts(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *productHandler) importProductsGrammage(c *fiber.Ctx) error {
	var (
		req = new(entity.ImportProductGrammageReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	file, err := c.FormFile("file")
	if err != nil {
		if err != fasthttp.ErrMissingFile {
			log.Warn().Err(err).Msg("handler::importProductsGrammage - Missing file")
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}
	}
	req.File = file

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::importProductsGrammage - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.ImportProductGrammage(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *productHandler) importProductTransactions(c *fiber.Ctx) error {
	var (
		req = new(entity.ImportProductTransactionsReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	file, err := c.FormFile("file")
	if err != nil {
		if err != fasthttp.ErrMissingFile {
			if err == fasthttp.ErrBodyTooLarge {
				log.Warn().Err(err).Msg("handler::importProductTransactions - File too large")
				return c.Status(fiber.StatusRequestEntityTooLarge).JSON(response.Error(err))
			}

			log.Warn().Err(err).Msg("handler::importProductTransactions - Missing file")
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
		}
	}
	req.File = file

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::importProductTransactions - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.ImportProductTransactions(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *productHandler) createProductTransaction(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateProductTransactionReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::createProductTransaction - Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::createProductTransaction - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.CreateProductTransaction(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *productHandler) getProducts(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductsReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetProducts - Invalid request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetProducts - Invalid request query")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	data, err := h.service.GetProducts(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(data, ""))
}

func (h *productHandler) getProductTransactions(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductTransactionsReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetProductTransactions - Invalid request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetProductTransactions - Invalid request query")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	data, err := h.service.GetProductTransactions(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(data, ""))
}

func (h *productHandler) getProductGrammages(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductGrammagesReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetProductGrammages - Invalid request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetProductGrammages - Invalid request query")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	data, err := h.service.GetProductGrammages(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(data, ""))
}
