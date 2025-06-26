package baredroid

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type InstallType string

const (
	PlayStore InstallType = "playstore"
	FDroidCL  InstallType = "fdroidcl"
	Sideload  InstallType = "sideload"
)

var FDroidCLIndexSynced bool = false

type Device struct {
	ctx		context.Context
	timeout time.Duration
}
//

func NewDevice(timeout time.Duration) *Device {
	return &Device{
		ctx:	 context.Background(),
		timeout: timeout,
	}
}

func (d *Device) execCommand(bin string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(d.ctx, d.timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %v\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func (d *Device) RemovePackage(pkgName string) (string, error) {
	out, err := d.execCommand("adb", "shell", "pm", "uninstall", "--user", "0", pkgName)

	if err != nil {
		return out, fmt.Errorf("could not remove package %s: %s", pkgName, err)
	}

	return out, nil
}

func (d *Device) InstallPackage(pkgName string, pkg string, source string, installType InstallType) error {
	
	//TODO: add install check here before calling anything
	switch installType {
	case PlayStore:
		return d.installFromPlayStore(pkg, source)
	case Sideload:
		return d.installFromAPK(pkgName, pkg, source)
	case FDroidCL:
		if !FDroidCLIndexSynced {
			_, err := d.execCommand("fdroidcl", "update")

			if err != nil {
				return fmt.Errorf("failed to update the fdroid package index: %v", err)
			}

			FDroidCLIndexSynced = true
		}

		return d.installFromFDroidCL(pkg)
	default:
		return fmt.Errorf("unsupported install type: %s", installType)
	}
}

func (d *Device) isPackageInstalled(pkg string) bool {
	out, err := d.execCommand("adb", "shell", "pm", "list", "packages", "-3")

	if err != nil {
		return false
	}

	for _, item := range strings.Split(out, "\n") {
		if strings.TrimSpace(item) == "package:" + pkg {
			return true
		}
	}

	return false
}