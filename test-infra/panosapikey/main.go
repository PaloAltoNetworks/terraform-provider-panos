package main

import (
	"fmt"
	"os"

	"github.com/PaloAltoNetworks/pango"
)

func main() {
	var (
		hostname, username, password string
		ok                           bool
	)

	if hostname, ok = os.LookupEnv("PANOS_HOSTNAME"); !ok {
		os.Stderr.WriteString("PANOS_HOSTNAME must be set\n")
		return
	}
	if username, ok = os.LookupEnv("PANOS_USERNAME"); !ok {
		os.Stderr.WriteString("PANOS_USERNAME must be set\n")
		return
	}
	if password, ok = os.LookupEnv("PANOS_PASSWORD"); !ok {
		os.Stderr.WriteString("PANOS_PASSWORD must be set\n")
		return
	}

	fw := &pango.Firewall{Client: pango.Client{
		Hostname: hostname,
		Username: username,
		Password: password,
		Logging:  pango.LogQuiet,
	}}
	if err := fw.Initialize(); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Failed initialize: %s\n", err))
		return
	}
	os.Stdout.WriteString(fmt.Sprintf("%s\n", fw.ApiKey))
}
