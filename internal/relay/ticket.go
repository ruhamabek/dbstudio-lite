package relay

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const TicketScheme = "dbshare://"


var (
	ErrInvalidTicketScheme = errors.New("invalid ticket format: must start with dbshare://")
	ErrEmptyPeerID         = errors.New("ticket is missing peer ID")
	ErrEmptyToken          = errors.New("ticket is missing security token")
	ErrInvalidPermission   = errors.New("ticket contains an invalid permission level")
)

type Ticket struct {
	PeerID     string     `json:"peerId"`
	Addrs      []string   `json:"addrs,omitempty"`
	Permission Permission `json:"permission"`
	Database   string     `json:"database,omitempty"`
	Token      string     `json:"token"`
}

func (t Ticket) Encode() (string, error) {
	if t.PeerID == "" {
		return "", ErrEmptyPeerID
	}
	if t.Token == "" {
		return "", ErrEmptyToken
	}
	data, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ticket: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	return TicketScheme + encoded, nil
}

func ParseTicket(raw string)(Ticket, error){
	if !strings.HasPrefix(raw, TicketScheme){
		return Ticket{}, ErrInvalidTicketScheme
	}

	payload := strings.TrimPrefix(raw, TicketScheme)
	if payload == "" {
		return Ticket{}, errors.New("ticket payload is empty")
	}

	data, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return Ticket{}, fmt.Errorf("invalid base64 payload: %w", err)
	}

	var t Ticket
	
	if err := json.Unmarshal(data, &t); err != nil {
		return Ticket{}, fmt.Errorf("failed to unmarshal ticket JSON: %w", err)
	}
	if t.PeerID == "" {
		return Ticket{}, ErrEmptyPeerID
	}
	if t.Token == "" {
		return Ticket{}, ErrEmptyToken
	}

	switch t.Permission {
		case PermissionView, PermissionRead, PermissionUpdate:
		default:
			return Ticket{}, fmt.Errorf("%w: %q", ErrInvalidPermission, t.Permission)
	}

	return t, nil
}