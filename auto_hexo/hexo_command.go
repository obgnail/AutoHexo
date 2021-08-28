package auto_hexo

import (
	"os/exec"
)

const (
	cmdHexoGenerate = "g"
	cmdHexoDeploy   = "d"
	cmdHexoServer   = "s"
)

type HexoCommand struct {
	hexoCmdPath string
}

func NewHexoCommand(hexoCmdPath string) *HexoCommand {
	return &HexoCommand{hexoCmdPath: hexoCmdPath}
}

func (hc *HexoCommand) ExecuteCmd(cmd string) error {
	cdCmd := exec.Command(hc.hexoCmdPath, cmd)
	if _, err := cdCmd.Output(); err != nil {
		return err
	}
	return nil
}

func (hc *HexoCommand) ExecuteHexoServer() error {
	return hc.ExecuteCmd(cmdHexoServer)
}

func (hc *HexoCommand) ExecuteHexoGenerate() error {
	return hc.ExecuteCmd(cmdHexoGenerate)
}

func (hc *HexoCommand) ExecuteHexoDeploy() error {
	return hc.ExecuteCmd(cmdHexoDeploy)
}
