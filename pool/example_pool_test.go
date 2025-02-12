package pool

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/nspcc-dev/neofs-sdk-go/client"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/user"
)

func ExampleNew_easiestWay() {
	// Signer generation, like example.
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	signer := user.NewSignerRFC6979(*pk)

	pool, _ := New(NewFlatNodeParams([]string{"grpc://localhost:8080", "grpcs://localhost:8081"}), signer, DefaultOptions())
	_ = pool

	// ...
}

func ExampleNew_adjustingParameters() {
	// Signer generation, like example.
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	signer := user.NewSignerRFC6979(*pk)

	opts := DefaultOptions()
	opts.SetErrorThreshold(10)

	pool, _ := New(NewFlatNodeParams([]string{"grpc://localhost:8080", "grpcs://localhost:8081"}), signer, opts)
	_ = pool

	// ...
}

func ExamplePool_ObjectPutInit_explicitAutoSessionDisabling() {
	var prm client.PrmObjectPutInit

	// If you don't provide the session manually with prm.WithinSession, the request will be executed without session.
	prm.IgnoreSession()
	// ...
}

func ExamplePool_ObjectPutInit_autoSessionDisabling() {
	var sess session.Object
	var prm client.PrmObjectPutInit

	// Auto-session disabled, because you provided session already.
	prm.WithinSession(sess)

	// ...
}
