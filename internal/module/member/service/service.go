package service

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/member/entity"
	"codebase-app/internal/module/member/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.MemberService = &memberService{}

type memberService struct {
	repo ports.MemberRepository
}

func NewMemberService(repo ports.MemberRepository) *memberService {
	return &memberService{
		repo: repo,
	}
}

func (s *memberService) ImportMembers(ctx context.Context, req *entity.ImportMembersReq) error {
	var v = adapter.Adapters.Validator
	if req.File == nil {
		log.Warn().Any("req", req).Msg("service::importMember - Missing file")
		return errmsg.NewCustomErrors(400).SetMessage("Missing file")
	}

	// Open the uploaded file
	file, err := req.File.Open()
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to open file")
		return err
	}
	defer file.Close()

	// convert multipart.File to *os.File
	f, err := os.Create(ulid.Make().String() + "_members.csv")
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to create file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to create file")
	}

	// insert the file content to the new file
	if _, err := file.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	if _, err := f.ReadFrom(file); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to read file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read file")
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}
	defer os.Remove(f.Name())

	// Read the file content
	members := make([]entity.Member, 0)

	if err := gocsv.UnmarshalFile(f, &members); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to unmarshal file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read CSV file")
	}

	req.Data = members

	for i, member := range members {
		if err := v.Validate(member); err != nil {
			log.Warn().Err(err).Any("req", req).Msg("service::importMember - Failed to validate member")
			return errmsg.NewCustomErrors(400).SetMessage("Invalid member data at row " + strconv.Itoa(i+2))
		}

		// Set default password if not provided
		if member.Pass == "" {
			members[i].Pass = generatePassword()
		}
	}

	// Call the repository to import members
	return s.repo.ImportMembers(ctx, req)

}

// Helper function to generate a password
func generatePassword() string {
	// Replace this with your preferred password generation logic
	return "$2y$10$7YY20k1zrYeOMBmGNu/3lufnHDgsdBkyQ.6bvzphC3ovG.3t5W4oO" // bcrypt hash of "password"
}

func (s *memberService) GetMembers(ctx context.Context, req *entity.GetMembersReq) (*entity.GetMembersResp, error) {
	return s.repo.GetMembers(ctx, req)
}
