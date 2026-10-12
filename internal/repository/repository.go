package repository

import (
	"github.com/paricorpdev/dashboard-pay2u-be/internal/helpers"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository[T]) FindByKeyValue(db *gorm.DB, entity *T, key string, value any) error {
	return db.Where(key+" = ?", value).Take(entity).Error
}

func (r *Repository[T]) ApplyIncludes(db *gorm.DB, includes []string) *gorm.DB {
	if len(includes) == 0 {
		return db
	}

	for _, include := range includes {
		db = db.Preload(helpers.StringToPascalCase(include))
	}
	return db
}

func (r *Repository[T]) ApplySorting(db *gorm.DB, sort map[string]interface{}) *gorm.DB {
	for column, direction := range sort {
		if dir, ok := direction.(string); ok {
			if dir == "asc" || dir == "desc" {
				db = db.Order(column + " " + dir)
			}
		}
	}
	return db
}
