package pushover

import (
	"os"
	"sync"
	"testing"
	"time"
)

const exampleToken = "KzGDORePK8gMaC0QOYAMyEEuzJnyUi"
const exampleUser = "uQiRzpo4DXghDmr9QzzfQu27cmVRsG"
const longString = "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFASDFA"

// IMPORTANT:
// Testing expects Pushover credentials to be found in the environemnt
// as POTTOKEN and POTUSER variables.

var setupOnce sync.Once
var testToken string
var testUser string

func setup() {
	testToken = os.Getenv("POTTOKEN")
	if testToken == "" {
		panic("POTUSER not set")
	}

	testUser = os.Getenv("POTUSER")
	if testUser == "" {
		panic("POTTOKEN not set")
	}
}

func TestNewPushOver(t *testing.T) {
	setupOnce.Do(setup)

	if _, err := NewPushover(exampleToken, exampleUser); err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	if _, err := NewPushover("invalid token", exampleUser); err != ErrInvalidToken {
		t.Errorf("expected 'ErrInvalidToken', got '%v'", err)
	}

	if _, err := NewPushover(exampleToken, "invalid user"); err != ErrInvalidUser {
		t.Errorf("expected 'ErrInvalidUser', got '%v'", err)
	}
}

func TestValidate(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	if err = po.Validate(); err != nil {
		t.Errorf("expected 'nil', got 'err'", err)
	}
}

func TestValidateInvalid(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(exampleToken, exampleUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	if err = po.Validate(); err == nil {
		t.Errorf("expected 'err', got 'nil'")
	}
}

func TestMessage(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	if err := po.Message(""); err != ErrBlankMessage {
		t.Errorf("expected 'ErrBlankMessage', got '%v'", err)
	}

	if err := po.Message(longString); err != ErrMessageTooLong {
		t.Errorf("expected, 'ErrMessageTooLong', got '%v'", err)
	}

	if err := po.Message("Test"); err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}
}

func TestMessagePush(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	m := new(Message)

	m.Message = "Message"
	m.Title = "Message Title"
	m.Device = "tdevice"
	m.Url = "https://twitter.com/superblock"
	m.UrlTitle = "Checkout Superblocks Pushover..."
	m.Priority = Low
	m.Timestamp = time.Now().Unix()
	m.Sound = "pushover"

	if _, _, err := po.Push(m); err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}
}

func TestMessageInvalidDevice(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	m := new(Message)

	m.Message = "Message"
	m.Title = "Message Title"
	m.Device = "!!!!"

	if _, _, err := po.Push(m); err != ErrInvalidDevice {
		t.Errorf("expected 'ErrInvalidDevice', got '%v'", err)
	}
}

func TestMessageUrlTooLong(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	m := new(Message)

	m.Message = "Message"
	m.Title = "Message Title"
	m.Url = longString

	if _, _, err := po.Push(m); err != ErrUrlTooLong {
		t.Errorf("expected 'ErrUrlTooLong', got '%v'", err)
	}
}

func TestMessageUrlTitleTooLong(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	m := new(Message)

	m.Message = "Message"
	m.Title = "Message Title"
	m.Url = "https://twitter.com/superblock"
	m.UrlTitle = longString

	if _, _, err := po.Push(m); err != ErrUrlTitleTooLong {
		t.Errorf("expected 'ErrUrlTitleTooLong', got '%v'", err)
	}
}

func TestMessageInvalidPriority(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	m := new(Message)

	m.Message = "Message"
	m.Title = "Message Title"
	m.Priority = 4711

	if _, _, err := po.Push(m); err != ErrInvalidPriority {
		t.Errorf("expected 'ErrInvalidPriority', got '%v'", err)
	}
}

func TestPushEmergencyReceipt(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	m := new(Message)

	m.Message = "Emergency Message"
	m.Priority = Emergency
	m.Retry = 60
	m.Expire = 3600

	_, receipt, err := po.Push(m)
	if err != nil {
		t.Errorf("expected 'nil' got '%v'", err)
	}

	if _, err := po.GetReceipt(receipt); err != nil {
		t.Errorf("expected 'nil' got '%v'", err)
	}
}

func TestGetReceiptInvalidReceipt(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	if _, err := po.GetReceipt("abc"); err != ErrInvalidReceipt {
		t.Errorf("expected 'ErrInvalidReceipt' got '%v'", err)
	}
}

func TestGetReceiptBogusReceipt(t *testing.T) {
	setupOnce.Do(setup)

	po, err := NewPushover(testToken, testUser)
	if err != nil {
		t.Errorf("expected 'nil', got '%v'", err)
	}

	_, err = po.GetReceipt("abasasdkljpoieksaqkwdlajsadaaa")
	if err == nil {
		t.Errorf("expected 'err' got 'nil'")
	}
}
