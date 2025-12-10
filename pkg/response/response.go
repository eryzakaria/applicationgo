package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedData struct {
	Items interface{} `json:"items"`
	Meta  Meta        `json:"meta"`
}

func Success(c *gin.Context, data interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessPaginated(c *gin.Context, data interface{}, page, limit int, total int64) {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: PaginatedData{
			Items: data,
			Meta: Meta{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}

func Error(c *gin.Context, statusCode int, message string, err error) {
	resp := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	c.JSON(statusCode, resp)
}
