package sshconnhelper

import (
	"context"
	"github.com/docker/cli/cli/connhelper"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"time"
)

// GetConnectionHelperBySshClient create ConnectionHelper by ssh Client
func GetConnectionHelperBySshClient(cli *ssh.Client) *connhelper.ConnectionHelper {
	return &connhelper.ConnectionHelper{
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return newConn(cli)
		},
		Host: "http://docker.example.com",
	}
}

type ssConn struct {
	*ssh.Session
	stdOut io.Reader
	stdIn  io.Writer
}

func (s *ssConn) Read(b []byte) (n int, err error) {
	return s.stdOut.Read(b)
}

func (s *ssConn) Write(b []byte) (n int, err error) {
	return s.stdIn.Write(b)
}

func (s *ssConn) LocalAddr() net.Addr {
	return &net.IPAddr{
		IP:   nil,
		Zone: "dummy",
	}
}

func (s *ssConn) RemoteAddr() net.Addr {
	return &net.IPAddr{
		IP:   nil,
		Zone: "dummy",
	}
}

func (s *ssConn) SetDeadline(t time.Time) error {
	logrus.Debugf("unimplemented call: SetDeadline(%v)", t)
	return nil
}

func (s *ssConn) SetReadDeadline(t time.Time) error {
	logrus.Debugf("unimplemented call: SetReadDeadline(%v)", t)
	return nil
}

func (s *ssConn) SetWriteDeadline(t time.Time) error {
	logrus.Debugf("unimplemented call: SetWriteDeadline(%v)", t)
	return nil
}

func (s *ssConn) CloseRead() error {
	return s.Close()
}

func (s *ssConn) CloseWrite() error {
	return s.Close()
}

func newConn(cli *ssh.Client) (net.Conn, error) {
	ss, err := cli.NewSession()
	c := ssConn{
		Session: ss,
	}
	c.stdOut, err = ss.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c.stdIn, err = ss.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = ss.Start("docker system dial-stdio")
	return &c, err
}

func (s *ssConn) Close() error {
	if err := s.Signal(ssh.SIGTERM); err != nil {
		return err
	}
	return s.Session.Close()
}
