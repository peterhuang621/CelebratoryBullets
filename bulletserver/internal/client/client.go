package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/proto/gen"
	"google.golang.org/grpc"
)

const GeneratingSpeed_Max = 10

type Client_gRPC_cl struct {
	gen.BulletServiceClient
}

type Client struct {
	Total   int
	rng     *rand.Rand
	GRPC_cl *Client_gRPC_cl
}

func (cl *Client) Init() {
	cl.Total = 0
	cl.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func NewClient(conn *grpc.ClientConn) *Client {
	cl := &Client{}
	cl.Init()
	cl.GRPC_cl = &Client_gRPC_cl{gen.NewBulletServiceClient(conn)}
	return cl
}

func (cl *Client) constructBullets(num int) (bullets []configs.Bullet) {
	log.Printf("Client Generating %d bullets...", num)
	for i := 0; i < num; i++ {
		bullets = append(bullets, configs.Bullet{
			DurationSecs: cl.rng.Intn(configs.DurationSecs_Max) + 1,
			Size:         cl.rng.Intn(configs.Size_Max) + 1,
			Color:        [4]float32{cl.rng.Float32(), cl.rng.Float32(), cl.rng.Float32(), cl.rng.Float32()},
			Position:     [3]int{cl.rng.Intn(configs.Position_Max + 1), cl.rng.Intn(configs.Position_Max + 1), cl.rng.Intn(configs.Position_Max + 1)},
		})
	}
	cl.Total += num
	return
}

func (cl *Client) SendingBullets(num int) {
	if num > GeneratingSpeed_Max {
		log.Fatalf("Bullets number request (%d) is out of limit (%d)! NO BULLETS will be sent!", num, GeneratingSpeed_Max)
		return
	}
	bullets := cl.constructBullets(num)
	serverurl := fmt.Sprintf("http://localhost:%s/%s", configs.Server_Port, configs.Server_Addr)

	sendingjsondata, err := json.Marshal(&bullets)
	pkg.WarnHandle(err, "Failed on Marshaling bullets into json format")

	resp, err := http.Post(serverurl, "application/json", bytes.NewBuffer(sendingjsondata))
	pkg.WarnHandle(err, "Failed on Posting to the server")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("server returned status: %d, body: %s", resp.StatusCode, string(body))
	}
}

func (cl *Client) Random_SendingBullets() {
	cl.SendingBullets(cl.rng.Intn(GeneratingSpeed_Max) + 1)
}
