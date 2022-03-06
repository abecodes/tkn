# :package: Tkn

**Tkn** is a cli tool that creates signed
[JSON Web Tokens](https://jwt.io/introduction) from a `JSON` or `YAML`
description.

The **JWT** can be signed via _secret_ or _RSA key_ and with different
algorithms.

The following algorithms are supported:

- HS256
- HS384
- HS512
- RS256
- RS384
- RS512
- PS256
- PS384
- PS512

## :floppy_disk: Install

```bash
go install github.com/abecodes/tkn
```

## :computer: Example

A **JWT** can be described via `JSON`:

```json
{
  "headers": { "KID": "01SwHRh9DRhSevAE" },
  "payload": {
    "iss": "admin",
    "sub": "1337",
    "exp": 16459032010,
    "iat": 1645903201,
    "Roles": ["ADMIN", "USER"],
    "Test": ["another", "array", "here"],
    "hello": "world"
  }
}
```

or `YAML`:

```yaml
headers:
  KID: 01SwHRh9DRhSevAE
payload:
  iss: admin
  sub: '1337'
  exp: 16459032010
  iat: 1645903201
  Roles:
    - ADMIN
    - USER
  Test:
    - another
    - array
    - here
  hello: world
```

### Signing with a secret

```bash
tkn -token token.json -secret '6w9z$C&F)J@NcRfTjWnZr4u7x!A%D*G-' -alg hs256
```

### Signing with a RSA key

```bash
tkn -t token.yaml -k path/to/keyfile.pem -a rs256
```

### Output

**Tkn** will print the generated token as well as _secret_ or _public key_ that
can be used to verify its integrity.

```
-----BEGIN JWT TOKEN-----
eyJLSUQiOiIwMVN3SFJoOURSaFNldkFFIiwiYWxnIjoiSFMyNTYiLCJ0eXAiOiJKV1QiLCJ5b2xvIjoxMjN9.eyJSb2xlcyI6WyJBRE1JTiIsIlVTRVIiXSwiVGVzdCI6WyJhbm90aGVyIiwiYXJyYXkiLCJoZXJlIl0sImV4cCI6MTY0NTkwMzIwMTAsImhlbGxvIjoid29ybGQiLCJpYXQiOjE2NDU5MDMyMDEsImlzcyI6ImFkbWluIiwic3ViIjoiMTMzNyJ9.thTxEWqPozH2WWmiBaVvHIi4tTTAYDBbt4GA3nhhzKY
-----END JWT TOKEN-----

-----BEGIN SECRET-----
6w9z$C&F)J@NcRfTjWnZr4u7x!A%D*G-
-----END SECRET-----
```

## :clipboard: Options

| Parameter     | Description                                  | Default      |
| ------------- | -------------------------------------------- | ------------ |
| -a<br>-alg    | algorithm used to sign the token             | RS256        |
| -k<br>-key    | path to a RSA private key                    | null         |
| -s<br>-secret | secret used to sign the token                | null         |
| -t<br>-token  | path to a JSON/YAML description of the token | ./token.yaml |
