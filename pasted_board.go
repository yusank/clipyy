package main

import (
	"context"
	"log"
	"time"

	"github.com/progrium/macdriver/cocoa"
)

var (
	default_pb = newPasteboard(cocoa.NSPasteboard_GeneralPasteboard())
)

type Pasteboard struct {
	pb       cocoa.NSPasteboard
	handlers []Handler
	onCopy   func(string)
	onHandle func(string)
}

func newPasteboard(pb cocoa.NSPasteboard) *Pasteboard {
	return &Pasteboard{
		pb:       pb,
		handlers: []Handler{},
	}
}

func (p *Pasteboard) AddHandler(h Handler) {
	p.handlers = append(p.handlers, h)
}

func (p *Pasteboard) setOnCopy(f func(string)) {
	p.onCopy = f
}

func (p *Pasteboard) setOnHandle(f func(string)) {
	p.onHandle = f
}

func (p *Pasteboard) Run(ctx context.Context) {
	var (
		ticker = time.NewTicker(time.Millisecond * 300)
		cc     = p.pb.ChangeCount()
	)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if p.pb.ChangeCount() == cc {
				continue
			}

			str := p.pb.StringForType(cocoa.NSPasteboardTypeString)
			log.Println("read:", str)
			if p.onCopy != nil {
				p.onCopy(str)
			}

			for _, h := range p.handlers {
				m := h.Match(context.Background(), str)
				if m.Match {
					rs := h.Convert(context.Background(), str, m)
					if rs != "" && rs != str {
						p.pb.ClearContents()
						p.pb.SetStringForType(rs, cocoa.NSPasteboardTypeString)
						log.Println("write:", rs)
						if p.onHandle != nil {
							p.onHandle(rs)
						}
					}
				}
			}
			cc = p.pb.ChangeCount()
		}
	}
}
