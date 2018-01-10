package handler

import (
	"github.com/labstack/echo"
	"io"
	"fmt"
	"github.com/beewit/file/global"
	"github.com/beewit/beekit/utils"
	path2 "path"
	"github.com/beewit/beekit/utils/convert"
	"strings"
	"io/ioutil"
)

var b float64 = 1024

/*
	文件上传统一管理，方便维护
*/
func UploadFile(c echo.Context) error {
	// Read form fields
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	//放置目录
	dir := c.FormValue("dir")
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorNull(c, "获取文件失败")
	}
	if file.Size <= 0 {
		return utils.ErrorNull(c, "空文件")
	}
	if convert.MustFloat64(file.Size) > global.MaxFileSize {
		return utils.ErrorNull(c, fmt.Sprintf("文件过大超出限制%vmb", fmt.Sprintf("%.2f", global.MaxFileSize/b/b)))
	}
	src, err := file.Open()
	if err != nil {
		return utils.ErrorNull(c, "打开文件失败")
	}
	defer src.Close()

	buf, err := ioutil.ReadAll(src)
	if err != nil {
		return utils.ErrorNull(c, "获取文件格式错误")
	}
	ext := path2.Ext(file.Filename)
	if !strings.Contains(fmt.Sprintf(",%v,", global.ExtFilter), fmt.Sprintf(",%v%v,", buf[0], buf[1])) {
		return utils.ErrorNull(c, fmt.Sprintf("%v文件格式错误", ext))
	}

	// Destination
	fileName := convert.ToString(utils.ID())
	path := getPath(dir, fileName, ext, acc)
	dst, err := utils.CreateFile(path)
	if err != nil {
		return utils.ErrorNull(c, "创建文件失败")
	}
	defer dst.Close()

	// Copy
	if _, err = dst.Write(buf); err != nil {
		return utils.ErrorNull(c, "保存文件失败")
	}
	go func() {
		_, err = global.DB.InsertMap("file_log", map[string]interface{}{
			"id":      utils.ID(),
			"name":    file.Filename,
			"path":    path,
			"size":    file.Size,
			"ext":     ext,
			"ct_time": utils.CurrentTime(),
			"ct_ip":   c.RealIP(),
		})
		if err != nil {
			global.Log.Error(fmt.Sprintf("保存上传文件日志失败，ERROR：%s", err.Error()))
		}
	}()
	return utils.SuccessNullMsg(c, map[string]interface{}{
		"id":   fileName,
		"size": file.Size,
		"path": path,
		"url":  getUrl(path),
		"name": file.Filename,
	})
}

func UploadMultipart(c echo.Context) error {
	// Read form fields
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	//放置目录
	dir := c.FormValue("dir")
	if dir == "" {
		dir = "default"
	}
	//------------
	// Read files
	//------------
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorNull(c, "获取文件失败")
	}
	var newFiles []map[string]interface{}

	files := form.File["files"]
	filesMap := []map[string]interface{}{}
	for _, file := range files {
		if file.Size <= 0 {
			continue
		}
		if convert.MustFloat64(file.Size) > global.MaxFileSize {
			return utils.ErrorNull(c, fmt.Sprintf("文件(%s)过大超出限制%vmb", file.Filename, fmt.Sprintf("%.2f", global.MaxFileSize/b/b)))
		}
		// Source
		src, err := file.Open()
		if err != nil {
			return utils.ErrorNull(c, "打开文件失败")
		}
		defer src.Close()

		buf := make([]byte, 2)
		_, err = src.Read(buf)
		if err != nil {
			return utils.ErrorNull(c, "类型判断错误")
		}
		ext := path2.Ext(file.Filename)
		if !strings.Contains(fmt.Sprintf(",%v,", global.ExtFilter), fmt.Sprintf(",%v%v,", buf[0], buf[1])) {
			return utils.ErrorNull(c, fmt.Sprintf("%v文件格式错误", ext))
		}

		fileName := convert.ToString(utils.ID())
		path := getPath(dir, fileName, ext, acc)
		dst, err := utils.CreateFile(path)
		if err != nil {
			return utils.ErrorNull(c, "创建文件失败")
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return utils.ErrorNull(c, "保存文件失败")
		}
		newFiles = append(newFiles, map[string]interface{}{
			"id":   fileName,
			"size": file.Size,
			"path": path,
			"url":  getUrl(path),
			"name": file.Filename,
		})

		filesMap = append(filesMap, map[string]interface{}{
			"id":      utils.ID(),
			"name":    file.Filename,
			"path":    path,
			"size":    file.Size,
			"ext":     ext,
			"ct_time": utils.CurrentTime(),
			"ct_ip":   c.RealIP(),
		})
	}
	if len(newFiles) <= 0 {
		return utils.ErrorNull(c, "空文件")
	}

	go func() {
		_, err = global.DB.InsertMapList("file_log", filesMap)
		if err != nil {
			global.Log.Error(fmt.Sprintf("保存批量上传文件日志失败，ERROR：%s", err.Error()))
		}
	}()

	return utils.SuccessNullMsg(c, newFiles)
}

func getPath(dir, fileName, suffix string, acc *global.Account) string {
	return fmt.Sprintf("%s/%s/%d/%s%s", global.FilesPath, dir, acc.ID, fileName, suffix)
	//now := time.Now()
	//return fmt.Sprintf("%s/%s/%d/%d/%d/%s%s", global.FilesPath, dir, now.Year(), now.Month(), now.Day(), fileName, suffix)
}

func getUrl(path string) string {
	return global.FilesDoMain + path
}
