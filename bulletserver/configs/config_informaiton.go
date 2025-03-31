package configs

const (
	GL_Drawingfile        = "GL_Drawing.txt"
	DrawingSpeed_Max      = 50
	DrawingSpeed_Heavy    = 30
	DrawingSpeed_Light    = 20
	DurationSecs_Max      = 5
	Size_Max              = 50
	Position_Max          = 500
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
	DurationSecs int        `json:"durationSecs"`
	Size         int        `json:"size"`
	Color        [4]float32 `json:"color"`
	Position     [3]int     `json:"position"`
}

func (bt *Bullet) Init() {
	bt.DurationSecs = 1
	bt.Size = 40
	bt.Color = [4]float32{0.5, 0.5, 0.5, 1.0}
	bt.Position = [3]int{0, 0, 0}
}
