package entity

import "codebase-app/pkg/types"

type GetProductsReq struct {
	Page     int `query:"page" validate:"required,numeric"`
	Paginate int `query:"paginate" validate:"required,numeric"`
}

func (r *GetProductsReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetProductsResp struct {
	Items []Product  `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type GetProductTransactionsReq struct {
	Page     int `query:"page" validate:"required,numeric"`
	Paginate int `query:"paginate" validate:"required,numeric"`
}

func (r *GetProductTransactionsReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetProductTransactionsResp struct {
	Items []ProductTransaction `json:"items"`
	Meta  types.Meta           `json:"meta"`
}

type GetProductGrammagesReq struct {
	Page     int `query:"page" validate:"required,numeric"`
	Paginate int `query:"paginate" validate:"required,numeric"`
}

func (r *GetProductGrammagesReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetProductGrammagesResp struct {
	Items []ProductGrammage `json:"items"`
	Meta  types.Meta        `json:"meta"`
}
