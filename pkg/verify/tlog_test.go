package verify_test

import (
	"testing"
	"time"

	"github.com/github/sigstore-verifier/pkg/testing/ca"
	"github.com/github/sigstore-verifier/pkg/verify"
	"github.com/stretchr/testify/assert"
)

func TestTlogVerifier(t *testing.T) {
	virtualSigstore, err := ca.NewVirtualSigstore()
	assert.NoError(t, err)

	verifier := verify.NewArtifactTransparencyLogVerifier(virtualSigstore, 1, false)
	statement := []byte(`{"_type":"https://in-toto.io/Statement/v0.1","predicateType":"customFoo","subject":[{"name":"subject","digest":{"sha256":"deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"}}],"predicate":{}}`)
	entity, err := virtualSigstore.Attest("foo@fighters.com", "issuer", statement)
	assert.NoError(t, err)

	_, err = verifier.Verify(entity)
	assert.NoError(t, err)

	virtualSigstore2, err := ca.NewVirtualSigstore()
	assert.NoError(t, err)

	verifier2 := verify.NewArtifactTransparencyLogVerifier(virtualSigstore2, 1, false)
	_, err = verifier2.Verify(entity)
	assert.Error(t, err) // different sigstore instance should fail to verify

	// Attempt to use tlog with integrated time outside certificate validity.
	//
	// This time was chosen assuming the Fulcio signing certificate expires
	// after 5 minutes, but while the TSA intermediate is still valid (2 hours).
	entity, err = virtualSigstore.AttestAtTime("foo@fighters.com", "issuer", statement, time.Now().Add(30*time.Minute))
	assert.NoError(t, err)

	_, err = verifier.Verify(entity)
	assert.Error(t, err)
}