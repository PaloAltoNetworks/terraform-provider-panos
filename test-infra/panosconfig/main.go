package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Various PANOS prompts.
var (
	stdPrompt      = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\._\-]+@[a-zA-Z][a-zA-Z0-9\._\-]+> `)
	cfgPrompt      = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\._\-]+@[a-zA-Z][a-zA-Z0-9\._\-]+# `)
	passwordPrompt = regexp.MustCompile(`(Enter|Confirm) password\s+:\s+?`)
)

func main() {
	if err := panosInit(); err != nil {
		fmt.Printf("\nFailed initial config: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nConfig initialization successful\n")
}

// Perform user initialization of PANOS.
func panosInit() error {
	// Load environment variables.
	hostname := os.Getenv("PANOS_HOSTNAME")
	username := os.Getenv("PANOS_USERNAME")
	password := os.Getenv("PANOS_PASSWORD")
	privateKey := os.Getenv("PANOS_SSH_PRIVATE_KEY")

	// Sanity check input.
	if (len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help")) || hostname == "" || username == "" || password == "" || privateKey == "" {
		u := []string{
			fmt.Sprintf("Usage: %s", os.Args[0]),
			"",
			"This will connect to a PANOS NGFW and perform initial config:",
			"",
			" * Adds the user as a superuser (if not the admin user)",
			" * Sets the user's password",
			" * Commit",
			"",
			"The following environment variables are required:",
			"",
			" * PANOS_HOSTNAME",
			" * PANOS_USERNAME",
			" * PANOS_PASSWORD",
			" * PANOS_SSH_PRIVATE_KEY",
		}
		for i := range u {
			fmt.Printf("%s\n", u[i])
		}
		os.Exit(0)
	}

	data := []byte(privateKey)
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return fmt.Errorf("Failed to parse private key: %s", err)
	}

	useSshKey := ssh.PublicKeys(signer)

	// Configure and open the ssh connection.
	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			useSshKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Printf("Connecting to %q ...\n", hostname)
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
	if err != nil {
		return fmt.Errorf("Failed dial: %s", err)
	}
	defer client.Close()
	fmt.Println("Connected.")

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = session.RequestPty("vt100", 80, 80, modes); err != nil {
		return fmt.Errorf("pty request failed: %s", err)
	}

	// Get input/output pipes for the ssh connection.
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("setup stdin err: %s", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("setup stdout err: %s", err)
	}

	// Invoke a shell on the remote host.
	fmt.Println("Starting session ...")
	if err = session.Start("/bin/sh"); err != nil {
		return fmt.Errorf("failed session.Start: %s", err)
	}

	// Perform initial config.
	ok := true
	commands := []struct {
		Send        string
		Expect      *regexp.Regexp
		Validation  string
		OmitIfAdmin bool
	}{
		{"", stdPrompt, "", false},
		{"set cli pager off", stdPrompt, "", false},
		{"show system info", stdPrompt, "", false},
		{"configure", cfgPrompt, "", false},
		{fmt.Sprintf("set mgt-config users %s permissions role-based superuser yes", username), cfgPrompt, "", true},
		{fmt.Sprintf("set mgt-config users %s password", username), passwordPrompt, "", false},
		{password, passwordPrompt, "", false},
		{password, cfgPrompt, "", false},
		{"commit description 'initial config'", cfgPrompt, "Configuration committed successfully", false},
		{"exit", stdPrompt, "", false},
		{"exit", nil, "", false},
	}

	for _, cmd := range commands {
		if cmd.OmitIfAdmin && username == "admin" {
			continue
		}
		if cmd.Send != "" {
			stdin.Write([]byte(cmd.Send + "\n"))
		}
		if cmd.Expect != nil {
			out, err := ReadTo(stdout, cmd.Expect)
			if err != nil {
				return fmt.Errorf("Error in %q: %s", cmd.Send, err)
			}
			if cmd.Validation != "" {
				ok = ok && strings.Contains(out, cmd.Validation)
			}
			// Delay slightly before sending passwords.
			if cmd.Expect == passwordPrompt {
				time.Sleep(1 * time.Second)
			}
		} else {
			fmt.Printf("exit\n")
			session.Wait()
		}
	}

	// Completed successfully.
	return nil
}

// ReadTo reads from stdout until the desired prompt is encountered.
func ReadTo(stdout io.Reader, prompt *regexp.Regexp) (string, error) {
	var i int
	var buf [65 * 1024]byte

	for {
		n, err := stdout.Read(buf[i:])
		if n > 0 {
			os.Stdout.Write(buf[i : i+n])
		}
		if err != nil {
			return "", err
		}
		i += n
		if prompt.Find(buf[:i]) != nil {
			return string(buf[:i]), nil
		}
	}
}
