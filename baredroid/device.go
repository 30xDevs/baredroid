package baredroid

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
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

func (d *Device) InstallPackage(pkg *PkgInstall) error {
	
	//TODO: add install check here before calling anything
	switch pkg.Type {
	case "playstore":
		if err := d.installFromPlayStore(pkg.Package, pkg.Source); err != nil {
			return err
		}
	case "sideload":
		if err := d.installFromAPK(pkg.Name, pkg.Package, pkg.Source); err != nil {
			return err
		}
	case "fdroidcl":
		if !FDroidCLIndexSynced {
			_, err := d.execCommand("fdroidcl", "update")

			if err != nil {
				return fmt.Errorf("failed to update the fdroid package index: %v", err)
			}

			FDroidCLIndexSynced = true
		}

		if err := d.installFromFDroidCL(pkg.Package); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported install type: %s", pkg.Type)
	}

	// Install children
	for _, child := range pkg.Children {
		if err := d.InstallPackage(&child); err != nil {
			return fmt.Errorf("could not install: %e", err)
		}
	}

	return nil
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