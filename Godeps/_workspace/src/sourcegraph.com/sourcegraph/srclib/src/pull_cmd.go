package src

import "log"

func init() {
	c, err := CLI.AddCommand("pull",
		"fetch remote build data",
		"The pull command fetches build data from the remote. It is currently an alias for 'src build-data fetch'.",
		&buildDataFetchCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
	setDefaultCommitIDOpt(c)
}
