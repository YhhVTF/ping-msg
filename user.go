package main

// Data about the client's user and other user data received from the server
type UserData struct {
    // The username of this client's user
    ThisUser string
    // -
    Users map[int]string
}
