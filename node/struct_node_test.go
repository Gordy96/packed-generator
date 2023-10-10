package node

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructNode_GenerateDefinition(t *testing.T) {
	var s = StructNode{
		Name:      "Sample",
		ByteOrder: "",
		Members: []StructFieldNode{
			{
				Name:  "FirstMember",
				Class: "uint32",
				Size:  32,
			}, {
				Name:  "SecondMember",
				Class: "string",
				Size:  42,
			},
		},
	}
	buf := bytes.NewBuffer(nil)
	err := s.GenerateDefinition(buf)
	assert.NoError(t, err)
	const exp = `type Sample struct {
	FirstMember uint32
	SecondMember string
}`
	assert.Equal(t, exp, buf.String())
}

func TestStructNode_GenerateMarshalBinary(t *testing.T) {
	var s = StructNode{
		Name:      "Sample",
		ByteOrder: "big-endian",
		Members: []StructFieldNode{
			{
				Name:  "FirstMember",
				Class: "uint32",
				Size:  4,
			}, {
				ByteOrder: "l",
				Name:      "SecondMember",
				Class:     "int64",
				Size:      8,
			}, {
				Name:  "ThirdMember",
				Class: "string",
				Size:  42,
			},
		},
	}
	buf := bytes.NewBuffer(nil)
	err := s.GenerateMarshalBinary(buf)
	assert.NoError(t, err)
	const exp = `func (s Sample) MarshalBinary() ([]byte, error) {
	var buf [54]byte
	binary.BigEndian.PutUint32(buf[0:4], s.FirstMember)
	binary.LittleEndian.PutInt64(buf[4:12], s.SecondMember)
	copy(buf[12:54], s.ThirdMember)
	return nil
}`
	assert.Equal(t, exp, buf.String())
}
