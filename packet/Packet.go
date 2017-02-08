package packet

import (
	"archive/tar"
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/sha3"

	"github.com/trusch/tatar"
)

// Packet contains all data of a software packet
type Packet struct {
	ControlInfo
	Data *tatar.Tar
}

// A List is sortable list of Packets
type List []*Packet

func (a List) Len() int      { return len(a) }
func (a List) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a List) Less(i, j int) bool {
	if a[i].Name < a[j].Name {
		return true
	} else if a[i].Name == a[j].Name {
		return len(a[i].Labels) < len(a[j].Labels)
	}
	return false
}

// ToData returns a tar archive containing control.tar.xz and data.tar.xz which contains the actual payload of this packet
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

// FromData parses an packet from data which needs to be a tar archive containing control.tar.xz and data.tar.xz
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

// Hash returnes the base64 encoded sha3 shake 256bit hash of the packet
func (packet *Packet) Hash() (string, error) {
	if packet.ControlInfo.Hash != "" {
		return packet.ControlInfo.Hash, nil
	}
	hash := make([]byte, 16)
	data, err := packet.ToData()
	if err != nil {
		return "", err
	}
	sha3.ShakeSum256(hash, data)
	str := hex.EncodeToString(hash)
	packet.ControlInfo.Hash = str
	return str, nil
}

// NewFromData returns a new packet parsed from data
func NewFromData(data []byte) (*Packet, error) {
	pack := &Packet{}
	return pack, pack.FromData(data)
}

// NewFromFile returns a new packet parsed from file
func NewFromFile(file string) (*Packet, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	pack := &Packet{}
	err = pack.FromData(bs)
	if err != nil {
		return nil, err
	}
	_, err = pack.Hash()
	if err != nil {
		return nil, err
	}
	return pack, nil
}

// NewFromDirectory parses a directory producing a packet
func NewFromDirectory(dir string) (*Packet, error) {
	data, err := tatar.NewFromDirectory(filepath.Join(dir, "data"))
	if err != nil {
		return nil, err
	}
	data.Compression = tatar.LZMA

	preInst, err := ioutil.ReadFile(filepath.Join(dir, "preinst"))
	if err != nil {
		return nil, err
	}

	postInst, err := ioutil.ReadFile(filepath.Join(dir, "postinst"))
	if err != nil {
		return nil, err
	}

	preRm, err := ioutil.ReadFile(filepath.Join(dir, "prerm"))
	if err != nil {
		return nil, err
	}

	postRm, err := ioutil.ReadFile(filepath.Join(dir, "postrm"))
	if err != nil {
		return nil, err
	}

	control, err := ioutil.ReadFile(filepath.Join(dir, "control"))
	if err != nil {
		return nil, err
	}

	info := &ControlInfo{}
	err = info.FromYaml(control)
	if err != nil {
		return nil, err
	}

	pack := &Packet{
		ControlInfo: ControlInfo{
			Name:   info.Name,
			Labels: info.Labels,
			Scripts: Scripts{
				PreInst:  string(preInst),
				PostInst: string(postInst),
				PreRm:    string(preRm),
				PostRm:   string(postRm),
			},
		},
		Data: data,
	}
	return pack, nil
}

// InitDirectory initializes a directory with the minimal set of files to render a valid jamesd packet
func InitDirectory(dir, name string, labels map[string]string) error {
	if err := checkDir(dir); err != nil {
		return err
	}
	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(dir, "data"), 0755)
	ioutil.WriteFile(filepath.Join(dir, "preinst"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(dir, "postinst"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(dir, "prerm"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(dir, "postrm"), []byte{}, 0755)
	ctrl := ControlInfo{
		Name:   name,
		Labels: labels,
	}
	return ioutil.WriteFile(filepath.Join(dir, "control"), ctrl.ToYaml(), 0755)
}

func checkDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil // no stat -> no dir -> safe to write
	}
	if !stat.IsDir() {
		return errors.New("target directory already exists (and is a file oO)")
	}
	f, err := os.Open(dir)
	if err != nil {
		return errors.New("target directory is not writeable")
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return nil // dir exists, but no content
	}
	return errors.New("target directory already exists and is not empty")
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
