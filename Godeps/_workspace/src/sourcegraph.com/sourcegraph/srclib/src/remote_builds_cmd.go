package src

import (
	"fmt"
	"log"
	"time"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func init() {
	_, err := CLI.AddCommand("builds",
		"list remote builds",
		"The builds command lists remote builds for the repository.",
		&buildsCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type BuildsCmd struct {
	N         int    `short:"n" description:"number of builds to show" default:"5"`
	CommitID  string `long:"commit" description:"filter builds by commit ID"`
	Queued    bool   `long:"queued"`
	Succeeded bool   `long:"succeeded"`
	Ended     bool   `long:"ended"`
	Failed    bool   `long:"failed"`
	Sort      string `long:"sort" default:"updated_at"`
	Direction string `long:"dir" default:"desc"`
}

var buildsCmd BuildsCmd

func (c *BuildsCmd) Execute(args []string) error {
	cl := NewAPIClientWithAuthIfPresent()

	rrepo, err := remoteCmd.getRemoteRepo(cl)
	if err != nil {
		return err
	}

	opt := &sourcegraph.BuildListOptions{
		Repo:        rrepo.URI,
		CommitID:    c.CommitID,
		Queued:      c.Queued,
		Succeeded:   c.Succeeded,
		Ended:       c.Ended,
		Failed:      c.Failed,
		Sort:        c.Sort,
		Direction:   c.Direction,
		ListOptions: sourcegraph.ListOptions{PerPage: c.N},
	}
	builds, _, err := cl.Builds.List(opt)
	if err != nil {
		return err
	}

	for _, b := range builds {
		if b.Success {
			fmt.Printf(green("#% 8d")+" succeeded % 9s ago", b.BID, ago(b.EndedAt.Time))
		} else if b.Failure {
			fmt.Printf(red("#% 8d")+" failed % 9s ago", b.BID, ago(b.EndedAt.Time))
		} else if b.StartedAt.Valid {
			fmt.Printf(cyan("#% 8d")+" started % 9s ago", b.BID, ago(b.StartedAt.Time))
		} else {
			fmt.Printf(gray("#% 8d")+" queued % 9s ago", b.BID, ago(b.CreatedAt))
		}
		fmt.Printf("\t%s\n", b.CommitID)
	}

	return nil
}

func ago(t time.Time) string {
	d := time.Since(t)
	d = (d / time.Second) * time.Second
	return d.String()
}
