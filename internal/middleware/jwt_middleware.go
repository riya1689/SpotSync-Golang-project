package middleware 

import (
	"net/http" 
	"os" 
	"strings" 

	"github.com/golang-jwt/jwt/v5" 
	"github.com/labstack/echo/v4" 
)

func JWTAuth() echo.MiddlewareFunc { 
	return func(next echo.HandlerFunc) echo.HandlerFunc { 
		return func(c echo.Context) error { 
			authHeader := c.Request().Header.Get("Authorization") 
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") { 
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{ 
					"success": false,
					"message": "Missing or invalid token", 
					"errors":  "Unauthorized", 
				})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ") 
			secret := os.Getenv("JWT_SECRET") 

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { 
				return []byte(secret), nil 
			})

			if err != nil || !token.Valid { 
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{ 
					"success": false,
					"message": "Invalid or expired token", 
					"errors":  err.Error(), 
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims) 
			if !ok { 
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{ 
					"success": false,
					"message": "Invalid token claims", 
				})
			}

			c.Set("user_id", uint(claims["id"].(float64))) 
			c.Set("role", claims["role"].(string)) 

			return next(c) 
		}
	}
}
