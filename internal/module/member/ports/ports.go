package ports

import (
	"codebase-app/internal/module/member/entity"
	"context"
)

type MemberRepository interface {
	ImportMembers(ctx context.Context, req *entity.ImportMembersReq) error

	GetMembers(ctx context.Context, req *entity.GetMembersReq) (*entity.GetMembersResp, error)
}

type MemberService interface {
	ImportMembers(ctx context.Context, req *entity.ImportMembersReq) error

	GetMembers(ctx context.Context, req *entity.GetMembersReq) (*entity.GetMembersResp, error)
}
