package internal

import (
	"encoding/hex"
	"fmt"
)

const (
	StartByte byte = 0x02
	EndByte   byte = 0x03
	BaudRate  int  = 4800
)

var CommandCodes = map[string]byte{
	"start_test_signs": '3',
	"stop_test_signs":  'C',
	"write_image":      '1',
}

func ToAsciiHex(value []byte) []byte {
	return []byte(hex.EncodeToString(value))
}

func ClosestLargerMultiple(value, base int) int {
	return ((value + base - 1) / base) * base
}

func CalculateChecksum(data []byte) byte {
	var sum byte
	for _, b := range data[1:] {
		sum += b
	}
	return (^sum + 1) & 0xFF
}

func FormatPacket(command byte, address byte, payload []byte) []byte {
	packet := []byte{StartByte, '2', '1', command, address}
	packet = append(packet, ToAsciiHex([]byte{byte(len(payload))})...)
	packet = append(packet, payload...)
	packet = append(packet, EndByte)
	checksum := CalculateChecksum(packet)
	packet = append(packet, []byte(fmt.Sprintf("%02X", checksum))...)
	return packet
}
