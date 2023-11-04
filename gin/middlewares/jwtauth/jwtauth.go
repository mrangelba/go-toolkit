package jwtauth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTAuth struct {
	alg       jwa.SignatureAlgorithm
	signKey   interface{} // private-key
	verifyKey interface{} // public-key, only used by RSA and ECDSA algorithms
	verifier  jwt.ParseOption
}

const (
	TokenCtxKey = "token"
	ErrorCtxKey = "error"
	UserCtxKey  = "user"
)

var (
	ErrUnauthorized = errors.New("token is unauthorized")
	ErrExpired      = errors.New("token is expired")
	ErrNBFInvalid   = errors.New("token nbf validation failed")
	ErrIATInvalid   = errors.New("token iat validation failed")
	ErrNoTokenFound = errors.New("no token found")
	ErrAlgoInvalid  = errors.New("algorithm mismatch")
)

type config struct {
	handler      gin.HandlerFunc
	firebaseAuth *auth.Client
}

func New(opts ...Option) gin.HandlerFunc {
	cfg := &config{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.firebaseAuth != nil {
		println("Using firebase auth")

		return FirebaseAuthVerify(cfg.firebaseAuth)
	}
	println("NO using firebase auth")

	publicKey, _ := jwtV5.ParseRSAPublicKeyFromPEM([]byte(strings.ReplaceAll(os.Getenv("PUBLIC_KEY"), "\\n", "\n")))
	jwtAuth := newJWTAuth("RS256", nil, publicKey)

	return Authenticator(jwtAuth, cfg)
}

func newJWTAuth(alg string, signKey interface{}, verifyKey interface{}) *JWTAuth {
	ja := &JWTAuth{alg: jwa.SignatureAlgorithm(alg), signKey: signKey, verifyKey: verifyKey}

	if ja.verifyKey != nil {
		ja.verifier = jwt.WithKey(ja.alg, ja.verifyKey)
	} else {
		ja.verifier = jwt.WithKey(ja.alg, ja.signKey)
	}

	return ja
}

func loadToken(c *gin.Context, ja *JWTAuth, findTokenFns ...func(r *http.Request) string) {
	token, err := verifyRequest(ja, c.Request, findTokenFns...)

	c.Set(TokenCtxKey, token)
	c.Set(ErrorCtxKey, err)
	c.Set(UserCtxKey, token.Subject())
}

func verifyRequest(ja *JWTAuth, r *http.Request, findTokenFns ...func(r *http.Request) string) (jwt.Token, error) {
	var tokenString string

	// Extract token string from the request by calling token find functions in
	// the order they where provided. Further extraction stops if a function
	// returns a non-empty string.
	for _, fn := range findTokenFns {
		tokenString = fn(r)
		if tokenString != "" {
			break
		}
	}
	if tokenString == "" {
		return nil, ErrNoTokenFound
	}

	return verifyToken(ja, tokenString)
}

func verifyToken(ja *JWTAuth, tokenString string) (jwt.Token, error) {
	// Decode & verify the token
	token, err := ja.Decode(tokenString)
	if err != nil {
		return token, ErrorReason(err)
	}

	if token == nil {
		return nil, ErrUnauthorized
	}

	if err := jwt.Validate(token); err != nil {
		return token, ErrorReason(err)
	}

	// Valid!
	return token, nil
}

func (ja *JWTAuth) Encode(claims map[string]interface{}) (t jwt.Token, tokenString string, err error) {
	t = jwt.New()
	for k, v := range claims {
		t.Set(k, v)
	}
	payload, err := ja.sign(t)
	if err != nil {
		return nil, "", err
	}
	tokenString = string(payload)
	return
}

func (ja *JWTAuth) Decode(tokenString string) (jwt.Token, error) {
	return ja.parse([]byte(tokenString))
}

func (ja *JWTAuth) sign(token jwt.Token) ([]byte, error) {
	return jwt.Sign(token, jwt.WithKey(ja.alg, ja.signKey))
}

func (ja *JWTAuth) parse(payload []byte) (jwt.Token, error) {
	// we disable validation here because we use jwt.Validate to validate tokens
	return jwt.Parse(payload, ja.verifier, jwt.WithValidate(false))
}

// ErrorReason will normalize the error message from the underlining
// jwt library
func ErrorReason(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired()), err == ErrExpired:
		return ErrExpired
	case errors.Is(err, jwt.ErrInvalidIssuedAt()), err == ErrIATInvalid:
		return ErrIATInvalid
	case errors.Is(err, jwt.ErrTokenNotYetValid()), err == ErrNBFInvalid:
		return ErrNBFInvalid
	default:
		return ErrUnauthorized
	}
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func Authenticator(ja *JWTAuth, cfg *config) gin.HandlerFunc {
	return func(c *gin.Context) {
		loadToken(c, ja, tokenFromHeader, TokenFromCookie)

		token, _, err := tokenFromContext(c)

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if cfg.handler != nil {
			cfg.handler(c)
		}

		// Token is authenticated, pass it through
		c.Next()
	}
}

func tokenFromContext(c *gin.Context) (jwt.Token, map[string]interface{}, error) {
	token, _ := c.Value(TokenCtxKey).(jwt.Token)

	var err error
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return token, nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	err, _ = c.Value(ErrorCtxKey).(error)

	return token, claims, err
}

// UnixTime returns the given time in UTC milliseconds
func UnixTime(tm time.Time) int64 {
	return tm.UTC().Unix()
}

// EpochNow is a helper function that returns the NumericDate time value used by the spec
func EpochNow() int64 {
	return time.Now().UTC().Unix()
}

// ExpireIn is a helper function to return calculated time in the future for "exp" claim
func ExpireIn(tm time.Duration) int64 {
	return EpochNow() + int64(tm.Seconds())
}

// Set issued at ("iat") to specified time in the claims
func SetIssuedAt(claims map[string]interface{}, tm time.Time) {
	claims["iat"] = tm.UTC().Unix()
}

// Set issued at ("iat") to present time in the claims
func SetIssuedNow(claims map[string]interface{}) {
	claims["iat"] = EpochNow()
}

// Set expiry ("exp") in the claims
func SetExpiry(claims map[string]interface{}, tm time.Time) {
	claims["exp"] = tm.UTC().Unix()
}

// Set expiry ("exp") in the claims to some duration from the present time
func SetExpiryIn(claims map[string]interface{}, tm time.Duration) {
	claims["exp"] = ExpireIn(tm)
}

// TokenFromCookie tries to retreive the token string from a cookie named
// "jwt".
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// tokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func tokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// TokenFromQuery tries to retreive the token string from the "jwt" URI
// query parameter.
//
// To use it, build our own middleware handler, such as:
//
//	func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
//		return func(next http.Handler) http.Handler {
//			return Verify(ja, TokenFromQuery, TokenFromHeader, TokenFromCookie)(next)
//		}
//	}
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}
