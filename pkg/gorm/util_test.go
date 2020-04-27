package gorm

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/require"
)

var connection *gorm.DB
var ctx = context.TODO()

func init() {
	var err error

	connection, err = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		panic(err)
	}

	connection = connection.AutoMigrate(&item{})
}

type item struct {
	ID   uint `gorm:"type:int AUTO_INCREMENT;PRIMARY_KEY"`
	Name string
}

func TestGorm(t *testing.T) {
	tx := connection.Begin()

	defer func() {
		tx.Rollback()
		connection.Close()
	}()

	myItem := &item{}

	t.Run("should create a record", func(t *testing.T) {
		err := Create(ctx, tx, myItem)
		require.NoError(t, err)
		require.NotEmpty(t, myItem.ID)
	})

	t.Run("should save a record", func(t *testing.T) {
		myItem.Name = "Jane"
		err := Save(ctx, tx, myItem)
		require.NoError(t, err)
	})

	t.Run("should update a record", func(t *testing.T) {
		err := Update(ctx, tx, myItem, map[string]interface{}{
			"name": "John",
		})
		require.NoError(t, err)
	})

	t.Run("should find a record", func(t *testing.T) {
		otherItem := &item{}
		err := First(ctx, tx.New().Where("id = ?", myItem.ID), otherItem)
		require.NoError(t, err)
		require.NotEmpty(t, otherItem.ID)
	})

	t.Run("should pull all the records", func(t *testing.T) {
		items := make([]*item, 0)
		err := Find(ctx, tx, &items)
		require.NoError(t, err)
		require.Equal(t, 1, len(items))
	})

	t.Run("should delte a record", func(t *testing.T) {
		err := Delete(ctx, tx, myItem)
		require.NoError(t, err)
	})
}
