package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type ButtonTextOptions struct {
    Label   string  `json:"label"`
}

type ButtonOptions struct {}

type DialogTextOptions struct {
    Buttons []ButtonTextOptions `json:"buttons"`
    Prompt  string              `json:"prompt"`
    Title   string              `json:"title"`
}

type DialogOptions struct {
    Size    []float32   `json:"size"`
}

type EntryTextOptions struct {
    Label   string  `json:"label"`
}

type EntryOptions struct {}

type GUITextOptions struct {
    ButtonAttach            ButtonTextOptions   `json:"button_attach"`
    ButtonOptions           ButtonTextOptions   `json:"button_options"`
    ButtonSend              ButtonTextOptions   `json:"button_send"`
    DialogConnIssues        DialogTextOptions   `json:"dialog_connection_issues"`
    DialogLogin             DialogTextOptions   `json:"dialog_login"`
    DialogLoginAltPrompt    string              `json:"dialog_login_altprompt"`
    EntryMessage            EntryTextOptions    `json:"entry_message"`
    EntryUsername           EntryTextOptions    `json:"entry_username"`
    Window                  WindowTextOptions   `json:"window"`
}

type GUIOptions struct {
    DialogConnIssues    DialogOptions   `json:"dialog_connection_issues"`
    DialogLogin         DialogOptions   `json:"dialog_login"`
    Language            string          `json:"language"`
    Window              WindowOptions   `json:"window"`
}

type NetOptions struct {
    InitialReconnCooldown   int `json:"initial_reconn_cooldown"`
    ReconnCooldown          int `json:"reconn_cooldown"`
}

type Options struct {
    GUI         GUIOptions
    GUIText     GUITextOptions
    Net         NetOptions
}

type WindowTextOptions struct {
    Title   string  `json:"title"`
}

type WindowOptions struct {
    Size    []float32   `json:"size"`
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

    // Open text gui options of the specified language
    f, err = os.OpenFile(
        fmt.Sprintf("%s/text/%s.json", pathToOptions, opt.GUI.Language), os.O_RDWR, 644,
    )
    if err != nil {
        Error.Printf("Failed to open text options: %s\n", err)
        return nil, err
    }

    // Decode text gui options
    decoder = json.NewDecoder(f)
    err = decoder.Decode(&opt.GUIText)
    if err != nil {
        Error.Printf("Failed to deocde text options: %s\n", err)
        return nil, err
    }

    return opt, nil
}
