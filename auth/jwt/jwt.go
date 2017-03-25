package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

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

type Signer struct {
	secret string
	method jwt.SigningMethod
}

func NewSigner(secret string, method jwt.SigningMethod) *Signer {
	return &Signer{
		secret: secret,
		method: method,
	}
}

func (s *Signer) Sign(claims jwt.Claims) (token string, err error) {
	t := jwt.NewWithClaims(s.method, claims)
	// Sign the token with the specified secret.
	token, err = t.SignedString([]byte(s.secret))

	return
}

func (s *Signer) SignWithStruct(claims jwt.Claims, context interface{}) {

}

type Parser struct {
	secret string
}

func NewParser(secret string) *Parser {
	return &Parser{
		secret: secret,
	}
}

func (p *Parser) Parse(token string, context interface{}) error {
	// Parse the token.
	t, err := jwt.Parse(token, secretFunc(p.secret))
	if err != nil {
		return err

		// Read the token if it's valid.
	} else if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {

		return mapstructure.Decode(claims, context)
		// Other errors.
	} else {
		return err
	}
}
