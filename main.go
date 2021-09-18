package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsClairton/Anny/base/discord"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	discord.Init(os.Getenv("DISCORD_TOKEN"))

	err = discord.Connect()
	if err != nil {
		panic(err)
	}

	logger.Info("Conexão com o Discord feita com Sucesso.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	discord.Disconnect()
}
