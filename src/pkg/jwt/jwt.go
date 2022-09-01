package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
)

// TokenExpireDuration Access Token有效时间
//const TokenExpireDuration = time.Minute * 10 // 这里不再写死,而是通过配置文件来读取

// 用于加盐的字符串
var mySecret = []byte("带雨云埋一半山")

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserID               int64  `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内嵌的声明
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {

	// 创建一个我们自己的声明
	claims := CustomClaims{
		userID,   // 自定义字段
		username, // 自定义字段
		jwt.RegisteredClaims{ // 指定内置声明的值,我们这里就简单的指定了两个内置声明，Issuer和ExpiresAt
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(viper.GetDuration("auth.jwt_expire")), // 过期时间来源于配置文件
			),
			Issuer: "bluebell", // 签发人/组织
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // HS256是一种加密算法

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析token字符串
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 用于接收解码后的对象
	var mc = new(CustomClaims)

	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		// 标准的Claim则可以直接使用Parse方法
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 校验token的签名是否正确(有没有被篡改)
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}

// 安装: go get -u github.com/golang-jwt/jwt/v4

// 每次调用GenToken时都会生成一个新的token
