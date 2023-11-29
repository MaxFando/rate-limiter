package cli

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
)

const (
	addCommand    = "add"
	removeCommand = "remove"
	getCommand    = "get"
)

func (c *CommandLineInterface) blackListHandler(ctx context.Context, setCommand []string) {
	switch setCommand[1] {
	case addCommand:
		if len(setCommand) != 4 {
			break
		}

		ipNetwork, err := network.NewIPNetwork(setCommand[2], setCommand[3])
		if err != nil {
			fmt.Printf("bootExecutor - blacklist add: %s", err.Error())
			return
		}

		c.addIPToBl(ctx, ipNetwork)
	case removeCommand:
		if len(setCommand) != 4 {
			break
		}

		ipNetwork, err := network.NewIPNetwork(setCommand[2], setCommand[3])
		if err != nil {
			fmt.Printf("bootExecutor - blacklist remove: %s", err.Error())
			return
		}
		c.removeIPToBl(ctx, ipNetwork)
	case getCommand:
		c.getIPListFromBl(ctx)
	default:
		fmt.Println("unknown command")
	}
}

func (c *CommandLineInterface) addIPToBl(ctx context.Context, ipNet network.IPNetwork) {
	err := c.blackListUseCase.AddIP(ctx, ipNet)
	if err != nil {
		fmt.Printf("addIPToBl: %s", err.Error())
		return
	}

	fmt.Printf("add address: %v to blacklist", ipNet)
}

func (c *CommandLineInterface) removeIPToBl(ctx context.Context, ipNet network.IPNetwork) {
	err := c.blackListUseCase.RemoveIP(ctx, ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from blacklist \n", ipNet)
}

func (c *CommandLineInterface) getIPListFromBl(ctx context.Context) {
	list, err := c.blackListUseCase.GetIPList(ctx)
	if err != nil {
		return
	}

	for _, _network := range list {
		fmt.Println(_network)
	}
}
