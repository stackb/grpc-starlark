package starlarkcrypto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type TlsConfig struct {
	*tls.Config
	*starlarkstruct.Struct
}

func newTlsConfig(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var certificates *starlark.List
	var rootCertPool *x509CertPool
	var clientAuth int
	var insecureSkipVerify bool
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"certificates?", &certificates,
		"client_auth?", &clientAuth,
		"insecure_skip_verify?", &insecureSkipVerify,
		"root_certificate_authorities?", &rootCertPool,
	); err != nil {
		return nil, err
	}

	var certs []tls.Certificate
	if certificates != nil {
		certs := make([]tls.Certificate, certificates.Len())
		for i := 0; i < certificates.Len(); i++ {
			value := certificates.Index(i)
			if cert, ok := value.(certificate); ok {
				certs = append(certs, *cert.Certificate)
			} else {
				return nil, fmt.Errorf("certificate list entry %d must be a tls.Certificate (got %T)", i, value)
			}
		}
	} else {
		certificates = starlark.NewList(nil)
	}

	var rootCas *x509.CertPool
	if rootCertPool != nil {
		rootCas = rootCertPool.CertPool
	}
	return TlsConfig{
		Config: &tls.Config{
			RootCAs:            rootCas,
			InsecureSkipVerify: insecureSkipVerify,
			Certificates:       certs,
			ClientAuth:         tls.ClientAuthType(clientAuth),
		},
		Struct: starlarkstruct.FromStringDict(
			starlarkutil.Symbol("tls.Config"),
			starlark.StringDict{
				"certificates":         certificates,
				"insecure_skip_verify": starlark.Bool(insecureSkipVerify),
			},
		),
	}, nil
}
