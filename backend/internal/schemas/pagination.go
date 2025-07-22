package schemas

type Pagination struct {
	Limit      int   `json:"limit"`
	Page       int   `json:"page"`
	Offset     int   `json:"-"`
	TotalPages int   `json:"total_pages"`
	TotalRows  int64 `json:"total_rows"`
}

func (p *Pagination) SetPagination(limit int, page int) {
	if limit <= 0 {
		p.Limit = 10
	} else {
		p.Limit = limit
	}

	if page <= 0 {
		p.Page = 1
	} else {
		p.Page = page
	}

	p.Offset = (p.Page - 1) * p.Limit
}

func (p *Pagination) SetTotalPages(totalRows int64) {
	p.TotalRows = totalRows
	if totalRows == 0 || p.Limit == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = int((totalRows + int64(p.Limit) - 1) / int64(p.Limit))
	}
}
