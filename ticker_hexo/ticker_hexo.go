package ticker_hexo

import (
	"github.com/obgnail/AutoHexo/auto_hexo"
	"log"
	"time"
)

const tickerTime = 24 * time.Hour

type TickerHexo struct {
	autoHexo   *auto_hexo.AutoHexo
	tickerTime time.Duration
}

func New(
	originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath string,
	tickerTime time.Duration,
) *TickerHexo {
	ah := auto_hexo.New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath)
	return &TickerHexo{autoHexo: ah, tickerTime: tickerTime}
}

func (th *TickerHexo) reset() {
	if err := th.autoHexo.DeleteAllBlogs(); err != nil {
		log.Println("[Error]: delete all blogs error:", err)
	}
	if err := th.autoHexo.HexoClean(); err != nil {
		log.Println("[Error]: hexo clean error:", err)
	}
	if err := th.autoHexo.Run(th.autoHexo.OriginMarkdownRootDir); err != nil {
		log.Println("[Error]: create all blogs error:", err)
	}
}

func (th *TickerHexo) AutoDeploy() {
	ticker := time.NewTicker(th.tickerTime)
	go func() {
		for range ticker.C {
			th.reset()
		}
	}()
	select {}
}
