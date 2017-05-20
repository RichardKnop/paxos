package main

import (
	"log"
	"os"

	"github.com/RichardKnop/paxos"
	"github.com/urfave/cli"
)

var cliApp *cli.App

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "paxos"
	cliApp.Usage = "Paxos"
	cliApp.Authors = []cli.Author{
		cli.Author{
			Name:  "Richard Knop",
			Email: "risoknop@gmail.com",
		},
	}
	cliApp.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "runs an agent (which acts as proposer, acceptor and learner)",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Usage: "TCP port to listen on",
				},
				cli.StringSliceFlag{
					Name:  "peers",
					Usage: "Peers to form cluster with",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Int("port") == 0 {
					return cli.NewExitError("Set port to listen on", 1)
				}
				if len(c.StringSlice("peers")) == 0 {
					return cli.NewExitError("Set at least one peer", 1)
				}
				agent := paxos.NewAgent(
					"", // ID
					"", // host
					c.Int("port"),
					c.StringSlice("peers"),
				)
				return agent.Run()
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
