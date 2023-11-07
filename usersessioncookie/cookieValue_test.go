package usersessioncookie

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"golang.org/x/crypto/ssh"
)

func Test_generateCookieValue(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}
	err = privateKey.Validate()
	if err != nil {
		t.Error(err)
	}
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Error(err)
	}
	manager := New("example.com", "/", signer)
	testValue := "TestValue"
	output, err := manager.generateCookieValue(testValue)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s == %s", testValue, output)
}
