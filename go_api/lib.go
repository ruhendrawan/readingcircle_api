package main

import (
	"fmt"
	"sort"
	"strconv"

	excelize "github.com/xuri/excelize/v2"
)

type Counts struct {
	UserDocNoCorrect     map[string]map[string]int
	UserDocSrcCorrect    map[string]map[string]int
	DocNoQuestionCounts  map[string]int
	DocSrcQuestionCounts map[string]int
}

func processQuizResults(data []map[string]string) Counts {
	userDocNoCorrect := make(map[string]map[string]int)
	userDocSrcCorrect := make(map[string]map[string]int)
	docNoQuestionCounts := make(map[string]int)
	docSrcQuestionCounts := make(map[string]int)

	for _, row := range data {
		usr := row["usr"]
		docNo := row["docno"]
		docSrc := row["docsrc"]

		correct, _ := strconv.Atoi(row["correct"])

		if _, ok := userDocNoCorrect[usr]; !ok {
			userDocNoCorrect[usr] = make(map[string]int)
		}

		if _, ok := userDocSrcCorrect[usr]; !ok {
			userDocSrcCorrect[usr] = make(map[string]int)
		}

		userDocNoCorrect[usr][docNo] += correct
		userDocSrcCorrect[usr][docSrc] += correct

		docNoQuestionCounts[docNo] = intMax(docNoQuestionCounts[docNo], userDocNoCorrect[usr][docNo])
		docSrcQuestionCounts[docSrc] = intMax(docSrcQuestionCounts[docSrc], userDocSrcCorrect[usr][docSrc])
	}

	return Counts{
		UserDocNoCorrect:     userDocNoCorrect,
		UserDocSrcCorrect:    userDocSrcCorrect,
		DocNoQuestionCounts:  docNoQuestionCounts,
		DocSrcQuestionCounts: docSrcQuestionCounts,
	}
}

func writeQuizResultsToExcel(counts Counts, outputFilePath string) error {
	f := excelize.NewFile()

	users := sortedKeys(counts.UserDocNoCorrect)
	docNos := sortedKeys(counts.DocNoQuestionCounts)
	docSrcs := sortedKeys(counts.DocSrcQuestionCounts)

	sheet1 := "User_DocNo"
	f.NewSheet(sheet1)
	// f.SetActiveSheet(sheet1)

	header := append([]string{"User"}, docNos...)
	writeRow(f, sheet1, 1, header)
	writeRow(f, sheet1, 2, append([]string{"Max"}, intMapValues(counts.DocNoQuestionCounts, docNos)...))

	for i, user := range users {
		userRow := append([]string{user}, intMapValues(counts.UserDocNoCorrect[user], docNos)...)
		writeRow(f, sheet1, i+3, userRow)
	}

	sheet2 := "User_DocSrc"
	f.NewSheet(sheet2)

	header = append([]string{"User"}, docSrcs...)
	writeRow(f, sheet2, 1, header)
	writeRow(f, sheet2, 2, append([]string{"Max"}, intMapValues(counts.DocSrcQuestionCounts, docSrcs)...))

	for i, user := range users {
		userRow := append([]string{user}, intMapValues(counts.UserDocSrcCorrect[user], docSrcs)...)
		writeRow(f, sheet2, i+3, userRow)

		// Add formula to compute the sum of columns B, C, and D
		cell1, _ := excelize.CoordinatesToCellName(len(header)+1, i+3)
		cell2, _ := excelize.CoordinatesToCellName(len(header), i+3)
		f.SetCellFormula(sheet2, cell1, fmt.Sprintf("=SUM(B%d:%s%d)", i+3, cell2, i+3))
		// f.SetCellStyle(sheet2, cell, cell, grayFillStyle(f))

		// Add formula to compute column E divided by $E$2
		cell3, _ := excelize.CoordinatesToCellName(len(header)+2, i+3)
		cell4, _ := excelize.CoordinatesToCellName(len(header)+1, i+3)
		f.SetCellFormula(sheet2, cell3, fmt.Sprintf("=ROUND(%s%d/%s$2*100, 2)", cell4, i+3))
		// f.SetCellStyle(sheet2, cell2, cell2, boldGrayFillStyle(f))
	}

	return f.SaveAs(outputFilePath)
}

// func grayFillStyle(f *excelize.File) int {
// 	style, _ := f.NewStyle(`{"fill":{"type":"pattern","patternType":"solid","fgColor":"D3D3D3"}}`)
// 	return style
// }

// func boldGrayFillStyle(f *excelize.File) int {
// 	style, _ := f.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","patternType":"solid","fgColor":"D3D3D3"}}`)
// 	return style
// }

func writeRow(f *excelize.File, sheet string, rowIndex int, values []string) {
	for colIndex, value := range values {
		cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
		f.SetCellValue(sheet, cell, value)
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func intMapValues(m map[string]int, keys []string) []string {
	values := make([]string, len(keys))
	for i, key := range keys {
		values[i] = strconv.Itoa(m[key])
	}
	return values
}

func intMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
