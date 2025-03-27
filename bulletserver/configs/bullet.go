package configs

const (
	GL_Drawingfile     = "GL_Drawing.txt"
	DrawingSpeed_Max   = 25
	DrawingSpeed_Heavy = 15
	DrawingSpeed_Light = 10
)

type Bullet struct {
	DurationSecs int
	Size         [2]int
	Color        [4]int
	Position     [3]int
}
