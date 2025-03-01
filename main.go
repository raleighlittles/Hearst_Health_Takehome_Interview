package main

import (
	"container/heap"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"log"
)

/// -------------------  Data structures provided below, do not modify this section ------------------- ///

type PriceUpdate struct {
	Retailer string  // the name of the retailer
	SKU      string  // assume retailers share a common SKU
	Price    float64 // always the price per unit
	URL      string  // product detail link, optional
}

func Receive(payload PriceUpdate) {

	// Insert the new price into the database

	insertStmt := "INSERT INTO products (dateAdded, sku, retailer, price, url) VALUES (date('now'), ?, ?, ?, ?)"
	_, err := productsDb.Exec(insertStmt, payload.SKU, payload.Retailer, payload.Price, payload.URL)
	if err != nil {
		log.Fatal("ERROR! Could not insert data into products table!", err)
	}

	if priceHeaps[payload.SKU] == nil {
		priceHeaps[payload.SKU] = &ProductHeap{}
		heap.Init(priceHeaps[payload.SKU])
	}

	heap.Push(priceHeaps[payload.SKU], ProductPrice(payload))
}

type ProductPrice struct {
	Retailer string
	SKU      string
	Price    float64
	URL      string
}

func findPrice(sku string) ProductPrice {

	// The instructions don't specifically say this, but based on the test cases,
	// it seems that we are supposed to remove the product from the heap after we find it.
	// I guess the assumption is that the consumer will purchase the product after finding the price,
	// and the retailer only has 1 quantity of that product in stock.
	//return (*priceHeaps[sku])[0]
	return heap.Pop(priceHeaps[sku]).(ProductPrice)
}

/// -------------------  Implementation of methods required to use `heap` package, below ------------------- ///

type ProductHeap []ProductPrice

func (h ProductHeap) Len() int           { return len(h) }
func (h ProductHeap) Less(i, j int) bool { return h[i].Price < h[j].Price }
func (h ProductHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ProductHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(ProductPrice))
}

func (h *ProductHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func SetupDb(DbName string, prePopulate bool) *sql.DB {

	// Open the database
	db, err := sql.Open("sqlite3", "./"+DbName+".db")
	if err != nil {
		log.Fatal("ERROR! Could not open SQL DB!", err)
	}
	//defer db.Close()

	// Corrected CREATE TABLE statement
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        dateAdded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		sku TEXT NOT NULL,
		retailer TEXT NOT NULL,
		price INTEGER NOT NULL, 
		url TEXT
	);`

	_, err = db.Exec(createProductsTable)
	if err != nil {
		log.Fatal("ERROR! Could not create products table!", err)
	}

	if prePopulate {
		insertDataIntoTable := `
        INSERT INTO products (sku, retailer, price, url) VALUES
        ('CLOCK', 'Walmart', 20, ''),
        ('BED', 'IKEA', 140, ''),
        ('CLOCK', 'Target', 15, ''),
        ('CLOCK', 'Target', 14, ''),
        ('CLOCK', 'Best Buy', 30, ''),
        ('BED', 'Wayfair', 120, ''),
        ('CLOCK', 'Target', 25, ''),
        ('CLOCK', 'Walmart', 27, ''),
        ('CLOCK', 'Costco', 12, ''),
        ('BED', 'IKEA', 100, ''),
        ('CLOCK', 'Costco', 13, '');`

		_, err = db.Exec(insertDataIntoTable)
		if err != nil {
			log.Fatal("ERROR! Could not insert data into products table!", err)
		}
	}

	fmt.Println("Finished populating database table!")

	return db
}

var priceHeaps = make(map[string]*ProductHeap)
var productsDb = SetupDb(dbName, false)
var dbName = "products"

func main() {

	Receive(PriceUpdate{"Walmart", "CLOCK", 20, ""})

	p1 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p1.Retailer, " Low price ", p1.Price)

	Receive(PriceUpdate{"IKEA", "BED", 140, ""})

	p2 := findPrice("BED")
	fmt.Println("Low retailer ", p2.Retailer, " Low price ", p2.Price)

	Receive(PriceUpdate{"Target", "CLOCK", 15, ""})

	p3 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p3.Retailer, " Low price ", p3.Price)

	Receive(PriceUpdate{"Target", "CLOCK", 14, ""})

	p4 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p4.Retailer, " Low price ", p4.Price)

	Receive(PriceUpdate{"Best Buy", "CLOCK", 30, ""})

	p5 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p5.Retailer, " Low price ", p5.Price)

	Receive(PriceUpdate{"Wayfair", "BED", 120, ""})
	p6 := findPrice("BED")
	fmt.Println("Low retailer ", p6.Retailer, " Low price ", p6.Price)

	Receive(PriceUpdate{"Target", "CLOCK", 25, ""})

	p7 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p7.Retailer, " Low price ", p7.Price)

	Receive(PriceUpdate{"Walmart", "CLOCK", 27, ""})

	p8 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p8.Retailer, " Low price ", p8.Price)

	Receive(PriceUpdate{"Costco", "CLOCK", 12, ""})

	p9 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p9.Retailer, " Low price ", p9.Price)

	Receive(PriceUpdate{"IKEA", "BED", 100, ""})

	p10 := findPrice("BED")
	fmt.Println("Low retailer ", p10.Retailer, " Low price ", p10.Price)

	Receive(PriceUpdate{"Costco", "CLOCK", 13, ""})

	p11 := findPrice("CLOCK")
	fmt.Println("Low retailer ", p11.Retailer, " Low price ", p11.Price)
}
