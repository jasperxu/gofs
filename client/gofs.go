package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GoFS 配置
type GoFS struct {
	URL string

	ReadKey  string
	WriteKey string
}

// Upload 将本地文件 file 上传到服务器 path
func (gofs GoFS) Upload(file, path string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filepath.Base(file))
	if err != nil {
		return "", err
	}

	//打开文件句柄操作
	fh, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}

	// 添加其他参数
	bodyWriter.WriteField("key", gofs.WriteKey)
	bodyWriter.WriteField("path", path)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(gofs.URL, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 将返回消息转为 input
	var input struct {
		Status  string
		Message string
		Data    struct {
			FileName string
			FullURL  string
			FileSize int64
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&input)
	if err != nil {
		return "", err
	}
	if input.Status != "success" {
		return "", errors.New(input.Message)
	}

	return input.Data.FullURL, nil
}

// Download 将服务器上的文件 file 下载到本地 path
func (gofs GoFS) Download(file, path string) error {
	// 获取文件地址
	fullURL := strings.TrimRight(gofs.URL, "/") + file
	if gofs.ReadKey != "" {
		fullURL = fullURL + "?" + gofs.ReadKey
	}

	res, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}

// Delete 将服务器上的 path 删除。可以是文件也可以是目录！
func (gofs GoFS) Delete(path string) error {
	// 获取文件地址
	fullURL := strings.TrimRight(gofs.URL, "/") + path
	if gofs.ReadKey != "" {
		fullURL = fullURL + "?" + gofs.ReadKey
	}

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 将返回消息转为 input
	var input struct {
		Status  string
		Message string
		Data    interface{}
	}
	err = json.NewDecoder(resp.Body).Decode(&input)
	if err != nil {
		return err
	}
	if input.Status != "success" {
		return errors.New(input.Message)
	}

	return nil
}
