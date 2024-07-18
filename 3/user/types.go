package user

//easyjson:json
type JsonUser struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company,nocopy"`
	Country  string   `json:"country,nocopy"`
	Email    string   `json:"email,nocopy"`
	Job      string   `json:"job,nocopy"`
	Name     string   `json:"name,nocopy"`
	Phone    string   `json:"phone,nocopy"`
}

//easyjson:json
//type JsonUser struct {
//	Browsers []string `json:"browsers"`
//	Company  string   `json:"company"`
//	Country  string   `json:"country"`
//	Email    string   `json:"email"`
//	Job      string   `json:"job"`
//	Name     string   `json:"name"`
//	Phone    string   `json:"phone"`
//}
