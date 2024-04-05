package model

import "avito_banners/internal/errs"

type Banner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

func (b *Banner) Validate() (ers []errs.Error) {
	if b.Text == "" {
		ers = append(ers, errs.GetErr(101))
	}

	if b.Title == "" {
		ers = append(ers, errs.GetErr(102))
	}

	if b.Url == "" {
		ers = append(ers, errs.GetErr(103))
	}

	return
}
