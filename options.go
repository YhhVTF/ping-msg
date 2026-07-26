package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type ButtonLabelOptions struct {
    Label   string  `json:"label"`
}

type ButtonOptions struct {}

type DialogLabelOptions struct {
    Buttons []ButtonLabelOptions    `json:"buttons"`
    Prompt  string                  `json:"prompt"`
    Title   string                  `json:"title"`
}

type DialogOptions struct {
    Size    []int                   `json:"size"`
}

type EntryLabelOptions struct {
    Label   string  `json:"label"`
}

type EntryOptions struct {}

type GUILabelOptions struct {
    ButtonAttach        ButtonLabelOptions   `json:"button_attach"`
    ButtonLabelOptions  ButtonLabelOptions   `json:"button_options"`
    ButtonSend          ButtonLabelOptions   `json:"button_send"`
    DialogConnIssues    DialogLabelOptions   `json:"dialog_connection_issues"`
    DialogLogin         DialogLabelOptions   `json:"dialog_login"`
    EntryMessage        EntryLabelOptions    `json:"entry_message"`
    EntryUsername       EntryLabelOptions    `json:"entry_username"`
    Window              WindowLabelOptions   `json:"window"`
}

type GUIOptions struct {
    DialogConnIssues    DialogOptions   `json:"dialog_connection_issues"`
    DialogLogin         DialogOptions   `json:"dialog_login"`
    Window              WindowOptions   `json:"window"`
}

type LanguageOptions struct {
    Language    string  `json:"language"`
}

type NetOptions struct {
    InitialReconnCooldown   int `json:"initial_reconn_cooldown"`
    ReconnCooldown          int `json:"reconn_cooldown"`
}

type Options struct {
    GUI         GUIOptions
    GUILabels   GUILabelOptions
    Language    LanguageOptions
    Net         NetOptions
}

type WindowLabelOptions struct {
    Title   string  `json:"title"`
}

type WindowOptions struct {
    Size    []int   `json:"size"`
}

func LoadOptions(pathToOptions string) (*Options, error) {
    Info.Printf("Loading options\n")

    opt := &Options{}

    // Open net options
    f, err := os.OpenFile(
        fmt.Sprintf("%s/net.json", pathToOptions), os.O_RDWR, 644,
    )
    if err != nil {
        Error.Printf("Failed to open net options: %s\n", err)
        return nil, err
    }

    // Decode net options
    decoder := json.NewDecoder(f)
    err = decoder.Decode(&opt.Net)
    if err != nil {
        Error.Printf("Failed to deocde net options: %s\n", err)
        return nil, err
    }

    // Open non text gui options
    f, err = os.OpenFile(
        fmt.Sprintf("%s/gui.json", pathToOptions), os.O_RDWR, 644,
    )
    if err != nil {
        Error.Printf("Failed to open gui options: %s\n", err)
        return nil, err
    }

    // Decode non text options
    decoder = json.NewDecoder(f)
    err = decoder.Decode(&opt.GUI)
    if err != nil {
        Error.Printf("Failed to deocde gui options: %s\n", err)
        return nil, err
    }

    // Open language options
    f, err = os.OpenFile(
        fmt.Sprintf("%s/language.json", pathToOptions), os.O_RDWR, 644,
    )
    if err != nil {
        Error.Printf("Failed to open language options: %s\n", err)
        return nil, err
    }

    // Decode languae options
    decoder = json.NewDecoder(f)
    err = decoder.Decode(&opt.Language)
    if err != nil {
        Error.Printf("Failed to deocde language options: %s\n", err)
        return nil, err
    }

    // Open text gui options of the specified language
    f, err = os.OpenFile(
        fmt.Sprintf("%s/text/%s.json", pathToOptions, opt.Language), os.O_RDWR, 644,
    )
    if err != nil {
        Error.Printf("Failed to open text options: %s\n", err)
        return nil, err
    }

    // Decode text gui options
    decoder = json.NewDecoder(f)
    err = decoder.Decode(&opt.GUILabels)
    if err != nil {
        Error.Printf("Failed to deocde text options: %s\n", err)
        return nil, err
    }

    return opt, nil
}
