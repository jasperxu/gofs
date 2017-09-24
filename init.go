package main

import "github.com/BurntSushi/toml"

//region **************** Config Model Begin ****************

// Config 配置
type Config struct {
	Project   string
	Company   string
	Developer string
	IsDebug   bool

	IsHTTPS bool

	URL string

	I18nNoAccessRights string
	I18nFileNotFound   string

	I18nUploadError   string
	I18nProhibitExt   string
	I18nCreateError   string
	I18nSaveError     string
	I18nUploadSuccess string

	I18nDeleteError   string
	I18nDeleteSuccess string

	ProhibitExt []string
	MaxFileSize int64

	ReadKey  string
	WriteKey string
}

//endregion

// ReadConfig 加载配置文件
func readConfig() error {
	if _, err := toml.DecodeFile("./config.toml", config); err != nil {
		return err
	}
	return nil
}
