package test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
)

func TestController(t *testing.T) {
	buf := new(bytes.Buffer)
	ctrl := goflipdot.NewController(buf)

	t.Run("AddSign", func(t *testing.T) {
		err := ctrl.AddSign("test", 1, 86, 7, false)
		if err != nil {
			t.Errorf("Failed to add sign: %v", err)
		}

		err = ctrl.AddSign("test", 2, 86, 7, false)
		if err == nil {
			t.Error("Expected error when adding duplicate sign, got nil")
		}
	})

	t.Run("StartTestSigns", func(t *testing.T) {
		buf.Reset()
		err := ctrl.StartTestSigns()
		if err != nil {
			t.Errorf("Failed to start test signs: %v", err)
		}
		expected := []byte{0x02, '3', '0', 0x03, '9', 'A'}
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("Unexpected output for StartTestSigns. Got %v, want %v", buf.Bytes(), expected)
		}
	})

	t.Run("StopTestSigns", func(t *testing.T) {
		buf.Reset()
		err := ctrl.StopTestSigns()
		if err != nil {
			t.Errorf("Failed to stop test signs: %v", err)
		}
		expected := []byte{0x02, 'C', '0', 0x03, '8', 'A'}
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("Unexpected output for StopTestSigns. Got %v, want %v", buf.Bytes(), expected)
		}
	})

	t.Run("DrawImage", func(t *testing.T) {
		buf.Reset()
		img := image.NewGray(image.Rect(0, 0, 86, 7))
		err := ctrl.DrawImage(img, "test")
		if err != nil {
			t.Errorf("Failed to draw image: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("Expected output for DrawImage, got empty buffer")
		}
	})
}
