package ticker_hexo

import (
	"github.com/obgnail/AutoHexo/auto_hexo"
	"log"
	"time"
)

const tickerTime = 24 * time.Hour

type TickerHexo struct {
	tickerTime time.Duration
	*auto_hexo.AutoHexo
}

func New(
	originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath string,
	tickerTime time.Duration,
) *TickerHexo {
	ah := auto_hexo.New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath)
	return &TickerHexo{tickerTime, ah}
}

func (th *TickerHexo) reset() {
	if err := th.DeleteAllBlogs(); err != nil {
		log.Println("[Error]: delete all blogs error:", err)
	}
	if err := th.CreateBlog(th.OriginMarkdownRootDir); err != nil {
		log.Println("[Error]: create all blogs error:", err)
	}
}

func (th *TickerHexo) Run() {
	ticker := time.NewTicker(th.tickerTime)
	go func() {
		for range ticker.C {
			th.reset()
		}
	}()
	select {}
}
