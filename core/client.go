package core

import (
	"context"
	"reflect"

	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/pkg/errors"
)

type Event struct {
	Handler interface{}
}

type Module struct {
	Name, Emote string

	Commands []*Command
	Events   []*Event

	OnInit, OnLogin func()
}

var (
	State *state.State
	Self  *discord.User

	Commands = make(map[string]*Command)
	Modules  = []*Module{}
)

func NewClient(token string) {
	State = state.NewWithIntents("Bot "+token, gateway.IntentGuilds, gateway.IntentGuildVoiceStates)

	State.AddHandler(InteractionEvent)
}

func Connect() error {
	err := State.Open(context.Background())

	if err == nil {
		Self, err = State.Me()
	}

	if err == nil {
		for _, module := range Modules {
			if module.OnLogin != nil {
				go module.OnLogin()
			}
		}
	}

	return err
}

func DeployCommands() error {
	app, err := State.CurrentApplication()
	if err != nil {
		return errors.Wrap(err, "unable to get application information")
	}

	previous, err := State.Commands(app.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to get Discord command list")
	}

	checked := make(map[string]interface{})
	for _, prevCmd := range previous {
		newCmd := Commands[prevCmd.Name]

		if newCmd == nil {
			logger.DebugF("Removendo comando \"%s\" do Discord.", prevCmd.Name)
			if err := State.DeleteCommand(app.ID, prevCmd.ID); err != nil {
				return errors.Wrapf(err, "failed to delete \"%s\" command", prevCmd.Name)
			}
		} else {
			if !reflect.DeepEqual(prevCmd.Options, newCmd.Options) || newCmd.Description != prevCmd.Description {
				logger.DebugF("Atualizando commando %s no Discord.", newCmd.Name)
				if _, err := State.EditCommand(app.ID, prevCmd.ID, newCmd.RAW()); err != nil {
					return errors.Wrapf(err, "failed to update \"%s\" command", newCmd.Name)
				}
			}
			checked[newCmd.Name] = true
		}
	}

	for _, command := range Commands {
		if checked[command.Name] == nil {
			logger.DebugF("Criando comando %s no Discord.", command.Name)
			if _, err := State.CreateCommand(app.ID, command.RAW()); err != nil {
				return errors.Wrapf(err, "failed to create \"%s\" command", command.Name)
			}
		}
	}

	return nil
}

func Close() {
	State.Close()
}

func AddModules(modules ...*Module) {
	for _, module := range modules {
		AddModule(module)
	}
}

func AddModule(module *Module) {
	for _, cmd := range module.Commands {
		cmd.Module = module
		Commands[cmd.Name] = cmd
	}

	for _, event := range module.Events {
		State.AddHandler(event.Handler)
	}

	if module.OnInit != nil {
		go module.OnInit()
	}

	Modules = append(Modules, module)
}
