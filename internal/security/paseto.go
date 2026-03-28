package security

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

const claimUserID = "user_id"

type TokenService struct {
	key      paseto.V4SymmetricKey
	implicit []byte // salt
}

func NewTokenService(secret string) (*TokenService, error) {
	// must be 256 bits (32 bytes) as per specs
	// https://github.com/paseto-standard/paseto-spec/blob/master/docs/01-Protocol-Versions/Version4.md
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("invalid secret key: %v", err)
	}

	return &TokenService{
		key:      key,
		implicit: []byte("imdb-unique-pseudo-token"),
	}, nil
}

func (s *TokenService) Generate(userID string, duration time.Duration) string {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	token.SetString(claimUserID, userID)

	return token.V4Encrypt(s.key, s.implicit)
}

func (s *TokenService) Verify(signedToken string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	token, err := parser.ParseV4Local(s.key, signedToken, s.implicit)
	if err != nil {
		return "", err
	}

	userID, err := token.GetString(claimUserID)
	if err != nil {
		return "", fmt.Errorf("user_id claim missing")
	}

	return userID, nil
}
