package user

import (
    "fyne.io/fyne/v2/data/binding"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/protocol"
)

// Data of a user
type User struct {
    Bio         string
    MemberOf    []int
    Pfp         []byte
    Username    string
    Visibility  prot.UserVisibility
}

// Bindings to data of a user
type UserBind struct {
    Pfp         binding.Bytes
    Username    binding.String
}

// Data about the client's user and other user data received from the server
type UserCache struct {
    // The bio of this client's user
    ThisBio         string
    // IDs of chats this client's user is apart of
    ThisMemberOf    []int
    // The profile picture of this client's user
    ThisPfp         []byte
    // The username of this client's user
    ThisUsername    string
    // Cache of other users' data
    //  Key (string) - Username
    //  Val (User) - User data
    Users           map[string]*User
    // Binding to cache of other users' data
    UsersBind       map[string]*UserBind
}

// Caches all avaliable data of a user
func (u *UserCache) CacheUserBack(userRaw prot.UserRaw) {
    log.Info.Printf("Caching back user data of %s\n", userRaw.Username)

    user, exists := u.Users[userRaw.Username]
    if !exists { 
        user = &User{}
        u.UsersBind[userRaw.Username] = &UserBind{}
        user.Username = userRaw.Username
        u.UsersBind[user.Username].Username = binding.BindString(&user.Username)
        u.Users[userRaw.Username] = user
    } else {
        u.UsersBind[user.Username].Username.Set(userRaw.Username)
    }
    user.Bio = userRaw.Bio
    user.MemberOf = userRaw.MemberOf
    user.Visibility = userRaw.Visibility
}

// Caches only the data of a user required to render a message from or card of them
func (u *UserCache) CacheUserFront(usernames []string) {
    for _, username := range usernames {
        if _, exists := u.Users[username]; !exists {
            log.Info.Printf("Caching front user data of %s\n", username)

            u.Users[username] = &User{Username: username}
            u.UsersBind[username] = &UserBind{
                Username: binding.BindString(&u.Users[username].Username),
            }
        }
    }
}

func NewUserCache() *UserCache {
    return &UserCache{
        Users:      make(map[string]*User),
        UsersBind:  make(map[string]*UserBind),
    }
}
