package ports

import (
	"codebase-app/internal/module/product/entity"
	"context"
)

type ProductRepository interface {
	ImportProducts(ctx context.Context, req *entity.ImportProductsReq) error
	ImportProductGrammage(ctx context.Context, req *entity.ImportProductGrammageReq) error
	ImportProductTransactions(ctx context.Context, req *entity.ImportProductTransactionsReq) error

	CreateProductTransaction(ctx context.Context, req *entity.CreateProductTransactionReq) error

	GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error)
	GetProductTransactions(ctx context.Context, req *entity.GetProductTransactionsReq) (*entity.GetProductTransactionsResp, error)
	GetProductGrammages(ctx context.Context, req *entity.GetProductGrammagesReq) (*entity.GetProductGrammagesResp, error)
}

type ProductService interface {
	ImportProducts(ctx context.Context, req *entity.ImportProductsReq) error
	ImportProductGrammage(ctx context.Context, req *entity.ImportProductGrammageReq) error
	ImportProductTransactions(ctx context.Context, req *entity.ImportProductTransactionsReq) error

	CreateProductTransaction(ctx context.Context, req *entity.CreateProductTransactionReq) error

	GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error)
	GetProductTransactions(ctx context.Context, req *entity.GetProductTransactionsReq) (*entity.GetProductTransactionsResp, error)
	GetProductGrammages(ctx context.Context, req *entity.GetProductGrammagesReq) (*entity.GetProductGrammagesResp, error)
}
