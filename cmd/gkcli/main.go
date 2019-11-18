package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/imorph/gate-keeper/pkg/api/gatekeeper"
	"github.com/imorph/gate-keeper/pkg/version"
)

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// ==================
	// Configuration
	var cfg struct {
		ServerHost string
	}

	// command line flags
	pfs := pflag.NewFlagSet(version.GetAppName(), pflag.ContinueOnError)
	pfs.StringVar(&cfg.ServerHost, "server-host", "127.0.0.1:10001", "ip:port of GATE-KEEPER gRPC server")
	versionFlag := pfs.BoolP("version", "v", false, "get version number")

	// parse flags
	err := pfs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		pfs.PrintDefaults()
		return err
	case *versionFlag:
		fmt.Printf("%s-%s\n", version.GetVersion(), version.GetRevision())
		os.Exit(0)
	}

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(cfg.ServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := gatekeeper.NewGateKeeperClient(conn)
	response, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "me", Password: "me1", Ip: "192.168.0.1"})
	if err != nil {
		log.Fatalf("Error when calling Check: %s", err)
	}
	log.Println("Response from server: ", response.Ok)

	return nil
}
