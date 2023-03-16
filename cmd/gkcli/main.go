package main

import (
	"fmt"
	"log"
	"time"

	"os"

	"github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	defaults.StringVar(&cfg.ServerHost, "server-host", "127.0.0.1:10001", "ip:port of GATE-KEEPER gRPC server")
	// versionFlag := defaults.BoolP("version", "v", false, "get version number")

	cmdCheck := pflag.NewFlagSet("check", pflag.ExitOnError)
	checkLogin := cmdCheck.String("login", "", "value of Login attempted")
	checkPass := cmdCheck.String("pass", "", "value of Password (hopefully hashed) attempted")
	checkIP := cmdCheck.String("ip", "", "value of IP from wich login was attempted")
	cmdCheck.AddFlagSet(defaults)

	cmdReset := pflag.NewFlagSet("reset", pflag.ExitOnError)
	resetLogin := cmdReset.String("login", "", "value of Login for wich attempts will be reseted")
	resetIP := cmdReset.String("ip", "", "value of IP for wich attempts will be reseted")
	cmdReset.AddFlagSet(defaults)

	cmdWhiteList := pflag.NewFlagSet("white-list", pflag.ExitOnError)
	whiteListSubNet := cmdWhiteList.String("subnet", "", `value of network to add/delete to white list, subnet in CIDR notation (RFC 4632 and RFC 4291): "IP/MASK" eg "192.0.2.0/24"`)
	whiteListAdd := cmdWhiteList.Bool("add", true, "Add/delete to/from whitelist")
	cmdWhiteList.AddFlagSet(defaults)

	cmdBlackList := pflag.NewFlagSet("black-list", pflag.ExitOnError)
	blackListSubNet := cmdBlackList.String("subnet", "", `value of network to add/delete to black list, subnet in CIDR notation (RFC 4632 and RFC 4291): "IP/MASK" eg "192.0.2.0/24"`)
	blackListAdd := cmdBlackList.Bool("add", true, "Add/delete to/from blacklist")
	cmdBlackList.AddFlagSet(defaults)

	cmdBench := pflag.NewFlagSet("simple-bench", pflag.ExitOnError)
	ipBase := cmdBench.String("ipbase", "192.168.10.", `base for IP like: "192.168.1."`)
	loginPassBase := cmdBench.String("lpbase", "test", `base for logins and passwords`)

	if len(os.Args) == 1 {
		fmt.Println("No subcomand given")
		fmt.Println("")
		fmt.Println("Valid subcomands are: check, reset, white-list, black-list")
		fmt.Println("")
		fmt.Println("Global settings:")
		defaults.PrintDefaults()
		fmt.Println("")
		fmt.Println("check settings:")
		cmdCheck.PrintDefaults()
		fmt.Println("")
		fmt.Println("reset settings:")
		cmdReset.PrintDefaults()
		fmt.Println("")
		fmt.Println("white-list settings:")
		cmdWhiteList.PrintDefaults()
		fmt.Println("")
		fmt.Println("black-list settings:")
		cmdBlackList.PrintDefaults()
		fmt.Println("")
		return fmt.Errorf("no subcomand given")
	}

	switch os.Args[1] {
	case "check":
		err := cmdCheck.Parse(os.Args[2:])
		if err != nil {
			return err
		}

		fmt.Println("Will CHECK login attempt for", "Login:", *checkLogin, "Pass:", *checkPass, "IP:", *checkIP)
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %s", err)
			return err
		}
		defer func(conn *grpc.ClientConn) {
			if err := conn.Close(); err != nil {
				log.Printf("Error when calling CLose: %s", err)
			}

		}(conn)
		c := gatekeeper.NewGateKeeperClient(conn)
		reply, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: *checkLogin, Password: *checkPass, Ip: *checkIP})
		if err != nil {
			log.Printf("Error when calling Check: %s", err)
		}
		fmt.Println("Response from server: ", reply.GetOk())
	case "reset":
		err := cmdReset.Parse(os.Args[2:])
		if err != nil {
			return err
		}

		fmt.Println("Will RESET login attempt counters for", "Login:", *resetLogin, "IP:", *resetIP)
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %s", err)
			return err
		}

		defer func(conn *grpc.ClientConn) {
			if err := conn.Close(); err != nil {
				log.Printf("Error when calling CLose: %s", err)
			}

		}(conn)

		c := gatekeeper.NewGateKeeperClient(conn)
		reply, err := c.Reset(context.Background(), &gatekeeper.ResetRequest{Login: *resetLogin, Ip: *resetIP})
		if err != nil {
			log.Printf("Error when calling Reset: %s", err)
		}
		fmt.Println("Response from server: ", reply.GetOk())
	case "white-list":
		err := cmdWhiteList.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		fmt.Println("Will include/exclude from WHITE-LIST subnet", "Sub-Network:", *whiteListSubNet, "ADD:", *whiteListAdd)
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %s", err)
			return err
		}
		defer func(conn *grpc.ClientConn) {
			if err := conn.Close(); err != nil {
				log.Printf("Error when calling CLose: %s", err)
			}

		}(conn)
		c := gatekeeper.NewGateKeeperClient(conn)
		reply, err := c.WhiteList(context.Background(), &gatekeeper.WhiteListRequest{Subnet: *whiteListSubNet, Isadd: *whiteListAdd})
		if err != nil {
			log.Printf("Error when calling White-list: %s", err)
		}
		fmt.Println("Response from server: ", reply.GetOk())
	case "black-list":
		err := cmdBlackList.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		fmt.Println("Will include/exclude from BLACK-LIST subnet", "Sub-Network:", *blackListSubNet, "ADD:", *blackListAdd)
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %s", err)
			return err
		}
		defer func(conn *grpc.ClientConn) {
			if err := conn.Close(); err != nil {
				log.Printf("Error when calling CLose: %s", err)
			}

		}(conn)
		c := gatekeeper.NewGateKeeperClient(conn)
		reply, err := c.BlackList(context.Background(), &gatekeeper.BlackListRequest{Subnet: *blackListSubNet, Isadd: *blackListAdd})
		if err != nil {
			log.Printf("Error when calling Black-list: %s", err)
		}
		fmt.Println("Response from server: ", reply.GetOk())
	case "simple-bench":
		fmt.Println("Will do simple single-threaded benchmark")
		fmt.Println("Version:", version.GetVersion(), "Revision:", version.GetRevision())
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %s", err)
			return err
		}
		defer func(conn *grpc.ClientConn) {
			if err := conn.Close(); err != nil {
				log.Printf("Error when calling CLose: %s", err)
			}

		}(conn)
		c := gatekeeper.NewGateKeeperClient(conn)
		start := time.Now()
		ips := 256
		lgs := 1000
		for i := 0; i < ips; i++ {
			for j := 0; j < lgs; j++ {
				_, err := c.Check(context.Background(), &gatekeeper.CheckRequest{
					Login:    fmt.Sprintf("%s-%d-%d", *loginPassBase, j, i),
					Password: fmt.Sprintf("%s-%d-%d", *loginPassBase, j, i),
					Ip:       fmt.Sprintf("%s%d", *ipBase, i),
				})
				if err != nil {
					log.Printf("Error when calling Check: %s", err)
				}
			}
		}
		fmt.Println(ips*lgs, " requests executed in ", time.Since(start))
		//fmt.Println("Response from server: ", reply.GetOk())
	default:
		fmt.Printf("%q is not valid subcommand.\n", os.Args[1])
		fmt.Println("")
		fmt.Println("Valid subcomands are: check, reset, white-list, black-list")
		fmt.Println("")
		fmt.Println("Global settings:")
		defaults.PrintDefaults()
		fmt.Println("")
		fmt.Println("check settings:")
		cmdCheck.PrintDefaults()
		fmt.Println("")
		fmt.Println("reset settings:")
		cmdReset.PrintDefaults()
		fmt.Println("")
		fmt.Println("white-list settings:")
		cmdWhiteList.PrintDefaults()
		fmt.Println("")
		fmt.Println("black-list settings:")
		cmdBlackList.PrintDefaults()
		fmt.Println("")
		fmt.Println("simple-bench settings:")
		cmdBench.PrintDefaults()
		fmt.Println("")
		return fmt.Errorf("invalid subcommand")
	}

	return nil
}
