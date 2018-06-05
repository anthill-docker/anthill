package commands

import (
	"os"
	"testing"

	"github.com/aenthill/aenthill/app/context"

	"github.com/aenthill/manifest"
	"github.com/spf13/afero"
)

func TestNoImagesToRemoveError(t *testing.T) {
	err := &noImagesToRemoveError{}
	if err.Error() != noImagesToRemoveErrorMessage {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), noImagesToRemoveErrorMessage)
	}
}

func TestRemoveCmd(t *testing.T) {
	t.Run("calling RunE without arguments", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{}
		cmd := NewRemoveCmd(m, ctx)
		if err := cmd.RunE(nil, nil); err == nil {
			t.Error("RunE should have thrown an error as there are no arguments")
		}
	})

	t.Run("calling RunE with a non-existing manifest file", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		ctx := &context.AppContext{}
		cmd := NewRemoveCmd(m, ctx)
		if err := cmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
			t.Error("RunE should have thrown an error as the manifest file does not exist")
		}
	})

	t.Run("calling RunE with a broken manifest file", func(t *testing.T) {
		m := manifest.New("../../tests/aenthill-broken.json", afero.NewOsFs())
		ctx := &context.AppContext{}
		cmd := NewRemoveCmd(m, ctx)
		if err := cmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
			t.Error("RunE should have thrown an error as the manifest file is broken")
		}
	})

	t.Run("calling RunE with a non-existing image as argument", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while trying to flush the given manifest: %s", err.Error())
		}
		ctx := &context.AppContext{}
		cmd := NewRemoveCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err == nil {
			t.Errorf("RunE should have thrown an error as the image %s should not exist in given manifest", image)
		}
	})

	t.Run("calling RunE with all parameters OK!", func(t *testing.T) {
		image := "aenthill/cassandra"
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Errorf("an unexpected error occurred while trying to flush the given manifest: %s", err.Error())
		}
		if err := m.AddAent(image); err != nil {
			t.Errorf("an unexpected error occurred while adding an aent: %s", err.Error())
		}
		ctx := &context.AppContext{ProjectDir: os.Getenv("HOST_PROJECT_DIR"), LogLevel: "DEBUG"}
		cmd := NewRemoveCmd(m, ctx)
		if err := cmd.RunE(nil, []string{image}); err != nil {
			t.Error("RunE should not have thrown an error as all parameters are OK")
		}
	})
}
