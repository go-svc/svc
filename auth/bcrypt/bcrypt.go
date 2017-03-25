package brcypt

import "golang.org/x/crypto/bcrypt"

// Encrypt 會將來源字串透過 bcrypt 加密，如果 `cost` 是 0 則會使用預設的計算花費時間。
func Encrypt(source string, cost int) (string, error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), cost)
	return string(hashedBytes), err
}

// Compare 會比對已經加密和純文字字串看是否吻合。
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
