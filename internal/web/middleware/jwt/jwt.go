package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kgosLj/opsvoid/config"
	"github.com/kgosLj/opsvoid/internal/model"
	"github.com/kgosLj/opsvoid/pkg/utils"
	"go.uber.org/zap"
	"strings"
	"time"
)

// GenerateToken 生成 token
func GenerateToken(user model.User) string {
	expireSeconds := config.Allconfig.Jwt.Expire
	expireDuration := expireSeconds * time.Second
	expirationTime := time.Now().Add(expireDuration)

	claim := &JWTClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Allconfig.Jwt.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(config.Allconfig.Jwt.Secret))
	if err != nil {
		zap.L().Error("token 转化成字符串的时候失败: ", zap.Error(err))
		panic(err)
	}
	return tokenString
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Allconfig.Jwt.Secret), nil
	})

	if err != nil {
		// **【已优化】**：根据不同的错误类型返回更具体的信息
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token 已过期")
		}
		return nil, errors.New("token 无效")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token 无效")
}

// JWTMiddleware JWT 中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在这里添加忽略的路径
		if c.Request.URL.Path == "/api/v1/user/login" {
			c.Next()
			return
		}

		Header := c.GetHeader("Authorization")
		if Header == "" {
			utils.RespondError(c, 401, "未登录或非法访问")
			c.Abort()
			return
		}

		tokenString := ""
		if strings.HasPrefix(Header, "Bearer ") {
			tokenString = Header[7:] // 去掉 "Bearer " 前缀
		} else {
			utils.RespondError(c, 401, "非法登录请求头")
			c.Abort()
			return
		}

		// 这一步是验证 token 是否有效
		claims, err := ParseToken(tokenString)
		if err != nil {
			utils.RespondError(c, 401, "token 解析失败")
			c.Abort()
			return
		}

		// 如果 token 未过期，那么就说明是登录状态，那么此时需要为他刷新一下 token
		// 因为 token 是有过期时间的，那么如果用户一直不操作，那么 token 就会过期
		// 那么就需要重新生成一个 token 给用户，这样用户就不需要重新登录了
		// 这样就可以保证用户的登录状态不会过期
		// 时间小于 token 过期时间的 1/2 时候就刷新 token
		if time.Until(claims.ExpiresAt.Time) < config.Allconfig.Jwt.Expire/2 {
			currentUser := model.User{
				Username: claims.Username,
			}
			newtoken := GenerateToken(currentUser)
			c.Header("Authorization", "Bearer "+newtoken)
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
