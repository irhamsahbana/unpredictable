package entity

import (
	"codebase-app/pkg/types"
	"mime/multipart"
)

type ImportMembersReq struct {
	File *multipart.FileHeader `form:"file" validate:"required"`

	Data []Member
}

type GetMembersReq struct {
	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`
}

func (r *GetMembersReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetMembersResp struct {
	Items []Member   `json:"data"`
	Meta  types.Meta `json:"meta"`
}

type Member struct {
	Id             string  `db:"id" csv:"MemberID" validate:"required" json:"id"`
	JoinDate       string  `db:"join_date" csv:"JoinDate" validate:"required,datetime=2006-01-02" json:"join_date"`
	DateOfBirth    *string `db:"date_of_birth" csv:"DateOfBirth,omitempty" validate:"omitempty,datetime=2006-01-02" json:"date_of_birth"`
	City           string  `db:"city" csv:"City" validate:"required" json:"city"`
	NoOfChild      int     `db:"no_of_child" csv:"NoOfChild" json:"no_of_child"`
	EldestKidDOB   string  `db:"eldest_kid_dob" csv:"EldestKidDOB" validate:"datetime=2006-01-02" json:"eldest_kid_dob"`
	YoungestKidDOB string  `db:"youngest_kid_dob" csv:"YoungestKidDOB" validate:"datetime=2006-01-02" json:"youngest_kid_dob"`

	Pass string `db:"password" json:"-"`
}
