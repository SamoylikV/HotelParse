package xlsxutils

import (
	"fmt"
	"testing"
)

func TestExportToExcel_LargeData(t *testing.T) {
	info := generateLargeTestData(100000)

	err := ExportToExcel(info)
	if err != nil {
		t.Fatalf("ExportToExcel failed: %v", err)
	}
	t.Log("ExportToExcel passed with large dataset")
}

func TestUniqueMaps(t *testing.T) {
	input := []map[string]string{
		{"A": "value1", "B": "value2"},
		{"A": "value1", "B": "value2"},
		{"A": "value3", "B": "value4"},
	}

	expected := 2
	output := UniqueMaps(input)
	if len(output) != expected {
		t.Errorf("expected %d unique maps, got %d", expected, len(output))
	}
}

func TestFormatDateToString(t *testing.T) {
	testCases := []struct {
		input      string
		expected   string
		shouldFail bool
	}{
		{"2023-10-15T12:45:00.000-0700", "2023-10-15 12:45:00", false},
		{"2023-10-15 12:45:00", "2023-10-15 12:45:00", false},
		{"15/10/2023", "2023-10-15 00:00:00", false},
		{"invalid_date", "", true},
	}

	for _, tc := range testCases {
		result, err := formatDateToString(tc.input)
		if tc.shouldFail {
			if err == nil {
				t.Errorf("expected failure for input: %s", tc.input)
			}
		} else {
			if err != nil || result != tc.expected {
				t.Errorf("expected %s, got %s (input: %s)", tc.expected, result, tc.input)
			}
		}
	}
}

func generateLargeTestData(count int) []map[string]string {
	var data []map[string]string
	for i := 0; i < count; i++ {
		data = append(data, map[string]string{
			"A": fmt.Sprintf("Hotel Name %d", i),
			"B": fmt.Sprintf("Type %d", i),
			"C": fmt.Sprintf("email%d@example.com", i),
			"D": fmt.Sprintf("+123456789%d", i),
			"E": fmt.Sprintf("City %d", i),
			"F": fmt.Sprintf("https://example%d.com", i),
			"G": fmt.Sprintf("Region Name %d", i),
			"H": fmt.Sprintf("Region Code %d", i),
			"I": fmt.Sprintf("License %d", i),
			"J": fmt.Sprintf("2023-10-15 12:45:00"),
			"K": fmt.Sprintf("2024-10-15 12:45:00"),
			"L": fmt.Sprintf("Stars Category %d", i),
			"M": fmt.Sprintf("Org %d", i),
			"N": fmt.Sprintf("Rooms %d", i),
		})
	}
	return data
}
