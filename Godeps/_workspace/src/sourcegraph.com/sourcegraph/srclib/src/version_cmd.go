package src

import (
	"fmt"
	"log"

	"github.com/inconshreveable/go-update/check"
)

// Version of src.
//
// For releases, this is set using the -X flag to `go tool ld`. See
// http://stackoverflow.com/a/11355611.
var Version string

const develVersion = "devel"

func init() {
	if Version == "" {
		Version = develVersion
	}
}

func init() {
	_, err := CLI.AddCommand("version",
		"show version",
		"Show version.",
		&versionCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type VersionCmd struct {
	NoCheck bool `long:"no-check-update" description:"don't check for updates"`
}

var versionCmd VersionCmd

func (c *VersionCmd) Execute(args []string) error {
	fmt.Printf("srclib %s\n", Version)

	// Only check for an update if we're running a released version.
	if Version != develVersion && !c.NoCheck {
		r, err := checkForUpdate()
		if err == check.NoUpdateAvailable {
			log.Println("\nYou are on the latest version of src.")
			return nil
		} else if err != nil {
			return err
		}

		if r != nil {
			log.Printf("\nA newer version of src is available: %s.", r.Version)
			log.Println("Run 'src selfupdate' to update.")
		}
	}

	return nil
}
