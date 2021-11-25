package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/events"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/joho/godotenv"

	_ "github.com/ItsClairton/Anny/interactions"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Warn("Arquivo .env não encontrado, utilizando variaveis de ambiente fornecidas via linha de comando.")
		} else {
			logger.Fatal("Um erro ocorreu ao carregar as variaveis de ambiente do .env.", err)
		}
	}

	if err := core.New(os.Getenv("DISCORD_TOKEN")); err != nil {
		logger.Fatal("Um erro ocorreu ao tentar se conectar ao Discord.", err)
	}

	core.AddHandler(events.OnReady)
	core.AddHandler(events.OnInteraction)
	core.AddHandler(events.VoiceServerUpdate)
	core.AddHandler(events.VoiceStateUpdate)

	if err := core.DeployInteractions(); err != nil {
		logger.Fatal("Um erro ocorreu ao fazer o deploy das interações para o Discord.", err)
	}

	logger.Info("Conexão com o Discord feita com Sucesso.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	core.Disconnect()
}
