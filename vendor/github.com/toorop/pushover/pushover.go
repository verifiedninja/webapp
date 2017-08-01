// Copyright (c) 2014 Thordur I. Bjornsson <thorduri@secnorth.net>
// All rights reserved. License can be found in the LICENSE file.

/*
	Package pushover implements a client interface to the pushover API.
	Peruse https://pushover.net/api for more information about the API.

	See comments in pushover_testing.go for important notices about testing.
*/
package pushover

import (
	"errors"
	"fmt"
	"regexp"

	"encoding/json"
	"net/http"
	"net/url"
)

const pushoverAPI = "https://api.pushover.net/1"
const pushoverAPIValidate = pushoverAPI + "/users/validate.json"
const pushoverAPIMessages = pushoverAPI + "/messages.json"
const pushoverAPIReceipt = pushoverAPI + "/receipts"

// MessageMaxLength is the maximum length of a message, including a title.
const MessageMaxLength = 512

// ErrMessageTooLong occurs when attempting to push a message that
// exceeds MessageMaxLength.
var ErrMessageTooLong = errors.New("pushover: message too long")

// UrlMaxLength is the maximum length of a supplementary URL.
const UrlMaxLength = 512

// ErrUrlTooLong occours when a supplementary URL exceeds UrlMaxLength
var ErrUrlTooLong = errors.New("pushover: supplementary URL too long")

// UrlTitleMaxLength is the maximum length of a supplementary URL title.
const UrlTitleMaxLength = 100

// ErrUrlTitleTooLong occours when a supplementary URL title exceeds UrlTitleMaxLength
var ErrUrlTitleTooLong = errors.New("pushover: supplementary URL title too long")

// ErrBlankMessage occours when attempting to push a blank message.
var ErrBlankMessage = errors.New("pushover: blank message")

// A Pushover encompasses interactions with the pushover API.
// Must be constructed with NewPushover.
type Pushover struct {
	user  string
	token string
}

var validToken = regexp.MustCompile(`^[A-Za-z0-9]{30}$`)

// ErrInvalidToken occurs when an invalid token is set.
var ErrInvalidToken = errors.New("pushover: invalid token")

var validUser = regexp.MustCompile(`^[A-Za-z0-9]{30}$`)

// ErrInvalidUser occurs when an invalid user is set.
var ErrInvalidUser = errors.New("pushover: invalid user")

var validDevice = regexp.MustCompile(`^[A-Za-z0-9_-]{1,25}$`)

// ErrInvalidDevice occurs when an invalid device is set.
var ErrInvalidDevice = errors.New("pushover: invalid device")

var validReceipt = regexp.MustCompile(`^[A-Za-z0-9]{30}$`)

// ErrInvalidReceipt occurs when an invalid receipt is set.
var ErrInvalidReceipt = errors.New("pushover: invalid receipt")

type errorResponse struct {
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors"`
}

type messagesResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
	Receipt string `json:"receipt"`
}

const (
	Low       = -1
	Normal    = 0
	High      = 1
	Emergency = 2
)

// ErrInvalidPriority occurs when an invalid message priority is set.
var ErrInvalidPriority = errors.New("pushover: invalid priority")

// A Message represents a message to be pushed over the pushover API.
type Message struct {
	// Required
	Message string

	// Optional
	Device    string
	Title     string
	Url       string
	UrlTitle  string
	Priority  int
	Retry     int
	Expire    int
	Timestamp int64
	Sound     string
}

// A Receipt represent a receipt from the pushover API for a pushed critical message.
type Receipt struct {
	Status          int   `json:"status"`
	Acknowledged    int   `json:"acknowledged"`
	AcknowledgedAt  int   `json:"acknowledged_at"`
	LastDeliveredAt int   `json:"last_delivered_at"`
	Expired         int   `json:"expired"`
	ExpiresAt       int64 `json:"expires_at"`
	CalledBack      int   `json:"called_back"`
	CalledBackAt    int64 `json:"called_back_at"`
}

// PushoverError occours when errors are made in the interaction with the pushover API.
type PushoverError struct {
	request string
	err     string
}

// Error returns the string representation of PushoverError.
func (e *PushoverError) Error() string {
	return fmt.Sprintf("pushover: request %s: %s", e.request, e.err)
}

