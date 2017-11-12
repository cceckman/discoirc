//
// lorem.go
// Copyright (C) 2017 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

package model

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
)

var lorem = strings.Split("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec bibendum consequat nibh vitae vestibulum. Praesent rutrum massa ac lorem bibendum, a facilisis est vehicula. Vivamus efficitur vehicula eros id mollis. In hac habitasse platea dictumst. Donec id est scelerisque, mollis nibh ut, tempus ligula. Donec interdum faucibus leo ac rutrum. Vestibulum gravida tempor dui, vitae vulputate arcu ullamcorper ac. Ut porttitor libero at ipsum mattis elementum. Nullam odio odio, lacinia non venenatis non, sagittis nec purus. Proin bibendum sollicitudin nibh at faucibus. Nunc ornare faucibus tortor. Praesent ut nibh vitae augue semper sodales. Praesent laoreet est eu posuere rhoncus. Suspendisse enim purus, tincidunt eu finibus sed, tristique sit amet diam. Integer nec euismod enim. Pellentesque bibendum arcu urna.", " ")

// EventGenerator generates a sequence of events for the Channel.
func EventGenerator(c *MockChannel) {
	go func() {
		glog.V(1).Info("[start] sending events")
		defer glog.V(1).Info("[done] sending events")

		connected := false

		initialNick := strings.Split("cceckman", "")
		sort.Strings(initialNick)
		nick := initialNick
		mode := strings.Split("one", "")

		for i := 0; true; i++ {
			time.Sleep(time.Millisecond * 1000)
			if i%11 == 0 && i/5 > 0 {
				n := (i / 5) % len(lorem)
				topic := strings.Join(lorem[0:n], " ")

				glog.V(1).Infof("generating state update: topic: %s", topic)
				c.updateState <- ChannelState{
					Topic: topic,
				}
				c.SendMessage("[system] topic updated")
			}
			if i%7 == 1 {
				newNick := strings.Join(nick, "")
				glog.V(1).Infof("generating state update: nick: %s", newNick)
				c.updateState <- ChannelState{
					Nick: newNick,
				}
				c.SendMessage("[system] nick updated")
				if !next(sort.StringSlice(nick)) {
					nick = initialNick
				}
			}
			if i%13 == 2 {
				connected = !connected
				glog.V(1).Infof("generating state update: connected: %v", connected)
				c.connected <- connected
				c.SendMessage("[system] connection state updated")
			}
			if i%17 == 3 {
				c.updateState <- ChannelState{
					Mode: strings.Join(append([]string{"+"}, mode[0:(i%3)]...), ""),
				}
				c.SendMessage("[system] mode updated")
			}

			msg := fmt.Sprintf("%d bottles of beer on the wall, %d bottles of beer...", i, i)
			glog.V(1).Infof("Chat/messages: [sending] : %v", msg)
			c.SendMessage(msg)
		}
	}()
}

// next returns false when it cannot permute any more
// http://en.wikipedia.org/wiki/Permutation#Generation_in_lexicographic_order
func next(data sort.Interface) bool {
	var k, l int
	for k = data.Len() - 2; ; k-- {
		if k < 0 {
			return false
		}
		if data.Less(k, k+1) {
			break
		}
	}
	for l = data.Len() - 1; !data.Less(k, l); l-- {
	}
	data.Swap(k, l)
	for i, j := k+1, data.Len()-1; i < j; i++ {
		data.Swap(i, j)
		j--
	}
	return true
}
