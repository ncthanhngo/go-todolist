package common

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"` // truyen form hay khong ko can quan tam, ko lay;
}

func (p *Paging) Process() {
	if p.Page < 1 {
		p.Page = 1

	}
	if p.Limit < 1 {
		p.Limit = 1
	}
	if p.Limit >= 100 {
		p.Limit = 100
	}

}
