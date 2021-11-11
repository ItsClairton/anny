package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsClairton/Anny/base/discord"
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

	discord.Init(os.Getenv("DISCORD_TOKEN"))
	discord.Session.AddHandler(events.InteractionsEvent)
	discord.Session.AddHandler(events.ReadyEvent)

	if err := discord.Connect(); err != nil {
		logger.Fatal("Um erro ocorreu ao tentar se conectar ao Discord.", err)
	}
	if err := discord.UpdateInteractions(); err != nil {
		logger.Fatal("Um erro ocorreu ao obter a lista de interações do Discord.", err)
	}

	logger.Info("Conexão com o Discord feita com Sucesso.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	discord.Disconnect()
}
