package main

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
    GUI GUIOptions  `json:"gui"`
    Net NetOptions  `json:"net"`
}

type WindowLabelOptions struct {
    Title   string  `json:"title"`
}

type WindowOptions struct {
    Size    []int   `json:"size"`
}
