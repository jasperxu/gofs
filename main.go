package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// config 配置信息
var config = new(Config)

func init() {
	// 修正Win下的目录问题
	file, _ := exec.LookPath(os.Args[0])
	dir, _ := path.Split(strings.Replace(file, "\\", "/", -1))
	os.Chdir(dir)
}

func main() {
	readConfig()

	http.HandleFunc("/", index)

	log.Println("[gofs] Start at ", config.URL)
	if config.IsHTTPS {
		http.ListenAndServeTLS(":"+config.Port, "./server.crt", "./server.key", nil)
	} else {
		http.ListenAndServe(":"+config.Port, nil)
	}
	log.Println("[gofs] Stopped.")
}

// index 页面
// Get, [key] 显示上传页面"/"
// Get, [key] 获取指定文件"/path/file"
// Delete, [key] 删除指定文件或文件夹"/path/file"或"/path"
// Post, [key], path, upload 为上传文件并保存到path路径
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method == "GET" {
		// 判断读权限
		if config.ReadKey != "" && config.ReadKey != r.URL.RawQuery {
			Output{Status: "error", Message: config.I18nNoAccessRights, Data: nil}.Writer2Response(w)
			return
		}

		if r.URL.Path == "/" {
			// 显示上传页面
			uploadPage(w, r)
		} else {
			// 获取路径的文件
			getPath(w, r)
		}
	} else {
		// 判断写权限
		if config.WriteKey != "" && (config.WriteKey != r.URL.RawQuery || config.WriteKey != r.FormValue("key")) {
			Output{Status: "error", Message: config.I18nNoAccessRights, Data: nil}.Writer2Response(w)
			return
		}

		if r.Method == "DELETE" {
			// 删除path路径的文件或文件夹
			deletePath(w, r)
			return
		}
		if r.Method == "POST" {
			// 上传文件并保存到path路径
			uploadFile(w, r)
			return
		}
	}
}

// uploadPage ,Get, [key] 为显示上传页面
func uploadPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<body>

<form action="/" method="POST" enctype="multipart/form-data">
key:<br /><input type="text" name="key"><br />
path:<br /><input type="text" name="path"><br />
file:<br /><input type="file" name="uploadfile"><br />
<input type="submit" value="Submit">
</form>

</body>
</html>`
	fmt.Fprintln(w, html)
}

// getPath ,Get, [key] 为获取path路径的文件
func getPath(w http.ResponseWriter, r *http.Request) {
	// 获取文件
	filePath := "./upload" + r.URL.Path
	info, err := os.Stat(filePath)
	if err != nil {
		// 路径不存在
		http.NotFound(w, r)
	} else if info.IsDir() {
		// 是文件夹
		http.NotFound(w, r)
	} else {
		// 文件存在,输出文件
		http.ServeFile(w, r, filePath)
	}
}

// deletePath ,Delete, [key] 为删除path路径的文件或文件夹
func deletePath(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	if os.RemoveAll("./upload"+filePath) != nil {
		Output{Status: "error", Message: config.I18nDeleteError, Data: nil}.Writer2Response(w)
		return
	}

	Output{Status: "success", Message: config.I18nDeleteSuccess, Data: nil}.Writer2Response(w)
}

// uploadFile ,Post, [key], path, upload 为上传文件并保存到path路径
func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(config.MaxFileSize)

	filePath := r.FormValue("path") // 必须以/开头
	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		Output{Status: "error", Message: config.I18nUploadError, Data: nil}.Writer2Response(w)
		return
	}
	defer file.Close()

	fileUploadExt := strings.ToLower(filepath.Ext(header.Filename))
	if check(fileUploadExt) == false {
		Output{Status: "error", Message: config.I18nProhibitExt, Data: nil}.Writer2Response(w)
		return
	}
	filePathExt := strings.ToLower(filepath.Ext(filePath))
	if check(filePathExt) == false {
		Output{Status: "error", Message: config.I18nProhibitExt, Data: nil}.Writer2Response(w)
		return
	}

	// 处理Dir
	fileFullName := "./upload" + filePath
	err = os.MkdirAll(filepath.Dir(fileFullName), 0777)

	// 创建文件并保存
	f, err := os.Create(fileFullName)
	if err != nil {
		Output{Status: "error", Message: config.I18nCreateError, Data: nil}.Writer2Response(w)
		return
	}
	size, err := io.Copy(f, file)
	if err != nil {
		Output{Status: "error", Message: config.I18nSaveError, Data: nil}.Writer2Response(w)
		return
	}
	f.Close()

	data := struct {
		FileName string
		FullURL  string
		FileSize int64
	}{
		FileName: filepath.Base(fileFullName),
		FullURL:  config.URL + filePath,
		FileSize: size,
	}

	Output{Status: "success", Message: config.I18nUploadSuccess, Data: data}.Writer2Response(w)
}

// check 判断后缀是否被禁止上传。name为小写
func check(name string) bool {
	for _, v := range config.ProhibitExt {
		if v == name {
			return false
		}
	}
	return true
}
