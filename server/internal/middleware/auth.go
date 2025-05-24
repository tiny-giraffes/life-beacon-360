// Life Beacon 360
// Copyright (C) 2025 Tim Yashin/tiny-giraffes
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// internal/middleware/auth.go

package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tiny-giraffes/life-beacon-360/server/config"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get authorization header
		authHeader := c.Request().Header.Get("Authorization")

		// Check if header is empty
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized - Missing authentication token",
			})
		}

		// Extract token - supports both "Bearer TOKEN" and plain "TOKEN" formats
		var token string
		if strings.HasPrefix(authHeader, "Bearer ") {
			// Extract token from "Bearer TOKEN" format
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Use the raw header value as token
			token = authHeader
		}

		// Validate token
		if token == "" || token != config.AppConfig.ApiToken {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized - Invalid authentication token",
			})
		}

		return next(c)
	}
}
