package user

type User struct {
	Title  string   `json:"title"`
	Amount float32  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"Tags"`
}

type Err struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}
