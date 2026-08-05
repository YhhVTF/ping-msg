package dialogs

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

type DialogLogin struct {
    ButtonLogin *widget.Button
    Dialog      *dialog.CustomDialog
    Title       string
    ContentBase *fyne.Container
    Prompt      *widget.Label
    Entry       *widget.Entry
}

func InitDialogLogin(
    w fyne.Window, s *screen.ScreenManager, u *user.UserCache, opt *options.Options,
) *DialogLogin {
	log.Info.Printf("Creating Dialog Login\n")

    d := &DialogLogin{}

    // Text telling the user what to do
	d.Prompt = widget.NewLabel(opt.GUIText.DialogLogin.Prompt)
    // Entry for username
    d.Entry = widget.NewEntry()
    d.Entry.SetPlaceHolder(opt.GUIText.EntryUsername.Label)

	// add them to a new vbox
    d.ContentBase = container.NewVBox(d.Prompt, d.Entry)

    // Create loging button
    d.ButtonLogin = widget.NewButton(opt.GUIText.DialogLogin.Buttons[0].Label, func() {
        // Set the text in the d.Entry as the username if it isn't a reserved username or empty
        if d.Entry.Text == "" { return }
        if d.Entry.Text == prot.NONE_STRING {
            d.Prompt.SetText(opt.GUIText.DialogLoginAltPrompt)
            return
        }
        u.ThisUsername = d.Entry.Text

        // Dismiss the dialog and set it as nil
        d.Dialog.Dismiss()
        d.Dialog.Hide()
        d.Dialog = nil

        // Set focus on message entry in chat screen now that dialog is dismissed
        s.ScreenChatFocusDefault()

        log.Info.Printf("Username set as %s\n", u.ThisUsername)
        log.Info.Printf("Dialog Login dismissed\n")
    })

	// Create a dialog with the vbox as the content
	d.Dialog = dialog.NewCustom(opt.GUIText.DialogLogin.Title, "", d.ContentBase, w)
	d.Dialog.SetButtons([]fyne.CanvasObject{ d.ButtonLogin })

	// Resize to default dialog size and show the dialog
	d.Dialog.Resize(fyne.NewSize(opt.GUI.DialogConnIssues.Size[0], opt.GUI.DialogConnIssues.Size[1]))
	d.Dialog.Show()

    d.Entry.OnSubmitted = func(text string) {
        // Set the text in the d.Entry as the username if it isn't a reserved username or empty
        if d.Entry.Text == "" { return }
        if d.Entry.Text == prot.NONE_STRING {
            d.Prompt.SetText(opt.GUIText.DialogLoginAltPrompt)
            return
        }
        u.ThisUsername = d.Entry.Text

        // Dismiss the dialog and set it as nil in the dialog table
        d.Dialog.Dismiss()
        d.Dialog.Hide()
        d.Dialog = nil

        // Set focus on message d.Entry now that dialog is dismissed
        s.ScreenChatFocusDefault()

        log.Info.Printf("Username set as %s\n", u.ThisUsername)
		log.Info.Printf("Dialog Login dismissed\n")
    }
    return d
}
