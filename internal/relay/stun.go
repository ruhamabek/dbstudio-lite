package relay

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	stunMagicCookie uint32 = 0x2112A442
	attrMappedAddr  uint16 = 0x0001
	attrXorMapped   uint16 = 0x0020
)

func DiscoverPublicIP(ctx context.Context, stunServer string) (string, error) {
	if stunServer == "" {
		stunServer = "stun.l.google.com:19302"
	}

	d := net.Dialer{Timeout: 3 * time.Second}
	conn, err := d.DialContext(ctx, "udp", stunServer)
	if err != nil {
		return "", fmt.Errorf("failed to reach STUN server %s: %w", stunServer, err)
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	}

	req := make([]byte, 20)
	binary.BigEndian.PutUint16(req[0:2], 0x0001)         // Message Type: Binding Request
	binary.BigEndian.PutUint16(req[2:4], 0x0000)         // Message Length: 0 (no body attributes)
	binary.BigEndian.PutUint32(req[4:8], stunMagicCookie) // Magic Cookie
	if _, err := rand.Read(req[8:20]); err != nil {      // 12-byte Transaction ID
		return "", err
	}

 	if _, err := conn.Write(req); err != nil {
		return "", fmt.Errorf("failed to send STUN request: %w", err)
	}

 	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read STUN response: %w", err)
	}

	if n < 20 {
		return "", errors.New("STUN response too short")
	}

 	msgType := binary.BigEndian.Uint16(buf[0:2])
	if msgType != 0x0101 {
		return "", fmt.Errorf("STUN server returned error/unexpected type: 0x%04x", msgType)
	}

	msgLen := int(binary.BigEndian.Uint16(buf[2:4]))
	if len(buf) < 20+msgLen {
		msgLen = n - 20
	}

 	offset := 20
	for offset+4 <= 20+msgLen {
		attrType := binary.BigEndian.Uint16(buf[offset : offset+2])
		attrLength := int(binary.BigEndian.Uint16(buf[offset+2 : offset+4]))
		offset += 4

		if offset+attrLength > n {
			break
		}

		val := buf[offset : offset+attrLength]

 		if attrType == attrXorMapped && len(val) >= 8 {
			family := val[1]
			if family == 0x01 {  
				ipBytes := make([]byte, 4)
				ipBytes[0] = val[4] ^ 0x21
				ipBytes[1] = val[5] ^ 0x12
				ipBytes[2] = val[6] ^ 0xA4
				ipBytes[3] = val[7] ^ 0x42
				return net.IP(ipBytes).String(), nil
			}
		}

 		if attrType == attrMappedAddr && len(val) >= 8 {
			family := val[1]
			if family == 0x01 {  
				return net.IPv4(val[4], val[5], val[6], val[7]).String(), nil
			}
		}

 		padding := (4 - (attrLength % 4)) % 4
		offset += attrLength + padding
	}

	return "", errors.New("no public IP address attribute found in STUN response")
}