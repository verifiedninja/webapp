package pushover

import (
	"errors"

	"github.com/toorop/pushover"
)

var (
	po                  PushoverInfo
	ErrPushoverDisabled = errors.New("Pushover is disabled")
)

// PushoverInfo has the details for Pushover
type PushoverInfo struct {
	Enabled bool
	Token   string
	UserKey string
}

// Configure adds the settings
func Configure(p PushoverInfo) {
	po = p
}

// ReadConfig returns the settings
func ReadConfig() PushoverInfo {
	return po
}

// New returns a new instance of Pushover
func New() (*pushover.Pushover, error) {
	if po.Enabled {
		p, err := pushover.NewPushover(po.Token, po.UserKey)
		return p, err
	}

	return nil, ErrPushoverDisabled
}
