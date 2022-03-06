package tkn_test

import (
	"errors"
	"testing"

	"github.com/abecodes/goutls/rsakys"
	"github.com/abecodes/tkn/pkg/tkn"
)

type unit struct {
	desc string
	fn   func() error
}

type test struct {
	desc string
	set  []unit
}

const (
	// Success and failure markers
	success = "\u2713" // check
	failed  = "\u2717" // cross
	// Test strings
	rsaAlg        = "RS256"
	secretAlg     = "HS256"
	invalidAlg    = "XXX123"
	validSecret   = "6w9z$C&F)J@NcRfTjWnZr4u7x!A%D*G-"
	invalidSecret = "yolo"
)

func TestTkn(t *testing.T) {
	validKey, err := rsakys.GetPrivateKey(1024)
	if err != nil {
		t.Error(err)
	}
	invalidKey, err := rsakys.GetPrivateKey(1024)
	if err != nil {
		t.Error(err)
	}
	hdrs := map[string]interface{}{
		"KID": "1234",
	}
	clms := map[string]interface{}{
		"iss": "admin",
		"sub": "1337",
		"exp": 16459032010,
		"iat": 1645903201,
		"Roles": []string{
			"ADMIN",
			"USER",
		},
	}
	invalidHdrs := map[string]interface{}{
		"alg": 1337,
	}
	invalidClms := map[string]interface{}{
		"sub": 1337,
	}

	tests := []test{
		{
			desc: "When creating a signed JWT token",
			set: []unit{
				{
					desc: "Should be able to generate a token string using a RSA key",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						_, err = t.Sign(validKey, rsaAlg)
						return err
					},
				},
				{
					desc: "Should be able to generate a token string using a secret",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						_, err = t.Sign([]byte(validSecret), secretAlg)
						return err
					},
				},
				{
					desc: "Should fail when trying to sign a RSA algorithm token with a secret string",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						_, err = t.Sign([]byte(validSecret), rsaAlg)
						if err == nil {
							return errors.New("token was signed when it shouldnt be")
						}
						return nil
					},
				},
				{
					desc: "Should fail when trying to sign a secret algorithm token with a RSA key",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						_, err = t.Sign(validKey, secretAlg)
						if err == nil {
							return errors.New("token was signed when it shouldnt be")
						}
						return nil
					},
				},
				{
					desc: "Should sign token properly with a RSA key",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						signed, err := t.Sign(validKey, rsaAlg)
						if err != nil {
							return err
						}

						_, err = tkn.Parse(signed, validKey, rsaAlg)
						return err
					},
				},
				{
					desc: "Should sign token properly with a secret",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						signed, err := t.Sign([]byte(validSecret), secretAlg)
						if err != nil {
							return err
						}

						_, err = tkn.Parse(signed, []byte(validSecret), secretAlg)
						return err
					},
				},
				{
					desc: "Should not verify a token signed with a RSA key if wrong key is used",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						signed, err := t.Sign(validKey, rsaAlg)
						if err != nil {
							return err
						}

						_, err = tkn.Parse(signed, invalidKey, rsaAlg)
						if err == nil {
							return errors.New("token was verified when it shouldnt be")
						}
						return nil
					},
				},
				{
					desc: "Should not verify a token signed with a secret if wrong secret is used",
					fn: func() error {
						t, err := tkn.New(
							tkn.WithClaims(clms),
							tkn.WithHeaders(hdrs),
						)
						if err != nil {
							return err
						}

						signed, err := t.Sign([]byte(validSecret), secretAlg)
						if err != nil {
							return err
						}

						_, err = tkn.Parse(signed, []byte(invalidSecret), secretAlg)
						if err == nil {
							return errors.New("token was verified when it shouldnt be")
						}
						return nil
					},
				},
				{
					desc: "Should not create token if the claims are invalid",
					fn: func() error {
						_, err := tkn.New(
							tkn.WithClaims(invalidClms),
						)
						if err == nil {
							return errors.New("token was created with invalid claims")
						}
						return nil
					},
				},
				{
					desc: "Should not create token if the headers are invalid",
					fn: func() error {
						_, err := tkn.New(
							tkn.WithHeaders(invalidHdrs),
						)
						if err == nil {
							return errors.New("token was created with invalid claims")
						}
						return nil
					},
				},
			},
		},
	}

	t.Log("Given the need to be able to create a signed JWT token")
	{
		for i, test := range tests {
			t.Logf("\tTest %d:\t%s", i+1, test.desc)
			{
				for _, x := range test.set {
					err := x.fn()
					if err != nil {
						t.Fatalf("Test %d:\t\t%s\t%s: %v", 1, failed, x.desc, err)
					}
					t.Logf("Test %d:\t\t%s\t%s", 1, success, x.desc)
				}
			}
		}
	}
}
