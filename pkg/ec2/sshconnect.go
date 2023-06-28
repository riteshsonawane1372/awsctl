package ec2

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func SshConnect(cmd *cobra.Command, args []string) error {

	publicIp, _ := cmd.Flags().GetString("ip")

	user := "ec2-user"
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("An error occurred", err)
	}
	privateKeyPath := filepath.Join(homedir, "/.awsctl/awsctl.pem")
	host := publicIp

	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(`
		Private key awsctl.pem not found at ~/.awsctl/ 
		Use -i option to pass the location of pem file 
		`)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		log.Fatal("Failed to connect to SSH server", err)
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	defer session.Close()

	originalState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set raw terminal mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), originalState)

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	termModes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, termModes); err != nil {
		log.Fatalf("Failed to request pseudo terminal: %v", err)
	}

	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %v", err)
	}

	if err := session.Wait(); err != nil {
		log.Fatalf("Failed to wait for session: %v", err)
	}

	return nil
}

