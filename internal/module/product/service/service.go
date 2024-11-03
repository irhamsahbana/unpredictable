package service

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductService = &productService{}

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) ImportProducts(ctx context.Context, req *entity.ImportProductsReq) error {
	v := adapter.Adapters.Validator

	if req.File == nil {
		log.Warn().Any("req", req).Msg("service::importProduct - Missing file")
		return errmsg.NewCustomErrors(400).SetMessage("Missing file")
	}

	file, err := req.File.Open()
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Failed to open file")
		return err
	}
	defer file.Close()

	// convert multipart.File to *os.File
	f, err := os.Create(ulid.Make().String() + "_products.csv")
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Failed to create file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to create file")
	}
	defer func() {
		f.Close()
		err = os.Remove(f.Name())
		if err != nil {
			log.Warn().Err(err).Any("file", f.Name()).Msg("service::importProduct - Failed to remove file")
		}
	}()

	// insert the file content to the new file
	if _, err := file.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	if _, err := f.ReadFrom(file); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Failed to read file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read file")
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	products := []*entity.Product{}

	if err := gocsv.UnmarshalFile(f, &products); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Failed to import product")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read CSV file")
	}

	req.Data = products

	for i, product := range products {
		if err := v.Validate(product); err != nil {
			log.Warn().Err(err).Any("req", req).Msg("service::importProduct - Invalid product data")
			return errmsg.NewCustomErrors(400).SetMessage(fmt.Sprintf("Invalid product data at row %d", i+2))
		}
	}

	return s.repo.ImportProducts(ctx, req)
}

func (s *productService) ImportProductGrammage(ctx context.Context, req *entity.ImportProductGrammageReq) error {
	v := adapter.Adapters.Validator

	if req.File == nil {
		log.Warn().Any("req", req).Msg("service::importProductGrammage - Missing file")
		return errmsg.NewCustomErrors(400).SetMessage("Missing file")
	}

	file, err := req.File.Open()
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductGrammage - Failed to open file")
		return err
	}
	defer file.Close()

	// convert multipart.File to *os.File
	f, err := os.Create(ulid.Make().String() + "_product_grammages.csv")
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductGrammage - Failed to create file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to create file")
	}
	defer func() {
		f.Close()
		err = os.Remove(f.Name())
		if err != nil {
			log.Warn().Err(err).Any("file", f.Name()).Msg("service::importProductGrammage - Failed to remove file")
		}
	}()

	// insert the file content to the new file
	if _, err := file.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductGrammage - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	if _, err := f.ReadFrom(file); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductGrammage - Failed to read file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read file")
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductGrammage - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	grammages := []*entity.ProductGrammage{}

	if err := gocsv.UnmarshalFile(f, &grammages); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductGrammage - Failed to import product grammage")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read CSV file")
	}

	req.Data = grammages

	for i, grammage := range grammages {
		if err := v.Validate(grammage); err != nil {
			log.Warn().Err(err).Any("payload", grammage).Msg("service::importProductGrammage - Invalid product grammage data")
			return errmsg.NewCustomErrors(400).SetMessage(fmt.Sprintf("Invalid product grammage data at row %d", i+2))
		}
	}

	return s.repo.ImportProductGrammage(ctx, req)
}

func (s *productService) ImportProductTransactions(ctx context.Context, req *entity.ImportProductTransactionsReq) error {
	v := adapter.Adapters.Validator

	if req.File == nil {
		log.Warn().Any("req", req).Msg("service::importProductTransactions - Missing file")
		return errmsg.NewCustomErrors(400).SetMessage("Missing file")
	}

	file, err := req.File.Open()
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductTransactions - Failed to open file")
		return err
	}
	defer file.Close()

	// convert multipart.File to *os.File
	f, err := os.Create(ulid.Make().String() + "_product_transactions.csv")
	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductTransactions - Failed to create file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to create file")
	}
	defer func() {
		f.Close()
		err = os.Remove(f.Name())
		if err != nil {
			log.Warn().Err(err).Any("file", f.Name()).Msg("service::importProductTransactions - Failed to remove file")
		}
	}()

	// insert the file content to the new file
	if _, err := file.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductTransactions - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	if _, err := f.ReadFrom(file); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductTransactions - Failed to read file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read file")
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductTransactions - Failed to seek file")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to seek file")
	}

	transactions := []*entity.ProductTransaction{}

	if err := gocsv.UnmarshalFile(f, &transactions); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("service::importProductTransactions - Failed to import product transactions")
		return errmsg.NewCustomErrors(400).SetMessage("Failed to read CSV file")
	}

	req.Data = transactions

	for i, transaction := range transactions {
		if err := v.Validate(transaction); err != nil {
			log.Warn().Err(err).Any("payload", transaction).Msg("service::importProductTransactions - Invalid product transaction data ")
			return errmsg.NewCustomErrors(400).SetMessage(fmt.Sprintf("Invalid product transaction data at row %d", i+2))
		}
	}

	return s.repo.ImportProductTransactions(ctx, req)
}

func (s *productService) CreateProductTransaction(ctx context.Context, req *entity.CreateProductTransactionReq) error {
	return s.repo.CreateProductTransaction(ctx, req)
}

func (s *productService) GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error) {
	return s.repo.GetProducts(ctx, req)
}

func (s *productService) GetProductTransactions(ctx context.Context, req *entity.GetProductTransactionsReq) (*entity.GetProductTransactionsResp, error) {
	return s.repo.GetProductTransactions(ctx, req)
}

func (s *productService) GetProductGrammages(ctx context.Context, req *entity.GetProductGrammagesReq) (*entity.GetProductGrammagesResp, error) {
	return s.repo.GetProductGrammages(ctx, req)
}
