package main

import (
	"flag"
	"fmt"

	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/wire"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("main")

var (
	version bool
)

func main() {
	flag.BoolVar(&version, "version", false, "version")
	flag.Parse()

	if version {
		fmt.Println(utils.CurrentVersion)
		return
	}

	s, err := wire.InitializeServer()
	if err != nil {
		log.Error("failed to initialize server", zap.Error(err))
		panic(err)
	}

	if err := s.Listen(); err != nil {
		log.Error("exit with error", zap.Error(err))
	}

	log.Info("bye.")
}
