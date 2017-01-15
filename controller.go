package kamemaru

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
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
