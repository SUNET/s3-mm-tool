package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/sunet/s3-mm-tool/pkg/api"
	"github.com/sunet/s3-mm-tool/pkg/manifest"
	"github.com/sunet/s3-mm-tool/pkg/storage"
)

var Log = logrus.New()

var helpFlag bool
var logLevelFlag string
var serverFlag bool

func usage(code int) {
	fmt.Println("usage: s3-mm-tool [-h] [-s]")
	os.Exit(code)
}

func is_not_tty() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func ConfigLoggers(logLevelFlag string) {
	configLogger(Log, logLevelFlag)
	configLogger(manifest.Log, logLevelFlag)
	configLogger(storage.Log, logLevelFlag)
	configLogger(api.Log, logLevelFlag)
}

func configLogger(log *logrus.Logger, ll string) {
	log.Out = os.Stdout

	if len(ll) > 0 {
		level, err := logrus.ParseLevel(logLevelFlag)
		if err != nil {
			log.Panicf("Unable to parse loglevel: %s", err.Error())
		}
		log.SetLevel(level)
	}
}

func main() {

	flag.Parse()
	if helpFlag {
		usage(0)
	}

	ConfigLoggers(logLevelFlag)

	defer func() {
		if r := recover(); r != nil {
			Log.Debug(r)
		}
	}()

	if serverFlag {
		api.Listen("0.0.0.0:3000")
	}

}

func init() {
	flag.BoolVar(&helpFlag, "h", false, "show help")
	flag.BoolVar(&serverFlag, "s", false, "run the s3-mm-tool api server on port 3000")
	flag.StringVar(&logLevelFlag, "loglevel", "info", "loglevel")
}
