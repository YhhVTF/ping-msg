package options

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/YhhVTF/ping-msg/log"
)

type ButtonTextOptions struct {
    Label   string  `json:"label"`
}

type ButtonOptions struct {}

type CardTextOptions struct {
    Subtitle    string  `json:"subtitle"`
    Title       string  `json:"title"`
}

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
    OptionCardLanguage      CardTextOptions     `json:"optioncard_language"`
    Window                  WindowTextOptions   `json:"window"`
    WindowOptions           WindowTextOptions   `json:"window_options"`
}

type GUIOptions struct {
    DialogConnIssues    DialogOptions   `json:"dialog_connection_issues"`
    DialogLogin         DialogOptions   `json:"dialog_login"`
    Language            string          `json:"language"`
    Window              WindowOptions   `json:"window"`
    WindowOptions       WindowOptions   `json:"window_options"`
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
    log.Info.Printf("Loading options\n")

    opt := &Options{}

    // Open net options
    f, err := os.OpenFile(
        fmt.Sprintf("%s/net.json", pathToOptions), os.O_RDWR, 644,
    )
    if err != nil {
        log.Error.Printf("Failed to open net options: %s\n", err)
        return nil, err
    }

    // Decode net options
    decoder := json.NewDecoder(f)
    err = decoder.Decode(&opt.Net)
    if err != nil {
        log.Error.Printf("Failed to deocde net options: %s\n", err)
        return nil, err
    }

    // Open non text gui options
    f, err = os.OpenFile(
        fmt.Sprintf("%s/gui.json", pathToOptions), os.O_RDWR, 644,
    )
    if err != nil {
        log.Error.Printf("Failed to open gui options: %s\n", err)
        return nil, err
    }

    // Decode non text options
    decoder = json.NewDecoder(f)
    err = decoder.Decode(&opt.GUI)
    if err != nil {
        log.Error.Printf("Failed to deocde gui options: %s\n", err)
        return nil, err
    }

    // Open text gui options of the specified language
    f, err = os.OpenFile(
        fmt.Sprintf("%s/text/%s.json", pathToOptions, opt.GUI.Language), os.O_RDWR, 644,
    )
    if err != nil {
        log.Error.Printf("Failed to open text options: %s\n", err)
        return nil, err
    }

    // Decode text gui options
    decoder = json.NewDecoder(f)
    err = decoder.Decode(&opt.GUIText)
    if err != nil {
        log.Error.Printf("Failed to deocde text options: %s\n", err)
        return nil, err
    }

    return opt, nil
}

func (opt *Options) SaveGUI(pathToOptions string) error {
    log.Info.Printf("Saving GUI options\n")

    // Open gui options
    f, err := os.OpenFile(fmt.Sprintf("%s/gui.json", pathToOptions), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 644)
    if err != nil {
        log.Error.Printf("Failed to open gui options: %s\n", err)
        return err
    }

    // Encode gui options and save them
    encoder := json.NewEncoder(f)
    err = encoder.Encode(opt.GUI)
    if err != nil {
        log.Error.Printf("Failed to encode and save gui options: %s\n", err)
        return err
    }

    return nil
}

func (opt *Options) SaveNet(pathToOptions string) error {
    log.Info.Printf("Saving Net options\n")

    // Open net options
    f, err := os.OpenFile(fmt.Sprintf("%s/net.json", pathToOptions), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 644)
    if err != nil {
        log.Error.Printf("Failed to open net options: %s\n", err)
        return err
    }

    // Encode net options and save them
    encoder := json.NewEncoder(f)
    err = encoder.Encode(opt.Net)
    if err != nil {
        log.Error.Printf("Failed to encode and save net options: %s\n", err)
        return err
    }

    return nil
}
