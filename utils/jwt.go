package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// 获取配置文件中的 JWT 密钥
var jwtSecret = []byte(viper.GetString("jwt.JwtSecret"))

// Claims 结构定义了 JWT token 中的声明
type Claims struct {
	Email              string `json:"email"`    // 邮箱
	Password           string `json:"password"` // 密码
	jwt.StandardClaims        // JWT 标准声明
}

// GenerateToken 生成 JWT token
func GenerateToken(email, password string) (string, error) {
	nowTime := time.Now()
	// 设置 token 过期时间为当前时间往后推 _ 小时
	expireTime := nowTime.Add(viper.GetDuration("jwt.Expiration") * time.Hour)

	// 创建声明
	claims := Claims{
		email,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			Issuer:    "gin-blog",        // 发行者
		},
	}

	// 使用声明创建 JWT token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析 JWT token
func ParseToken(token string) (*Claims, error) {
	// 解析 JWT token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// 验证 token 是否有效
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
