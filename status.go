package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusBar *walk.StatusBarItem
var statusConnectivity string

func StatusUpdate(connectivity string) {
	if statusBar != nil {
		statusBar.SetText(connectivity)
	}
	statusConnectivity = connectivity
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		{
			AssignTo: &statusBar,
			Text:     statusConnectivity,
			Icon:     ICON_Status,
			Width:    300,
		},
	}
}
