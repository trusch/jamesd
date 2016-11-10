package installer

import (
	"archive/tar"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(archive *tar.Reader, installRoot, preInstallScript, postInstallScript string) error {
	if preInstallScript != "" {
		if err := execScript(preInstallScript); err != nil {
			return err
		}
	}
	if err := installTar(archive, installRoot); err != nil {
		return err
	}
	if postInstallScript != "" {
		if err := execScript(postInstallScript); err != nil {
			return err
		}
	}

	return nil
}

func Uninstall(archive *tar.Reader, installRoot, preRemoveScript, postRemoveScript string) error {
	if preRemoveScript != "" {
		if err := execScript(preRemoveScript); err != nil {
			return err
		}
	}
	if err := uninstallTar(archive, installRoot); err != nil {
		return err
	}
	if postRemoveScript != "" {
		if err := execScript(postRemoveScript); err != nil {
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
