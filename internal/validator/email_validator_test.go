package validator

import "testing"

func TestVerifyEmail(t *testing.T) {
	testCases := []struct {
		email    string
		expected bool
	}{
		{"valid@email.com", true},
		{"invalid", false},
	}

	for _, tc := range testCases {
		isValid, _ := VerifyEmail(tc.email)
		if isValid != tc.expected {
			t.Errorf("Email %s should be %t", tc.email, tc.expected)
		}
	}
}
