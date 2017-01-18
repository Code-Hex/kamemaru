package util

var imageList = map[string]func([]byte) bool{
	"jpg": IsJPEG,
	"png": IsPNG,
	"gif": IsGIF,
	"bmp": IsBMP,
}

func IsJPEG(buf []byte) bool {
	return len(buf) > 2 &&
		buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF
}

func IsPNG(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47
}

func IsGIF(buf []byte) bool {
	return len(buf) > 2 &&
		buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46
}

func IsBMP(buf []byte) bool {
	return len(buf) > 1 &&
		buf[0] == 0x42 && buf[1] == 0x4D
}

func IsImage(buf []byte) bool {
	for _, k := range imageList {
		if imageList[k](buf) {
			return true
		}
	}
	return false
}
