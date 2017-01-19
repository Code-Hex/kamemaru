package kamemaru

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"image"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Code-Hex/kamemaru/internal/database"
	"github.com/Code-Hex/kamemaru/internal/util"
	"github.com/Code-Hex/saltissimo"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/disintegration/imaging"
	"github.com/labstack/echo"
)

func ReadUploadedFiles(ctx context.Context, c echo.Context, fdata chan<- filedata, fh <-chan *multipart.FileHeader) func() error {
	return func() error {
		defer close(fdata)
		for file := range fh {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				src, err := file.Open()
				if err != nil {
					return c.JSON(http.StatusConflict, whyError(err))
				}
				defer src.Close()

				buf, err := ioutil.ReadAll(src)
				if err != nil {
					return c.JSON(http.StatusConflict, whyError(err))
				}
				t, ok := util.IsImage(buf)
				if !ok {
					return c.JSON(http.StatusConflict, whyError(fmt.Errorf("A file other than image files are included")))
				}
				fdata <- filedata{buf, file.Filename, t}
			}
		}

		return nil
	}
}

func (k *kamemaru) StoreImageFiles(ctx context.Context, c echo.Context, fdata <-chan filedata) func() error {
	return func() error {
		for f := range fdata {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				rand, err := saltissimo.RandomBytes(16)
				if err != nil {
					return c.JSON(http.StatusConflict, whyError(err))
				}

				hash := hex.EncodeToString(rand)
				path := filepath.Join("images", hash+"."+f.extension)
				dst, err := os.Create(path)
				if err != nil {
					return c.JSON(http.StatusConflict, whyError(err))
				}
				dst.Write(f.buf)
				dst.Close()

				user := c.Get("user").(*jwt.Token)
				claims := user.Claims.(jwt.MapClaims)

				img, _, err := image.Decode(bytes.NewReader(f.buf))
				if err != nil {
					return c.JSON(http.StatusConflict, whyError(err))
				}
				img400 := imaging.Fill(img, 400, 400, imaging.Center, imaging.Lanczos)
				path400 := filepath.Join("images", "400", hash+"."+f.extension)
				if err := imaging.Save(img400, path400); err != nil {
					return c.JSON(http.StatusConflict, whyError(err))
				}

				var imgdb database.Image
				imgdb.UserID = uint(claims["id"].(float64))
				split := strings.Split(f.filename, ".")
				if len(split) > 0 {
					imgdb.Name = split[0]
				} else {
					imgdb.Name = f.filename
				}
				imgdb.Hash = hash
				imgdb.Ext = f.extension
				imgdb.OriginalURL = "/" + path
				imgdb.Resize400URL = "/" + path400

				if err := k.DB.Create(&imgdb).Error; err != nil {
					return c.JSON(http.StatusConflict, whyError(fmt.Errorf("Failed to create user:%s", err.Error())))
				}
			}
		}
		return nil
	}
}
