package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cekkustomer.com/dtos"
	"github.com/lib/pq"
)

func GetAllKec(db *sql.DB) ([]string, error) {
	query := `
	SELECT tablename FROM pg_catalog.pg_tables WHERE tablename like 'dpt_%'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []string
	for rows.Next() {
		var each string
		var err = rows.Scan(
			&each,
		)
		if err != nil {
			log.Println("record not found")
			return nil, err
		}
		result = append(result, each)
	}

	return result, nil
}

func CheckData(db *sql.DB, tableName, concatCustomer string) ([]dtos.CheckDPT, error) {
	query := `
	select 
	COALESCE(nama,''),
	COALESCE(kodepos,''),
	COALESCE(kec,''),
	COALESCE(kel,'')
	from ` + tableName + `
	WHERE concat = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{
		concatCustomer,
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var results []dtos.CheckDPT
	for rows.Next() {
		var each = dtos.CheckDPT{}
		var err = rows.Scan(
			&each.Nama,
			&each.Kodepos,
			&each.Kecamatan,
			&each.Kelurahan,
		)
		if err != nil {
			log.Println("record not found")
			return nil, err
		}

		results = append(results, each)

	}

	return results, nil
}

func (customer *ImportCustomerXls) CompareCustomer(db *sql.DB, filePath, agenciesName, concatCustomer string) (bool, error) {
	query := `
	select card_number from customer
	where files = $1 and agencies = $2 and concat_customer = $3
	`

	args := []interface{}{
		filePath,
		agenciesName,
		concatCustomer,
	}

	var exists bool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println()
	}

	defer rows.Close()
	for rows.Next() {
		var each = dtos.DataPreview{}
		var err = rows.Scan(
			&each.CardNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil
			}
			log.Println(err.Error())
			return false, err
		}
	}

	return exists, nil

	//	rows, err := db.QueryContext(ctx, query, args...)
	//	if err != nil {
	//		log.Println(err.Error())
	//		return false
	//	}
	//	defer rows.Close()
	//
	//	var resultFound bool
	//
	//	for rows.Next() {
	//		resultFound = true
	//	}

}

func (customer *ImportCustomerXls) GetCustomer(db *sql.DB, filePath, agenciesName string) ([]dtos.DataPreview, error) {
	query := `
	SELECT distinct card_number,
	first_name,
	collector,
	home_address_3 address_3,
	home_address_4 address_4,
	home_zip_code zipcode
	FROM customer 
	WHERE files = $1 AND agencies = $2
	`

	args := []interface{}{
		filePath,
		agenciesName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []dtos.DataPreview
	for rows.Next() {
		var each = dtos.DataPreview{}
		var err = rows.Scan(
			&each.CardNumber,
			&each.FirstName,
			&each.Collector,
			&each.Address3,
			&each.Address4,
			&each.ZipCode,
		)
		if err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil

}

func (customer *ImportCustomerXls) GetAll(db *sql.DB, tableName, agenciesName, filePath string) ([]dtos.CheckDPT, error) {

	query := fmt.Sprintf(`
	select distinct on(t1.card_number) card_number,
	t1.first_name first_name,
	t1.collector collector,
	t1.agencies agencies,
	t1.home_address_3 address_3,
	t1.home_address_4 address_4,
	t1.home_zip_code zipcode,
	t2.kodepos kodepos,
	t2.nama nama,
	t2.kel kelurahan,
	t2.kec kecamatan from customer t1 
	JOIN %s t2 ON t1.concat_customer = t2.concat
	WHERE t1.files = $1 AND t1.agencies = $2
	`, pq.QuoteIdentifier(tableName))

	args := []interface{}{
		filePath,
		agenciesName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []dtos.CheckDPT
	for rows.Next() {
		var each = dtos.CheckDPT{}
		var err = rows.Scan(
			&each.CardNumber,
			&each.FirstName,
			&each.Collector,
			&each.Agencies,
			&each.Address3,
			&each.Address4,
			&each.ZipCode,
			&each.Kodepos,
			&each.Nama,
			&each.Kecamatan,
			&each.Kelurahan,
		)
		if err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)

	}

	return result, nil
}

func (customer *ImportCustomerXls) InsertCustomer(db *sql.DB) error {
	query := `
	INSERT INTO customer(
	card_number,
	first_name,
	collector,
	agencies,
	home_address_3,
	home_address_4,
	home_zip_code,
	concat_customer,
	files,
	created)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	args := []interface{}{
		customer.CardNumber,
		customer.FirstName,
		customer.Collector,
		customer.Agencies,
		customer.Address3,
		customer.Address4,
		customer.ZipCode,
		customer.ConcatCustomer,
		customer.Files,
		customer.Created,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (customer *ImportCustomerXls) DeleteCustomer(db *sql.DB, filePath, agenciesName string) error {
	query := `
	DELETE FROM customer WHERE files = $1 AND agencies = $2
	`

	args := []interface{}{
		filePath,
		agenciesName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
