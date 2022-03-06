package tkn

//revive:disable:unexported-return
type option func(t *Token) error

// WithClaims adds a list of claims to the token
func WithClaims(clms map[string]interface{}) option {
	return func(t *Token) (err error) {
		for k, v := range clms {
			err = t.token.Set(k, v)
			if err != nil {
				return err
			}
		}
		return err
	}
}

// WithHeaders adds a list of headers to the token
func WithHeaders(hdrs map[string]interface{}) option {
	return func(t *Token) error {
		for k, v := range hdrs {
			err := t.headers.Set(k, v)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

//revive:enable:unexported-return
