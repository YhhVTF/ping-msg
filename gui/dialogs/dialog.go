package dialogs

// All dialogs to be used by chat screen
type DialogTable struct {
	// Informs the user that there are issues with connecting to the server
	ConnectionIssues *DialogConnIssues
    // Asks the user for information to log in
    Login *DialogLogin
}
