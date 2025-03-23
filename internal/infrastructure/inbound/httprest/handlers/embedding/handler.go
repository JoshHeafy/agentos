package alegracategory

import (
	"agentos/internal/domain"
	"agentos/internal/domain/ports/app"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase app.EmbeddingUseCase
}

func New(useCase app.EmbeddingUseCase) Handler {
	return Handler{useCase: useCase}
}

func (h Handler) Upload(c echo.Context) error {
	service := domain.EmbeddingProvider(c.QueryParam("service"))
	if service == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service parameter is required"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.useCase.Upload(c.Request().Context(), service, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Insert completed successfully"})
}

func (h Handler) Search(c echo.Context) error {
	service := domain.EmbeddingProvider(c.QueryParam("service"))
	if service == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "service parameter is required"})
	}

	result, err := h.useCase.Search(c.Request().Context(), service, c.QueryParam("query"), 3)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
