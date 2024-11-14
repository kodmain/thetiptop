package transfert

import "github.com/kodmain/thetiptop/api/internal/infrastructure/data"

type Ticket struct {
	ID       *string `json:"id" xml:"id" form:"id"`
	Prize    *string `json:"prize" xml:"prize" form:"prize"`
	ClientID *string `json:"client_id" xml:"client_id" form:"client_id"`
	Token    *string `json:"token" xml:"token" form:"token"`
}

func (c *Ticket) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"id":        c.ID,
		"prize":     c.Prize,
		"client_id": c.ClientID,
		"token":     c.Token,
	})
}
