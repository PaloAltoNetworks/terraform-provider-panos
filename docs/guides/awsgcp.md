---
page_title: "AWS / GCP Considerations"
---

# AWS / GCP Considerations

There are [a few types](https://aws.amazon.com/marketplace/seller-profile?id=0ed48363-5064-4d47-b41b-a53f7c937314)
of PAN-OS VMs available to bring up in AWS.  Both these VMs as well as the ones
that can be deployed in Google Cloud Platform are different in that
the `admin` password is unset, but it has an SSH key associated with it.  As
the panos Terraform provider package authenticates via username/password, an
initialization step of configuring a password using the given SSH key is
required.  Right now, this initialization step requires manual intervention;
the user must download this SSH key, at which point the following may be used
to automate this initialization:

```go
package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
)

// Various prompts.
var (
    P1 *regexp.Regexp
    P2 *regexp.Regexp
    P3 *regexp.Regexp
)

func init() {
    P1 = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\._\-]+@[a-zA-Z][a-zA-Z0-9\._\-]+> `)
    P2 = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\._\-]+@[a-zA-Z][a-zA-Z0-9\._\-]+# `)
    P3 = regexp.MustCompile(`(Enter|Confirm) password\s+:\s+?`)
}

// Globals to handle I/O.
var (
    stdin io.Writer
    stdout io.Reader
    buf [65 * 1024]byte
)

// ReadTo reads from stdout until the desired prompt is encountered.
func ReadTo(prompt *regexp.Regexp) (string, error) {
    var i int

    for {
        n, err := stdout.Read(buf[i:])
        if n > 0 {
            os.Stdout.Write(buf[i:i + n])
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

// Perform user initialization.
func panosInit() error {
    var err error

    // Load environment variables.
    hostname := os.Getenv("PANOS_HOSTNAME")
    username := os.Getenv("PANOS_USERNAME")
    password := os.Getenv("PANOS_PASSWORD")

    // Sanity check input.
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" || hostname == "" || username == "" || password == "" {
        u := []string{
            fmt.Sprintf("Usage: %s <key_file>", os.Args[0]),
            "",
            "This will connect to a PAN-OS NGFW and perform initial config:",
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
        }
        for i := range u {
            fmt.Printf("%s\n", u[i])
        }
        os.Exit(0)
    }

    // Read in the ssh key file.
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        return fmt.Errorf("Failed to read SSH key file %q: %s", os.Args[1], err)
    }

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

    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
    if err != nil {
        return fmt.Errorf("Failed dial: %s", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return fmt.Errorf("Failed to create session: %s", err)
    }
    defer session.Close()

    modes := ssh.TerminalModes{
        ssh.ECHO: 0,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }

    if err = session.RequestPty("vt100", 80, 80, modes); err != nil {
        return fmt.Errorf("pty request failed: %s", err)
    }

    // Get input/output pipes for the ssh connection.
    stdin, err = session.StdinPipe()
    if err != nil {
        return fmt.Errorf("setup stdin err: %s", err)
    }

    stdout, err = session.StdoutPipe()
    if err != nil {
        return fmt.Errorf("setup stdout err: %s", err)
    }

    // Invoke a shell on the remote host.
    if err = session.Start("/bin/sh"); err != nil {
        return fmt.Errorf("failed session.Start: %s", err)
    }

    // Perform initial config.
    ok := true
    commands := []struct{
        Send string
        Expect *regexp.Regexp
        Validation string
        OmitIfAdmin bool
    }{
        {"", P1, "", false},
        {"set cli pager off", P1, "", false},
        {"show system info", P1, "", false},
        {"configure", P2, "", false},
        {fmt.Sprintf("set mgt-config users %s permissions role-based superuser yes", username), P2, "", true},
        {fmt.Sprintf("set mgt-config users %s password", username), P3, "", false},
        {password, P3, "", false},
        {password, P2, "", false},
        {"commit description 'initial config'", P2, "Configuration committed successfully", false},
        {"exit", P1, "", false},
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
            out, err := ReadTo(cmd.Expect)
            if err != nil {
                return fmt.Errorf("Error in %q: %s", cmd.Send, err)
            }
            if cmd.Validation != "" {
                ok = ok && strings.Contains(out, cmd.Validation)
            }
            // Delay slightly before sending passwords.
            if cmd.Expect == P3 {
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

func main() {
    if err := panosInit(); err != nil {
        fmt.Printf("\nFailed initial config: %s\n", err)
        os.Exit(1)
    }
    fmt.Printf("\nConfig initialization successful")
}
```

Compile the above, put it somewhere in your `$PATH` (such as `$HOME/bin`),
then invoke it after the device is accessible in AWS:

```bash
$ go get golang.org/x/crypto/ssh
$ go build panos_init.go
$ mv panos_init ~/bin
$ panos_init my_ssh_key.pem
```

The private SSH key file is expected to be given as the first param, while the
hostname is retrieved from the following environment variable:

* `PANOS_HOSTNAME`

The username and password are expected to be in the following environment
variables:

* `PANOS_USERNAME`
* `PANOS_PASSWORD`

If `PANOS_USERNAME` is set to `admin`, then the above will skip the step that
creates the account, as the `admin` account already exists.
