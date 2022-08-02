package geoip

import (
	"net"
	"strconv"
	"strings"
)

func IP2Uint(ip string) uint32 {
	if ip == "" {
		return 0
	}
	bytes := strings.Split(ip, ".")
	b0, _ := strconv.Atoi(bytes[0])
	b1, _ := strconv.Atoi(bytes[1])
	b2, _ := strconv.Atoi(bytes[2])
	b3, _ := strconv.Atoi(bytes[3])

	var sum uint32
	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)
	return sum
}

func Uint2IP(ip uint32) string {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)
	ipv4 := net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
	return  ipv4.String()
}
