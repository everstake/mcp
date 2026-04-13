package dashboard

type BodyError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e BodyError) Error() string {
	return e.Description
}
