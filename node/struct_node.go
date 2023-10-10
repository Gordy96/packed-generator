package node

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type StructFieldNode struct {
	Name      string
	Class     string
	Size      int
	Offset    int
	ByteOrder string
}

type StructNode struct {
	Name      string
	ByteOrder string
	Members   []StructFieldNode
}

func (s StructNode) GenerateDefinition(w io.Writer) error {
	_, err := fmt.Fprintf(w, "type %s struct {\n", s.Name)
	if err != nil {
		return err
	}
	for _, m := range s.Members {
		_, err = fmt.Fprintf(w, "\t%s %s\n", m.Name, m.Class)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprint(w, "}")
	return err
}

func getBO(bo string) string {
	switch bo {
	case "b", "big", "big-endian":
		return "binary.BigEndian"
	case "l", "little", "little-endian":
		return "binary.LittleEndian"
	}
	return "binary.LittleEndian"
}

func (s StructNode) GenerateMarshalBinary(w io.Writer) error {
	callerName := strings.ToLower(s.Name[:1])
	_, err := fmt.Fprintf(w, "func (%s %s) MarshalBinary() ([]byte, error) {\n", callerName, s.Name)
	if err != nil {
		return err
	}
	var offset int
	var structByteOrder = getBO(s.ByteOrder)
	var buf = bytes.NewBuffer(nil)
	for _, m := range s.Members {
		var byteOrder string
		if len(m.ByteOrder) > 0 {
			byteOrder = getBO(m.ByteOrder)
		} else {
			byteOrder = structByteOrder
		}
		offset += m.Offset
		if inSlice(m.Class, []string{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64"}) {
			bits := parseInt(m.Class)
			op := "Int"
			if m.Class[0] == 'u' {
				op = "Uint"
			}
			_, err = fmt.Fprintf(buf, "\t%s.Put%s%d(buf[%d:%d], %s.%s)\n", byteOrder, op, bits, offset, offset+m.Size, callerName, m.Name)
		} else if m.Class == "string" {
			_, err = fmt.Fprintf(buf, "\tcopy(buf[%d:%d], %s.%s)\n", offset, offset+m.Size, callerName, m.Name)
		}
		if err != nil {
			return err
		}
		offset += m.Size
	}

	_, err = fmt.Fprintf(w, "\tvar buf [%d]byte\n", offset)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "\treturn nil\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, "}")
	return err
}

func inSlice(e string, ops []string) bool {
	for _, o := range ops {
		if e == o {
			return true
		}
	}
	return false
}

func parseInt(s string) int {
	var res int
	for _, c := range s {
		if c > '9' || c < '0' {
			continue
		}
		res = res*10 + int(c-'0')
	}
	return res
}
