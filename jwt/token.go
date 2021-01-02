package token

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrEmptyToken empty token
	ErrEmptyToken = errors.New("token: empty token")
)

// M payload in jwt
type M map[string]interface{}

var conf *Config

// Config jwt config
type Config struct {
	Secret string
	// LookupMethod methods to lookup token
	// Here are support 3 methods:
	// ["query-<param>", "header-<key>", "cookie-<key>"]
	// eg: ["query-token", "header-token", "cookie-token"]
	// Will start looking up the token from the beginning of the array, and stop until it is found
	LookupMethod  []string
	lookupMethods [][]string
}

// Init -
func (c *Config) Init() error {
	c.lookupMethods = nil
	for _, method := range c.LookupMethod {
		method = strings.TrimSpace(method)
		if method == "" {
			continue
		}
		parts := strings.SplitN(method, "-", 2)
		if len(parts) < 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		c.lookupMethods = append(c.lookupMethods, []string{k, v})
	}
	return nil
}

// Init init config
func Init(c *Config) error {
	if c == nil {
		return errors.New("nil config")
	}

	conf = c
	return conf.Init()
}

// ParseFromRequest parse token from request
func ParseFromRequest(r *http.Request) (payload M, err error) {
	token, err := getToken(r)
	if err != nil {
		return
	}
	return Parse(token, conf.Secret)
}

func getToken(r *http.Request) (string, error) {
	var token string
	var err error
	for _, method := range conf.lookupMethods {
		if token != "" {
			break
		}
		k := method[0]
		v := method[1]
		switch k {
		case "header":
			token, err = jwtFromHeader(r, v)
		case "query":
			token, err = jwtFromQuery(r, v)
		case "cookie":
			token, err = jwtFromCookie(r, v)
		}
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func jwtFromHeader(r *http.Request, key string) (string, error) {
	auth := r.Header.Get(key)
	if auth == "" {
		return "", ErrEmptyToken
	}

	var t string
	// Parse the header to get the token part.
	fmt.Sscanf(auth, "Bearer %s", &t)
	if t == "" {
		return "", ErrEmptyToken
	}
	return t, nil
}

func jwtFromQuery(r *http.Request, key string) (string, error) {
	if values, ok := r.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], nil
	}

	return "", ErrEmptyToken
}

func jwtFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	token, _ := url.QueryUnescape(cookie.Value)
	if token == "" {
		return "", ErrEmptyToken
	}

	return token, nil
}

// Parse validates the token with the specialized secret,
func Parse(tokenString string, secret string) (payload M, err error) {
	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return M(claims), nil
	}

	return
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

// GenerateToken generate a token with given payload
func GenerateToken(payload M) (string, error) {
	claims := jwt.MapClaims(payload)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign the token with the specified secret.
	return token.SignedString([]byte(conf.Secret))
}
