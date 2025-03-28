package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/internal/client"
)

func isNumber(s string) (int, bool) {
	num, err := strconv.Atoi(s)
	return num, err == nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl := client.NewClient()
	go func() {
		defer cancel()
		scanner := bufio.NewScanner(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Print("Enter command (Enter=Random, x=Send x bullets, other=Exit): ")
				if !scanner.Scan() {
					return
				}
				input := strings.TrimSpace(scanner.Text())
				switch input {
				case "":
					cl.Random_SendingBullets()

				default:
					x, isnum := isNumber(input)
					if isnum {
						cl.SendingBullets(x)
					} else {
						fmt.Println("Exiting...")
						return
					}
				}
			}
		}
	}()
	<-ctx.Done()
	fmt.Println("Program terminated.")
}
