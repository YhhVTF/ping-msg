package shortcuts

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/driver/desktop"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

var ShortcutDebugNewChat =
    &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt}

func DoDebugNewChat(
    chatsListW *widget.List, c *chat.ChatCache, u *user.UserCache,
) {
    u.ThisUser.MemberOf = []int{ 1 }
    c.Chats[1] = chat.NewChat(&prot.ChatMetadataRaw{ID: 1})
    chatsListW.Refresh()
}
