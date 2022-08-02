package main

import (
	"fmt"
	"github.com/civet148/geoip"
	"github.com/civet148/log"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
)

const (
	Version      = "1.0.0"
	PROGRAM_NAME = "geoip"
)

var (
	BuildTime = "2022-08-01"
	GitCommit = ""
)

const (
	CMD_NAME_TEST   = "test"
	CMD_NAME_IMPORT = "import"
)

const (
	CMD_FLAG_NAME_DEBUG   = "debug"
	CMD_FLAG_NAME_DSN     = "dsn"
	CMD_FLAG_NAME_TABLE   = "table"
	CMD_FLAG_NAME_EXTRACT = "extract"
)

const (
	DEFAULT_TABLE_NAME = "ip_info"
)

func init() {
	log.SetLevel("info")
}

func grace() {
	//capture signal of Ctrl+C and gracefully exit
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		for {
			select {
			case s := <-sigChannel:
				{
					if s != nil && s == os.Interrupt {
						fmt.Printf("Ctrl+C signal captured, program exiting...\n")
						close(sigChannel)
						os.Exit(0)
					}
				}
			}
		}
	}()
}

func main() {

	grace()

	local := []*cli.Command{
		testCmd,
		importCmd,
	}
	app := &cli.App{
		Name:     PROGRAM_NAME,
		Version:  fmt.Sprintf("v%s %s commit %s", Version, BuildTime, GitCommit),
		Flags:    []cli.Flag{},
		Commands: local,
		Action:   nil,
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit in error %s", err)
		os.Exit(1)
		return
	}
}

var testCmd = &cli.Command{
	Name:      CMD_NAME_TEST,
	Usage:     "test",
	ArgsUsage: "",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  CMD_FLAG_NAME_DEBUG,
			Usage: "open debug mode",
		},
	},
	Action: func(cctx *cli.Context) error {
		strIP := "218.76.88.19"
		uip := geoip.IP2Uint(strIP)
		log.Infof("IP [%s] to uint32 [%d]", strIP, uip)
		strIP = geoip.Uint2IP(uip)
		log.Infof("uint32 [%d] to IP [%s]", uip, strIP)
		geo, err := geoip.NewGeoIP("ip.dat")
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		loc := geo.Find(strIP)
		log.Infof("IP [%s] location [%+v]", strIP, loc)
		return nil
	},
}

var importCmd = &cli.Command{
	Name:      CMD_NAME_IMPORT,
	Usage:     "import to MySQL database",
	ArgsUsage: "[/path/to/ip.dat]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  CMD_FLAG_NAME_DEBUG,
			Usage: "open debug mode",
		},
		&cli.StringFlag{
			Name:     CMD_FLAG_NAME_DSN,
			Usage:    "data source name of database",
			Required: true,
		},
		&cli.StringFlag{
			Name:  CMD_FLAG_NAME_TABLE,
			Usage: "table name",
			Value: DEFAULT_TABLE_NAME,
		},
	},
	Action: func(cctx *cli.Context) error {

		return nil
	},
}
