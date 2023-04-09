package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var banner = `
-------------------------------------------
-       Denko Inventory Management        -
-------------------------------------------`

type Record struct {
	Item   string `gorm:"primary key,unique"`
	Amount int
}

func addRecord(ctx *cli.Context) error {
	item := ctx.Args().Get(0)
	amount, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	record := Record{
		Item:   item,
		Amount: amount,
	}

	db.Create(&record)

	return nil
}

func getRecords(ctx *cli.Context) error {

	var records []Record
	db.Find(&records)

	fmt.Println(banner)
	for _, i := range records {
		fmt.Printf("- %-20s| %-4d          -\n", i.Item, i.Amount)
	}

	fmt.Println()

	return nil
}

func updateRecord(ctx *cli.Context) error {
	newItem := ctx.Args().Get(0)
	newAmount, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		return err
	}

	newRecord := Record{
		Item:   newItem,
		Amount: newAmount,
	}

	db.Model(&newRecord).Where("item = ?", newRecord.Item).Update("amount", newRecord.Amount)

	return nil
}

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("inventory.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Record{})
	if err != nil {
		fmt.Println("error migrating DB")
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "create",
				Usage:  "Adds an item to track",
				Action: addRecord,
			},
			{
				Name:   "ls",
				Usage:  "lists the inventory",
				Action: getRecords,
			},
			{
				Name:   "update",
				Usage:  "updates inventory record",
				Action: updateRecord,
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		panic(err)
	}
}
