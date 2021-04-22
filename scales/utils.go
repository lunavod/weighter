package scales

import (
	"encoding/binary"
	"reflect"
)

func CRC16(crc uint16, buf []byte, len uint16) uint16 {
	var bits, k, a, temp uint16
	crc = 0
	for k=0; k < len; k++ {
		a = 0
		temp = (crc >> 8) << 8
		for bits = 0; bits < 8; bits++ {
			if (temp ^ a) & 0x8000 != 0 {
				a=(a<<1) ^ 0x1021
			} else {
				a<<=1
			}
			temp<<=1
		}
		crc = a ^ (crc<<8) ^ (binary.LittleEndian.Uint16([]byte{buf[k], 0x0}) & 0xFF)
	}
	return crc
}

func FillResponseStruct(raw RawResponse, st reflect.Value) {
	raw.Offset = 6
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if !field.CanInterface() || field.Kind() == reflect.Struct { continue }

		var tmp reflect.Value

		switch field.Kind() {
		case reflect.Uint32:
			tmp = reflect.ValueOf(binary.LittleEndian.Uint32(raw.Get(4)))
			break
		case reflect.Uint16:
			tmp = reflect.ValueOf(binary.LittleEndian.Uint16(raw.Get(2)))
			break
		case reflect.Uint8:
			t := raw.Get(1)[0]
			tmp = reflect.ValueOf(t)
			break
		case reflect.Bool:
			t := raw.Get(1)[0] == 1
			tmp = reflect.ValueOf(t)
			break
		case reflect.String:
			pair := raw.Get(2)
			var slice []byte
			for !reflect.DeepEqual(pair, []byte{0x0D, 0x0A}) {
				slice = append(slice, pair[0])
				pair[0] = pair[1]
				pair[1] = raw.Get(1)[0]
			}
			tmp = reflect.ValueOf(string(slice))
			break
		case reflect.Array:
			slice := raw.Get(field.Cap())
			tmp = reflect.New(field.Type()).Elem()
			for i := 0; i < field.Cap(); i++ {
				tmp.Index(i).Set(reflect.ValueOf(slice[i]))
			}
			break
		}

		field.Set(tmp)
	}
}

func BuildRequest(data []byte) []byte {
	var result []byte

	result = append(result, []byte{0xF8, 0x55, 0xCE}...)

	lenBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBytes, uint16(len(data)))
	result = append(result, lenBytes...)

	result = append(result, data...)

	crcBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcBytes, CRC16(0, data, uint16(len(data))))
	result = append(result, crcBytes...)

	return result
}
