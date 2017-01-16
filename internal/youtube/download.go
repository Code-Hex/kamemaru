package youtube

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rylio/ytdl"
)

func Download(url string, status chan<- float64) error {
	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return err
	}

	file, err := os.Create(vid.Title + ".mp4")
	if err != nil {
		return err
	}
	defer file.Close()

	return download(vid, file, status)
}

func download(vid *ytdl.VideoInfo, dest *os.File, status chan<- float64) error {
	format := vid.Formats.Worst(ytdl.FormatVideoEncodingKey)[0]
	u, err := vid.GetDownloadURL(format)
	if err != nil {
		return err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}
	_, err = io.Copy(dest, newReaderWithPercent(resp, status))
	close(status)

	return err
}
