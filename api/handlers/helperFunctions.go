package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func parseInt64Param(paramName string, c echo.Context) (int64, error) {
	paramStr := c.Param(paramName)
	value, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		log.Printf("error parsing %s: %v", paramName, err)
		return 0, fmt.Errorf("invalid %s: %w", paramName, err)
	}
	return value, nil
}
