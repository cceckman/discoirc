// Package view provides UI views for discoirc.
package view

import (
	"log"

	"github.com/cceckman/discoirc/model"
	"github.com/jroimartin/gocui"
)

// Context provides data necessary for all Windows.
type Context struct {
	Gui     *gocui.Gui
	Log     *log.Logger
	Backend model.Client
}

// Window is a top-level view, e.g. Channel or Session.
type Window interface {
	// Start replaces the Gui with this Window, or returns an error.
	Start() error
}
