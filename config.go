package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	AutoType              string `json:"auto_type"`
	OriginMarkdownRootDir string `json:"origin_markdown_root_dir"`
	BlogMarkdownRootDir   string `json:"blog_markdown_root_dir"`
	BlogResourceRootDir   string `json:"blog_resource_root_dir"`
	HexoCmdPath           string `json:"hexo_cmd_path"`
	WaitingWindows        int    `json:"waiting_windows_minute"`
	TickerTime            int    `json:"ticker_time_hour"`
}

func ReadConfig(filePath string) *Config {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("read config.json failed")
	}
	cfg := &Config{}
	if err := json.Unmarshal(file, cfg); err != nil {
		log.Fatalln("unmarshal config.json failed", err)
	}
	return cfg
}
