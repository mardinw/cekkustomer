package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
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

func (customer *ImportCustomerXls) GetCustomerByName(db *sql.DB, filePath, agenciesName, firstName string) ([]dtos.DataPreview, error) {
	query := `
	SELECT distinct card_number,
	first_name,
	collector,
	home_address_3 address_3,
	home_address_4 address_4,
	home_zip_code zipcode
	FROM customer 
	WHERE files = $1 AND agencies = $2 AND first_name like $3
	`

	firstNameUpper := "%" + strings.ToUpper(firstName) + "%"

	args := []interface{}{
		filePath,
		agenciesName,
		firstNameUpper,
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
		if err := rows.Scan(
			&each.CardNumber,
			&each.FirstName,
			&each.Collector,
			&each.Address3,
			&each.Address4,
			&each.ZipCode,
		); err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (customer *ImportCustomerXls) GetAllConcatByName(db *sql.DB, tableName, agenciesName, firstName, filePath string) ([]dtos.CheckDPT, error) {

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
	WHERE t1.files = $1 AND t1.agencies = $2 AND t2.nama like $3`,
		pq.QuoteIdentifier(tableName))

	firstNameUpper := "%" + strings.ToUpper(firstName) + "%"

	args := []interface{}{
		filePath,
		agenciesName,
		firstNameUpper,
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

func (customer *ImportCustomerXls) GetAllConcat(db *sql.DB, tableName, agenciesName, filePath string) ([]dtos.CheckDPT, error) {

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
	nik,
	first_name,
	collector,
	agencies,
	home_address_3,
	home_address_4,
	home_zip_code,
	concat_customer,
	files,
	created)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	args := []interface{}{
		customer.CardNumber,
		customer.NIK,
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

func (customer *ImportCustomerXls) GetCustomerNIK(db *sql.DB, filePath, agenciesName string) ([]dtos.DataPreviewNIK, error) {
	query := `
	SELECT distinct card_number,
	nik,
	first_name,
	collector,
	home_address_3 address_3,
	home_address_4 address_4,
	home_zip_code zipcode,
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

	var result []dtos.DataPreviewNIK
	for rows.Next() {
		var each = dtos.DataPreviewNIK{}
		if err := rows.Scan(
			&each.CardNumber,
			&each.NIK,
			&each.FirstName,
			&each.Collector,
			&each.Address3,
			&each.Address4,
			&each.ZipCode,
		); err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (customer *ImportCustomerXls) GetCustomerNIKByName(db *sql.DB, agenciesName, firstName, filePath string) ([]dtos.DataPreviewNIK, error) {
	query := `
	SELECT distinct card_number,
	nik,
	first_name,
	collector,
	home_address_3 address_3,
	home_address_4 address_4,
	home_zip_code zipcode,
	from customer
	WHERE files = $1 AND agencies = $2 AND first_name like $3
	`

	firstNameUpper := "%" + strings.ToUpper(firstName) + "%"

	args := []interface{}{
		filePath,
		agenciesName,
		firstNameUpper,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []dtos.DataPreviewNIK
	for rows.Next() {
		var each = dtos.DataPreviewNIK{}
		if err := rows.Scan(
			&each.CardNumber,
			&each.NIK,
			&each.FirstName,
			&each.Collector,
			&each.Address3,
			&each.Address4,
			&each.ZipCode,
		); err != nil {
			log.Println("record not found")
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}
