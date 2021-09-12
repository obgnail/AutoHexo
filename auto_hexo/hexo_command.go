package auto_hexo

import (
	"fmt"
	"os/exec"
)

const (
	cmdHexoGenerate = "g"
	cmdHexoDeploy   = "d"
	cmdHexoServer   = "s"
)

type HexoCommand struct {
	hexoCmdPath string
	execDir     string
}

func NewHexoCommand(hexoCmdPath, execDir string) *HexoCommand {
	return &HexoCommand{hexoCmdPath: hexoCmdPath, execDir: execDir}
}

func (hc *HexoCommand) ExecuteCmd(cmd string) error {
	execCmd := exec.Command(hc.hexoCmdPath, cmd)
	execCmd.Dir = hc.execDir
	out, err := execCmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
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
