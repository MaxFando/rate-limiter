package cli

import (
	"context"
	"fmt"
	"github.com/MaxFando/rate-limiter/internal/domain/network"
)

func (c *CommandLineInterface) bucketHandler(ctx context.Context, setCommand []string) {
	if len(setCommand) != 4 {
		return
	}
	if setCommand[1] == "reset" {
		request, err := network.NewRequest(
			setCommand[2],
			"",
			setCommand[3],
		)
		if err != nil {
			fmt.Printf("bootExecutor - bucket reset: %w", err)
			return
		}

		c.resetBucket(ctx, request)
	} else {
		fmt.Println("unknown command")
	}
}

func (c *CommandLineInterface) resetBucket(ctx context.Context, request network.Request) {
	isLoginReset, isIpReset, err := c.bucketUseCase.Reset(ctx, request.Login, request.Ip.String())
	if err != nil {
		fmt.Printf("service error: %v \n", err)
		return
	}

	if !isLoginReset {
		fmt.Printf("login: %v has not been reseted\n", request.Login)
	} else {
		fmt.Printf("login: %v has been reseted\n", request.Login)
	}

	if !isIpReset {
		fmt.Printf("ip: %v has not been reseted\n", request.Ip)
	} else {
		fmt.Printf("ip: %v has been reseted\n", request.Ip)
	}
}
