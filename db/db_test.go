package db

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/spec"
)

func TestPacket(t *testing.T) {
	db, err := New("mongodb://localhost/test-db")
	defer func() {
		assert.NoError(t, db.Drop())
	}()
	assert.NoError(t, err)

	var originalPacket *packet.Packet
	for i := 0; i < 20; i++ {
		labels := map[string]string{
			"n": strconv.Itoa(i),
		}
		if i%2 == 0 {
			labels["even"] = "true"
			labels["odd"] = "false"
		} else {
			labels["even"] = "false"
			labels["odd"] = "true"
		}
		packet.InitDirectory("./test", "test-packet", labels)
		pack, _ := packet.NewFromDirectory("./test")
		os.RemoveAll("./test")
		err = db.SavePacket(pack)
		assert.NoError(t, err)
		if i == 3 {
			originalPacket = pack
		}
	}
	info, err := db.GetBestInfo("test-packet", map[string]string{
		"n":      "3",
		"even":   "false",
		"odd":    "true",
		"doesnt": "exist",
	})
	assert.NoError(t, err)
	restoredPacket, err := db.GetPacket(info.Hash)
	assert.NoError(t, err)
	originalPacket.Hash()
	restoredPacket.Hash()
	assert.Equal(t, originalPacket, restoredPacket)
	err = db.DeletePacket(restoredPacket.ControlInfo.Hash)
	assert.NoError(t, err)
	infos, err := db.GetInfos("test-packet")
	assert.NoError(t, err)
	assert.Equal(t, 19, len(infos))

	names, err := db.GetPacketNames()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, "test-packet", names[0])
}

func TestSpec(t *testing.T) {
	db, err := New("mongodb://localhost/test-db")
	defer func() {
		assert.NoError(t, db.Drop())
	}()
	assert.NoError(t, err)

	err = db.SaveSpec(&spec.Spec{ID: "foo", Target: map[string]string{"a": "a"}, Apps: []*spec.App{&spec.App{Name: "foo"}}})
	assert.NoError(t, err)
	err = db.SaveSpec(&spec.Spec{ID: "bar", Target: map[string]string{"b": "b"}, Apps: []*spec.App{&spec.App{Name: "bar"}}})
	assert.NoError(t, err)
	err = db.SaveSpec(&spec.Spec{ID: "baz", Target: map[string]string{"c": "c"}, Apps: []*spec.App{&spec.App{Name: "baz"}}})
	assert.NoError(t, err)

	s, err := db.GetSpec("foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", s.ID)
	assert.Equal(t, map[string]string{"a": "a"}, s.Target)

	specs, err := db.GetSpecs()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(specs))

	s, err = db.GetMergedSpec(map[string]string{"a": "a", "b": "b"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(s.Apps))
	assert.Equal(t, "foo", s.Apps[0].Name)
	assert.Equal(t, "bar", s.Apps[1].Name)

	err = db.DeleteSpec("foo")
	assert.NoError(t, err)

	specs, err = db.GetSpecs()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(specs))

	_, err = db.GetSpec("foo")
	assert.Error(t, err)

	s, err = db.GetMergedSpec(map[string]string{"a": "a", "b": "b"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(s.Apps))
	assert.Equal(t, "bar", s.Apps[0].Name)

}
