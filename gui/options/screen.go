package sopt

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/opt"
)

type ContainerTableOptions struct {
    // Base container of the options sceen, makes the child container scrollable
    Base *container.Scroll
    // Contains widgets for changing options
    VBox *fyne.Container
}

type ScreenOptions struct {
    Containers      ContainerTableOptions
    Widgets         WidgetTableOptions
    Window          fyne.Window
}

type WidgetTableOptions struct {
    CardLanguage    *widget.Card
    SelectLanguage  *widget.Select
}

func InitScreenOptions(
    w fyne.Window, opt *options.Options,
) *ScreenOptions {
    g := &ScreenOptions{}
    g.Window = w

    g.Widgets.SelectLanguage = createSelectLanguage(opt)
    g.Widgets.CardLanguage = 
        widget.NewCard("", "Language", g.Widgets.SelectLanguage)

    g.Containers.VBox = container.NewVBox(g.Widgets.CardLanguage)
    g.Containers.Base = container.NewVScroll(g.Containers.VBox)

    return g
}
