package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/proto/gen"
	"google.golang.org/grpc"
)

type BulletgRPCServiceServer struct {
	gen.UnimplementedBulletServiceServer
}

func initgRPC() {
	lis, err := net.Listen("tcp", configs.GRPC_Addr)
	pkg.FailOnError(err, fmt.Sprintf("(gRPC address: %s)", configs.GRPC_Addr))

	grpc_server = grpc.NewServer()

	gen.RegisterBulletServiceServer(grpc_server, &BulletgRPCServiceServer{})

	go func() {
		if err := grpc_server.Serve(lis); err != nil {
			pkg.FailOnError(err, "Failed to serve gRPC")
		}
	}()
}

func (s *BulletgRPCServiceServer) DirectDrawBullets(ctx context.Context, in *gen.BulletList) (*gen.Ack, error) {
	var bullets []configs.Bullet
	for id, b := range in.Bullets {
		if len(b.Color) != 4 || len(b.Position) != 3 {
			return &gen.Ack{
				Success: false,
				Message: fmt.Sprintf("Invaild bullets at #%d", id),
			}, fmt.Errorf("invalid bullet data received")
		}

		bullets = append(bullets, configs.Bullet{
			DurationSecs: b.DurationSecs,
			Size:         b.Size,
			Color:        [4]float32{b.Color[0], b.Color[1], b.Color[2], b.Color[3]},
			Position:     [3]float32{b.Position[0], b.Position[1], b.Position[2]},
		})
	}
	body, err := json.Marshal(bullets)
	pkg.FailOnError(err, "Failed to Marshal the input data on gRPC server")
	SendToMQwithRoutingKey(body, mqqueues[0])

	return &gen.Ack{
		Success: true,
		Message: "Successfully received bullets by gRPC server"}, nil
}
