package jwt

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	key []byte
}

type User struct {
	Address common.Address
	ChainID int
	Expires int64 // Unix time
}

func New(key string) (*JWT, error) {
	if len(key) == 0 {
		return nil, errors.New("lack of key")
	}

	return &JWT{
		key: []byte(key),
	}, nil
}

func (d *JWT) ParseUser(authToken string) (*User, error) {
	if authToken == "" {
		return nil, errors.New("empty auth token")
	}

	// Parse the JWT string and store the result in `*User`.
	user := &User{}
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return d.key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.Address = common.HexToAddress(claims["address"].(string))
		user.ChainID = int(claims["chain_id"].(float64))
		user.Expires = int64(claims["exp"].(float64))
	} else {
		return nil, errors.New("invalid token")
	}

	return user, nil
}

func (d *JWT) SignToken(user *User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"address":  user.Address.Hex(),
		"chain_id": user.ChainID,
		"exp":      user.Expires,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(d.key)
}
