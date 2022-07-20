package main

import (
	"context"
	"runtime"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

func main() {
	runtime.LockOSThread()

	cocoa.TerminateAfterWindowsClose = false
	var (
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()
	app := cocoa.NSApp_WithDidLaunch(func(_ objc.Object) {
		// watch pasted board
		default_pb.AddHandler(&timeHandler{})
		go default_pb.Run(ctx)

		// status bar
		sb := NewStatusBar("🩹Clipyy")
		sb.Init()

		default_pb.setOnCopy(sb.onCopy)
		default_pb.setOnHandle(sb.onConvert)
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}
