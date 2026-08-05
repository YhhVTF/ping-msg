package screen

import (
)

const S_CHAT_FLOAT          = "CHAT_FLOAT"
const S_CHAT_FOCUS_DEFAULT  = "CHAT_FOCUS_DEFAULT"
const S_CHAT_FULL           = "CHAT_FULL"
const S_OPTIONS_FLOAT       = "OPTIONS_FLOAT"
const S_OPTIONS_FULL        = "OPTIONS_FULL"

type ScreenManager struct {
    Chan chan string
}

func NewScreenManager() *ScreenManager {
    return &ScreenManager{
        Chan: make(chan string),
    }
}

// This should be used ONLY when the code using it does not have access to chat screen's widgets, since it is slower than just doing it directly with fyne.Window.Canvas().Focus()
func (s *ScreenManager) ScreenChatFocusDefault() {
    s.Chan <- S_CHAT_FOCUS_DEFAULT
}

func (s *ScreenManager) ScreenChatFull() {
    s.Chan <- S_CHAT_FULL
}

func (s *ScreenManager) ScreenOptionsFloat() {
    s.Chan <- S_OPTIONS_FLOAT
}
