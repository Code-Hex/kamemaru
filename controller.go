package kamemaru

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Code-Hex/kamemaru/internal/database"
	"github.com/Code-Hex/saltissimo"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/sync/errgroup"
)

type (
	User struct {
		Username string `validate:"required"`
		Password string `validate:"min=8,max=16"`
	}

	Image struct {
		Offset int `query:"offset"`
	}

	filedata struct {
		buf       []byte
		filename  string
		extension string
	}
)

func (k *kamemaru) imgfetch(c echo.Context) error {
	img := new(Image)
	if err := c.Bind(img); err != nil {
		return c.JSON(http.StatusOK, whyError(err))
	}

	if img.Offset <= 0 {
		img.Offset = 1
	}

	imgs := database.GetImages(k.DB, 20, img.Offset)

	return c.JSON(http.StatusOK, imgs)
}

// success: 201 failed: 409
func (k *kamemaru) register(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusInternalServerError, whyError(err))
	}
	if err = c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, whyError(err))
	}

	username, password := u.Username, u.Password

	if database.IsExistUser(k.DB, username) {
		return c.JSON(http.StatusConflict, whyError(fmt.Errorf("Already exist user:%s", username)))
	}

	var udb database.User
	udb.Pass, udb.Salt, err = saltissimo.HexHash(sha256.New, password)
	if err != nil {
		return c.JSON(http.StatusConflict, whyError(fmt.Errorf("Failed to create hash:%s", err.Error())))
	}
	udb.Name = username

	if err = k.DB.Create(&udb).Error; err != nil {
		return c.JSON(http.StatusConflict, whyError(fmt.Errorf("Failed to create user:%s", err.Error())))
	}

	t, err := k.createToken(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, whyError(err))
	}
	return c.JSON(http.StatusCreated, echo.Map{"status": "success", "token": t})
}

func (k *kamemaru) login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusInternalServerError, whyError(err))
	}
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, whyError(err))
	}

	username, password := u.Username, u.Password

	dbu := database.PickUser(k.DB, username)

	isSame, err := saltissimo.CompareHexHash(sha256.New, password, dbu.Pass, dbu.Salt)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, whyError(err))
	}

	if isSame {
		t, err := k.createToken(username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, whyError(err))
		}
		return c.JSON(http.StatusOK, echo.Map{"status": "success", "token": t})
	}
	return c.JSON(http.StatusUnauthorized, whyError(fmt.Errorf("invalid user")))
}

func (k *kamemaru) Upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, whyError(err))
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, whyError(fmt.Errorf("The request form is empty")))
	}

	g, ctx := errgroup.WithContext(context.Background())

	fdata := make(chan filedata)
	fh := make(chan *multipart.FileHeader)

	g.Go(ReadUploadedFiles(ctx, c, fdata, fh))
	g.Go(k.StoreImageFiles(ctx, c, fdata))
	for _, file := range files {
		fh <- file
	}
	close(fh)

	if err := g.Wait(); err != nil {
		// return json error
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "success"})
}

func (k *kamemaru) createToken(username string) (string, error) {
	dbu := database.PickUser(k.DB, username)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = dbu.ID
	claims["user"] = dbu.Name
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 36).Unix()

	return token.SignedString(k.JWTSecret)
}

func whyError(err error) echo.Map {
	return echo.Map{
		"status": "failed",
		"reason": err.Error(),
	}
}
