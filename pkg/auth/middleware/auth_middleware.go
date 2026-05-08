package middleware

import (
	"go_monolith_sample/internal/domain/common"
	"go_monolith_sample/pkg/auth/jwt"
	"go_monolith_sample/pkg/cache/redis"
	apperror "go_monolith_sample/pkg/error"
	"go_monolith_sample/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authorize kiểm tra xem User có quyền thực hiện hành động này không
func Authorize(allowedRoles ...common.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy Role từ Context (Giả sử đã được set vào ở bước Authenticate)
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, apperror.Unauthorized("Bạn chưa đăng nhập", nil))
			c.Abort()
			return
		}

		// 2. Kiểm tra xem Role hiện tại có nằm trong danh sách cho phép không
		role := common.UserRole(userRole.(string))
		isAllowed := false
		for _, r := range allowedRoles {
			if role == r {
				isAllowed = true
				break
			}
		}

		// 3. Nếu không có quyền -> Chặn luôn
		if !isAllowed {
			response.Error(c, apperror.Forbidden("Bạn không có quyền thực hiện hành động này", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

func Authenticate(redisService *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy Token từ Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperror.Unauthorized("Vui lòng đăng nhập", nil))
			c.Abort()
			return
		}

		// 2. Validate định dạng Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, apperror.Unauthorized("Định dạng Token không đúng", nil))
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. KIỂM TRA REDIS BLACKLIST TRƯỚC (Để tiết kiệm CPU giải mã nếu token đã bị hủy)
		isBlacklisted, _ := redisService.IsTokenBlacklisted(c, tokenString)
		if isBlacklisted {
			response.Error(c, apperror.Unauthorized("Phiên đăng nhập đã hết hạn", nil))
			c.Abort()
			return
		}

		// 4. GIẢI MÃ TOKEN
		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			response.Error(c, apperror.Unauthorized("Token không hợp lệ hoặc đã hết hạn", err))
			c.Abort()
			return
		}

		// 5. ĐẨY THÔNG TIN VÀO CONTEXT
		// Các hàm xử lý sau này (Handler, Service) sẽ lấy user_id từ đây
		c.Set("user_id", claims.UserID)
		c.Set("user_role", string(claims.Role))
		c.Set("username", claims.Username)

		c.Next()
	}
}
