package sopt

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/opt"
)

func selectLanguageOnChanged(language string, opt *options.Options) {
    opt.GUI.Language = language
}

func createSelectLanguage(opt *options.Options) *widget.Select {
    sel := widget.NewSelect([]string{
        "english", "entish", "gakotolo",
    }, func(language string) {
        selectLanguageOnChanged(language, opt)
    })
    sel.PlaceHolder = opt.GUI.Language
    return sel
}
