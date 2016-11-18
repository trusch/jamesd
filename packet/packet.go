package packet

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/trusch/tatar"
	"gopkg.in/yaml.v2"
)

type Packet struct {
	ControlInfo
	Data *tatar.Tar
}

type ControlInfo struct {
	Name    string
	Tags    []string
	Depends []string
	Scripts `yaml:"-"`
}

type Scripts struct {
	PreInst  string
	PostInst string
	PreRm    string
	PostRm   string
}

func (info *ControlInfo) ToYaml() []byte {
	d, _ := yaml.Marshal(info)
	return d
}

func (info *ControlInfo) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, info)
}

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

func (packet *Packet) ToData() ([]byte, error) {
	mainTar := &bytes.Buffer{}
	controlData, err := packet.ControlInfo.ToData()
	if err != nil {
		return nil, err
	}
	mainTarWriter := tar.NewWriter(mainTar)
	err = addDataToTarWriter(mainTarWriter, controlData, "control.tar.xz")
	if err != nil {
		return nil, err
	}
	dataBytes, err := packet.Data.ToData()
	if err != nil {
		return nil, err
	}
	err = addDataToTarWriter(mainTarWriter, dataBytes, "data.tar.xz")
	if err != nil {
		return nil, err
	}
	err = mainTarWriter.Close()
	if err != nil {
		return nil, err
	}
	return mainTar.Bytes(), nil
}

func (packet *Packet) FromData(data []byte) error {
	t, err := tatar.NewFromData(data, tatar.NO_COMPRESSION)
	if err != nil {
		return err
	}
	return t.ForEach(func(header *tar.Header, reader io.Reader) error {
		if header.Name == "control.tar.xz" {
			buf := &bytes.Buffer{}
			_, err = io.Copy(buf, reader)
			if err != nil {
				return err
			}
			return packet.ControlInfo.FromData(buf.Bytes())
		}
		if header.Name == "data.tar.xz" {
			d := &tatar.Tar{Compression: tatar.LZMA}
			_, err = d.Load(reader)
			if err != nil {
				return err
			}
			packet.Data = d
		}
		return nil
	})
}

func NewFromData(data []byte) (*Packet, error) {
	pack := &Packet{}
	return pack, pack.FromData(data)
}

func addDataToTarWriter(t *tar.Writer, data []byte, name string) error {
	err := t.WriteHeader(&tar.Header{
		Name: name,
		Mode: 0600,
		Size: int64(len(data)),
	})
	if err != nil {
		return err
	}
	_, err = t.Write(data)
	if err != nil {
		return err
	}
	return nil
}
