package main

type GUIOptions struct {
    Window WidgetOptions   `json:"window"`
}

type WidgetOptions struct {
    Label   string  `json:"label"`
    Size    []int   `json:"size"`
}
