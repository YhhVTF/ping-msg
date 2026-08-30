package dialogs

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"

    "errors"
    "strconv"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
)

// Dialog prompting the user for a chat ID. The user will then be made a member of that chat
type DialogJoinChat struct {
    ButtonCancel    *widget.Button
    ButtonJoin      *widget.Button
    Dialog          *dialog.CustomDialog
    Entry           *widget.Entry
    Prompt          *widget.Label
}

func createJoinChatButtonCancel(
    d *DialogJoinChat, done chan int, opt *options.Options,
) *widget.Button {
    return widget.NewButton("Cancel", func() {
        defer d.Dialog.Dismiss()
        done <- prot.NONE_INT
    })
}

func createJoinChatButtonJoin(
    d *DialogJoinChat, entryW *widget.Entry, opt *options.Options,
) *widget.Button {
    return widget.NewButton("Join", func() {
        entryW.OnSubmitted(entryW.Text)
    })
}

func createJoinChatEntry(
    d *DialogJoinChat, done chan int, promptW *widget.Label, opt *options.Options,
) *widget.Entry {
    entryW := widget.NewEntry()
    entryW.OnSubmitted = func(text string) {
        if err := entryW.Validate(); err != nil {
            promptW.SetText(err.Error())
            return
        }
        chatID, _ := strconv.Atoi(text)
        done <- chatID
        d.Dialog.Dismiss()
    }
    entryW.Validator = func(text string) error {
        _, err := strconv.Atoi(text)
        if err != nil {
            return errors.New("Not a valid chat ID, try again")
        }
        return nil
    }
    return entryW
}

func createJoinChatPrompt(opt *options.Options) *widget.Label {
    return widget.NewLabel("Enter the ID of the chat you want to join")
}

func InitDialogJoinChat(
    w fyne.Window, done chan int, opt *options.Options,
) *DialogJoinChat {
    log.Info.Printf("Creating dialog JoinChat\n")
    d := &DialogJoinChat{}

    d.Prompt = createJoinChatPrompt(opt)

    d.Entry = createJoinChatEntry(d, done, d.Prompt, opt)

    d.ButtonJoin = createJoinChatButtonJoin(d, d.Entry, opt)

    d.ButtonCancel = createJoinChatButtonCancel(d, done, opt)

    d.Dialog =
        dialog.NewCustom("Join chat", "", container.NewVBox(d.Prompt, d.Entry), w)

    d.Dialog.SetButtons([]fyne.CanvasObject{
        d.ButtonJoin, d.ButtonCancel,
    })
    d.Dialog.Show()
    return d
}
