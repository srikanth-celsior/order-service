package middleware

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
)

func JWTMiddleware(ctx iris.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.StopWithStatus(401)
		return
	}
	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		ctx.StopWithStatus(401)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	ctx.Values().Set("userID", claims["user_id"])
	ctx.Next()
}
