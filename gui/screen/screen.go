package screen

import (
)

const S_CHAT_FLOAT     = "CHAT_FLOAT"
const S_CHAT_FULL      = "CHAT_FULL"
const S_OPTIONS_FLOAT  = "OPTIONS_FLOAT"
const S_OPTIONS_FULL   = "OPTIONS_FULL"

type ScreenManager struct {
    Chan chan string
}

func NewScreenManager() *ScreenManager {
    return &ScreenManager{
        Chan: make(chan string),
    }
}

func (s *ScreenManager) ScreenChatFull() {
    s.Chan <- S_CHAT_FULL
}

func (s *ScreenManager) ScreenOptionsFloat() {
    s.Chan <- S_OPTIONS_FLOAT
}
