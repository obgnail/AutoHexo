package auto_hexo

import (
	"log"
	"os"
	"strings"

	"github.com/obgnail/AutoHexo/auto_hexo/strategy"
	handlerConfig "github.com/obgnail/MarkdownResouceCollecter/config"
	"github.com/obgnail/MarkdownResouceCollecter/handler"
)

const (
	markdownFileSuffix     = ".md"
	localPictureUseAbsPath = true
)

type AutoHexo struct {
	OriginMarkdownRootDir string // 监控的md目录
	BlogMarkdownRootDir   string // 用于生成blog的目录
	BlogResourceRootDir   string // blog的资源目录
	hexoCommandPath       string // hexo命令的位置

	hexoCommand *HexoCommand
}

func New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCommandPath string) *AutoHexo {
	return &AutoHexo{
		OriginMarkdownRootDir: originMarkdownRootDir,
		BlogMarkdownRootDir:   blogMarkdownRootDir,
		BlogResourceRootDir:   blogResourceRootDir,
		hexoCommandPath:       hexoCommandPath,
		hexoCommand:           NewHexoCommand(hexoCommandPath, blogResourceRootDir),
	}
}

func (ah *AutoHexo) DeleteBlog(deleteFilePath string) error {
	if err := os.Remove(deleteFilePath); err != nil {
		return err
	}
	log.Println("[INFO] file remove:", deleteFilePath)
	if err := ah.HexoAutoDeploy(); err != nil {
		return err
	}
	return nil
}

func (ah *AutoHexo) DeleteAllBlogs() error {
	return ah.DeleteBlog(ah.BlogMarkdownRootDir)
}

func (ah *AutoHexo) Run(changedFilePath string) error {
	h := ah.newHandler(
		changedFilePath, ah.OriginMarkdownRootDir, ah.BlogResourceRootDir, ah.BlogMarkdownRootDir,
		markdownFileSuffix, localPictureUseAbsPath,
	)
	if err := h.Run(); err != nil {
		return err
	}
	if err := ah.HexoAutoDeploy(); err != nil {
		return err
	}
	return nil
}

func (ah *AutoHexo) AutoDeploy() {
	if err := ah.Run(ah.OriginMarkdownRootDir); err != nil {
		log.Println("[Error]: deploy blogs error:", err)
	}
}

func (ah *AutoHexo) HexoGenerate() error {
	log.Println("[INFO] hexo Generate...")
	return ah.hexoCommand.ExecuteHexoGenerate()
}

func (ah *AutoHexo) HexoDeploy() error {
	log.Println("[INFO] hexo Deploying...")
	return ah.hexoCommand.ExecuteHexoDeploy()
}

func (ah *AutoHexo) HexoClean() error {
	log.Println("[INFO] hexo Cleaning...")
	return ah.hexoCommand.ExecuteHexoClean()
}

func (ah *AutoHexo) HexoAutoDeploy() error {
	log.Println("[INFO] hexo AutoDeploying...")
	if err := ah.HexoGenerate(); err != nil {
		return err
	}
	if err := ah.HexoDeploy(); err != nil {
		return err
	}
	log.Println("[INFO] hexo AutoDeployed")
	return nil
}

func (ah *AutoHexo) newHandler(
	changedFilePath, originMarkdownRootDir, blogResourceRootDir, blogMarkdownRootDir string,
	markdownFileSuffix string, localPictureUseAbsPath bool,
) *handler.BaseHandler {
	blogMarkdownPath := strings.Replace(changedFilePath, originMarkdownRootDir, blogMarkdownRootDir, 1)
	cfg := handlerConfig.InitConfig(
		changedFilePath, blogResourceRootDir, blogMarkdownPath,
		nil, markdownFileSuffix, localPictureUseAbsPath,
	)
	h := handler.New(cfg)
	h.AppendStrategy(&handler.CollectNetWorkPictureStrategy{})
	h.AppendStrategy(&handler.CollectLocalPictureStrategy{})
	h.AppendStrategy(&strategy.InsertHexoHeaderStrategy{BlogMarkdownRootDir: blogMarkdownRootDir})
	h.AppendStrategy(&strategy.FixHexoPicturePathStrategy{})
	return h
}
