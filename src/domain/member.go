package domain

type Member struct {
	ID   int    `json:"id"`
	Club Club   `json:"club"`
	Name string `json:"name"`
	Debt int    `json:"debt"`
}
