package cli

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
)

func (c *CommandLineInterface) blackListHandler(ctx context.Context, setCommand []string) {
	switch setCommand[1] {
	case "add":
		if len(setCommand) != 4 {
			break
		}

		ipNetwork, err := network.NewIpNetwork(setCommand[2], setCommand[3])
		if err != nil {
			fmt.Printf("bootExecutor - blacklist add: %w", err)
			return
		}

		c.addIpToBl(ctx, ipNetwork)
	case "remove":
		if len(setCommand) != 4 {
			break
		}

		ipNetwork, err := network.NewIpNetwork(setCommand[2], setCommand[3])
		if err != nil {
			fmt.Printf("bootExecutor - blacklist remove: %w", err)
			return
		}
		c.removeIpToBl(ctx, ipNetwork)
	case "get":
		c.getIpListFromBl(ctx)
	default:
		fmt.Println("unknown command")
	}
}

func (c *CommandLineInterface) addIpToBl(ctx context.Context, ipNet network.IpNetwork) {
	err := c.blackListUseCase.AddIP(ctx, ipNet)
	if err != nil {
		fmt.Printf("addIpToBl: %w", err)
		return
	}

	fmt.Printf("add address: %v to blacklist", ipNet)
}

func (c *CommandLineInterface) removeIpToBl(ctx context.Context, ipNet network.IpNetwork) {
	err := c.blackListUseCase.RemoveIP(ctx, ipNet)
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}
	fmt.Printf("remove address: %v from blacklist \n", ipNet)
}

func (c *CommandLineInterface) getIpListFromBl(ctx context.Context) {
	list, err := c.blackListUseCase.GetIPList(ctx)
	if err != nil {
		return
	}

	for _, _network := range list {
		fmt.Println(_network)
	}
}
