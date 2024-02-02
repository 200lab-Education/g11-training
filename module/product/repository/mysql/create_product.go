package productmysql

import (
	"context"
	"my-app/module/product/domain"
)

func (repo MysqlRepository) CreateProduct(ctx context.Context, prod *productdomain.ProductCreationDTO) error {
	if err := repo.db.Table(prod.TableName()).Create(&prod).Error; err != nil {
		return err
	}

	return nil
}
