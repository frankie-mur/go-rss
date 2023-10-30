package validator

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

var verifier = emailverifier.NewVerifier()

func VerifyEmail(email string) (bool, error) {
	ret, err := verifier.Verify(email)
	if err != nil {
		return false, err
	}
	if !ret.Syntax.Valid {
		return false, nil
	}
	return true, nil
}
