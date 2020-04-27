package main

import (
	"fmt"
	"log"
	"os"

	"generate/pkg"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lapix-com-co/dataloader/pkg/pagination"
	"golang.org/x/net/context"
)

var err error
var connection *gorm.DB

func init() {
	connection, err = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))
	check(err)
	connection.AutoMigrate(&pkg.Pet{})
}

func main() {
	tx := connection.Begin()
	defer func() {
		tx.Rollback()
		connection.Close()
	}()

	loader := pkg.NewPetMySQLDataLoader(tx)
	ctx := context.TODO()

	// Creates an item.
	myPet := &pkg.Pet{Name: "Bunny"}
	err = loader.Create(ctx, myPet)
	check(err)

	// Updates an item.
	myPet.Name = "Lucas"
	err = loader.Update(ctx, myPet)
	check(err)

	// Find an item.
	petsFromDB, err := loader.Find(ctx, []uint{myPet.ID})
	check(err)
	if len(petsFromDB) == 0 {
		log.Fatal("Find() method did not pull any data")
	}

	// Pull all the items.
	page, err := loader.All(ctx, pagination.Input{First: u(1)})
	check(err)
	if len(page.Nodes) == 0 {
		log.Fatal("All() method did not pull any data")
	}

	if page.Nodes[0].ID != petsFromDB[0].ID {
		log.Fatal("All() and Find() pull different items")
	}

	// Deletes the item.
	if err = loader.Delete(ctx, petsFromDB[0]); err != nil {
		log.Fatal(err)
	}
}

func u(i int) *uint16 {
	var x = uint16(i)
	return &x
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
