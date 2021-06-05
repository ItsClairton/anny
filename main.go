package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/commands/anime"
	"github.com/ItsClairton/Anny/commands/image"
	"github.com/ItsClairton/Anny/commands/misc"
	"github.com/ItsClairton/Anny/listeners"
	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		logger.ErrorAndExit("Um erro ocorreu ao carregar o arquivo .env de configurações. (%s)", err.Error())
	}

	err = base.Init(os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		logger.ErrorAndExit("Um erro ocorreu ao criar o cliente do discord. (%s)", err.Error())
	}

	base.AddHandler(listeners.MessageCreateListener)

	base.AddCommand(&misc.PingCommand)

	base.AddCommand(&image.CatCommand)
	base.AddCommand(&image.NekoCommand)

	base.AddCommand(&anime.SceneCommand)
	base.AddCommand(&anime.AnimeCommand)

	err = base.Connect()

	if err != nil {
		logger.ErrorAndExit("Um erro ocorreu ao tentar se conectar ao Discord. (%s)", err.Error())
	}

	logger.Info("Conexão com o Discord feita com sucesso, Yay.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	base.Disconnect()
	logger.Info("Processo finalizado com sucesso, Yay.")
}
