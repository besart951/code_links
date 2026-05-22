package productauth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Licenses []string `json:"licenses"`
	jwt.RegisteredClaims
}

func (c Claims) HasLicense(productID string) bool {
	for _, license := range c.Licenses {
		if license == productID {
			return true
		}
	}

	return false
}
