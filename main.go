package main

import (
	"github.com/hsn723/rdap-exporter/cmd"
	_ "golang.org/x/crypto/x509roots/fallback"
)

func main() {
	cmd.Execute()
}
