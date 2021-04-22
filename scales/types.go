package scales

import (
	"encoding/binary"
	"reflect"
)

var Commands = map[string]byte{
	"CMD_GET_NAME":     0x20,
	"CMD_SET_NAME":     0x22,
	"CMD_GET_ETHERNET": 0x2D,
	"CMD_SET_ETHERNET": 0x39,
	"CMD_GET_WIFI":     0x3A,
	"CMD_SET_WIFI":     0x3C,
	"CMD_GET_MASSA":    0x23,
	"CMD_SET_TARE":     0xA3,
}

type ResponseMeta struct {
	Valid       bool
	Len         uint16
	CommandCode byte
	CommandName string
	CRC         uint16
}

type RawResponse struct {
	Bytes  []byte
	Offset int
	Meta   ResponseMeta
}

func (r *RawResponse) Get(n int) []byte {
	result := r.Bytes[r.Offset : r.Offset+n]
	r.Offset += n
	return result
}

func (r *RawResponse) FillMeta() {
	header := []byte{0xF8, 0x55, 0xCE}
	responses := map[byte]string{
		0x21: "CMD_ACK_NAME",
		0x27: "CMD_ACK_SET",
		0x2E: "CMD_ACK_ETHERNET",
		0x34: "CMD_ACK_WIFI_IP",
		0x24: "CMD_ACK_MASSA",
		0x12: "CMD_ACK_SET_TARE",
		0x15: "CMD_NACK_TARE",

		0xF0: "CMD_NACK",
		0x28: "CMD_ERROR",
	}

	r.Meta.Valid = true

	if !reflect.DeepEqual(r.Bytes[:3], header) {
		r.Meta.Valid = false
		return
	}

	r.Meta.Len = binary.LittleEndian.Uint16(r.Bytes[3:5])

	r.Meta.CommandCode = r.Bytes[5]
	r.Meta.CommandName = responses[r.Meta.CommandCode]

	r.Meta.CRC = binary.LittleEndian.Uint16(r.Bytes[6+r.Meta.Len-1 : 6+r.Meta.Len+1])

	if r.Meta.CRC != CRC16(0, r.Bytes[5:6+r.Meta.Len-1], r.Meta.Len) {
		r.Meta.Valid = false
	}

	return
}

type Response struct {
	raw RawResponse
}

func (r *Response) Raw() RawResponse {
	return r.raw
}

type Connection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

func (r *Response) Read(s Connection) error {
	buf := make([]byte, 128)
	n, err := s.Read(buf)
	if err != nil {
		return err
	}

	r.raw.Bytes = make([]byte, n)
	copy(r.raw.Bytes, buf)

	r.raw.FillMeta()

	return nil
}
