package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui/screen"
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
    ChatCards []*ChatCard
    ChatsList *widget.List
}

func InitScreenSidebar(w fyne.Window, s *screen.ScreenManager, c *chat.ChatCache) *ScreenSidebar {
    g := &ScreenSidebar{}
    g.Window = w
    g.ScreenManager = s

    g.Widgets.ChatCards = make([]*ChatCard, 0)
    g.Widgets.ChatCards = append(g.Widgets.ChatCards, createChatCard(c.Chats[1]))

    g.Widgets.ChatsList = widget.NewList(
        func() int {
            return len(g.Widgets.ChatCards)
        },
        func() fyne.CanvasObject {
            return container.NewVBox()
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            o.(*fyne.Container).Add(g.Widgets.ChatCards[i].Base)
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
