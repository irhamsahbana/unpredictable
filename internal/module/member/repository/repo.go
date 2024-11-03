package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/member/entity"
	"codebase-app/internal/module/member/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.MemberRepository = &memberRepo{}

type memberRepo struct {
	db *sqlx.DB
}

func NewMemberRepository() *memberRepo {
	return &memberRepo{
		db: adapter.Adapters.Postgres,
	}
}

func (r *memberRepo) ImportMembers(ctx context.Context, req *entity.ImportMembersReq) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}

	// Prepare the insert query
	query := `
		INSERT INTO members (id, join_date, date_of_birth, city, no_of_child, eldest_kid_dob, youngest_kid_dob, password)
		VALUES (:id, :join_date, :date_of_birth, :city, :no_of_child, :eldest_kid_dob, :youngest_kid_dob, :password)
	`

	// Iterate over each member in req.Data and execute the insert query
	for _, member := range req.Data {
		_, err := tx.NamedExecContext(ctx, query, member)
		if err != nil {
			log.Error().Err(err).Any("member", member).Msg("failed to insert member")
			tx.Rollback() // Roll back transaction on error
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction")
		return err
	}

	return nil
}

func (r *memberRepo) GetMembers(ctx context.Context, req *entity.GetMembersReq) (*entity.GetMembersResp, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.Member
	}
	var (
		res  = new(entity.GetMembersResp)
		data = make([]dao, 0)
	)
	res.Items = make([]entity.Member, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			join_date,
			date_of_birth,
			city,
			no_of_child,
			eldest_kid_dob,
			youngest_kid_dob,
			password
		FROM members
		LIMIT ? OFFSET ?
		`

	if err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Paginate, (req.Page-1)*req.Paginate); err != nil {
		log.Error().Err(err).Msg("failed to get members")
		return nil, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.Member)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
