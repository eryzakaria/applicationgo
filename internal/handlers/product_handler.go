package handlers

import (
	"net/http"

	"suitemedia/internal/models"
	"suitemedia/internal/service"
	"suitemedia/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) List(c *gin.Context) {
	var params models.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	products, total, err := h.productService.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch products", err)
		return
	}

	response.SuccessPaginated(c, products, params.Page, params.Limit, total)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found", err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	product, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	response.Success(c, product, http.StatusCreated)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	response.Success(c, gin.H{"message": "Product deleted successfully"})
}
