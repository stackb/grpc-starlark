package starlarkcrypto

import (
	"crypto/x509"
	"fmt"

	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type x509CertPool struct {
	*x509.CertPool
	*starlarkstruct.Struct
}

func newX509CertPool(name string, pool *x509.CertPool) *x509CertPool {
	return &x509CertPool{
		CertPool: pool,
		Struct: starlarkstruct.FromStringDict(
			starlarkutil.Symbol(name),
			starlark.StringDict{
				"add":    starlark.NewBuiltin("add", addCert(pool)),
				"append": starlark.NewBuiltin("append", appendCertsFromPem(pool)),
			},
		),
	}
}

func newCertPool(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs); err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()

	return newX509CertPool(fn.Name(), pool), nil
}

func newSystemCertPool(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs); err != nil {
		return nil, err
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("accessing system cert pool: %w", err)
	}

	return newX509CertPool(fn.Name(), pool), nil
}

func addCert(pool *x509.CertPool) func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var cert *x509.Certificate
		if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
			"cert", &cert,
		); err != nil {
			return nil, err
		}
		pool.AddCert(cert)
		return starlark.None, nil
	}
}

func appendCertsFromPem(pool *x509.CertPool) func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var pem string
		if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
			"pem", &pem,
		); err != nil {
			return nil, err
		}
		ok := pool.AppendCertsFromPEM([]byte(pem))
		return starlark.Bool(ok), nil
	}
}
