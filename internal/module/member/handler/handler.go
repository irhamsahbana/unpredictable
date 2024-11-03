package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/member/entity"
	"codebase-app/internal/module/member/ports"
	"codebase-app/internal/module/member/repository"
	"codebase-app/internal/module/member/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type memberHandler struct {
	service ports.MemberService
}

func NewMemberHandler() *memberHandler {
	var (
		repo    = repository.NewMemberRepository()
		service = service.NewMemberService(repo)
		handler = new(memberHandler)
	)
	handler.service = service

	return handler
}

func (h *memberHandler) Register(router fiber.Router) {
	router.Get("/", h.getMembers)
	router.Post("/import", h.importMembers)
	router.Get("/data", h.getMembers)
}

func (h *memberHandler) importMembers(c *fiber.Ctx) error {
	var (
		req = new(entity.ImportMembersReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	file, err := c.FormFile("file")
	if err != nil {
		if err != fasthttp.ErrMissingFile {
			log.Warn().Err(err).Msg("handler::importMember - Missing file")
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
		}
	}
	req.File = file

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::importMember - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.ImportMembers(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *memberHandler) getMembers(c *fiber.Ctx) error {
	var (
		req = new(entity.GetMembersReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::getMembers - Failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::getMembers - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetMembers(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(resp, ""))
}
