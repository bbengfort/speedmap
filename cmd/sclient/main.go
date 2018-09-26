package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bbengfort/speedmap"
	"github.com/bbengfort/speedmap/server"
	"github.com/bbengfort/speedmap/server/pb"
	"github.com/urfave/cli"
)

var (
	addr, identity string
	client         *server.Client
)

func main() {

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "sclient"
	app.Version = speedmap.Version
	app.Usage = "client to access a speedmap server"
	app.Before = initClient
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "a, addr",
			Usage: "address of the server to connect on",
			Value: "localhost:3264",
		},
		cli.StringFlag{
			Name:  "i, identity",
			Usage: "unique identity of the client",
		},
	}

	// Define commands available to application
	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "get a key from the speedmap server",
			Action: get,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "k, key",
					Usage: "specify the key to get",
				},
			},
		},
		{
			Name:   "put",
			Usage:  "put a value to a key on the speedmap server",
			Action: put,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "k, key",
					Usage: "specify the key to put to",
				},
				cli.StringFlag{
					Name:  "v, val",
					Usage: "value to put to the key",
				},
			},
		},
		{
			Name:   "del",
			Usage:  "delte a key from the speedmap server",
			Action: del,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "k, key",
					Usage: "specify the key to delete",
				},
				cli.BoolFlag{
					Name:  "f, force",
					Usage: "force the deletion of the key",
				},
			},
		},
		{
			Name:   "blast",
			Usage:  "run blast throughput benchmark against speedmap server",
			Action: blast,
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "r, requests",
					Usage: "number of requests issued per client",
					Value: 1000,
				},
				cli.UintFlag{
					Name:  "s, size",
					Usage: "number of bytes per value",
					Value: 32,
				},
				cli.DurationFlag{
					Name:  "d, delay",
					Usage: "wait specified time before starting benchmark",
				},
				cli.IntFlag{
					Name:  "i, indent",
					Usage: "indent the results by specified number of spaces",
				},
			},
		},
	}

	// Run the CLI program
	app.Run(os.Args)

}

func initClient(c *cli.Context) error {
	addr, identity = c.String("addr"), c.String("identity")

	client = server.NewClient(identity)
	if err := client.Connect(addr); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

//===========================================================================
// Client Methods
//===========================================================================

func get(c *cli.Context) (err error) {
	var rep *pb.ClientReply
	if rep, err = client.Get(c.String("key")); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if !rep.Success {
		return cli.NewExitError(rep.Error, 2)
	}

	fmt.Println(rep.Pair)
	return nil
}

func put(c *cli.Context) (err error) {
	var rep *pb.ClientReply
	if rep, err = client.Put(c.String("key"), []byte(c.String("val"))); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if !rep.Success {
		return cli.NewExitError(rep.Error, 2)
	}

	return nil
}

func del(c *cli.Context) (err error) {
	var rep *pb.ClientReply
	if rep, err = client.Del(c.String("key"), c.Bool("force")); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if !rep.Success {
		return cli.NewExitError(rep.Error, 2)
	}

	return nil
}

func blast(c *cli.Context) error {
	if delay := c.Duration("delay"); delay > 0 {
		time.Sleep(delay)
	}

	// Close the client
	client.Close()

	// Create options for benchmark
	N := c.Uint("requests")
	S := c.Uint("size")

	// Create and run benchmark
	benchmark := &server.Blast{}
	if err := benchmark.Run(addr, N, S); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// Print the results
	results, err := benchmark.JSON(c.Int("indent"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(string(results))
	return nil
}
