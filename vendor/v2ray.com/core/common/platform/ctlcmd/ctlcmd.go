package ctlcmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform"
)

//go:generate errorgen

func Run(args []string, input io.Reader) (buf.MultiBuffer, error) {
	v2ctl := platform.GetToolLocation("v2ctl")
	current_file, _ := exec.LookPath(os.Args[0])
	v2ray, _ := filepath.Abs(current_file)

	var errBuffer buf.MultiBufferContainer
	var outBuffer buf.MultiBufferContainer

	ctlargs := []string{"-ctl"}
	ctlargs = append(ctlargs, args[:]...)

	cmd := exec.Command(v2ray, ctlargs...)
	if _, err := os.Stat(v2ctl); err == nil {
		cmd = exec.Command(v2ctl, args...)
	}

	cmd.Stderr = &errBuffer
	cmd.Stdout = &outBuffer
	cmd.SysProcAttr = getSysProcAttr()
	if input != nil {
		cmd.Stdin = input
	}

	if err := cmd.Start(); err != nil {
		return nil, newError("failed to start v2ctl").Base(err)
	}

	if err := cmd.Wait(); err != nil {
		msg := "failed to execute" + v2ray
		if errBuffer.Len() > 0 {
			msg += ": " + errBuffer.MultiBuffer.String()
		}
		return nil, newError(msg).Base(err)
	}

	return outBuffer.MultiBuffer, nil
}