// NewPushover returns a new Pushover structure suitable for interaction with the pushover API.
func NewPushover(token, user string) (*Pushover, error) {
	po := new(Pushover)

	if validToken.MatchString(token) == false {
		return nil, ErrInvalidToken
	}
	po.token = token

	if validUser.MatchString(user) == false {
		return nil, ErrInvalidUser
	}
	po.user = user

	return po, nil
}

func (po *Pushover) push(url string, message url.Values) (request, receipt string, err error) {
	resp, err := http.PostForm(url, message)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	jd := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse

		if err := jd.Decode(&errResp); err != nil {
			return "", "", err
		}

		return errResp.Request, "", &PushoverError{
			request: errResp.Request,
			err:     errResp.Errors[0],
		}
	}

	// Does the status in the message body really matter vis-a-vis
	// the HTTP status code ?

	if url == pushoverAPIMessages {
		var msgResp messagesResponse
		if err := jd.Decode(&msgResp); err != nil {
			return "", "", err
		}

		return msgResp.Request, msgResp.Receipt, nil
	}

	return "", "", nil
}

// Validate verifies that the credentials and devices are valid.
func (po *Pushover) Validate() error {
	message := url.Values{}

	message.Add("token", po.token)
	message.Add("user", po.user)

	_, _, err := po.push(pushoverAPIValidate, message)
	if err != nil {
		return err
	}

	return nil
}

// Push pushes a Message over the pushover API.
func (po *Pushover) Push(m *Message) (request, receipt string, err error) {
	message := url.Values{}

	message.Add("token", po.token)
	message.Add("user", po.user)

	if m.Message == "" {
		return "", "", ErrBlankMessage
	} else if len(m.Message)+len(m.Title) > MessageMaxLength {
		return "", "", ErrMessageTooLong
	}
	message.Add("message", m.Message)

	if m.Device != "" {
		if validDevice.MatchString(m.Device) == false {
			return "", "", ErrInvalidDevice
		}
		message.Add("device", m.Device)
	}

	if m.Title != "" {
		message.Add("title", m.Title)
	}

	if m.Url != "" {
		if len(m.Url) > UrlMaxLength {
			return "", "", ErrUrlTooLong
		}
		message.Add("url", m.Url)
	}

	if m.UrlTitle != "" {
		if len(m.UrlTitle) > UrlTitleMaxLength {
			return "", "", ErrUrlTitleTooLong
		}
		message.Add("url_title", m.UrlTitle)
	}

	switch m.Priority {
	case Low:
		fallthrough
	case Normal:
		fallthrough
	case High:
		fallthrough
	case Emergency:
		priostr := fmt.Sprintf("%d", m.Priority)
		message.Add("priority", priostr)
		if m.Priority == Emergency {
			rstr := fmt.Sprintf("%d", m.Retry)
			message.Add("retry", rstr)

			estr := fmt.Sprintf("%d", m.Expire)
			message.Add("expire", estr)
		}
	default:
		return "", "", ErrInvalidPriority
	}

	if m.Timestamp != 0 {
		tsstr := fmt.Sprintf("%d", m.Timestamp)
		message.Add("timestamp", tsstr)
	}

	if m.Sound != "" {
		message.Add("sound", m.Sound)
	}

	request, receipt, err = po.push(pushoverAPIMessages, message)
	if err != nil {
		return "", "", err
	}

	return request, receipt, nil
}

// Message pushes a simple message with default parameters over the pushover API.
// For more granular control over pushed messages see Message and Push.
func (po *Pushover) Message(msg string) error {
	m := new(Message)

	m.Message = msg
	if _, _, err := po.Push(m); err != nil {
		return err
	}

	return nil
}

// GetReceipt retrives the receipt data from a receipt of a pushed message
// of Emergency priority.
func (po *Pushover) GetReceipt(receipt string) (*Receipt, error) {
	if validReceipt.MatchString(receipt) == false {
		return nil, ErrInvalidReceipt
	}

	receiptURL := fmt.Sprintf("%s/%s.json?token=%s", pushoverAPIReceipt, receipt, po.token)
	resp, err := http.Get(receiptURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jd := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse

		if err := jd.Decode(&errResp); err != nil {
			return nil, err
		}

		return nil, &PushoverError{
			request: errResp.Request,
			err:     errResp.Errors[0],
		}
	}

	// Does the status in the message body really matter vis-a-vis
	// the HTTP status code ?

	r := new(Receipt)
	if err := jd.Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
