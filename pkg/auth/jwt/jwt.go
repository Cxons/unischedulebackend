package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type CustomClaims struct {
	User_id string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}


var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


func GenerateToken(user_id string, email string, role string, typeOfToken string)(string,jwt.NumericDate ,error){
	var expiresAt *jwt.NumericDate;
	refereshTokenDuration,_ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_IN_HOURS"))
	accessTokenDuration,_ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_IN_SECONDS"))
	now := time.Now();
	convertedRefreshTime := jwt.NewNumericDate(now.Add(time.Hour * time.Duration(refereshTokenDuration)))
	convertedAccessTime := jwt.NewNumericDate(now.Add(time.Second * time.Duration(accessTokenDuration)))

	if typeOfToken == "refresh_token" {
		expiresAt = convertedRefreshTime
	} else if typeOfToken == "access_token" {
		expiresAt = convertedAccessTime
	} else {
		return "",jwt.NumericDate{},errors.New("correct type of token not specified")
	}

	claims := &CustomClaims{
		User_id: user_id,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "unischeduleserver",
			Subject: user_id,
			Audience: jwt.ClaimStrings{"http://localhost:5173"},
			ExpiresAt: expiresAt,
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil{
		return "",jwt.NumericDate{},err
	}
	return tokenString,*expiresAt,nil
}

func VerifyToken(tokenStr string)(*CustomClaims,error){
	token,err := jwt.ParseWithClaims(tokenStr, &CustomClaims{},func(token *jwt.Token) (interface{}, error) {
		 if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,errors.New("siging method does not match")
        }
		return jwtSecret,nil
	})
	if err != nil || !token.Valid{
		return &CustomClaims{},err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return &CustomClaims{},errors.New("invalid claims type")
	}

	return claims,nil
}