package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// Signer 呈現了一個簽署者。
type Signer struct {
	secret string
	method jwt.SigningMethod
}

// Parser 呈現了一個解析者。
type Parser struct {
	secret string
}

// secretFunc 會驗證暗號的格式。
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// 確保 `alg` 是我們期望的。
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

// NewSigner 會建立並回傳一個新的簽署者，並以此來簽發之後的 JSON Web Token。
func NewSigner(secret string, method jwt.SigningMethod) *Signer {
	return &Signer{
		secret: secret,
		method: method,
	}
}

// Sign 會簽發一個接收到的內容，並回傳 JSON Web Token 格式。
func (s *Signer) Sign(claims jwt.Claims) (token string, err error) {
	t := jwt.NewWithClaims(s.method, claims)
	// 以指定的密碼簽署內容。
	token, err = t.SignedString([]byte(s.secret))
	return
}

// SignWithStruct 會簽發一個接收到的建構體，並回傳 JSON Web Token 格式。
func (s *Signer) SignWithStruct(context interface{}) (token string, err error) {
	t := jwt.NewWithClaims(s.method, jwt.MapClaims(structs.Map(context)))
	// 以指定的密碼簽署內容。
	token, err = t.SignedString([]byte(s.secret))
	return
}

// NewParser 會建立並回傳一個新的解析器。
func NewParser(secret string) *Parser {
	return &Parser{
		secret: secret,
	}
}

// Parse 會解析接收到的 JSON Web Token，並將其資料帶入到傳入的建構體中。
func (p *Parser) Parse(token string, context interface{}) error {
	t, err := jwt.Parse(token, secretFunc(p.secret))
	if err != nil {
		return err

		// 如果 Token 沒有問題則繼續。
	} else if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		// 將 Claims 從 Map 轉換、映照成建構體。
		return mapstructure.Decode(claims, context)

		// 其他錯誤。
	} else {
		return err
	}
}
