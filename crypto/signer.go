package neofscrypto

import (
	"errors"
	"fmt"

	"github.com/nspcc-dev/neofs-api-go/v2/refs"
)

// ErrIncorrectSigner is returned from function when the signer passed to it
// is incompatible with the function requirements. This variable is intended
// to be used as documentation and for [errors.Is] purposes and MUST NOT be
// changed.
var ErrIncorrectSigner = errors.New("incorrect signer")

// Scheme represents digital signature algorithm with fixed cryptographic hash function.
//
// Negative values are reserved and depend on context (e.g. unsupported scheme).
type Scheme int32

//nolint:revive
const (
	_ Scheme = iota - 1

	ECDSA_SHA512               // ECDSA with SHA-512 hashing (FIPS 186-3)
	ECDSA_DETERMINISTIC_SHA256 // Deterministic ECDSA with SHA-256 hashing (RFC 6979)
	ECDSA_WALLETCONNECT        // Wallet Connect signature scheme
)

// String implements fmt.Stringer.
func (x Scheme) String() string {
	return refs.SignatureScheme(x).String()
}

// maps Scheme to blank PublicKey constructor.
var publicKeys = make(map[Scheme]func() PublicKey)

// RegisterScheme registers a function that returns a new blank PublicKey
// instance for the given Scheme. This is intended to be called from the init
// function in packages that implement signature schemes.
//
// RegisterScheme panics if function for the given Scheme is already registered.
//
// Note that RegisterScheme isn't tread-safe.
func RegisterScheme(scheme Scheme, f func() PublicKey) {
	_, ok := publicKeys[scheme]
	if ok {
		panic(fmt.Sprintf("scheme %v is already registered", scheme))
	}

	publicKeys[scheme] = f
}

// Signer is an interface of entities that can be used for signing operations
// in NeoFS. Unites secret and public parts. For example, an ECDSA private key
// or external auth service.
//
// See also PublicKey.
type Signer interface {
	// Scheme returns corresponding signature scheme.
	Scheme() Scheme

	// Sign signs digest of the given data. Implementations encapsulate data
	// hashing that depends on Scheme. For example, if scheme uses SHA-256, then
	// Sign signs SHA-256 hash of the data.
	Sign(data []byte) ([]byte, error)

	// Public returns the public key corresponding to the Signer.
	Public() PublicKey
}

// PublicKey represents a public key using fixed signature scheme supported by
// NeoFS.
//
// See also Signer.
type PublicKey interface {
	// MaxEncodedSize returns maximum size required for binary-encoded
	// public key.
	//
	// MaxEncodedSize MUST NOT return value greater than any return of
	// Encode.
	MaxEncodedSize() int

	// Encode encodes public key into buf. Returns number of bytes
	// written.
	//
	// Encode MUST panic if buffer size is insufficient and less than
	// MaxEncodedSize (*). Encode MUST return negative value
	// on any failure except (*).
	//
	// Encode is a reverse operation to Decode.
	//
	// [PublicKeyBytes] may be used to skip explicit buffer allocation.
	Encode(buf []byte) int

	// Decode decodes binary public key.
	//
	// Decode is a reverse operation to Encode.
	Decode([]byte) error

	// Verify checks signature of the given data. True means correct signature.
	Verify(data, signature []byte) bool
}

// StaticSigner emulates real sign and contains already precalculated hash.
// Provides neofscrypto.Signer interface.
type StaticSigner struct {
	scheme Scheme
	sig    []byte
	pubKey PublicKey
}

// NewStaticSigner creates new StaticSigner.
func NewStaticSigner(scheme Scheme, sig []byte, pubKey PublicKey) *StaticSigner {
	return &StaticSigner{
		scheme: scheme,
		sig:    sig,
		pubKey: pubKey,
	}
}

// Scheme returns neofscrypto.ECDSA_DETERMINISTIC_SHA256.
// Implements neofscrypto.Signer.
func (s *StaticSigner) Scheme() Scheme {
	return s.scheme
}

// Sign returns precalculated hash.
// Implements neofscrypto.Signer.
func (s *StaticSigner) Sign(_ []byte) ([]byte, error) {
	return s.sig, nil
}

// Public returns neofscrypto.PublicKey.
// Implements neofscrypto.Signer.
func (s *StaticSigner) Public() PublicKey {
	return s.pubKey
}
