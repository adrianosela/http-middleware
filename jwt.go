package middleman

//
// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"strings"
//
// 	jwt "github.com/dgrijalva/jwt-go"
// )
//
// // JWTAuthorizer is a Middleware which runs a given function to authenticate
// // requests with basic authentication in the appropriate header
// type JWTAuthorizer struct {
// 	validateJWT ValidateJWTFunc
// }
//
// // ValidateJWTFunc authorizes a user with a Bearer token
// type ValidateJWTFunc func(tk string) (jwt.Claims, error)
//
// const (
// 	jwtCtxKey = requestCtxKey("jwt-auther-claims")
// )
//
// // NewJWTAuther is the constructor for a JWTAuthorizer Middleware
// func NewJWTAuther(vf ValidateJWTFunc) *JWTAuthorizer {
// 	return &JWTAuthorizer{validateJWT: vf}
// }
//
// // Wrap makes JWTAuthorizer implement the Middleware interface
// func (ja *JWTAuthorizer) Wrap(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var authHeader string
// 		if authHeader = r.Header.Get("Authorization"); authHeader == "" {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("no token in request header"))
// 			return
// 		}
// 		authHeader = strings.TrimPrefix(authHeader, "Bearer")
// 		claims, err := ja.validateJWT(authHeader)
// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte(fmt.Sprintf("user could not be authorized. %s", err)))
// 			return
// 		}
// 		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), jwtCtxKey, claims)))
// 	})
// }
//
// // GetClaims returns the authorized claims
// func GetClaims(r *http.Request) jwt.Claims {
// 	return r.Context().Value(jwtCtxKey).(jwt.Claims)
// }
