package client

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// Session is struct representing ssh Session
type Session struct {
	config  *ssh.ClientConfig
	conn    *ssh.Client
	session *ssh.Session
}

// NewSession returns new Session instance
func NewSession(ip, port, user, privateKey string) (*Session, error) {
	buf, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return &Session{
		config:  config,
		conn:    conn,
		session: session,
	}, nil
}

// Close close the session & connection
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}

// Get executes the command
func (s *Session) Get(abs string) ([]byte, error) {
	cmd := fmt.Sprintf("cat %s\n", abs)
	return s.session.Output(cmd)
}
