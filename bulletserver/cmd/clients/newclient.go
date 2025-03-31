package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/internal/client"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func isNumber(s string) (int, bool) {
	num, err := strconv.Atoi(s)
	return num, err == nil
}

var conn *grpc.ClientConn
var err error

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err = grpc.NewClient(configs.GRPC_Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	pkg.FailOnError(err, fmt.Sprintf("Failed to dial %s to get a NewClient", configs.GRPC_Addr))
	defer conn.Close()

	cl := client.NewClient(conn)
	bulletList := &gen.BulletList{
		Bullets: []*gen.Bullet{
			{
				DurationSecs: 10,
				Size:         5,
				Color:        []float32{1.0, 0.0, 0.0, 1.0},
				Position:     []int32{10, 20, 30},
			},
		},
	}

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
						if x != 0 {
							cl.SendingBullets(x)
						} else {
							fmt.Println("Sending a dedicate bullet by gRPC")
							ack, err := cl.GRPC_cl.DirectDrawBullets(
								context.Background(),
								bulletList,
							)
							pkg.FailOnError(err, "Error message from Sever-gPRC")
							fmt.Printf("Server-gPRC response: %v\n", ack.Message)
						}

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
