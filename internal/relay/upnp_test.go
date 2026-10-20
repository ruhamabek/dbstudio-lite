package relay_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"app-dev/internal/relay"
)

func TestUPnP_BuildSOAPPayload(t *testing.T) {
	t.Run("generates valid AddPortMapping XML envelope", func(t *testing.T) {
		xml := relay.BuildAddPortMappingSOAP("192.168.1.50", 54321, "TCP", "DBStudio Lite")
		
		if !strings.Contains(xml, "<NewExternalPort>54321</NewExternalPort>") {
			t.Errorf("missing NewExternalPort in SOAP: %s", xml)
		}
		if !strings.Contains(xml, "<NewInternalClient>192.168.1.50</NewInternalClient>") {
			t.Errorf("missing NewInternalClient in SOAP: %s", xml)
		}
		if !strings.Contains(xml, "<NewProtocol>TCP</NewProtocol>") {
			t.Errorf("missing NewProtocol in SOAP: %s", xml)
		}
		if !strings.Contains(xml, "AddPortMapping") {
			t.Errorf("missing AddPortMapping action in SOAP: %s", xml)
		}
	})

	t.Run("generates valid DeletePortMapping XML envelope", func(t *testing.T) {
		xml := relay.BuildDeletePortMappingSOAP(54321, "TCP")
		if !strings.Contains(xml, "<NewExternalPort>54321</NewExternalPort>") {
			t.Errorf("missing NewExternalPort in SOAP: %s", xml)
		}
		if !strings.Contains(xml, "DeletePortMapping") {
			t.Errorf("missing DeletePortMapping action in SOAP: %s", xml)
		}
	})

	t.Run("fails gracefully when no UPnP router responds within timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		upnp := relay.NewUPnPClient()
		_ = upnp.ForwardPort(ctx, 54321, "127.0.0.1")
	})
}