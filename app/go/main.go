package main

import (
	"database/sql"
	"net/http"
	"problem1/configs"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "minimal_sns_app")
	})

	e.GET("/get_friend_list", getFriendList)

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		// FIXME
		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}

func getFriendList(c echo.Context) error {
	// ページと1ページあたりの項目数をクエリパラメータから取得
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "id is required",
		})
	}
	// クエリを実行

	q := `
		SELECT
			fl.id,
			fl.user1_id,
			u1.name AS user1_name,
			fl.user2_id,
			u2.name AS user2_name
		FROM
			friend_link AS fl
		JOIN
			users AS u1 ON fl.user1_id = u1.user_id
		JOIN
			users AS u2 ON fl.user2_id = u2.user_id
		WHERE
			fl.user1_id = ? OR fl.user2_id = ?;
		`
	rows, err := db.Query(q, id, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error",
		})
	}

	// 結果をJSONで返す
	return c.JSON(http.StatusOK, map[string]interface{}{
		"friend_list": rows,
	})
}
