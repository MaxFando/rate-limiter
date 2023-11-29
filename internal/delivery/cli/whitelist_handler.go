package cli

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
)

func (c *CommandLineInterface) whiteListHandler(ctx context.Context, setCommand []string) {
	switch setCommand[1] {
	case "add":
		if len(setCommand) != 4 {
			break
		}

		ipNetwork, err := network.NewIPNetwork(setCommand[2], setCommand[3])
		if err != nil {
			fmt.Printf("bootExecutor - whitelist add: %s", err.Error())
			return
		}

		c.addIPToWl(ctx, ipNetwork)
	case "remove":
		if len(setCommand) != 4 {
			break
		}

		ipNetwork, err := network.NewIPNetwork(setCommand[2], setCommand[3])
		if err != nil {
			fmt.Printf("bootExecutor - whitelist remove: %s", err.Error())
			return
		}
		c.removeIPToWl(ctx, ipNetwork)
	case "get":
		c.getIPListFromWl(ctx)
	default:
		fmt.Println("unknown command")
	}
}

func (c *CommandLineInterface) addIPToWl(ctx context.Context, ipNet network.IPNetwork) {
	err := c.whiteListUseCase.AddIP(ctx, ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("add address: %v to whitelist \n", ipNet)
}

func (c *CommandLineInterface) removeIPToWl(ctx context.Context, ipNet network.IPNetwork) {
	err := c.whiteListUseCase.RemoveIP(ctx, ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from whitelist \n", ipNet)
}

func (c *CommandLineInterface) getIPListFromWl(ctx context.Context) {
	list, err := c.whiteListUseCase.GetIPList(ctx)
	if err != nil {
		return
	}
	for _, network := range list {
		fmt.Println(network)
	}
}
