package starlarkcrypto

import (
	"crypto/tls"
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type certificate struct {
	*tls.Certificate
	*starlarkstruct.Module // using Module instead of struct to prevent default printing of private_key
}

func newCertificate(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var publicKey string
	var privateKey string
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"public_key", &publicKey,
		"private_key", &privateKey,
	); err != nil {
		return nil, err
	}

	certPEMBlock := []byte(publicKey)
	keyPEMBlock := []byte(privateKey)

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("creating x509 key pair: %w", err)
	}

	tlsCert := certificate{
		Certificate: &cert,
		Module: &starlarkstruct.Module{
			Name: "tls.Certificate",
			Members: starlark.StringDict{
				"public_key":  starlark.String(publicKey),
				"private_key": starlark.String(privateKey),
			},
		},
	}
	return tlsCert, nil
}
