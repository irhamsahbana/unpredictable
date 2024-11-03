package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepo{}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository() *productRepo {
	return &productRepo{
		db: adapter.Adapters.Postgres,
	}
}

func (r *productRepo) ImportProducts(ctx context.Context, req *entity.ImportProductsReq) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}

	// Prepare the insert query
	query := `
		INSERT INTO products (id, name, category, level)
		VALUES (:id, :name, :category, :level)
	`

	// Iterate over each product in req.Data and execute the insert query
	for _, product := range req.Data {
		_, err := tx.NamedExecContext(ctx, query, product)
		if err != nil {
			log.Error().Err(err).Any("product", product).Msg("failed to insert product")
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

func (r *productRepo) ImportProductGrammage(ctx context.Context, req *entity.ImportProductGrammageReq) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}

	// Prepare the insert query
	query := `
		INSERT INTO product_grammages (id, name, point, price)
		VALUES (:id, :name, :point, :price)
	`

	// Iterate over each product in req.Data and execute the insert query
	for _, productGrammage := range req.Data {
		_, err := tx.NamedExecContext(ctx, query, productGrammage)
		if err != nil {
			log.Error().Err(err).Any("productGrammage", productGrammage).Msg("failed to insert product grammage")
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

func (r *productRepo) ImportProductTransactions(ctx context.Context, req *entity.ImportProductTransactionsReq) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}

	// Prepare the insert query
	query := `
		INSERT INTO product_transactions (
			id,
			member_id,
			product_id,
			product_grammage_id,
			source,
			qty,
			price_per_unit,
			created_at,
			is_training_data
		) VALUES (
		 	:id,
			:member_id,
			:product_id,
			:product_grammage_id,
			:source,
			:qty,
			:price_per_unit,
			:created_at,
			TRUE
		) ON CONFLICT (id) DO UPDATE SET
			member_id = EXCLUDED.member_id,
			product_id = EXCLUDED.product_id,
			product_grammage_id = EXCLUDED.product_grammage_id,
			source = EXCLUDED.source,
			qty = EXCLUDED.qty,
			price_per_unit = EXCLUDED.price_per_unit,
			created_at = EXCLUDED.created_at
	`

	// Iterate over each product in req.Data and execute the insert query
	for _, productTransaction := range req.Data {
		_, err := tx.NamedExecContext(ctx, query, productTransaction)
		if err != nil {
			log.Error().Err(err).Any("productTransaction", productTransaction).Msg("failed to insert product transaction")
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

func (r *productRepo) CreateProductTransaction(ctx context.Context, req *entity.CreateProductTransactionReq) error {
	// Implement the logic to create a product transaction

	id := uuid.New()                                     // Generate a new UUID
	idNoDash := strings.ReplaceAll(id.String(), "-", "") // Remove dashes

	query := `
		INSERT INTO product_transactions (
			id,
			member_id,
			product_id,
			product_grammage_id,
			source,
			qty,
			price_per_unit,
			created_at,
			is_training_data
		) VALUES (
			:id,
			:member_id,
			:product_id,
			:product_grammage_id,
			:source,
			:qty,
			:price_per_unit,
			:created_at,
			FALSE
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":                  idNoDash,
		"member_id":           req.MemberId,
		"product_id":          req.ProductId,
		"product_grammage_id": req.ProductGrammageId,
		"source":              req.Source,
		"qty":                 req.Qty,
		"price_per_unit":      req.PricePerUnit,
		"created_at":          req.CreatedAt,
	})
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("failed to create product transaction")
		return err
	}

	return nil
}

func (r *productRepo) GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.Product
	}

	var (
		data = make([]dao, 0)
		res  = new(entity.GetProductsResp)
	)
	res.Items = make([]entity.Product, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			p.id,
			p.name,
			p.category,
			p.level
		FROM products p
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Paginate, (req.Page-1)*req.Paginate)
	if err != nil {
		log.Error().Err(err).Msg("failed to get products")
		return nil, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.Product)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}

func (r *productRepo) GetProductGrammages(ctx context.Context, req *entity.GetProductGrammagesReq) (*entity.GetProductGrammagesResp, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductGrammage
	}
	var (
		data = make([]dao, 0)
		res  = new(entity.GetProductGrammagesResp)
	)
	res.Items = make([]entity.ProductGrammage, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			pg.id,
			pg.name,
			pg.point,
			pg.price
		FROM product_grammages pg
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Paginate, (req.Page-1)*req.Paginate)
	if err != nil {
		log.Error().Err(err).Msg("failed to get product grammages")
		return nil, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.ProductGrammage)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}

func (r *productRepo) GetProductTransactions(ctx context.Context, req *entity.GetProductTransactionsReq) (*entity.GetProductTransactionsResp, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductTransaction
	}
	var (
		data = make([]dao, 0)
		res  = new(entity.GetProductTransactionsResp)
	)
	res.Items = make([]entity.ProductTransaction, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			pt.id,
			pt.member_id,
			pt.product_id,
			pt.product_grammage_id,
			pt.source,
			pt.qty,
			pt.price_per_unit,
			pt.created_at,
			pt.is_training_data
		FROM product_transactions pt
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.Paginate, (req.Page-1)*req.Paginate)
	if err != nil {
		log.Error().Err(err).Msg("failed to get product transactions")
		return nil, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.ProductTransaction)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
