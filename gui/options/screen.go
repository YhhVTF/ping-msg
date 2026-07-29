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
    Widgets         WidgetTableOptions
    Window          fyne.Window
}

type WidgetTableOptions struct {
    SelectLanguage *widget.Select
}

func InitScreenOptions(
    w fyne.Window, opt *options.Options,
) *ScreenOptions {
    g := &ScreenOptions{}
    g.Window = w

    g.Widgets.SelectLanguage = createSelectLanguage(opt)

    g.Containers.Base = container.NewVBox(g.Widgets.SelectLanguage)

    return g
}
