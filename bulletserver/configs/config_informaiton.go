package configs

import "fmt"

const (
	GL_Drawingfile        = "GL_DRAWING.txt"
	DrawingSpeed_Max      = 50
	DrawingSpeed_Heavy    = 30
	DrawingSpeed_Light    = 20
	DurationSecs_Max      = 5.0
	Size_Max              = 50.0
	Position_Max          = 500.0
	Server_Port           = "8866"
	Server_Addr           = "bullets"
	Client_Addr           = "bullets_line"
	MQ_Addr               = "amqp://guest:guest@localhost:5672/"
	MQ_QueueNumber        = 5
	GRPC_Addr             = "localhost:50051"
	Redis_Addr            = "localhost:6379"
	Redis_Key             = "bullets_count_num"
	Redis_Key_DefaultExpr = -1
)

type Bullet struct {
	DurationSecs float32    `json:"durationSecs"`
	Size         float32    `json:"size"`
	Color        [4]float32 `json:"color"`
	Position     [3]float32 `json:"position"`
}

func (bt *Bullet) Init() {
	bt.DurationSecs = 1
	bt.Size = 40
	bt.Color = [4]float32{0.5, 0.5, 0.5, 1.0}
	bt.Position = [3]float32{0, 0, 0}
}

func (b *Bullet) String() string {
	return fmt.Sprintf("%.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f %.3f",
		float32(b.DurationSecs),
		float32(b.Size),
		b.Color[0], b.Color[1], b.Color[2], b.Color[3],
		b.Position[0], b.Position[1], b.Position[2],
	)
}
