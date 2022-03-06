package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/abecodes/goutls/rsakys"
	"github.com/abecodes/tkn/pkg/tkn"
	"github.com/abecodes/tkn/pkg/validate"

	"gopkg.in/yaml.v3"
)

type token struct {
	Headers map[string]interface{} `yaml:"headers"`
	Payload map[string]interface{} `yaml:"payload"`
}

const (
	// ===== JWT descripton flag =====
	jwtDefault    = "./token.yaml"
	jwtLong       = "token"
	jwtLongUsage  = "path to a JSON/YAML description of the token"
	jwtShort      = "t"
	jwtShortUsage = "-token shorthand"
	// ===== Algorithm flag =====
	algDefault    = "RS256"
	algLong       = "alg"
	algLongUsage  = "algorithm used to sign the token: HS256, HS384, HS512, RS256, RS384, RS512, PS256, PS384, PS512"
	algShort      = "a"
	algShortUsage = "-alg shorthand"
	// ===== Private Key flag =====
	keyDefault    = ""
	keyLong       = "key"
	keyLongUsage  = "path to a RSA private key"
	keyShort      = "k"
	keyShortUsage = "-key shorthand"
	// ===== Secret flag =====
	scrtDefault    = ""
	scrtLong       = "secret"
	scrtLongUsage  = "secret used to sign the token"
	scrtShort      = "s"
	scrtShortUsage = "-secret shorthand"
)

var errUnableToParse = errors.New("unable to parse the token description")

// Execute runs the cli tool
func Execute() {
	var jwtPth string
	var keyPth string
	var alg string
	var scrt string

	flag.StringVar(&jwtPth, jwtLong, jwtDefault, jwtLongUsage)
	flag.StringVar(&jwtPth, jwtShort, jwtDefault, jwtShortUsage)
	flag.StringVar(&keyPth, keyLong, keyDefault, keyLongUsage)
	flag.StringVar(&keyPth, keyShort, keyDefault, keyShortUsage)
	flag.StringVar(&alg, algLong, algDefault, algLongUsage)
	flag.StringVar(&alg, algShort, algDefault, algShortUsage)
	flag.StringVar(&scrt, scrtLong, scrtDefault, scrtLongUsage)
	flag.StringVar(&scrt, scrtShort, scrtDefault, scrtShortUsage)

	flag.Parse()

	alg = strings.ToUpper(alg)
	err := validate.Params(alg, scrt, keyPth, jwtPth)
	if err != nil {
		log.Fatalln(err)
	}

	desc, err := os.Open(jwtPth)
	if err != nil {
		log.Fatalln(err)
	}
	defer desc.Close()

	descCntnt, _ := ioutil.ReadAll(desc)
	desc.Close()

	var tk token
	err = yaml.Unmarshal(descCntnt, &tk)
	if err != nil {
		log.Fatalln(errUnableToParse, "\n", err)
	}

	var t *tkn.Token

	switch {
	case tk.Headers == nil && tk.Payload == nil:
		log.Fatalln(errors.New("no valid token description"))
	case tk.Headers != nil && tk.Payload != nil:
		t, err = tkn.New(
			tkn.WithClaims(tk.Payload),
			tkn.WithHeaders(tk.Headers),
		)
	default:
		t, err = tkn.New(
			tkn.WithClaims(tk.Payload),
		)
	}

	if err != nil {
		log.Fatalln(err)
	}

	var key interface{}
	var vrfy string

	if validate.NeedsSecret(alg) {
		key = []byte(scrt)
		vrfy = fmt.Sprintf("-----BEGIN SECRET-----\n%s\n-----END SECRET-----\n", scrt)
	} else {
		prv, kerr := rsakys.ReadPrivate(keyPth)
		if kerr != nil {
			log.Fatalln(kerr)
		}

		pub, perr := rsakys.GetPKIXPublicKeyString(&prv.PublicKey)
		if perr != nil {
			log.Fatalln(perr)
		}

		key = prv
		vrfy = string(pub)
	}

	sign, err := t.Sign(key, alg)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("-----BEGIN JWT TOKEN-----")
	fmt.Println(string(sign))
	fmt.Println("-----END JWT TOKEN-----")
	fmt.Println()
	fmt.Println(vrfy)
}
