package main

type GUIOptions struct {
    ButtonAttach        ButtonOptions   `json:"button_attach"`
    ButtonOptions       ButtonOptions   `json:"button_options"`
    ButtonSend          ButtonOptions   `json:"button_send"`
    DialogConnIssues    DialogOptions   `json:"dialog_connection_issues"`
    DialogLogin         DialogOptions   `json:"dialog_login"`
    EntryMessage        EntryOptions    `json:"entry_message"`
    EntryUsername       EntryOptions    `json:"entry_username"`
    Window              WindowOptions   `json:"window"`
}

type ButtonOptions struct {
    Label   string  `json:"label"`
}

type DialogOptions struct {
    Buttons []ButtonOptions `json:"buttons"`
    Prompt  string          `json:"prompt"`
    Size    []int           `json:"size"`
    Title   string          `json:"title"`
}

type EntryOptions struct {
    Label   string  `json:"label"`
}

type WindowOptions struct {
    Size    []int   `json:"size"`
    Title   string  `json:"title"`
}
