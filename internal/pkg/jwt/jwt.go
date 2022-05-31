package jwt

// https://jwt.io
// JWT token looks like xxxxx.yyyyy.zzzzz, consists of 3 parts: header, payload, signature
// https://restfulapi.cn/page/jwt
// realize and test inspire from https://zhuanlan.zhihu.com/p/113376580

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/def"
	"github.com/gorilla/websocket"
)

var (
	ErrTokenMalformed        error = errors.New("token is malformed")
	ErrTokenUnverifiable     error = errors.New("token could not be verified because of signing problems")
	ErrTokenSignatureInvalid error = errors.New("signature validation failed")
	ErrTokenExpired          error = errors.New("token is expired")
	ErrTokenNotValidYet      error = errors.New("token not active yet")
	ErrTokenInvalid          error = errors.New("token is invalid")

	SignKey string = "tuo"
)

// self-define claim in payload
type CustomClaims struct {
	UserID uint `json:"userID"`
	// StandardClaims contains some Registered Claim Names, detailes at https://datatracker.ietf.org/doc/html/rfc7519#section-4.1,
	// and also realizes jwt.Claim interface
	jwt.StandardClaims
}

type JWT struct {
	// declare sign key
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(SignKey),
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // use signing algorithm jwt.SigningMethodHS256
	return token.SignedString(j.SigningKey)
}

// parse and validate token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		// type assert
		if ve, ok := err.(*jwt.ValidationError); ok { // invalid token
			// ve.Errors is one of below constants
			// detailes at https://pkg.go.dev/github.com/dgrijalva/jwt-go@v3.2.0+incompatible#pkg-constants
			if (ve.Errors & jwt.ValidationErrorMalformed) != 0 { // Token is malformed
				return nil, ErrTokenMalformed
			} else if (ve.Errors & jwt.ValidationErrorUnverifiable) != 0 { // Token could not be verified because of signing problems
				return nil, ErrTokenUnverifiable
			} else if (ve.Errors & jwt.ValidationErrorSignatureInvalid) != 0 { // Signature validation failed
				return nil, ErrTokenSignatureInvalid
			} else if (ve.Errors & jwt.ValidationErrorExpired) != 0 { // EXP validation failed
				return nil, ErrTokenExpired
			} else if (ve.Errors & jwt.ValidationErrorNotValidYet) != 0 { // NBF validation failed
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}

	// assert token.Claim into CustomClaim form
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// custom a JWT middleware for gin
// If the authorization fails (ex: the password does not match),
// call Abort to ensure the remaining handlers for this request are not called.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token is in http header
		token := c.Request.Header.Get(def.HEADER_TOKEN_KEY)
		// if the req header is set to "Token:abc", the token var here is "abc", it's ok to parse
		// SPECIAL CASE:
		// if the req header is set to "Token:[\"abc\"]" like what http.Header.Set() do, the token var here is "\"abc\""
		token = strings.Trim(token, "\"")

		if token == "" {
			c.JSON(http.StatusUnauthorized, def.BaseResp{
				Error: "request with no token, authorization fails",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// parse payload of token
		claims, err := j.ParseToken(token)
		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, def.BaseResp{
				Error: err.Error(),
			})
			c.Abort()
			return
		}

		// put valid claim into gin.Context
		c.Set("claims", claims)
	}
}

func WSJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// use key Sec-WebSocket-Protocol of header to deliver token
		// if this key is used in websocket req, the resp of server must contains it
		token := c.Request.Header.Get(def.HEADER_SecWebSocketProtocol_KEY)
		// if the req header is set to "Token:abc", the token var here is "abc", it's ok to parse
		// SPECIAL CASE:
		// if the req header is set to "Token:[\"abc\"]" like what http.Header.Set() do, the token var here is "\"abc\""
		token = strings.Trim(token, "\"")

		var upgrader = websocket.Upgrader{
			// solve cross origin domain
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols: []string{token}, // https://www.jianshu.com/p/7b1deb1e0a07
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		//defer conn.Close()

		if token == "" {
			conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "request with no token, authorization fails"),
				time.Now().Add(time.Second))
			conn.Close() // 1006 abnormal closure, last choice
			c.Abort()
			return
		}

		j := NewJWT()
		// parse payload of token
		claims, err := j.ParseToken(token)
		if err != nil || claims == nil {
			// explicitly send close message to inform client to close
			// close codes meaning: https://www.rfc-editor.org/rfc/rfc6455.html#section-7.4
			conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error()),
				time.Now().Add(time.Second))
			conn.Close() // 1006 abnormal closure, last choice
			c.Abort()
			return
		}

		// put valid claim into gin.Context
		c.Set("claims", claims)
		c.Set("wsConn", conn)
	}
}
