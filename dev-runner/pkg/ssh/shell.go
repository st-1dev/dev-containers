package ssh

import (
	"fmt"
	"strconv"

	"github.com/blacknon/go-sshlib"

	"golang.org/x/crypto/ssh"
)

func Shell(host string, port int, user string, password string) (err error) {
	con := &sshlib.Connect{ForwardX11: true}
	auth := sshlib.CreateAuthMethodPassword(password)
	err = con.CreateClient(host, strconv.Itoa(port), user, []ssh.AuthMethod{auth})
	if err != nil {
		return fmt.Errorf("cannot create ssh client: %w", err)
	}

	var session *ssh.Session
	session, err = con.CreateSession()
	if err != nil {
		return fmt.Errorf("cannot create ssh session to '%s@%s:%p': %w", user, host, port, err)
	}

	err = con.Shell(session)
	if err != nil {
		return fmt.Errorf("cannot run ssh shell on '%s@%s:%p': %w", user, host, port, err)
	}

	return nil
}
