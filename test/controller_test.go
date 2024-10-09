package test

import (
	"bytes"
	"errors"
	"image"
	"testing"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
)

func TestNewController(t *testing.T) {
	buf := new(bytes.Buffer)
	ctrl, err := goflipdot.NewController(buf)
	if err != nil {
		t.Fatalf("Failed to create controller: %v", err)
	}
	if ctrl == nil {
		t.Fatal("Controller is nil")
	}
}

func TestAddSign(t *testing.T) {
	buf := new(bytes.Buffer)
	ctrl, _ := goflipdot.NewController(buf)

	t.Run("ValidSign", func(t *testing.T) {
		err := ctrl.AddSign("test1", 1, 86, 7)
		if err != nil {
			t.Errorf("Failed to add valid sign: %v", err)
		}
	})

	t.Run("DuplicateSign", func(t *testing.T) {
		err := ctrl.AddSign("test1", 2, 86, 7)
		if err == nil {
			t.Error("Expected error when adding duplicate sign, got nil")
		}
	})

	t.Run("InvalidDimensions", func(t *testing.T) {
		err := ctrl.AddSign("test2", 3, 0, 0)
		if err == nil {
			t.Error("Expected error when adding sign with invalid dimensions, got nil")
		}
	})
}

func TestStartTestSigns(t *testing.T) {
	buf := new(bytes.Buffer)
	ctrl, _ := goflipdot.NewController(buf)

	err := ctrl.StartTestSigns()
	if err != nil {
		t.Errorf("Failed to start test signs: %v", err)
	}
	expected := []byte{0x02, '3', '0', 0x03, '9', 'A'}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Unexpected output for StartTestSigns. Got %v, want %v", buf.Bytes(), expected)
	}
}

func TestStopTestSigns(t *testing.T) {
	buf := new(bytes.Buffer)
	ctrl, _ := goflipdot.NewController(buf)

	err := ctrl.StopTestSigns()
	if err != nil {
		t.Errorf("Failed to stop test signs: %v", err)
	}
	expected := []byte{0x02, 'C', '0', 0x03, '8', 'A'}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Unexpected output for StopTestSigns. Got %v, want %v", buf.Bytes(), expected)
	}
}

func TestDrawImage(t *testing.T) {
	buf := new(bytes.Buffer)
	ctrl, _ := goflipdot.NewController(buf)
	ctrl.AddSign("test", 1, 86, 7)

	t.Run("ValidImage", func(t *testing.T) {
		buf.Reset()
		img := image.NewGray(image.Rect(0, 0, 86, 7))
		err := ctrl.DrawImage(img, "test")
		if err != nil {
			t.Errorf("Failed to draw valid image: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("Expected output for DrawImage, got empty buffer")
		}
	})

	t.Run("InvalidImageSize", func(t *testing.T) {
		buf.Reset()
		img := image.NewGray(image.Rect(0, 0, 100, 100))
		err := ctrl.DrawImage(img, "test")
		if err == nil {
			t.Error("Expected error when drawing invalid image size, got nil")
		}
	})

	t.Run("NonexistentSign", func(t *testing.T) {
		buf.Reset()
		img := image.NewGray(image.Rect(0, 0, 86, 7))
		err := ctrl.DrawImage(img, "nonexistent")
		if err == nil {
			t.Error("Expected error when drawing to nonexistent sign, got nil")
		}
	})
}

func TestNetworkError(t *testing.T) {
	errorBuf := &ErrorBuffer{err: errors.New("network error")}
	ctrl, _ := goflipdot.NewController(errorBuf)

	err := ctrl.StartTestSigns()
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}

type ErrorBuffer struct {
	err error
}

func (e *ErrorBuffer) Write(p []byte) (n int, err error) {
	return 0, e.err
}

func (e *ErrorBuffer) Read(p []byte) (n int, err error) {
	return 0, e.err
}
