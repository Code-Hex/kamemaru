package kamemaru

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Code-Hex/kamemaru/internal/youtube"
	"github.com/labstack/echo"
	"golang.org/x/sync/errgroup"
)

// コマンドラインからイベントの記録を投稿する api
/*
 * Method: POST /api/v1/list
 * Param:  {
 *             text:"Today's [link](http://example.com)!!",
 *             tags:["golang", "perl"],
 *         }
 */
func (k *kamemaru) List(c echo.Context) error {
	var list List
	json.NewDecoder(c.Request().Body).Decode(&list)
	return c.JSON(http.StatusOK, list)
}

type (
	YRequest struct {
		URL string `json:"url"`
	}
	Youtube struct {
		Reason string `json:"reason"`
	}
	Download struct {
		Percent float64 `json:"percent"`
	}
)

func (k *kamemaru) YoutubeDownload(c echo.Context) error {
	var (
		d    Download
		data Youtube
	)

	currentSize := make(chan float64)
	g, _ := errgroup.WithContext(context.Background())

	var req YRequest
	json.NewDecoder(c.Request().Body).Decode(&req)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	g.Go(func() error {
		return youtube.Download(req.URL, currentSize)
	})

	g.Go(func() error {
		for d.Percent = range currentSize {
			if err := json.NewEncoder(c.Response()).Encode(d); err != nil {
				return err
			}
			c.Response().Flush()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		data.Reason = err.Error()
		return c.JSON(http.StatusInternalServerError, data)
	}

	return c.JSON(http.StatusOK, data)
}
