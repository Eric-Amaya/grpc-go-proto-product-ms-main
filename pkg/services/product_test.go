package services

import (
	"context"
	"testing"

	"grpc-go-proto-product-ms-main/pkg/db"
	pb "grpc-go-proto-product-ms-main/pkg/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.CreateProductRequest{
		Name:        "Test Product",
		Sku:         "test-sku",
		Category:    "test-category",
		Description: "test-description",
		Stock:       10,
		Price:       100.0,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"products\"").
		WithArgs(req.Name, req.Sku, req.Category, req.Description, req.Stock, req.Price).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	resp, err := server.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, int32(201), resp.Status)
	assert.Equal(t, int32(1), resp.Id)
}

func TestFindOne(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.FindOneRequest{
		Id: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "sku", "category", "description", "stock", "price"}).
		AddRow(1, "Test Product", "test-sku", "test-category", "test-description", 10, 100)

	mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."id" = \$1 ORDER BY "products"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	resp, err := server.FindOne(context.Background(), req)

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		assert.Equal(t, int32(200), resp.Status)
		assert.Equal(t, int32(1), resp.Data.Id)
		assert.Equal(t, "Test Product", resp.Data.Name)
		assert.Equal(t, "test-sku", resp.Data.Sku)
		assert.Equal(t, "test-category", resp.Data.Category)
		assert.Equal(t, "test-description", resp.Data.Description)
		assert.Equal(t, int32(10), resp.Data.Stock)
		assert.Equal(t, int32(100), resp.Data.Price)
	} else {
		t.Fatal("Response is nil")
	}
}

func TestUpdateProduct(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.UpdateProductRequest{
		ProductId: 1,
		Product: &pb.Product{
			Id:          1,
			Name:        "Test Product",
			Sku:         "test-sku",
			Category:    "test-category",
			Description: "test-description",
			Stock:       5,
			Price:       200,
		},
	}

	mock.ExpectQuery(`^SELECT \* FROM "products" WHERE "products"."id" = \$1 ORDER BY "products"."id" LIMIT \$2$`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "sku", "category", "description", "stock", "price"}).
			AddRow(1, "Test Product", "test-sku", "test-category", "test-description", 10, 100))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "products" SET "name" = \$1, "sku" = \$2, "category" = \$3, "description" = \$4, "stock" = \$5, "price" = \$6 WHERE "id" = \$7`).
		WithArgs(req.Product.Name, req.Product.Sku, req.Product.Category, req.Product.Description, req.Product.Stock, req.Product.Price, req.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp, err := server.UpdateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, int32(200), resp.Status)
}

func TestFindByCategory(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.FindByCategoryRequest{
		Category: "test-category",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "sku", "category", "description", "stock", "price"}).
		AddRow(1, "Test Product", "test-sku", "test-category", "test-description", 10, 100)

	mock.ExpectQuery(`SELECT \* FROM "products" WHERE category = \$1`).
		WithArgs("test-category").
		WillReturnRows(rows)

	resp, err := server.FindByCategory(context.Background(), req)

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		assert.Equal(t, int32(200), resp.Status)
		assert.Equal(t, 1, len(resp.Products))
		assert.Equal(t, int32(1), resp.Products[0].Id)
		assert.Equal(t, "Test Product", resp.Products[0].Name)
		assert.Equal(t, "test-sku", resp.Products[0].Sku)
		assert.Equal(t, "test-category", resp.Products[0].Category)
		assert.Equal(t, "test-description", resp.Products[0].Description)
		assert.Equal(t, int32(10), resp.Products[0].Stock)
		assert.Equal(t, int32(100), resp.Products[0].Price)
	} else {
		t.Fatal("Response is nil")
	}
}

func TestFindByName(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.FindByNameRequest{
		Name: "Test Product",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "sku", "category", "description", "stock", "price"}).
		AddRow(1, "Test Product", "test-sku", "test-category", "test-description", 10, 100)

	mock.ExpectQuery(`SELECT \* FROM "products" WHERE name = \$1`).
		WithArgs("Test Product").
		WillReturnRows(rows)

	resp, err := server.FindByName(context.Background(), req)

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		assert.Equal(t, int32(200), resp.Status)
		assert.Equal(t, 1, len(resp.Products))
		assert.Equal(t, int32(1), resp.Products[0].Id)
		assert.Equal(t, "Test Product", resp.Products[0].Name)
		assert.Equal(t, "test-sku", resp.Products[0].Sku)
		assert.Equal(t, "test-category", resp.Products[0].Category)
		assert.Equal(t, "test-description", resp.Products[0].Description)
		assert.Equal(t, int32(10), resp.Products[0].Stock)
		assert.Equal(t, int32(100), resp.Products[0].Price)
	} else {
		t.Fatal("Response is nil")
	}
}

func TestFindAll(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.FindAllRequest{}

	rows := sqlmock.NewRows([]string{"id", "name", "sku", "category", "description", "stock", "price"}).
		AddRow(1, "Test Product", "test-sku", "test-category", "test-description", 10, 100)

	mock.ExpectQuery(`SELECT \* FROM "products"`).
		WillReturnRows(rows)

	resp, err := server.FindAll(context.Background(), req)

	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		assert.Equal(t, int32(200), resp.Status)
		assert.Equal(t, 1, len(resp.Products))
		assert.Equal(t, int32(1), resp.Products[0].Id)
		assert.Equal(t, "Test Product", resp.Products[0].Name)
		assert.Equal(t, "test-sku", resp.Products[0].Sku)
		assert.Equal(t, "test-category", resp.Products[0].Category)
		assert.Equal(t, "test-description", resp.Products[0].Description)
		assert.Equal(t, int32(10), resp.Products[0].Stock)
		assert.Equal(t, int32(100), resp.Products[0].Price)
	} else {
		t.Fatal("Response is nil")
	}
}

func TestDeleteProduct(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	server := &Server{
		H: db.Handler{DB: gormDB},
	}

	req := &pb.DeleteProductRequest{
		Id: 1,
	}

	mock.ExpectQuery(`^SELECT \* FROM "products" WHERE "products"."id" = \$1 ORDER BY "products"."id" LIMIT \$2$`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "sku", "category", "description", "stock", "price"}).
			AddRow(1, "Test Product", "test-sku", "test-category", "test-description", 10, 100))

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "products" WHERE "products"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp, err := server.DeleteProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, int32(200), resp.Status)
}
