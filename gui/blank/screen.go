package sblank

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func InitScreenBlank() *fyne.Container {
    silly := widget.NewLabel(":P")
    silly.Importance = widget.LowImportance
    return container.NewCenter(silly)
}
