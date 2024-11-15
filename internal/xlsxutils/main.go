package xlsxutils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

func ExportToExcel(info []map[string]string) error {
	headers := []string{"Название гостиницы", "Вид", "Электронная почта", "Телефон", "Город",
		"Сайт", "Название региона", "Код региона", "Лицензия", "Дата выдачи лицензии", "Дата окончания лицензии", "Категория звезд",
		"Краткое название аккредитационной организации", "Количество номеров", "Номер в реестре",
	}
	f := excelize.NewFile()
	sheetName := "Hotel info"
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()
	styleIdx, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"93fac0"},
			Pattern: 1,
		},
	})
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	if err = f.SetCellStyle(sheetName, "A1", "N1", styleIdx); err != nil {
		return err
	}
	for i := 0; i < len(headers); i++ {
		cell := fmt.Sprintf("%c%d", 'A'+i, 1)
		err = f.SetCellValue(sheetName, cell, headers[i])
		if err != nil {
			return err
		}
	}

	uniqueInfo := uniqueMaps(info)
	for i, dict := range uniqueInfo {
		for key, value := range dict {
			cell := key + strconv.Itoa(i+2)
			var cellValue interface{}
			var err error
			if key == "J" || key == "K" {
				cellValue, err = formatDateToString(value)
			} else {
				cellValue = value
			}
			if err != nil {
				return err
			}
			if err := f.SetCellValue(sheetName, cell, cellValue); err != nil {
				return err
			}
		}
	}
	err = f.AutoFilter(sheetName, "A1:N1", []excelize.AutoFilterOptions{})
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)
	if err = f.SaveAs("HotelInfo.xlsx"); err != nil {
		fmt.Println(err)
	}
	return nil
}

func formatDateToString(dateStr string) (string, error) {
	dateLayout := "2006-01-02T15:04:05.000-0700"
	t, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02 15:04:05"), nil
}

func uniqueMaps(infoArray []map[string]string) []map[string]string {
	seen := make(map[string]struct{})
	var result []map[string]string
	for _, dict := range infoArray {
		if len(dict) == 0 {
			continue
		}
		hasEmptyKey := false
		for key := range dict {
			if key == "" {
				hasEmptyKey = true
				break
			}
		}
		if hasEmptyKey {
			continue
		}
		dictStr := fmt.Sprintf("%v", dict)
		if _, exists := seen[dictStr]; exists {
			continue
		}
		seen[dictStr] = struct{}{}
		result = append(result, dict)
	}
	return result
}
