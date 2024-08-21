package service

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/TheTNB/go-web-skeleton/internal/app"
	"github.com/TheTNB/go-web-skeleton/internal/http/request"
)

// SuccessResponse 通用成功响应
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Message string `json:"message"`
}

// Success 响应成功
func Success(c fiber.Ctx, data any) error {
	return c.JSON(&SuccessResponse{
		Message: "success",
		Data:    data,
	})
}

// Error 响应错误
func Error(c fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(&ErrorResponse{
		Message: message,
	})
}

// ErrorSystem 响应系统错误
func ErrorSystem(c fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(&ErrorResponse{
		Message: http.StatusText(http.StatusInternalServerError),
	})
}

// Bind 验证并绑定请求参数
func Bind[T any, PT request.Request[T]](c fiber.Ctx) (*T, error) {
	req := PT(new(T))

	// 绑定参数
	if err := c.Bind().URI(req); err != nil {
		return nil, err
	}
	if err := c.Bind().Query(req); err != nil {
		return nil, err
	}
	if slices.Contains([]string{"POST", "PUT", "PATCH"}, strings.ToUpper(c.Method())) {
		if err := c.Bind().Body(req); err != nil {
			return nil, err
		}
	}

	// 验证参数
	if err := req.PrepareForValidation(c); err != nil {
		return nil, err
	}
	if err := app.Validator.Struct(req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, e := range errs {
				return nil, errors.New(e.Translate(*app.Translator))
			}
		}
		return nil, err
	}

	return req, nil
}

// Paginate 取分页条目
func Paginate[T any](c fiber.Ctx, allItems []T) (pagedItems []T, total uint) {
	req, err := Bind[request.Paginate](c)
	if err != nil {
		req.Page = 1
		req.Limit = 10
	}
	total = uint(len(allItems))
	startIndex := (req.Page - 1) * req.Limit
	endIndex := req.Page * req.Limit

	if total == 0 {
		return []T{}, 0
	}
	if startIndex > total {
		return []T{}, total
	}
	if endIndex > total {
		endIndex = total
	}

	return allItems[startIndex:endIndex], total
}

// removeTopStruct 移除验证器返回中的顶层结构
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
