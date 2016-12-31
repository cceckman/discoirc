//
// model.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

package mvvm

import (
	"fmt"
	"context"
	"strings"
	"unicode"
)

// Model is an object that interacts with the ModelView.
type Model interface {
	Run(context.Context, *ModelView)
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
