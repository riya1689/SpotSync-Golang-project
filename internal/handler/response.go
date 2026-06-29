package handler

import (
	"github.com/labstack/echo/v4" 
	"spotSync-golang-Project/internal/dto"
)

func JSONSuccess(c echo.Context, status int, message string, data interface{}) error { 
	return c.JSON(status, dto.StandardResponse{ 
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONError(c echo.Context, status int, message string, errors interface{}) error { 
	return c.JSON(status, dto.StandardResponse{ 
		Success: false, 
		Message: message,
		Errors:  errors, 
	})
}
