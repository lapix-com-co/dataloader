package gorm

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/lapix-com-co/dataloader/pkg"
)

// Create save a item but does not update the associations.
func Create(ctx context.Context, db *gorm.DB, i interface{}) error {
	return db.Create(i).Error
}

// Save save an item but does not update the associations.
func Save(ctx context.Context, db *gorm.DB, i interface{}) error {
	return db.Set("gorm:association_autoupdate", false).Save(i).Error
}

// Update updates an item but does not update the associations.
func Update(ctx context.Context, db *gorm.DB, model interface{}, i map[string]interface{}) error {
	return db.Set("gorm:association_autoupdate", false).
		Model(model).
		Update(i).
		Error
}

// First returns the first item that match the given query
func First(ctx context.Context, db *gorm.DB, i interface{}) error {
	err := db.First(i).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return pkg.ErrRecordNotFound
	}

	return err
}

// Find returns all the items that match the given query.
func Find(ctx context.Context, db *gorm.DB, i interface{}) error {
	return db.Find(i).Error
}

// Count return the total amount of items that match the given query.
func Count(ctx context.Context, db *gorm.DB) (total uint32, err error) {
	err = db.Count(&total).Error
	return total, err
}

// Delete removes an item.
func Delete(ctx context.Context, db *gorm.DB, i interface{}) error {
	return db.Delete(i).Error
}
