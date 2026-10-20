package relay

import (
 	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func BuildAddPortMappingSOAP(internalIP string, port int, proto, desc string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:AddPortMapping xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1">
      <NewRemoteHost></NewRemoteHost>
      <NewExternalPort>%d</NewExternalPort>
      <NewProtocol>%s</NewProtocol>
      <NewInternalPort>%d</NewInternalPort>
      <NewInternalClient>%s</NewInternalClient>
      <NewEnabled>1</NewEnabled>
      <NewPortMappingDescription>%s</NewPortMappingDescription>
      <NewLeaseDuration>3600</NewLeaseDuration>
    </u:AddPortMapping>
  </s:Body>
</s:Envelope>`, port, proto, port, internalIP, desc)
}

func BuildDeletePortMappingSOAP(port int, proto string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:DeletePortMapping xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1">
      <NewRemoteHost></NewRemoteHost>
      <NewExternalPort>%d</NewExternalPort>
      <NewProtocol>%s</NewProtocol>
    </u:DeletePortMapping>
  </s:Body>
</s:Envelope>`, port, proto)
}

type UPnPClient struct {
	controlURL string
}

func NewUPnPClient() *UPnPClient {
	return &UPnPClient{}
}


func (u *UPnPClient) ForwardPort(ctx context.Context, port int, internalIP string) error {
	if u.controlURL == "" {
		ctrl, err := u.discoverControlURL(ctx)
		if err != nil {
			return err
		}
		u.controlURL = ctrl
	}

	body := BuildAddPortMappingSOAP(internalIP, port, "TCP", "DBStudio Lite")
	return u.sendSOAP(ctx, "urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping", body)
}

func (u *UPnPClient) DeletePort(ctx context.Context, port int) error {
	if u.controlURL == "" {
		return nil
	}
	body := BuildDeletePortMappingSOAP(port, "TCP")
	return u.sendSOAP(ctx, "urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping", body)
}

func (u *UPnPClient) sendSOAP(ctx context.Context, action, body string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", u.controlURL, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Set("SOAPAction", fmt.Sprintf("\"%s\"", action))

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("router rejected UPnP request (HTTP %d): %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

func (u *UPnPClient) discoverControlURL(ctx context.Context) (string, error) {
	ssdpAddr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return "", err
	}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	msg := "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 1\r\n\r\n"

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
	}

	if _, err := conn.WriteTo([]byte(msg), ssdpAddr); err != nil {
		return "", err
	}

	buf := make([]byte, 2048)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		return "", errors.New("no UPnP router responded")
	}

	locRegex := regexp.MustCompile(`(?i)LOCATION:\s*([^\r\n]+)`)
	match := locRegex.FindSubmatch(buf[:n])
	if len(match) < 2 {
		return "", errors.New("invalid SSDP response: no location header")
	}

	locationURL := strings.TrimSpace(string(match[1]))
	parsedLoc, err := url.Parse(locationURL)
	if err != nil {
		return "", fmt.Errorf("invalid router location URL: %w", err)
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", locationURL, nil)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	ctrlRegex := regexp.MustCompile(`(?s)<serviceType>[^<]*WAN(?:IP|PPP)Connection:[^<]*</serviceType>.*?<controlURL>([^<]+)</controlURL>`)
	ctrlMatch := ctrlRegex.FindSubmatch(bodyBytes)
	if len(ctrlMatch) < 2 {
		return "", errors.New("no WANIPConnection control URL found in router description")
	}

	ctrlPath := strings.TrimSpace(string(ctrlMatch[1]))
	ctrlURL, err := parsedLoc.Parse(ctrlPath)
	if err != nil {
		return "", err
	}

	return ctrlURL.String(), nil
}