package schat

import (
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"fmt"

	"github.com/YhhVTF/ping-msg/chat"
	"github.com/YhhVTF/ping-msg/log"
	options "github.com/YhhVTF/ping-msg/opt"
	"github.com/YhhVTF/ping-msg/user"
)

// A widget representing a reply to a message. Is a component of Message widgets that reply to other messages
type RepliedMessage struct {
	Base        *fyne.Container
	Icon        *widget.Icon
	MessageID   int
	UsernameLbl *widget.Label
	Text        *widget.RichText
}

func createRepliedMessage(
	repliedMsg *chat.Message, repliedID int, chat *chat.Chat, u *user.UserCache, opt *options.Options,
) *RepliedMessage {
	log.Info.Printf("Creating replied message widget for message %d\n", repliedID)

	repliedMsgWidget := &RepliedMessage{}
	repliedMsgWidget.MessageID = repliedID

	formatContent := func(raw string) string {
		re := regexp.MustCompile(`(?m)^#{1,6}\s*`)
		cleanedContent := re.ReplaceAllString(raw, "")
		return "**" + strings.TrimSpace(cleanedContent) + "**"
	}

	var usernameStr, contentStr string
	if repliedMsg == nil {
		usernameStr = ""
		contentStr = opt.GUIText.Greeting.PlaceholderUnloaded
	} else {
		usernameStr = fmt.Sprintf("%s:", *repliedMsg.Username)
		contentStr = formatContent(repliedMsg.Content)
	}

	repliedMsgWidget.Icon = widget.NewIcon(theme.Current().Icon(theme.IconNameMailReply))

	//Username container label
	repliedMsgWidget.UsernameLbl = widget.NewLabel(usernameStr)
	repliedMsgWidget.UsernameLbl.TextStyle.Bold = true // MAKE USERNAME BOLD :3

	// Richtext content
	repliedMsgWidget.Text = widget.NewRichTextFromMarkdown(contentStr)
	repliedMsgWidget.Text.Wrapping = fyne.TextWrapWord

	if chat != nil && chat.MessagesBind[repliedID] != nil {
		chat.MessagesBind[repliedID].AddListener(binding.NewDataListener(func() {
			updatedContent, _ := chat.MessagesBind[repliedID].Get()

			formattedUpdate := formatContent(updatedContent)

			repliedMsgWidget.Text.ParseMarkdown(formattedUpdate)
			repliedMsgWidget.Text.Refresh()
		}))
	}

	repliedMsgWidget.Base = container.NewHBox(
		repliedMsgWidget.Icon,
		repliedMsgWidget.UsernameLbl,
		repliedMsgWidget.Text,
	)

	return repliedMsgWidget
}
