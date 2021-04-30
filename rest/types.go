package rest

// User data input
type User struct {
	// Username of the user
	Username string `json:"username" example:"lehrer1234"`
	// Password of the user
	Password string `json:"password" example:"password1234"`
}

// TokenPair consists of an access and a refresh token
type TokenPair struct {
	// the access token
	AccessToken string `json:"access_token" example:"<jwt-token>"`
	// the refresh token
	RefreshToken string `json:"refresh_token" example:"<jwt-token>"`
}

// Error maps an error message
type Error struct {
	// the message that should be sent
	Message string `json:"error" example:"couldn't convert token"`
}

// Information maps an information message
type Information struct {
	// the message that should be sent
	Message string `json:"info" example:"updated teacher successfully"`
}

// TeacherInformation contains changeable information of the teacher
type TeacherInformation struct {
	// Degree of the Teacher
	Degree string `json:"degree" example:"DI"`
	// Title of the Teacher
	Title string `json:"title" example:"Prof"`
	// The Staffnr of the regarding teacher
	Staffnr int `json:"staffnr" example:"938503154"`
	// The Group number
	Group int `json:"group" example:"1"`
	// The StartingAddresses of the teacher
	StartingAddresses []string `json:"starting_addresses" example:"Zuhause 1, Zuhause 2"`
	// The TripGoals the teacher visited before
	TripGoals []string `json:"trip_goals" example:"Karl Hönck Heim, PH Wien, Landesgericht St. Pölten"`
	// The Departments this teacher belongs to
	Departments []string `json:"departments" example:"HIT, HBG"`
}

// RefreshToken maps an refresh token
type RefreshToken struct {
	// Token is the refresh token
	Token string `json:"refresh_token" example:"<jwt-token>"`
}

// Permissions lists the permissions of a teacher
type Permissions struct {
	// SuperUser permission
	SuperUser bool `json:"super_user" example:"true"`
	// Administration permission
	Administration bool `json:"administration" example:"true"`
	// AV permission
	AV bool `json:"av" example:"true"`
	// PEK permission
	PEK bool `json:"pek" example:"true"`
}

// News is a news object for applications
type News struct {
	// UUID of the application
	UUID string `json:"uuid" example:"3fcf7f67-e0ed-4339-99b4-a6765aaa3dc4"`
	// Title of the application
	Title string `json:"title" example:"Sommersportwoche"`
	// State of the application
	State int `json:"state" example:"3"`
	// LastChanged is the date of last changes of the application
	LastChanged string `json:"last_changed" example:"2009-11-10 23:00:00 +0000 UTC m=+0.000000001"`
}

// PDF represents a pdf file
type PDF struct {
	// Content is the content of this file
	Content string `json:"pdf" example:"<base64>"`
}

// PDFs is a wrapper for a single pdf
type PDFs struct {
	Files []PDF `json:"files"`
}

// Excel represents an excel output
type Excel struct {
	// Content is the content of the excel file
	Content string `json:"excel" example:"<base64>"`
}
