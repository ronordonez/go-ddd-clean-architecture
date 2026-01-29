package http

import (
	"github.com/gofiber/fiber/v2"
	"go-architecture/internal/product/application"
	"go-architecture/internal/shared/logger"
)

type ProductHandler struct {
	service *application.ProductService
	log     *logger.Logger
}

func NewProductHandler(service *application.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		log:     log,
	}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var dto application.CreateProductDTO

	if err := c.BodyParser(&dto); err != nil {
		h.log.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	product, err := h.service.Create(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": product,
	})
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": product,
	})
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	var filters application.ProductListFiltersDTO

	if err := c.QueryParser(&filters); err != nil {
		h.log.Error("Failed to parse query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	products, err := h.service.GetAll(c.Context(), filters)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data":  products,
		"count": len(products),
	})
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var dto application.UpdateProductDTO
	if err := c.BodyParser(&dto); err != nil {
		h.log.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	product, err := h.service.Update(c.Context(), id, dto)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": product,
	})
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
