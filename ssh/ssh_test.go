package client

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSession(t *testing.T) {
	type data struct {
		TestName   string
		IP         string
		Port       string
		User       string
		PrivateKey string
		Err        error
	}

	privateKey := filepath.Clean(os.Getenv("CONFIGCTL_TEST_PRIVATE_KEY"))

	tests := []data{
		{"Success", "127.0.0.1", "2222", "vagrant", privateKey, nil},
		{"Invalid PrivateKey path", "127.0.0.1", "2222", "vagrant", "./noop_rsa", errors.New("open ./noop_rsa: no such file or directory")},
		{"Invalid PrivateKey", "127.0.0.1", "2222", "vagrant", "./testdata/dummy_rsa", errors.New("ssh: no key found")},
		{"Invalid conn", "127.0.0.1", "23456", "vagrant", privateKey, errors.New("dial tcp 127.0.0.1:23456: getsockopt: connection refused")},
	}

	for i, test := range tests {
		session, err := NewSession(test.IP, test.Port, test.User, test.PrivateKey)

		if test.Err == nil {
			if err != test.Err {
				t.Errorf("Test #%d %s: expected '%s', got '%s'", i, test.TestName, test.Err, err)
			}
		} else {
			if err.Error() != test.Err.Error() {
				t.Errorf("Test #%d %s: expected '%s', got '%s'", i, test.TestName, test.Err, err)
			}
		}
		if session != nil {
			session.Close()
		}
	}

}
