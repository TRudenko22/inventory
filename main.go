package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/TRudenko22/inventory/data"
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

	db.Model(&Record{}).Where("item = ?", newItem).Update("amount", newAmount)

	return nil
}

func removeRecord(ctx *cli.Context) error {
	item := ctx.Args().Get(0)

	db.Where("item = ?", item).Delete(&Record{})

	return nil
}

func getEntries(ctx *cli.Context) error {
	var records []Record
	db.Find(&records)

	fmt.Printf("Total entries tracked %-3d\n", len(records))

	return nil
}

func decreaseAmount(ctx *cli.Context) error {
	var record Record
	item := ctx.Args().Get(0)

	db.Where("item = ?", item).First(&record).Update("amount", record.Amount-1)

	return nil
}

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open(string(data.MustAsset("data/inventory.db"))), &gorm.Config{})
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
				Name:    "list",
				Aliases: []string{"l", "ls"},
				Usage:   "Lists the inventory",
				Action:  getRecords,
			},
			{
				Name:   "update",
				Usage:  "Updates inventory record",
				Action: updateRecord,
			},
			{
				Name:    "remove",
				Aliases: []string{"rm", "rem"},
				Usage:   "Removes a tracked inventory item",
				Action:  removeRecord,
			},
			{
				Name:   "entries",
				Usage:  "Prints the total amount of items tracked",
				Action: getEntries,
			},
			{
				Name:    "decrease",
				Usage:   "Decreases the amount of an item by 1",
				Aliases: []string{"d", "dec"},
				Action:  decreaseAmount,
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		panic(err)
	}
}
