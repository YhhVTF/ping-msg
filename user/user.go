package user

// Data of a user
type User struct {
    Bio         string
    ID          int
    Pfp         []byte
    Username    string
}

// Data about the client's user and other user data received from the server
type UserCache struct {
    // The bio of this client's user
    ThisBio         string
    // The profile picture of this client's user
    ThisPfp         []byte
    // The user ID of this client's user
    ThisUserID      int
    // The username of this client's user
    ThisUsername    string
    // Cache of other users' data
    //  Key (int) - User ID
    //  Val (User) - User data
    Users           map[int]User
}
