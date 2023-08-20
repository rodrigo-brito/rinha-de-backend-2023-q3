package entity

type User struct {
	ID    string   `json:"id"`
	Nick  string   `json:"apelido"`
	Name  string   `json:"nome"`
	Birth string   `json:"nascimento"`
	Stack []string `json:"stack"`
}
