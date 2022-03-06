package validate

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
)

var secretAlgos = []string{
	string(jwa.HS256),
	string(jwa.HS384),
	string(jwa.HS512),
}

var keyAlgos = []string{
	string(jwa.PS256),
	string(jwa.PS384),
	string(jwa.PS512),
	string(jwa.RS256),
	string(jwa.RS384),
	string(jwa.RS512),
}

// Params validates the passed params
func Params(alg, scrt, keyPth, jwtPth string) error {
	if !algo(alg) {
		return fmt.Errorf("%s is not a valid signing algorithm", alg)
	}

	if !filepath(jwtPth) {
		return fmt.Errorf("%s is not a valid path to a token description", jwtPth)
	}

	if NeedsSecret(scrt) {
		if scrt == "" {
			return fmt.Errorf("%s requires a secret string to sign the token", alg)
		}
	} else {
		if keyPth == "" {
			return fmt.Errorf("%s requires a path to a keyfile to sign the token", alg)
		}

		if !filepath(keyPth) {
			return fmt.Errorf("%s is not a valid path to a keyfile", keyPth)
		}
	}
	return nil
}

// NeedsSecret evaluates if a given algorithm requires a secret string
func NeedsSecret(str string) bool {
	for _, alg := range secretAlgos {
		if str == alg {
			return true
		}
	}

	return false
}

// Algo validates if a string equals a valid algorithm
func algo(str string) bool {
	for _, alg := range append(keyAlgos, secretAlgos...) {
		if str == alg {
			return true
		}
	}

	return false
}

// Filepath validates if a passed path string is in a useable format
func filepath(str string) bool {
	str = path.Clean(str)
	if path.IsAbs(str) {
		return fs.ValidPath(strings.TrimPrefix(str, string(os.PathSeparator)))
	}
	return fs.ValidPath(str)
}
