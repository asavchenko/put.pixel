package characters

const BTW_CH_WIDTH = 3
const BTW_CH_HEIGHT = 10

type Text struct {
	text            string
	color           byte
	wH              int
	wW              int
	screenPositionX int
	screenPositionY int
	characters      [][]*Chr
}
