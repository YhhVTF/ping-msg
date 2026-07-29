package sopt

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/opt"
)

type ContainerTableOptions struct {
    Base *fyne.Container
}

type ScreenOptions struct {
    Containers      ContainerTableOptions
    InnerWindows    *container.MultipleWindows
    Widgets         WidgetTableOptions
    Window          fyne.Window
}

type WidgetTableOptions struct {
    SelectLanguage *widget.SelectEntry
}

func InitScreenOptions(
    w fyne.Window, wInner *container.MultipleWindows, opt *options.Options,
) *ScreenOptions {
    g := &ScreenOptions{}
    g.Window = w
    g.InnerWindows = wInner

    g.Containers.Base = container.NewHBox()

    return g
}
