package screen

import (
)

const SS_CHAT_FLOAT     = "CHAT_FLOAT"
const SS_CHAT_FULL      = "CHAT_FULL"
const SS_OPTIONS_FLOAT  = "OPTIONS_FLOAT"
const SS_OPTIONS_FULL   = "OPTIONS_FULL"

type ScreenManager struct {
    Chan chan string
}

func NewScreenManager() *ScreenManager {
    return &ScreenManager{
        Chan: make(chan string),
    }
}

func (ss *ScreenManager) ScreenOptionsFloat() {
    ss.Chan <- SS_OPTIONS_FLOAT
}
