package main

import (
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
		logger.Fatal("Um erro ocorreu ao carregar o arquivo .env.", err)
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
