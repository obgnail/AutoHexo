package main

import (
	"log"
	"time"

	"github.com/obgnail/AutoHexo/auto_hexo"
	"github.com/obgnail/AutoHexo/notify_hexo"
	"github.com/obgnail/AutoHexo/ticker_hexo"
)

func main() {
	config := ReadConfig("config.json")
	var hexoDeployer HexoDeployer
	switch config.AutoType {
	case "manual":
		hexoDeployer = auto_hexo.New(
			config.OriginMarkdownRootDir,
			config.BlogMarkdownRootDir,
			config.BlogResourceRootDir,
			config.HexoCmdPath,
		)
	case "notify":
		hexoDeployer = notify_hexo.New(
			config.OriginMarkdownRootDir,
			config.BlogMarkdownRootDir,
			config.BlogResourceRootDir,
			config.HexoCmdPath,
			time.Duration(config.WaitingWindows)*time.Minute,
		)
	case "ticker":
		hexoDeployer = ticker_hexo.New(
			config.OriginMarkdownRootDir,
			config.BlogMarkdownRootDir,
			config.BlogResourceRootDir,
			config.HexoCmdPath,
			time.Duration(config.TickerTime)*time.Hour,
		)
	default:
		log.Fatalln("no such hexo deployer")
	}
	hexoDeployer.AutoDeploy()
}
