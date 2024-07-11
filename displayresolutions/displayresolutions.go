package displayresolutions

import (
	"errors"
	"image"
)

var ErrDisplayIsNotSupported = errors.New("The display is not supported")

const GridInaccuracy = 5

const (
	DisplayWidth1920 = 1920
	DisplayWidth2560 = 2560
	DisplayWidth3440 = 3440
	DisplayWidth3840 = 3840
)

const (
	DisplayHeight1080 = 1080
	DisplayHeight1440 = 1440
	DisplayHeight2160 = 2160
)

type GridType int

const (
	gridType64  GridType = 64
	gridType84  GridType = 84
	gridType126 GridType = 126
)

// GetGridType размер грида, по скрину,
func GetGridType(point image.Point) (GridType, bool) {
	sizeGrid := gridType64 // 2560x1080 1920x1080
	switch true {
	case point.X == DisplayWidth1920 && point.Y == DisplayHeight1080:
		return gridType64, true
	case point.X == DisplayWidth2560 && point.Y == DisplayHeight1080:
		return gridType64, true
	case point.X == DisplayWidth2560 && point.Y == DisplayHeight1440:
		return gridType84, true
	case point.X == DisplayWidth3440 && point.Y == DisplayHeight1440:
		return gridType84, true
	case point.X == DisplayWidth3840 && point.Y == DisplayHeight2160:
		return gridType126, true
	}

	return sizeGrid, false
}

func GetHeightShortText(point image.Point) int {
	switch true {
	case point.X == DisplayWidth1920 && point.Y == DisplayHeight1080:
		return 16
	case point.X == DisplayWidth2560 && point.Y == DisplayHeight1080:
		return 16
	case point.X == DisplayWidth2560 && point.Y == DisplayHeight1440:
		return 19
	case point.X == DisplayWidth3440 && point.Y == DisplayHeight1440:
		return 19
	case point.X == DisplayWidth3840 && point.Y == DisplayHeight2160:
		return 26
	default:
		return 16
	}
}
