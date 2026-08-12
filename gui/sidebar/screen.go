package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    //"fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/gui/screen"
)

type ScreenSidebar struct {
    Base            *container.AppTabs
    ScreenManager   *screen.ScreenManager
    Window          fyne.Window
}

func InitScreenSidebar(w fyne.Window, s *screen.ScreenManager) *ScreenSidebar {
    g := &ScreenSidebar{}
    g.Window = w
    g.ScreenManager = s

    g.Base = container.NewAppTabs(
        container.NewTabItemWithIcon("", 
            theme.Icon(theme.IconNameMailCompose), &fyne.Container{},
        ),
        container.NewTabItemWithIcon("", 
            theme.Icon(theme.IconNameGrid), &fyne.Container{},
        ),
    )
    return g
}
