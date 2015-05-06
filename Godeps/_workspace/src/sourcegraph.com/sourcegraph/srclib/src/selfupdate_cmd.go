package src

import (
	"fmt"
	"log"

	"github.com/inconshreveable/go-update"
	"github.com/inconshreveable/go-update/check"
)

func init() {
	_, err := CLI.AddCommand("selfupdate",
		"update the 'src' program",
		"Checks for updates. If an update is available, downloads and installs it. The next invocation of 'src' will use the updated program.",
		&selfupdateCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	updateUri       = "https://api.equinox.io/1/Updates"
	updatePublicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2zGxdsE2sVXJlMgeOjXn
LY43+PTZ4wCOw4GZ+PKS7lr0kyeCyn4veKZYHzbdNE1rByO+IZmutqn1ylbI92ck
LKEBe/hKP5M6fnXVuVTPEjO/O27d/pJ9OmC9h3IqorV+G/1dJWlfyhZtqxbCCCbM
d3XRztnRqpx/HI4HgM7Xjjt3U1Qn9pOB9+gcBFBBIVSePWhVkD/uoxP+tqaRUDmp
fcUebuSxG4rsluJCEPoXhzwyJb1l79omTocevjRmyoo7ILxu5Et3w+t5qZ5tY7kL
9EK3l6CPtp//uI8rrMgqsKejekUnNX6I8wd1zJZ9u9Aj/KXXid4AoebFCMf+tjFh
twIDAQAB
-----END RSA PUBLIC KEY-----`
)

type SelfupdateCmd struct {
}

var selfupdateCmd SelfupdateCmd

func (c *SelfupdateCmd) Execute(args []string) error {
	log.Printf("Current: src %s.", Version)

	r, err := checkForUpdate()
	if err == check.NoUpdateAvailable {
		fmt.Println("No updates available.")
		return nil
	} else if err != nil {
		return fmt.Errorf("checking for update: %s.", err)
	}

	log.Printf("Updating to src %s...", r.Version)

	// apply update
	err, errRecover := r.Update()
	if err != nil {
		if errRecover != nil {
			return fmt.Errorf("update recovery failed: %s (%s) -- you may need to recover manually!", errRecover, err)
		}
		return fmt.Errorf("update failed: %s", err)
	}

	fmt.Printf("Updated to src %s.\n", r.Version)
	return nil
}

func checkForUpdate() (*check.Result, error) {
	params := check.Params{
		AppVersion: Version,
		AppId:      "ap_BQxVz1iWMxmjQnbVGd85V58qz6",
		Channel:    "stable",
	}

	up := update.New()
	up, err := up.VerifySignatureWithPEM([]byte(updatePublicKey))
	if err != nil {
		return nil, fmt.Errorf("parse public key for updates: %s.", err)
	}

	return params.CheckForUpdate(updateUri, up)
}
