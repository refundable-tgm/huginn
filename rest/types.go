package rest

// User data input
type User struct {
	// Username of the user
	Username string `json:"username"`
	// Password of the user
	Password string `json:"password"`
}

// TokenPair consists of an access and a refresh token
type TokenPair struct {
	// the access token
	AccessToken string `json:"access_token"`
	// the refresh token
	RefreshToken string `json:"refresh_token"`
}

// Error maps an error message
type Error struct {
	// the message that should be sent
	Message string `json:"error"`
}

// Information maps an information message
type Information struct {
	// the message that should be sent
	Message string `json:"info"`
}

// RefreshToken maps an refresh token
type RefreshToken struct {
	// Token is the refresh token
	Token string `json:"refresh_token"`
}

// Permissions lists the permissions of a teacher
type Permissions struct {
	// SuperUser permission
	SuperUser bool `json:"super_user"`
	// Administration permission
	Administration bool `json:"administration"`
	// AV permission
	AV bool `json:"av"`
	// PEK permission
	PEK bool `json:"pek"`
}

// News is a news object for applications
type News struct {
	// UUID of the application
	UUID string `json:"uuid"`
	// Title of the application
	Title string `json:"title"`
	// State of the application
	State int `json:"state"`
	// LastChanged is the date of last changes of the application
	LastChanged string `json:"last_changed"`
}

// PDF represents a pdf file
type PDF struct {
	// Content is the content of this file
	Content string `json:"pdf"`
}

// PDFs is a wrapper for a single pdf
type PDFs struct {
	Files []PDF `json:"files"`
}

// Excel represents an excel output
type Excel struct {
	// Content is the content of the excel file
	Content string `json:"excel"`
}
