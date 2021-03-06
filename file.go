package main

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/file/global"
	"github.com/beewit/file/handler"
	"github.com/labstack/echo"
)

/**
文件上传服务器
*/
func main() {
	e := echo.New()

	e.Static("/page", "page")
	e.Static("/files", "files")
	e.File("/.well-known/pki-validation/fileauth.txt", "fileauth.txt")
	e.File("/MP_verify_3Z6AKFClzM8nQt3q.txt", "MP_verify_3Z6AKFClzM8nQt3q.txt")

	e.POST("/upload", handler.UploadFile, handler.Filter)
	e.POST("/upload/multi", handler.UploadMultipart, handler.Filter)

	utils.Open(global.Host)
	port := ":" + convert.ToString(global.Port)
	e.Logger.Fatal(e.Start(port))
}
