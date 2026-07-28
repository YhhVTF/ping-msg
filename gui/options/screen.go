package sopt

import (
    "fyne.io/fyne/v2"

    "github.com/YhhVTF/ping-msg/opt"
)

type ContainerTableOptions struct {}

type ScreenOptions struct {
    Containers  ContainerTableOptions
    Widgets     WidgetTableOptions
    Window      fyne.Window
}

type WidgetTableOptions struct {}

func InitScreenOptions(w fyne.Window, opt *options.Options) *ScreenOptions {
    return nil
}
