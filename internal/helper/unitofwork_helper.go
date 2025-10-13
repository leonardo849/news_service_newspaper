package helper

import (
	"news_service/internal/unitofwork"

	"gorm.io/gorm"
)

func CreateUnitOfWork(db *gorm.DB) *unitofwork.UnitOfWork {
	return unitofwork.CreateUnitOfWork(
		createUserRepository(db),
		db,
		CreateImageRepository(db),
		createNewsRepository(db),
		CreateBlockRepository(db),
	)
}