package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsClairton/Anny/base"
	"github.com/ItsClairton/Anny/commands/image"
	"github.com/ItsClairton/Anny/commands/miscellaneous"
	"github.com/ItsClairton/Anny/commands/utilities"
	"github.com/ItsClairton/Anny/i18n"
	"github.com/ItsClairton/Anny/listeners"
	"github.com/ItsClairton/Anny/logger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		logger.ErrorAndExit("Um erro ocorreu ao carregar o arquivo .env de configurações. (%s)", err.Error())
	}

	err = i18n.Load("./locales")
	if err != nil {
		logger.ErrorAndExit("Um erro ocorreu ao carregar os arquivos de tradução. (%s)", err.Error())
	}

	err = base.Init(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		logger.ErrorAndExit("Um erro ocorreu ao criar o cliente do discord. (%s)", err.Error())
	}

	base.AddHandler(listeners.MessageCreateListener)

	base.AddCategory(image.Category)
	base.AddCategory(miscellaneous.Category)
	base.AddCategory(utilities.Category)

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
