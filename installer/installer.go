package installer

import (
	"archive/tar"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/trusch/jamesd/packet"
)

// Install installs a packet to a given root directory
func Install(pack *packet.Packet, installRoot string) error {
	if err := os.MkdirAll(installRoot, 0755); err != nil {
		return err
	}
	if pack.ControlInfo.PreInst != "" {
		if err := execScript(pack.ControlInfo.PreInst); err != nil {
			return err
		}
	}
	if err := installTar(pack.Data.GetReader(), installRoot); err != nil {
		return err
	}
	if pack.ControlInfo.PostInst != "" {
		if err := execScript(pack.ControlInfo.PostInst); err != nil {
			return err
		}
	}
	return nil
}

// Uninstall uninstalls a packet from a given root directory
func Uninstall(pack *packet.Packet, installRoot string) error {
	if pack.ControlInfo.PreRm != "" {
		if err := execScript(pack.ControlInfo.PreRm); err != nil {
			return err
		}
	}
	if err := uninstallTar(pack.Data.GetReader(), installRoot); err != nil {
		return err
	}
	if pack.ControlInfo.PostRm != "" {
		if err := execScript(pack.ControlInfo.PostRm); err != nil {
			return err
		}
	}
	return nil
}

func execScript(script string) error {
	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installTar(archive *tar.Reader, installRoot string) error {
	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.FileInfo().IsDir() {
			err = os.MkdirAll(filepath.Join(installRoot, hdr.Name), os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
		} else {
			path := filepath.Join(installRoot, hdr.Name)
			f, e := os.Create(path)
			if e != nil {
				return e
			}
			if _, e = io.Copy(f, archive); e != nil {
				f.Close()
				return e
			}
			e = f.Chmod(os.FileMode(hdr.Mode))
			if e != nil {
				f.Close()
				return e
			}
			f.Close()
		}
	}
	return nil
}

func uninstallTar(archive *tar.Reader, installRoot string) error {
	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			break
		}
		if !hdr.FileInfo().IsDir() {
			path := filepath.Join(installRoot, hdr.Name)
			e := os.Remove(path)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
