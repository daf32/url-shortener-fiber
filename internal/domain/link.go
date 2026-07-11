package domain

type Link struct {
	Code     string `json:"code"`
	Original string `json:"original"`
}

func NewLink(code string, original string) Link {
	return Link{
		Code:     code,
		Original: original,
	}
}
