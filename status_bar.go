package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

type StatusBar struct {
	sb    cocoa.NSStatusBar
	title string
	si    cocoa.NSStatusItem
	menu  *menu
}

func NewStatusBar(title string) *StatusBar {
	return &StatusBar{
		sb:    cocoa.NSStatusBar_System(),
		title: title,
	}
}

func (s *StatusBar) Init() {
	s.si = s.sb.StatusItemWithLength_(cocoa.NSVariableStatusItemLength)
	s.si.Retain()
	s.si.Button().SetTitle(s.title)

	s.menu = initMenu
	s.menu.init(false)
	s.si.SetMenu(s.menu.m)
}

func (s *StatusBar) onCopy(text string) {
	var st = text
	// if text length greater than 7, use three dots
	if len(text) > 7 {
		st = text[:7] + "..."
	}
	// update status bar title
	s.si.Button().SetTitle(fmt.Sprintf("%s.Copied: %s", s.title, st))

	// update menu item
	s.menu.items[0].item.SetTitle(fmt.Sprintf("Last Copied: %s", text))
}

func (s *StatusBar) onConvert(text string) {
	var st = text
	// if text length greater than 7, use three dots
	if len(text) > 7 {
		st = text[:7] + "..."
	}
	// update status bar title
	s.si.Button().SetTitle(fmt.Sprintf("%s.Converted: %s", s.title, st))

	// update menu item
	s.menu.items[1].item.SetTitle(fmt.Sprintf("Last Converted: %s", text))
}

type menu struct {
	m     cocoa.NSMenu
	items []*menuItem
}

func (m *menu) init(sub bool) {
	m.m = cocoa.NSMenu_New()

	for _, item := range m.items {
		item.parent = m
		item.init()
		m.m.AddItem(item.item)
	}
}

func (m *menu) onSelectedChanged(item *menuItem) {
	for _, it := range m.items {
		if it.actionName == item.actionName {
			it.selectedChanged(true)
			continue
		}

		it.selectedChanged(false)
	}
}

type actionFunc func(parent *menu, item *menuItem)

var (
	doNoting = func(parent *menu, item *menuItem) {}
	tz       = time.Local
	tzAction = func(parent *menu, item *menuItem) {
		if item.title == "Local" {
			tz = time.Local
		}

		if item.title == "UTC" {
			tz = time.UTC
		}

		parent.onSelectedChanged(item)
	}
	// sec/ms
	secType   = "sec"
	secAction = func(parent *menu, item *menuItem) {
		if item.title == "sec" {
			secType = "sec"
		}

		if item.title == "ms" {
			secType = "ms"
		}

		parent.onSelectedChanged(item)
	}

	initMenu = &menu{
		items: []*menuItem{
			newMenuItem("Last Copied", nil, false),
			newMenuItem("Last Convert", nil, false),
			newMenuItem("Set Timezone", doNoting, false).setSubMenu(
				newMenuItem("Local", tzAction, true),
				newMenuItem("UTC", tzAction, false),
			),
			newMenuItem("Switch sec/ms", doNoting, false).setSubMenu(
				newMenuItem("sec", secAction, true),
				newMenuItem("ms", secAction, false),
			),
			newMenuItem("Quit", nil, false),
			newMenuItem("github.com/yusank", doNoting, false),
		},
	}
)

type menuItem struct {
	parent     *menu
	title      string
	action     actionFunc
	actionName string
	item       cocoa.NSMenuItem
	selected   bool

	subMenu *menu
}

func newMenuItem(title string, action actionFunc, isSelected bool) *menuItem {
	i := &menuItem{
		title:      title,
		action:     action,
		actionName: "action:" + strings.ReplaceAll(title, " ", "") + ":",
		selected:   isSelected,
	}

	i.item = cocoa.NSMenuItem_New()
	if isSelected {
		title = "✓ " + title
	}
	i.item.SetTitle(title)
	if i.title == "Quit" {
		i.action = nil
		i.actionName = "terminate:"
	}

	i.item.SetAction(objc.Sel(i.actionName))
	return i
}

func (i *menuItem) selectedChanged(ch bool) {
	if i.selected == ch {
		return
	}

	i.selected = ch
	title := i.title
	if i.selected {
		title = "✓ " + title
	}

	i.item.SetTitle(title)
}

func (i *menuItem) init() {
	if i.subMenu != nil {
		i.subMenu.init(true)
		i.item.SetSubmenu(i.subMenu.m)
	}

	if i.action != nil {
		cocoa.DefaultDelegateClass.AddMethod(i.actionName, func(notification objc.Object) {
			i.action(i.parent, i)
		})
	}
}

func (i *menuItem) setSubMenu(items ...*menuItem) *menuItem {
	i.subMenu = &menu{
		items: items,
	}

	return i
}
