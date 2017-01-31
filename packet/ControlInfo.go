package packet

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/trusch/tatar"
	yaml "gopkg.in/yaml.v2"
)

// ControlInfo contains the metadata of a packet
type ControlInfo struct {
	Name    string
	Labels  map[string]string
	Hash    string `yaml:"-"`
	Scripts `yaml:"-"`
}

// Scripts is a wrapper for install/deinstall related scripts
type Scripts struct {
	PreInst  string
	PostInst string
	PreRm    string
	PostRm   string
}

// ToYaml dumps the controlinfo as yaml
func (info *ControlInfo) ToYaml() []byte {
	d, _ := yaml.Marshal(info)
	return d
}

// FromYaml parses the controlinfo from yaml data
func (info *ControlInfo) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, info)
}

// ToData returns a lzma compressed tar archive containing controldata and all install/deinstall scripts
func (info *ControlInfo) ToData() ([]byte, error) {
	controlTar := &bytes.Buffer{}
	controlTarWriter := tar.NewWriter(controlTar)
	controlData := info.ToYaml()
	err := addDataToTarWriter(controlTarWriter, controlData, "control")
	if err != nil {
		return nil, err
	}
	err = addDataToTarWriter(controlTarWriter, []byte(info.Scripts.PreInst), "preinst")
	if err != nil {
		return nil, err
	}
	err = addDataToTarWriter(controlTarWriter, []byte(info.Scripts.PostInst), "postinst")
	if err != nil {
		return nil, err
	}
	err = addDataToTarWriter(controlTarWriter, []byte(info.Scripts.PreRm), "prerm")
	if err != nil {
		return nil, err
	}
	err = addDataToTarWriter(controlTarWriter, []byte(info.Scripts.PostRm), "postrm")
	if err != nil {
		return nil, err
	}

	err = controlTarWriter.Close()
	if err != nil {
		return nil, err
	}
	controlTatar := &tatar.Tar{}
	_, err = controlTatar.Load(controlTar)
	if err != nil {
		return nil, err
	}
	controlTatar.Compression = tatar.LZMA
	return controlTatar.ToData()
}

// FromData parses an lzma compressed tar archive containing the controlinfo and the scripts
func (info *ControlInfo) FromData(data []byte) error {
	t, err := tatar.NewFromData(data, tatar.LZMA)
	if err != nil {
		return err
	}
	return t.ForEach(func(header *tar.Header, reader io.Reader) error {
		switch header.Name {
		case "control":
			{
				buf := &bytes.Buffer{}
				_, err := io.Copy(buf, reader)
				if err != nil {
					return err
				}
				err = info.FromYaml(buf.Bytes())
				if err != nil {
					return err
				}
			}
		case "preinst":
			{
				buf := &bytes.Buffer{}
				_, err := io.Copy(buf, reader)
				if err != nil {
					return err
				}
				info.Scripts.PreInst = string(buf.Bytes())
			}
		case "postinst":
			{
				buf := &bytes.Buffer{}
				_, err := io.Copy(buf, reader)
				if err != nil {
					return err
				}
				info.Scripts.PostInst = string(buf.Bytes())
			}
		case "prerm":
			{
				buf := &bytes.Buffer{}
				_, err := io.Copy(buf, reader)
				if err != nil {
					return err
				}
				info.Scripts.PreRm = string(buf.Bytes())
			}
		case "postrm":
			{
				buf := &bytes.Buffer{}
				_, err := io.Copy(buf, reader)
				if err != nil {
					return err
				}
				info.Scripts.PostRm = string(buf.Bytes())
			}
		}
		return nil
	})
}
