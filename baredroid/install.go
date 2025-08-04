package baredroid

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func (d *Device) installFromPlayStore( ctx context.Context, pkg string, url string) error {

    if !d.isPackageInstalled(ctx, pkg) {
         _, err := d.execCommand(
            ctx,
            "adb",
            "shell", 
            "am", 
            "start", 
            "-a", 
            "android.intent.action.VIEW", 
            "-d", 
            url,
            "-n",
            "com.aurora.store/.MainActivity",
        )
        if err != nil {
            return fmt.Errorf("failed to open Play Store: %w", err)
        }
    } else {
        return nil
    }

	// Wait and verify installation
	return d.waitForPackageInstall(ctx, pkg)
}

func (d *Device) installFromAPK(ctx context.Context, pkgName string, pkg string, source string) error {

    apkPath := "./"+pkgName+".apk"

    if !d.isPackageInstalled(ctx, pkg) {
        
        if _, err := os.Stat(apkPath); os.IsNotExist(err) {
            downloadFile(apkPath, source)
        }

        // Execute streamed install with adb
        _, err := d.execCommand(
            ctx,
            "adb",
            "install",
            apkPath,
        )

        if err != nil {
            return fmt.Errorf("error installing apk %s: %v", apkPath, err)
        }
    }

    return nil
}

func (d *Device) installFromFDroidCL(ctx context.Context, pkg string) error {

    // TODO: Perhaps move this to device.go and check before calling
    if !d.isPackageInstalled(ctx, pkg) {
        _, err := d.execCommand(
            ctx,
            "fdroidcl",
            "install",
            pkg,
        )

        if err != nil {
            return fmt.Errorf("error installing %s with fdroidcl: %v", pkg, err)
        }
    }

    return nil
}

func (d *Device) waitForPackageInstall(ctx context.Context, pkg string) error {
    deadline := time.Now().Add(2 * time.Minute)
    for time.Now().Before(deadline) {
        installed:= d.isPackageInstalled(ctx, pkg)
        if installed {
            return nil
        }
        time.Sleep(5 * time.Second)
    }
    return fmt.Errorf("timeout waiting for package installation")
}

func downloadFile(filepath, url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }

    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}