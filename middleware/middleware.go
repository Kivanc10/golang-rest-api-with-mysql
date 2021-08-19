package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySignInKey = "jdnfksdmfksd"

//var deletedToken []string

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		woutBearer := r.Header.Get("Authorization")
		if !strings.Contains(woutBearer, "Bearer") {
			ctx := context.WithValue(r.Context(), "props", jwt.MapClaims{"user_name": ""})
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			fmt.Printf("authheader -> %s and len -> %d\n", authHeader, len(authHeader))
			if len(authHeader) != 2 || authHeader[0] == "null" {
				//fmt.Println("Malformed token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
				log.Fatal("Malformed token")
			} else {
				jwtToken := authHeader[1]
				token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(mySignInKey), nil
				})

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					ctx := context.WithValue(r.Context(), "props", claims)
					// Access context values in handlers like this
					//props, _ := r.Context().Value("props").(jwt.MapClaims)

					next.ServeHTTP(w, r.WithContext(ctx))

				} else {
					fmt.Println("token err -> ", err)
					//r.Header.Set("ExpiredToken", jwtToken)
					//DelTokenIfExpired(jwtToken)
					// usernameInter := claims["user_name"]
					// if username, ok := usernameInter.(fmt.Stringer); ok {
					// 	person := dbop.GetPersonToDelToken(username.String())
					// 	dbop.DeleteTokenIfExpired(person)
					// }

					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("you are Unauthorized or your token is expired"))
				}
			}
		}

	})
}

// func DelTokenIfExpired(token string) []string {
// 	deletedToken = append(deletedToken, token)
// 	return deletedToken
// }

func CreateToken(userId uint64, name string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", mySignInKey) //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["user_name"] = name
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", errors.New("an error occured during the create token")
	}
	fmt.Println("jwt map --> ", atClaims)
	return token, nil
}

// func ParseMapClaims(myMap jwt.MapClaims, tokenStr string) jwt.Claims {
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(mySignInKey), nil
// 	})
// 	if err != nil {
// 		log.Fatal("an error occured during the parse jwt ,,,")
// 	}
// 	claims := token.Claims.(jwt.MapClaims)
// 	return claims
// }
