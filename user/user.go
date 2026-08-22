package user

import (
    "fyne.io/fyne/v2/data/binding"
)

// Data of a user
type User struct {
    Bio         string
    Pfp         []byte
    Username    string
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

func (u *UserCache) CacheUsernames(usernames []string) {
    for _, username := range usernames {
        if user, exists := u.Users[username]; !exists {
            user = &User{Username: username}
            u.UsersBind[username] = &UserBind{
                Username: binding.BindString(&user.Username),
            }
        }
    }
}
