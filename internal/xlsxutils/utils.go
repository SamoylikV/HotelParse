package xlsxutils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"sort"
	"strings"
	"time"
)

func ExportToExcel(info []map[string]string) error {
	headers := []string{
		"Название гостиницы", "Вид", "Электронная почта", "Телефон", "Город",
		"Сайт", "Название региона", "Код региона", "Лицензия", "Дата выдачи лицензии",
		"Дата окончания лицензии", "Категория звезд", "Краткое название аккредитационной организации", "Количество номеров",
	}

	f := excelize.NewFile()

	sheetName := "Hotel info"
	f.SetSheetName("Sheet1", sheetName)
	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	headerRow := make([]interface{}, len(headers))
	for i, header := range headers {
		headerRow[i] = header
	}
	if err := streamWriter.SetRow("A1", headerRow); err != nil {
		return err
	}

	uniqueInfo := UniqueMaps(info)
	for i, dict := range uniqueInfo {
		row := make([]interface{}, len(headers))
		for j := range headers {
			col := string(rune('A' + j))
			row[j] = dict[col]
		}
		rowIdx := i + 2
		if err := streamWriter.SetRow(fmt.Sprintf("A%d", rowIdx), row); err != nil {
			return err
		}
	}

	if err := streamWriter.Flush(); err != nil {
		return err
	}

	if err := f.SaveAs("HotelInfo.xlsx"); err != nil {
		return err
	}

	return nil
}

func formatDateToString(dateStr string) (string, error) {
	layouts := []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02 15:04:05",
		"02/01/2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006-01-02 15:04:05"), nil
		}
	}
	return "", fmt.Errorf("не удалось распознать формат даты: %s", dateStr)
}

func UniqueMaps(infoArray []map[string]string) []map[string]string {
	seen := make(map[string]bool)
	var result []map[string]string
	for _, dict := range infoArray {
		if len(dict) == 0 {
			continue
		}
		var keyParts []string
		var keys []string
		for k := range dict {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			keyParts = append(keyParts, dict[k])
		}
		dictKey := strings.Join(keyParts, "|")
		if seen[dictKey] {
			continue
		}
		seen[dictKey] = true
		result = append(result, dict)
	}
	return result
}
