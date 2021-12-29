package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsClairton/Anny/core"
	"github.com/ItsClairton/Anny/misc"
	"github.com/ItsClairton/Anny/music"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Fatal("Um erro ocorreu ao carregar as variaveis de ambiente.", err)
		}
		logger.Warn("Utilizando variaveis de ambiente fornecidas por linha de comando.")
	}

	if err := core.NewClient(os.Getenv("DISCORD_TOKEN")); err != nil {
		logger.Fatal("Um erro ocorreu ao criar um cliente do Discord.", err)
	}

	core.AddModules(music.Module, misc.Module)

	if err := core.Connect(); err != nil {
		logger.Fatal("Um erro ocorreu ao tentar se autenticar com o Discord.")
	}

	if err := core.DeployCommands(); err != nil {
		logger.Fatal("Um erro ocorreu ao fazer deploy dos comandos para o Discord.", err)
	}

	logger.Info("Conexão com o Discord feita com sucesso.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	core.Close()
}
