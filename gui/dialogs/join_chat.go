package dialogs

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"

    "strconv"

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
    d *DialogJoinChat, done chan int, entryBind binding.String,
    promptBind binding.String, opt *options.Options,
) *widget.Button {
    return widget.NewButton("Join", func() {
        text, _ := entryBind.Get()
        chatID, err := strconv.Atoi(text)
        if err != nil {
            promptBind.Set("Not a valid chat ID, try again")
            return
        }
        done <- chatID
        d.Dialog.Dismiss()
    })
}

func createJoinChatEntry(
    d *DialogJoinChat, done chan int, entryBind binding.String, opt *options.Options,
) *widget.Entry {
    entryW := widget.NewEntryWithData(entryBind)
    entryW.OnSubmitted = func(text string) {
        defer d.Dialog.Dismiss()
        chatID, err := strconv.Atoi(text)
        if err != nil {
            done <- prot.NONE_INT
            return
        }
        done <- chatID
    }
    return entryW
}

func createJoinChatPrompt(promptBind binding.String, opt *options.Options) *widget.Label {
    promptBind.Set("Enter the ID of the chat you want to join")
    return widget.NewLabelWithData(promptBind)
}

func InitDialogJoinChat(
    w fyne.Window, done chan int, opt *options.Options,
) *DialogJoinChat {
    d := &DialogJoinChat{}

    promptText := ""
    promptBind := binding.BindString(&promptText)
    d.Prompt = createJoinChatPrompt(promptBind, opt)

    entryText := ""
    entryBind := binding.BindString(&entryText)
    d.Entry = createJoinChatEntry(d, done, entryBind, opt)

    d.ButtonJoin = createJoinChatButtonJoin(d, done, entryBind, promptBind, opt)

    d.Dialog =
        dialog.NewCustom("Join chat", "", container.NewVBox(d.Prompt, d.Entry), w)

    d.Dialog.SetButtons([]fyne.CanvasObject{
        d.ButtonJoin, d.ButtonCancel,
    })
    return d
}
