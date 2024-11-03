package entity

import "mime/multipart"

type ImportProductsReq struct {
	File *multipart.FileHeader `form:"file" validate:"required"`

	Data []*Product
}

type Product struct {
	ProductId       int    `db:"id" csv:"productID" validate:"numeric" json:"id"`
	ProductName     string `db:"name" csv:"ProductName" validate:"required" json:"name"`
	ProductCategory string `db:"category" csv:"ProductCategory" validate:"required" json:"category"`
	ProductLevel    string `db:"level" csv:"ProductLevel" validate:"required" json:"level"`
}

type ImportProductGrammageReq struct {
	File *multipart.FileHeader `form:"file" validate:"required"`

	Data []*ProductGrammage
}

type ProductGrammage struct {
	Id    int      `db:"id" csv:"prodgramID" validate:"numeric" json:"id"`
	Name  *string  `db:"name" csv:"GrammageName,omitempty" json:"name"`
	Point *int64   `db:"point" csv:"Point,omitempty" json:"point"`
	Price *float64 `db:"price" csv:"Price,omitempty" json:"price"`
}

type ImportProductTransactionsReq struct {
	File *multipart.FileHeader `form:"file" validate:"required"`

	Data []*ProductTransaction
}

type CreateProductTransactionReq struct {
	MemberId          string   `db:"member_id" csv:"MemberID" json:"member_id" validate:"exist=members.id"`
	ProductId         int64    `db:"product_id" csv:"FK_PRODUCT_ID" json:"product_id" validate:"required"`
	ProductGrammageId int64    `db:"product_grammage_id" csv:"FK_PROD_GRAM_ID" json:"product_grammage_id" validate:"required"`
	Source            string   `db:"source" csv:"Source" json:"source" validate:"required"`
	Qty               int      `db:"qty" csv:"Qty" json:"qty" validate:"required"`
	PricePerUnit      *float64 `db:"price_per_unit" csv:"PricePerUnit,omitempty" json:"price_per_unit" validate:"required"`
	CreatedAt         string   `db:"created_at" csv:"TransactionDatetime" json:"created_at" validate:"required,datetime=2006-01-02 15:04:05 MST"`
}

type ProductTransaction struct {
	Id                string   `db:"id" csv:"TransactionID" json:"id"`
	MemberId          string   `db:"member_id" csv:"MemberID" json:"member_id"`
	ProductId         int64    `db:"product_id" csv:"FK_PRODUCT_ID" json:"product_id"`
	ProductGrammageId int64    `db:"product_grammage_id" csv:"FK_PROD_GRAM_ID" json:"product_grammage_id"`
	Source            string   `db:"source" csv:"Source" json:"source"`
	Qty               int      `db:"qty" csv:"Qty" json:"qty"`
	PricePerUnit      *float64 `db:"price_per_unit" csv:"PricePerUnit,omitempty" json:"price_per_unit"`
	CreatedAt         string   `db:"created_at" csv:"TransactionDatetime" json:"created_at" validate:"required,datetime=2006-01-02 15:04:05 MST"`
	IsTrainingData    bool     `db:"is_training_data" json:"is_training_data"`
}
