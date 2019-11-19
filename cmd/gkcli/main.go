package main

import (
	"fmt"
	"log"
	//"net"
	"os"

	"github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/imorph/gate-keeper/pkg/api/gatekeeper"
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

	// // command line flags
	// pfs := pflag.NewFlagSet(version.GetAppName(), pflag.ContinueOnError)
	// pfs.StringVar(&cfg.ServerHost, "server-host", "127.0.0.1:10001", "ip:port of GATE-KEEPER gRPC server")
	// versionFlag := pfs.BoolP("version", "v", false, "get version number")

	// // parse flags
	// err := pfs.Parse(os.Args[1:])
	// switch {
	// case err == pflag.ErrHelp:
	// 	os.Exit(0)
	// case err != nil:
	// 	pfs.PrintDefaults()
	// 	return err
	// case *versionFlag:
	// 	fmt.Printf("%s-%s\n", version.GetVersion(), version.GetRevision())
	// 	os.Exit(0)
	// }

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
		os.Exit(2)
	}

	switch os.Args[1] {
	case "check":
		err := cmdCheck.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		//cmdCheck.PrintDefaults()
		//defaults.PrintDefaults()
		// _, _, err = net.SplitHostPort(cfg.ServerHost)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// ipTMP := net.ParseIP(*checkIP)
		// if ipTMP == nil {
		// 	log.Fatal(*checkIP, " is not valid IP")
		// }
		fmt.Println("Will CHECK login attempt for", "Login:", *checkLogin, "Pass:", *checkPass, "IP:", *checkIP)
	case "reset":
		err := cmdReset.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		//fmt.Println(*cmdReset)
		//defaults.PrintDefaults()
		// _, _, err = net.SplitHostPort(cfg.ServerHost)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// ipTMP := net.ParseIP(*checkIP)
		// if ipTMP == nil {
		// 	log.Fatal(*resetIP, " is not valid IP")
		// }
		fmt.Println("Will RESET login attempt counters for", "Login:", *resetLogin, "IP:", *resetIP)
	case "white-list":
		err := cmdWhiteList.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		//cmdWhiteList.PrintDefaults()
		//defaults.PrintDefaults()
		// _, _, err = net.SplitHostPort(cfg.ServerHost)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// _, _, err = net.ParseCIDR(*whiteListSubNet)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		fmt.Println("Will include/exclude from WHITE-LIST subnet", "Sub-Network:", *whiteListSubNet, "ADD:", *whiteListAdd)
	case "black-list":
		err := cmdBlackList.Parse(os.Args[2:])
		if err != nil {
			return err
		}
		//cmdBlackList.PrintDefaults()
		//defaults.PrintDefaults()
		// _, _, err = net.SplitHostPort(cfg.ServerHost)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// _, _, err = net.ParseCIDR(*blackListSubNet)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		fmt.Println("Will include/exclude from BLACK-LIST subnet", "Sub-Network:", *blackListSubNet, "ADD:", *blackListAdd)
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
		os.Exit(2)
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(cfg.ServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := gatekeeper.NewGateKeeperClient(conn)
	reply, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: *checkLogin, Password: *checkPass, Ip: *checkIP})
	if err != nil {
		log.Printf("Error when calling Check: %s", err)
	}
	log.Println("Response from server: ", reply.GetOk())

	return nil
}
