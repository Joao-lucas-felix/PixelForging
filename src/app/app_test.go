package app

import (
	"testing"
)

func TestGenApp(t *testing.T) {
	app := GenAPP()

	if app.Name != "Pixel Forging" ||
		app.Usage != "A CLI to processing images" {
		t.Errorf("app name o usage descroption incorrets!")
	}
}
