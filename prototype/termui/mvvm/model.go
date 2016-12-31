//
// model.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

package mvvm

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Model is an object that interacts with the ModelView.
type Model interface {
	Run(context.Context, *ModelView)
}

// ModelFunc wraps a function as a Model.
type ModelFunc func(context.Context, *ModelView)
func (f ModelFunc) Run(c context.Context, mv *ModelView) {
	f(c, mv)
}

// DemuxModel reads the input. It strings that start with a bang to Notices, and sends other strings to Messages.
type DemuxModel struct{}
var _ Model = &DemuxModel{}

func (d *DemuxModel) Run(ctx context.Context, mv *ModelView) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-mv.UserInput():
			// TODO: I'm being unfriendly; RTL should absolutely be supported by this app.
			tr := strings.TrimRightFunc(input, unicode.IsSpace)
			if len(tr) == 0 {
				continue
			}

			if tr[0] == '!' {
				mv.Notice(tr[1:])
			} else {
				mv.Message(fmt.Sprintf("\"%s\"", tr))
			}
		}
	}
}

// PrintingModel just writes output to Message.
func PrintingModel(ctx context.Context, mv *ModelView) {
	ticker := time.NewTicker(time.Second)

	var noticeDone <-chan struct{} = nil
	isDone := true
	for i := 0; true; i++{
		select {
		case <-ctx.Done():
			return
		case <-noticeDone:
			isDone = true
		case <-ticker.C:
			mv.Message(fmt.Sprintf("%d; notice up? %t", i, isDone))
			if i % 10 == 0 {
				noticeDone = mv.Notice(fmt.Sprintf("Happy %dth anniversary!", i))
				isDone = false
			}
		}
	}
}
