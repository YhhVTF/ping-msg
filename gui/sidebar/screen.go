package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

type ScreenSidebar struct {
    // Tabs for DMs, Chats, and Friends in that order
    Base                        *container.AppTabs
    OutgoingReqChatMetadata     chan prot.ChatMetadataRequest
    OutgoingReqMember           chan prot.MemberRequest
    ScreenManager               *screen.ScreenManager
    // All widgets used in the sidebar
    Widgets                     WidgetTableSidebar
    Window                      fyne.Window
}

type WidgetTableSidebar struct {
    // Sends an outgoing Create chat metadata request in order to create a new chat
    ButtonCreateChat    *widget.Button
    // Prompts the user to enter the ID of a chat they want to join
    ButtonJoinChat      *widget.Button
    // List of chat cards which display info about a chat and open that chat when clicked
    ChatsList           *widget.List
}

func InitScreenSidebar(w fyne.Window, s *screen.ScreenManager, c *chat.ChatCache, u *user.UserCache, opt *options.Options) *ScreenSidebar {
    g := &ScreenSidebar{}
    g.Window = w
    g.ScreenManager = s

    g.OutgoingReqChatMetadata = make(chan prot.ChatMetadataRequest)
    g.OutgoingReqMember = make(chan prot.MemberRequest)

    g.Widgets.ButtonJoinChat = createButtonJoin(g, w, u, opt)

    g.Widgets.ButtonCreateChat = createButtonCreate(g, u, opt)

    // Initialize app tabs widget
    g.Base = container.NewAppTabs(
        // DMs tab
        container.NewTabItemWithIcon("", 
            theme.Icon(theme.IconNameMailCompose), &fyne.Container{},
        ),
        // Chats tab
        container.NewTabItemWithIcon("", 
            theme.Icon(theme.IconNameGrid), container.NewBorder(
                nil,
                container.NewStack(container.NewHBox(
                    container.NewStack(g.Widgets.ButtonJoinChat),
                    container.NewStack(g.Widgets.ButtonCreateChat),
                )),
                nil, nil, nil,
            ),
        ),
    )
    g.Base.OnSelected = func(_ *container.TabItem) {
        s.ScreenChatFocusDefault()
    }
    return g
}
