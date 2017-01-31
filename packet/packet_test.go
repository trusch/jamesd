package packet

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDirectory(t *testing.T) {
	err := InitDirectory("./test", "test-packet", map[string]string{"a": "label"})
	assert.NoError(t, err)
	defer os.RemoveAll("./test")
	bs, _ := ioutil.ReadFile("./test/control")
	expected := `name: test-packet
labels:
  a: label
`
	assert.Equal(t, expected, string(bs))
}

func TestNewFromDirectory(t *testing.T) {
	InitDirectory("./test", "test-packet", map[string]string{"a": "label"})
	defer os.RemoveAll("./test")
	ioutil.WriteFile("./test/data/foo", []byte("bar"), 0755)
	pack, err := NewFromDirectory("./test")
	assert.NoError(t, err)
	reader := pack.Data.GetReader()
	header, err := reader.Next()
	assert.NoError(t, err)
	assert.Equal(t, "foo", header.FileInfo().Name())
	buff := make([]byte, 32)
	n, err := reader.Read(buff)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "bar", string(buff[:n]))
	_, err = reader.Next()
	assert.Equal(t, io.EOF, err)
}

func TestToDataFromData(t *testing.T) {
	InitDirectory("./test", "test-packet", map[string]string{"a": "label"})
	defer os.RemoveAll("./test")
	ioutil.WriteFile("./test/data/foo", []byte("bar"), 0755)
	pack, _ := NewFromDirectory("./test")
	data, err := pack.ToData()
	assert.NoError(t, err)
	restoredPacket, err := NewFromData(data)
	assert.NoError(t, err)
	assert.Equal(t, pack.Name, restoredPacket.Name)
	assert.Equal(t, pack.Labels, restoredPacket.Labels)
	reader := restoredPacket.Data.GetReader()
	header, err := reader.Next()
	assert.NoError(t, err)
	assert.Equal(t, "foo", header.FileInfo().Name())
	buff := make([]byte, 32)
	n, err := reader.Read(buff)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "bar", string(buff[:n]))
	_, err = reader.Next()
	assert.Equal(t, io.EOF, err)
}

func TestHash(t *testing.T) {
	InitDirectory("./test", "test-packet", map[string]string{"a": "label"})
	defer os.RemoveAll("./test")
	ioutil.WriteFile("./test/data/foo", []byte("bar"), 0755)
	pack, _ := NewFromDirectory("./test")
	hash, err := pack.Hash()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	hash2, err := pack.Hash()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash2)
	assert.Equal(t, hash, hash2)
	data, _ := pack.ToData()
	restored, _ := NewFromData(data)
	hash3, err := restored.Hash()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash3)
	assert.Equal(t, hash, hash3)
	fmt.Println(hash)
}
