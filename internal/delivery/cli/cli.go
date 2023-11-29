package cli

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/providers"
	authUC "github.com/MaxFando/rate-limiter/internal/usecase/auth"
	blacklistUC "github.com/MaxFando/rate-limiter/internal/usecase/blacklist"
	bucketUC "github.com/MaxFando/rate-limiter/internal/usecase/bucket"
	whiteListUC "github.com/MaxFando/rate-limiter/internal/usecase/whitelist"
	"github.com/c-bata/go-prompt"
	"log"
	"os"
	"strings"
)

var suggestions = []prompt.Suggest{
	{Text: "blacklist add [ip_address] [mask]", Description: "Add ip net to blacklist"},
	{Text: "blacklist remove [ip_address] [mask]", Description: "Remove ip net to blacklist"},
	{Text: "blacklist get", Description: "Get ip list from blacklist"},
	{Text: "whitelist add [ip_address] [mask]", Description: "Add ip net to whitelist"},
	{Text: "whitelist remove [ip_address] [mask]", Description: "Remove ip net to whitelist"},
	{Text: "whitelist get", Description: "Get ip list from whitelist"},
	{Text: "bucket remove [login] [ip_address]", Description: "Remove login and ip address from bucket"},
	{Text: "help", Description: "Display list of commands"},
	{Text: "exit", Description: "Exit anti bruteforce app"},
}

type CommandLineInterface struct {
	authUseCase      *authUC.UseCase
	blackListUseCase *blacklistUC.UseCase
	whiteListUseCase *whiteListUC.UseCase
	bucketUseCase    *bucketUC.UseCase
}

func New(ctx context.Context) *CommandLineInterface {
	useCaseProvider := ctx.Value(providers.UseCaseProviderKey).(*providers.UseCaseProvider)

	return &CommandLineInterface{
		authUseCase:      useCaseProvider.AuthUseCase,
		blackListUseCase: useCaseProvider.BlackListUseCase,
		whiteListUseCase: useCaseProvider.WhiteListUseCase,
		bucketUseCase:    useCaseProvider.BucketUseCase,
	}
}

func (c *CommandLineInterface) Run(ctx context.Context, interruptCh chan os.Signal) {
	executor := prompt.Executor(func(s string) {
		s = strings.TrimSpace(s)
		setCommand := strings.Split(s, " ")

		switch setCommand[0] {
		case "blacklist":
			c.blackListHandler(ctx, setCommand)
		case "whitelist":
			c.whiteListHandler(ctx, setCommand)
		case "bucket":
			c.bucketHandler(ctx, setCommand)
		case "exit":
			interruptCh <- os.Interrupt
			return
		case "help":
			for _, suggestion := range suggestions {
				fmt.Println("Command:", suggestion.Text, "Description:", suggestion.Description)
			}

		default:
			fmt.Println("unknown command")
		}
	})

	completer := prompt.Completer(func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}
		return prompt.FilterHasPrefix(suggestions, w, true)
	})

	defer func() {
		if a := recover(); a != nil {
			log.Println("Command line interface not available. Please run container with tty mode")
		}
	}()
	prompt.New(executor, completer).Run()
}
