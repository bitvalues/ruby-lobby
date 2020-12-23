package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bitvalues/ruby-lobby/src/config"
	"github.com/bitvalues/ruby-lobby/src/servers/auth"
	"github.com/bitvalues/ruby-lobby/src/servers/data"
	"github.com/bitvalues/ruby-lobby/src/servers/view"
	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/sirupsen/logrus"
)

var signalChannel chan os.Signal
var log *logrus.Logger

func init() {
	// setup a channel for receiving interrupt signals
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	// generate a new seed
	rand.Seed(time.Now().Unix())

	// parse the log level from CLI
	logLevel := flag.Int("v", 0, "the level of logging verbosity: 0 = errors, 1 = warnings, 2 = info, 3+ = debug")
	flag.Parse()

	// setup a new logger
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// set the logging level based on the CLI flag value
	switch *logLevel {
	case 0:
		log.SetLevel(logrus.ErrorLevel)
	case 1:
		log.SetLevel(logrus.WarnLevel)
	case 2:
		log.SetLevel(logrus.InfoLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	// get the configuration
	cfg := config.GetConfig()

	// create a new session manager
	sessionManager := sessions.CreateNewSessionManager(log)

	// create all of our servers
	authServer := auth.NewServer(cfg, sessionManager, log)
	dataServer := data.NewServer(cfg, sessionManager, log)
	viewServer := view.NewServer(cfg, sessionManager, log)

	// make sure everything shuts down when we're done
	defer authServer.Shutdown()
	defer dataServer.Shutdown()
	defer viewServer.Shutdown()

	// start everything up
	go sessionManager.StartCleanup()
	go authServer.Startup()
	go dataServer.Startup()
	go viewServer.Startup()

	// wait for an interrupt signal
	<-signalChannel
}
