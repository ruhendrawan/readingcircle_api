package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "localhost"
	database = "readingcircle_dev"
	user     = "root"
	password = ""
)

func main() {
	r := gin.Default()

	r.GET("/raw_quiz_results", rawQuizResults)
	r.GET("/raw_reading_activities", rawReadingActivities)
	// r.GET("/xls_quiz_results", xlsQuizResults)

	err := r.Run()
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func getDBConnection() (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database)
	return sql.Open("mysql", connString)
}

func getQuizResults(grp string) (header []string, results [][]string, err error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	query := `
	SELECT * FROM submitted_answers s
	INNER JOIN questions q ON s.idquestions = q.idquestions
	INNER JOIN document d ON d.docid = q.docid
	WHERE s.grp = ?
	ORDER BY s.time DESC;
	`
	rows, err := db.Query(query, grp)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Get the header dynamically from the column names
	header, err = rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		cols, err := rows.ColumnTypes()
		if err != nil {
			return nil, nil, err
		}

		data := make([]string, len(cols))

		resultRow := make([]string, len(cols))
		for i, v := range data {
			resultRow[i] = fmt.Sprintf("x%s", v)

		}

		results = append(results, resultRow)
	}

	return header, results, nil
}

func getReadingActivities(grp string) (header []string, results [][]string, err error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	query := `
	SELECT * FROM tracking
	WHERE grp = ?;
	`
	rows, err := db.Query(query, grp)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Get the header dynamically from the column names
	header, err = rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		cols, err := rows.ColumnTypes()
		if err != nil {
			return nil, nil, err
		}

		data := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))

		for i := range data {
			ptrs[i] = &data[i]
		}

		err = rows.Scan(ptrs...)
		if err != nil {
			return nil, nil, err
		}

		resultRow := make([]string, len(cols))
		for i, v := range data {
			resultRow[i] = fmt.Sprintf("%s", v)
		}

		results = append(results, resultRow)
	}

	return header, results, nil
}

func toCSV(header []string, results [][]string) (string, error) {
	var csvData [][]string
	csvData = append(csvData, header)
	csvData = append(csvData, results...)

	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)
	err := writer.WriteAll(csvData)
	if err != nil {
		return "", err
	}

	return csvBuilder.String(), nil
}

func rawQuizResults(c *gin.Context) {
	grp := c.Query("grp")
	if grp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a 'grp' parameter."})
		return
	}

	header, results, err := getQuizResults(grp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	csvOutput, err := toCSV(header, results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment;filename=quiz_results.csv")
	c.Data(http.StatusOK, "text/csv", []byte(csvOutput))
}

func rawReadingActivities(c *gin.Context) {
	grp := c.Query("grp")
	if grp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a 'grp' parameter."})
		return
	}

	header, results, err := getReadingActivities(grp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	csvOutput, err := toCSV(header, results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment;filename=reading_activities.csv")
	c.Data(http.StatusOK, "text/csv", []byte(csvOutput))
}

// func xlsQuizResults(c *gin.Context) {
// 	grp := c.Query("grp")
// 	if grp == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a 'grp' parameter."})
// 		return
// 	}

// 	header, results, err := getQuizResults(grp)

// 	data := make([]map[string]interface{}, len(results))
// 	for i, row := range results {
// 		data[i] = make(map[string]interface{})
// 		for j, value := range row {
// 			data[i][header[j]] = value
// 		}
// 	}

// 	userDocnoCorrect, userDocsrcCorrect, docnoQuestionCounts, docsrcQuestionCounts := processQuizResults(data)

// 	outputFilePath := "ISD_rc_results.xlsx"
// 	err := writeQuizResultsToExcel(userDocnoCorrect, userDocsrcCorrect, docnoQuestionCounts, docsrcQuestionCounts, outputFilePath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file."})
// 		return
// 	}

// 	c.Header("Content-Description", "File Transfer")
// 	c.Header("Content-Disposition", "attachment; filename=ISD_rc_results.xlsx")
// 	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	c.File(outputFilePath)
// }
