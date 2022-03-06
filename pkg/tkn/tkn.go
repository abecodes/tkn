package tkn

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

// Token describes the internal token before signing
type Token struct {
	token   jwt.Token
	headers jws.Headers
}

// New creates a new token instance
func New(options ...option) (*Token, error) {
	t := &Token{
		token:   jwt.New(),
		headers: jws.NewHeaders(),
	}

	for _, opt := range options {
		err := opt(t)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Sign returns []byte representation of the token signed with the given key and algorithm
func (t *Token) Sign(key interface{}, alg string) ([]byte, error) {
	return jwt.Sign(t.token, jwa.SignatureAlgorithm(alg), key, jwt.WithHeaders(t.headers))
}

// Parse parses a token []byte representation and verifies it against the given key and algorithm
func Parse(tkn []byte, key interface{}, alg string) (jwt.Token, error) {
	return jwt.Parse(tkn, jwt.WithVerify(jwa.SignatureAlgorithm(alg), key))
}
