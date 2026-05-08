package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go_monolith_sample/internal/domain/common"
	"go_monolith_sample/pkg/auth/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Allowed - User has ADMIN role for ADMIN route", func(t *testing.T) {
		r := gin.New()
		// Giả lập middleware Authenticate đã set user_role vào context
		r.Use(func(c *gin.Context) {
			c.Set("user_role", string(common.RoleAdmin))
			c.Next()
		})

		r.GET("/admin-only", middleware.Authorize(common.RoleAdmin), func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("Forbidden - Staff tries to access Admin route", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_role", string(common.RoleStaff))
			c.Next()
		})

		r.DELETE("/medicines/:id", middleware.Authorize(common.RoleAdmin), func(c *gin.Context) {
			c.String(http.StatusOK, "Deleted")
		})

		req := httptest.NewRequest(http.MethodDelete, "/medicines/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Phải trả về 403 Forbidden
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Allowed - Manager accesses Manager route", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_role", string(common.RoleManager))
			c.Next()
		})

		// Route này cho phép Admin HOẶC Manager
		r.PUT("/refund/:id", middleware.Authorize(common.RoleAdmin, common.RoleManager), func(c *gin.Context) {
			c.String(http.StatusOK, "Refunded")
		})

		req := httptest.NewRequest(http.MethodPut, "/refund/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Refunded", w.Body.String())
	})

	t.Run("Unauthorized - No role in context", func(t *testing.T) {
		r := gin.New()
		// Không set gì vào context cả

		r.GET("/secure", middleware.Authorize(common.RoleAdmin), func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/secure", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
