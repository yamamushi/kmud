package session

import (
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type mapBuilder struct {
	width    int
	height   int
	depth    int
	data     [][][]mapTile
	userRoom types.Room
}

type mapTile struct {
	char  rune
	color types.Color
}

func (t *mapTile) toString() string {
	if t.char == ' ' {
		return string(t.char)
	}

	return types.Colorize(t.color, string(t.char))
}

func newMapBuilder(width int, height int, depth int) mapBuilder {
	var builder mapBuilder

	// Double the X/Y axis to account for extra space to draw exits
	width *= 2
	height *= 2

	builder.data = make([][][]mapTile, depth)

	for z := 0; z < depth; z++ {
		builder.data[z] = make([][]mapTile, height)
		for y := 0; y < height; y++ {
			builder.data[z][y] = make([]mapTile, width)
		}
	}

	builder.width = width
	builder.height = height
	builder.depth = depth

	for z := 0; z < depth; z++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				builder.data[z][y][x].char = ' '
			}
		}
	}

	return builder
}

func (b *mapBuilder) setUserRoom(room types.Room) {
	b.userRoom = room
}

func (b *mapBuilder) addRoom(room types.Room, x int, y int, z int) {
	x = x * 2
	y = y * 2

	addIfExists := func(dir types.Direction, x int, y int) {
		if x < 0 || y < 0 {
			return
		}

		if room.HasExit(dir) {
			b.data[z][y][x].addExit(dir)
		}
	}

	if b.userRoom.GetId() == room.GetId() {
		b.data[z][y][x].char = '*'
		b.data[z][y][x].color = types.ColorRed
	} else {
		b.data[z][y][x].color = types.ColorMagenta
		if room.HasExit(types.DirectionUp) && room.HasExit(types.DirectionDown) {
			b.data[z][y][x].char = '+'
		} else if room.HasExit(types.DirectionUp) {
			b.data[z][y][x].char = '^'
		} else if room.HasExit(types.DirectionDown) {
			b.data[z][y][x].char = 'v'
		} else {
			char := '#'

			/*
				count := len(model.CharactersIn(room.GetId()))
				if count < 10 {
					char = rune(strconv.Itoa(count)[0])
				}
			*/

			b.data[z][y][x].char = char
			b.data[z][y][x].color = types.ColorWhite
		}
	}

	addIfExists(types.DirectionNorth, x, y-1)
	addIfExists(types.DirectionNorthEast, x+1, y-1)
	addIfExists(types.DirectionEast, x+1, y)
	addIfExists(types.DirectionSouthEast, x+1, y+1)
	addIfExists(types.DirectionSouth, x, y+1)
	addIfExists(types.DirectionSouthWest, x-1, y+1)
	addIfExists(types.DirectionWest, x-1, y)
	addIfExists(types.DirectionNorthWest, x-1, y-1)
}

func (b *mapBuilder) toString() string {
	str := ""

	for z := 0; z < b.depth; z++ {
		var rows []string
		for y := 0; y < b.height; y++ {
			row := ""
			for x := 0; x < b.width; x++ {
				tile := b.data[z][y][x].toString()
				row = row + tile
			}
			rows = append(rows, row)
		}

		rows = utils.TrimLowerRows(rows)

		if b.depth > 1 {
			divider := types.Colorize(types.ColorWhite, "================================================================================\r\n")
			rows = append(rows, divider)
		}

		for _, row := range rows {
			str = str + row + "\r\n"
		}
	}

	return str
}

func (t *mapTile) addExit(dir types.Direction) {
	combineChars := func(r1 rune, r2 rune, r3 rune) {
		if t.char == r1 {
			t.char = r2
		} else {
			t.char = r3
		}
	}

	t.color = types.ColorBlue

	switch dir {
	case types.DirectionNorth:
		combineChars('|', '|', '|')
	case types.DirectionNorthEast:
		combineChars('\\', 'X', '/')
	case types.DirectionEast:
		combineChars('-', '-', '-')
	case types.DirectionSouthEast:
		combineChars('/', 'X', '\\')
	case types.DirectionSouth:
		combineChars('|', '|', '|')
	case types.DirectionSouthWest:
		combineChars('\\', 'X', '/')
	case types.DirectionWest:
		combineChars('-', '-', '-')
	case types.DirectionNorthWest:
		combineChars('/', 'X', '\\')
	default:
		panic("Unexpected direction given to mapTile::addExit()")
	}
}
