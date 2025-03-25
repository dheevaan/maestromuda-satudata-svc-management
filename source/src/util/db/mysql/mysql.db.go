package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"data-management/src/config"

	_ "github.com/go-sql-driver/mysql"
)

var MySQL_LiveConnection int

type MySQLDbUtil struct {
	dsn string
	db  *sql.DB
}

func NewMySQLDbUtil() (*MySQLDbUtil, error) {
	dsn := os.Getenv(config.ENV_KEY_MYSQL_DSN) // "user:password@tcp(localhost:3306)/dbname"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLDbUtil{dsn: dsn, db: db}, nil
}

// Connect ke MySQL
func (this *MySQLDbUtil) Connect() error {
	err := this.db.Ping()
	if err != nil {
		log.Println("Gagal koneksi ke MySQL:", err)
		return err
	}
	MySQL_LiveConnection += 1
	fmt.Println("Koneksi ke MySQL berhasil!")
	return nil
}

// Disconnect dari MySQL
func (this *MySQLDbUtil) Disconnect() {
	if this.db != nil {
		this.db.Close()
		MySQL_LiveConnection -= 1
	}
}

// Insert Data
func (this *MySQLDbUtil) Insert(table string, data map[string]interface{}) (int64, error) {
	columns := ""
	values := ""
	var args []interface{}

	for key, value := range data {
		columns += key + ","
		values += "?,"
		args = append(args, value)
	}
	columns = columns[:len(columns)-1] // Hapus koma terakhir
	values = values[:len(values)-1]

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columns, values)
	result, err := this.db.Exec(query, args...)
	if err != nil {
		log.Println("Gagal insert:", err)
		return 0, err
	}

	lastInsertID, _ := result.LastInsertId()
	return lastInsertID, nil
}

// Update Data
func (this *MySQLDbUtil) Update(table string, data map[string]interface{}, condition string, condValues ...interface{}) (int64, error) {
	setClause := ""
	var args []interface{}

	for key, value := range data {
		setClause += key + "=?,"
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-1] // Hapus koma terakhir
	args = append(args, condValues...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, setClause, condition)
	result, err := this.db.Exec(query, args...)
	if err != nil {
		log.Println("Gagal update:", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (this *MySQLDbUtil) Update(table string, data map[string]interface{}, condition string, condValues ...interface{}) (int64, error) {
	setClause := ""
	var args []interface{}

	for key, value := range data {
		setClause += key + "=?,"
		args = append(args, value)
	}
	setClause = setClause[:len(setClause)-1] // Hapus koma terakhir
	args = append(args, condValues...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, setClause, condition)
	result, err := this.db.Exec(query, args...)
	if err != nil {
		log.Println("Gagal update:", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// Find One Data
func (this *MySQLDbUtil) FindOne(table, condition string, condValues ...interface{}) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", table, condition)
	row := this.db.QueryRow(query, condValues...)

	columns, _ := this.db.Query(fmt.Sprintf("SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s'", table))
	defer columns.Close()

	columnNames := []string{}
	for columns.Next() {
		var col string
		columns.Scan(&col)
		columnNames = append(columnNames, col)
	}

	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err := row.Scan(valuePtrs...)
	if err != nil {
		log.Println("Gagal mengambil data:", err)
		return nil, err
	}

	result := make(map[string]interface{})
	for i, col := range columnNames {
		result[col] = values[i]
	}

	return result, nil
}

// Find Many Data
func (this *MySQLDbUtil) Find(table, condition string, condValues ...interface{}) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condition)
	rows, err := this.db.Query(query, condValues...)
	if err != nil {
		log.Println("Gagal mengambil data:", err)
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	result := []map[string]interface{}{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)
		rowData := make(map[string]interface{})
		for i, col := range columns {
			rowData[col] = values[i]
		}
		result = append(result, rowData)
	}

	return result, nil
}

// Soft Delete (Update Status)
func (this *MySQLDbUtil) SoftDelete(table, condition string, condValues ...interface{}) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET status='archive' WHERE %s", table, condition)
	result, err := this.db.Exec(query, condValues...)
	if err != nil {
		log.Println("Gagal menghapus data:", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// Hard Delete (Delete Permanen)
func (this *MySQLDbUtil) Delete(table, condition string, condValues ...interface{}) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, condition)
	result, err := this.db.Exec(query, condValues...)
	if err != nil {
		log.Println("Gagal menghapus data:", err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}
