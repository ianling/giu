package giu

import (
	"fmt"

	"github.com/ianling/imgui-go"
)

func SingleWindow(title string) *WindowWidget {
	size := Context.platform.DisplaySize()
	return Window(title).
		Flags(
			WindowFlagsNoTitleBar|
				WindowFlagsNoCollapse|
				WindowFlagsNoScrollbar|
				WindowFlagsNoMove|
				WindowFlagsNoResize).
		Size(size[0], size[1])
}

func SingleWindowWithMenuBar(title string) *WindowWidget {
	size := Context.platform.DisplaySize()
	return Window(title).
		Flags(
			WindowFlagsNoTitleBar|
				WindowFlagsNoCollapse|
				WindowFlagsNoScrollbar|
				WindowFlagsNoMove|
				WindowFlagsMenuBar|
				WindowFlagsNoResize).Size(size[0], size[1])
}

type windowState struct {
	hasFocus bool
}

func (s *windowState) Dispose() {
	// noop
}

type WindowWidget struct {
	title         string
	open          *bool
	flags         WindowFlags
	x, y          float32
	width, height float32
	bringToFront  bool
	pos           imgui.Vec2
	size          imgui.Vec2
}

func Window(title string) *WindowWidget {
	open := true

	return &WindowWidget{
		title: title,
		open:  &open,
	}
}

func (w *WindowWidget) IsOpen(open *bool) *WindowWidget {
	w.open = open
	return w
}

func (w *WindowWidget) Flags(flags WindowFlags) *WindowWidget {
	w.flags = flags
	return w
}

func (w *WindowWidget) Size(width, height float32) *WindowWidget {
	w.width, w.height = width, height
	return w
}

func (w *WindowWidget) Pos(x, y float32) *WindowWidget {
	w.x, w.y = x, y
	return w
}

func (w *WindowWidget) Layout(widgets ...Widget) {
	if widgets == nil {
		return
	}

	ws := w.getState()

	if w.flags&WindowFlagsNoMove != 0 && w.flags&WindowFlagsNoResize != 0 {
		imgui.SetNextWindowPos(imgui.Vec2{X: w.x, Y: w.y})
		imgui.SetNextWindowSize(imgui.Vec2{X: w.width, Y: w.height})
	} else {
		imgui.SetNextWindowPosV(imgui.Vec2{X: w.x, Y: w.y}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})
		imgui.SetNextWindowSizeV(imgui.Vec2{X: w.width, Y: w.height}, imgui.ConditionFirstUseEver)
	}

	widgets = append(widgets, Custom(func() {
		hasFocus := IsWindowFocused()
		if !hasFocus && ws.hasFocus {
			unregisterWindowShortcuts()
		}

		ws.hasFocus = hasFocus
		w.pos = imgui.WindowPos()
		w.size = imgui.WindowSize()
	}))

	if w.bringToFront {
		w.bringToFront = false
		imgui.SetNextWindowFocus()
	}

	showed := imgui.BeginV(w.title, w.open, imgui.WindowFlags(w.flags))

	if showed {
		Layout(widgets).Build()
	}

	imgui.End()
}

func (w *WindowWidget) CurrentPosition() (x, y float32) {
	return w.pos.X, w.pos.Y
}

func (w *WindowWidget) CurrentSize() (width, height float32) {
	return w.size.X, w.size.Y
}

func (w *WindowWidget) BringToFront() {
	w.bringToFront = true
}

func (w *WindowWidget) HasFocus() bool {
	return w.getState().hasFocus
}

func (w *WindowWidget) RegisterKeyboardShortcuts(s ...WindowShortcut) *WindowWidget {
	if w.HasFocus() {
		for _, shortcut := range s {
			RegisterKeyboardShortcuts(Shortcut{
				Key:      shortcut.Key,
				Modifier: shortcut.Modifier,
				Callback: shortcut.Callback,
				IsGlobal: LocalShortcut,
			})
		}
	}

	return w
}

func (w *WindowWidget) getStateID() string {
	return fmt.Sprintf("%s_windowState", w.title)
}

// returns window state
func (w *WindowWidget) getState() (state *windowState) {
	s := Context.GetState(w.getStateID())

	if s != nil {
		state = s.(*windowState)
	} else {
		state = &windowState{}

		Context.SetState(w.getStateID(), state)
	}

	return state
}
