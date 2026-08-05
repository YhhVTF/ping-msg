package dialogs

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
)

type DialogConnIssues struct {
    Dialog                  *dialog.CustomDialog
    Title                   string
    ContentBase             *fyne.Container
    LabelTechnicalErrMsg    *widget.Label
    LabelUXErrMsg           *widget.Label
}

// DialogConnectionIssues: Creates and shows a dialog set to the default size that informs the user that there are connection issues, with a user friendly message and a technical message
// Parameters:
//
//	err (error) - The error that occurred. This will be used as the technical error message
//  opt (*options.Options) - Options/settings
func InitDialogConnIssues(w fyne.Window, err error, opt *options.Options) *DialogConnIssues {
	log.Info.Printf("Creating Dialog ConnIssues\n")

    d := &DialogConnIssues{}

	// Create the user friendly error message as a label
	d.LabelUXErrMsg = widget.NewLabel(opt.GUIText.DialogConnIssues.Prompt)
	// Create the technical error message as a label, make it selectable and a low importance widget
	d.LabelTechnicalErrMsg = widget.NewLabel(err.Error())
	d.LabelTechnicalErrMsg.Selectable = true
	d.LabelTechnicalErrMsg.Importance = widget.LowImportance
    d.LabelTechnicalErrMsg.Wrapping = fyne.TextWrapWord
	// add then to a new vbox
	d.ContentBase = container.NewVBox(d.LabelUXErrMsg, d.LabelTechnicalErrMsg)

	// Create a dialog with the vbox as the content
	d.Dialog =
        dialog.NewCustom(opt.GUIText.DialogConnIssues.Title, "", d.ContentBase, w)

    //d.Dialog.SetButtons([]fyne.CanvasObject{})

	// Resize to default dialog size and show the dialog
	d.Dialog.Resize(fyne.NewSize(opt.GUI.DialogConnIssues.Size[0], opt.GUI.DialogConnIssues.Size[1]))
	d.Dialog.Show()

    return d
}
