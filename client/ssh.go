package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// Session is struct representing ssh Session
type Session struct {
	config    *ssh.ClientConfig
	conn      *ssh.Client
	session   *ssh.Session
	StdinPipe io.WriteCloser
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

// Get is func that get file contents
func (s *Session) Get(abs string) ([]byte, error) {
	cmd := fmt.Sprintf("cat %s\n", abs)
	return s.session.Output(cmd)
}

// Scp is func copy file to abs
func (s *Session) Scp(content, abs string) error {
	w, err := s.session.StdinPipe()
	if err != nil {
		return err
	}

	filename := filepath.Base(abs)
	dir := filepath.Dir(abs)

	go func() {
		defer w.Close()
		fmt.Fprintln(w, "C0644", len(content), filename)
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00")
	}()

	return s.session.Run(fmt.Sprintf("/usr/bin/scp -tr %s", dir))
}
