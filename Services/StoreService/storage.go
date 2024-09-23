package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func InitPostgresDB() {
	env := os.Getenv("ENV")

	var dsn string
	if env == "test" {
		dsn = "host=localhost port=5432 user=postgres dbname=postgres password=password sslmode=disable"
	} else {
		var (
			host     = os.Getenv("DB_HOST")
			port     = os.Getenv("DB_PORT")
			dbUser   = os.Getenv("DB_USER")
			dbName   = os.Getenv("DB_NAME")
			password = os.Getenv("DB_PASSWORD")
		)
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host,
			port,
			dbUser,
			dbName,
			password,
		)
	}

	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(Store{}, Product{})
	Prepare()
}

func CloseConn() {
	db.Close()
}

func Prepare() {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	s := Store{}
	res := tx.First(&s)

	if res.RowsAffected != 0 {
		return
	}

	storeOne := Store{
		Name: "Store One",
		Products: []Product{
			{
				Name:     "Sugar",
				Price:    50,
				Quantity: 200,
			},
			{
				Name:     "Milk",
				Price:    100,
				Quantity: 300,
			},
			{
				Name:     "Beef Steak",
				Price:    2000,
				Quantity: 50,
			},
		},
	}
	tx.Create(&storeOne)

	storeTwo := Store{
		Name: "Store Two",
		Products: []Product{
			{
				Name:     "Sugar",
				Price:    50,
				Quantity: 250,
			},
			{
				Name:     "Milk",
				Price:    100,
				Quantity: 400,
			},
			{
				Name:     "Potato pack.",
				Price:    100,
				Quantity: 2000,
			},
		},
	}
	tx.Create(&storeTwo)

	storeThree := Store{
		Name: "Store Three",
		Products: []Product{
			{
				Name:     "Toilet Paper",
				Price:    250,
				Quantity: 100,
			},
			{
				Name:     "Milk",
				Price:    100,
				Quantity: 200,
			},
			{
				Name:     "Potato pack.",
				Price:    100,
				Quantity: 1500,
			},
		},
	}
	tx.Create(&storeThree)
	tx.Commit()
}

type Store struct {
	gorm.Model
	Name     string
	Products []Product `gorm:"foreignKey:StoreID"`
}

type Product struct {
	gorm.Model
	Name     string
	Price    uint
	Quantity uint
	StoreID  uint
}

func GetProduct(storeId, id uint) (*Product, error) {
	var product Product
	res := db.First(&product, "id = ? AND store_id = ?", id, storeId)
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("product with id %d not found", id)
	}
	return &product, nil
}

func DeleteProduct(storeId, id uint) error {
	var product Product
	res := db.First(&product, "id = ? AND store_id = ?", id, storeId)
	if res.RowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	product.Quantity--

	db.Save(&product)

	return nil
}
