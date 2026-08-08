package ping

import (
    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/user"
)

// Cache and data for each loaded chat
var ChatCache *chat.ChatCache

// Whether or not Ping is currently connected ot the server
var Connected bool

// All tweakable values in Ping
var Options *options.Options

// Whether or not the user has quit Ping
var Quit bool

// Create and exit screens
var ScreenManager *screen.ScreenManager

// Cache of user data
var UserCache *user.UserCache
