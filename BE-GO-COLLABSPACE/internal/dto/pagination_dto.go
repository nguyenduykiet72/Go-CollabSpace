package dto

type PaginationReq struct {
	Page  int `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *PaginationReq) GetOffset() int {
	page := p.Page
	if page <= 0 {
		page = 1
	}

	limit := p.GetLimit()

	return (page - 1) * limit
}

func (p *PaginationReq) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	return p.Limit
}
