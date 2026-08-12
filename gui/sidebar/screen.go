package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/log"
)

type ScreenSidebar struct {
    // Tabs for DMs, Chats, and Friends in that order
    Base            *container.AppTabs
    ScreenManager   *screen.ScreenManager
    // All widgets used in the sidebar
    Widgets         WidgetTableSidebar
    Window          fyne.Window
}

type WidgetTableSidebar struct {
    ChatsList *widget.List
}

func InitScreenSidebar(w fyne.Window, s *screen.ScreenManager, c *chat.ChatCache) *ScreenSidebar {
    g := &ScreenSidebar{}
    g.Window = w
    g.ScreenManager = s

    log.Info.Printf("%d chat cards in sidebar\n", len(c.Chats))

    g.Widgets.ChatsList = widget.NewList(
        func() int {
            return len(c.Chats)
        },
        func() fyne.CanvasObject {
            return createChatCardTemplate()
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            updateChatCard(o.(*widget.Card), c.Chats[1])
        },
    )

    // Initialize app tabs widget
    g.Base = container.NewAppTabs(
        // DMs tab
        container.NewTabItemWithIcon("", 
            theme.Icon(theme.IconNameMailCompose), &fyne.Container{},
        ),
        // Chats tab
        container.NewTabItemWithIcon("", 
            theme.Icon(theme.IconNameGrid), container.NewStack(g.Widgets.ChatsList),
        ),
    )
    return g
}
