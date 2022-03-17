package helper

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

var APPLICATION_NAME = "lesBe"
var LOGIN_EXPIRATION_DURATION = time.Duration(240) * time.Hour
var RESET_TOKEN_EXPIRATION_DURATION = time.Duration(10) * time.Minute
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var SECRET_APP_KEY = "dsad5435nbnfgher"
var SECRET_PASSWORD_RESET = "dsadsa34543nfgher"

type MyClaims struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

// CReateToken to create token for auth handler
func CreateToken(userId int) (string, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		Id: userId,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString([]byte(SECRET_APP_KEY))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// CReateResePasswordToken to create token for reset password
func CreateResePasswordToken(userId int) (string, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(RESET_TOKEN_EXPIRATION_DURATION).Unix(),
		},
		Id: userId,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString([]byte(SECRET_PASSWORD_RESET))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Parse takes the token string and a function for looking up the key. The latter is especially
// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
// head of the token to identify which key to use, but the parsed token (head and claims) is provided
// to the callback, providing flexibility.
func parseTokenJwt(tokenString, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signing method invalid")
		} else if method != JWT_SIGNING_METHOD {
			return nil, errors.New("signing method invalid")
		}

		return []byte(secretKey), nil
	})

	return token, err
}

// ParseToken to parse token in auth handler
func ParseToken(tokenString string) (int, error) {
	token, err := parseTokenJwt(tokenString, SECRET_APP_KEY)

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token invalid")
	}

	// look the containts of claims
	id := int(claims["id"].(float64))
	expires_at := int(claims["exp"].(float64))

	// convert expires_at to time.Time
	expires_at_time := time.Unix(int64(expires_at), 0)

	// cek if token expired
	if time.Now().Unix() > expires_at_time.Unix() {
		return 0, errors.New("token expired")
	}

	return int(id), nil
}

// ParseResePasswordToken to parse token for reset password
func ParseResetPasswordToken(tokenString string) (int, error) {
	token, err := parseTokenJwt(tokenString, SECRET_PASSWORD_RESET)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token invalid")
	}

	// look the containts of claims
	id := int(claims["id"].(float64))
	expires_at := int(claims["exp"].(float64))

	// convert expires_at to time.Time
	expires_at_time := time.Unix(int64(expires_at), 0)

	// cek if token expired
	if time.Now().Unix() > expires_at_time.Unix() {
		return 0, errors.New("token expired")
	}

	return int(id), nil
}
