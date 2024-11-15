package parse

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func Id(number int) (string, error) {
	url1 := "https://tor.knd.gov.ru/ext/search/simpleRegistries"
	body1 := map[string]interface{}{
		"search": map[string]interface{}{
			"search": []map[string]interface{}{
				{
					"field":    "registryType.id",
					"operator": "eq",
					"value":    "63ef2fc7a445e900072d7e10",
				},
				{
					"field":    "status",
					"operator": "in",
					"value":    []string{"active"},
				},
				{
					"field":    "number",
					"operator": "eq",
					"value":    number,
				},
			},
		},
		"prj": "simpleRegistriesList",
	}
	jsonData, _ := json.Marshal(body1)
	req, _ := http.NewRequest("POST", url1, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if res.Body != nil {
			err := res.Body.Close()
			if err != nil {
				return
			}
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", err
	}

	var response Response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return "", err
	}

	if len(response.Content) == 0 {
		return "None", nil
	}

	return response.Content[0].ID, nil
}

func Info(resultID string) (map[string]string, error) {
	url2 := "https://tor.knd.gov.ru/ext/search/simpleRegistryItems"
	body2 := map[string]interface{}{
		"search": map[string]interface{}{
			"search": []map[string]interface{}{
				{
					"field":    "registryId",
					"operator": "eq",
					"value":    resultID,
				},
				{
					"field":    "status",
					"operator": "neq",
					"value":    "draft",
				},
			},
		},
		"prj":  "externalSimpleRegistryItems",
		"sort": "dateLastModification,ASC",
	}
	jsonData2, _ := json.Marshal(body2)
	req2, _ := http.NewRequest("POST", url2, bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res2, _ := client.Do(req2)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res2.Body)

	var response2 Response2
	bodyBytes2, _ := io.ReadAll(res2.Body)
	err := json.Unmarshal(bodyBytes2, &response2)
	if err != nil {
		return nil, err
	}

	var roomCount int
	if len(response2.Content) == 0 {
		return nil, err
	}
	object := response2.Content[0].Objects[0]
	xsdData := object.XsdData
	classInfo := xsdData.ClassificationInfo
	for _, room := range xsdData.InformationRooms {
		roomNum, _ := strconv.Atoi(room.InformationRoomsBlock.NumberRooms)
		roomCount += roomNum
	}
	//return map[string]string{
	//	"rooms":                      strconv.Itoa(roomCount),
	//	"name":                       object.HotelName,
	//	"view":                       xsdData.View,
	//	"email":                      xsdData.Email,
	//	"phone":                      xsdData.Phone,
	//	"city":                       xsdData.City,
	//	"site":                       xsdData.SiteUrl,
	//	"region name":                object.Region.Name,
	//	"region code":                object.Region.Code,
	//	"license number":             classInfo.Order.LicenseNumber,
	//	"license date issued":        classInfo.Order.DateIssued,
	//	"license date end":           classInfo.Order.DateEnd,
	//	"category stars":             classInfo.CategoryStars,
	//	"accredit organization name": classInfo.InfoAccredOrganization.AccredOrganizationShortName,
	//}
	return map[string]string{
		"A": object.HotelName,
		"B": xsdData.View,
		"C": xsdData.Email,
		"D": xsdData.Phone,
		"E": xsdData.City,
		"F": xsdData.SiteUrl,
		"G": object.Region.Name,
		"H": object.Region.Code,
		"I": classInfo.Order.LicenseNumber,
		"J": classInfo.Order.DateIssued,
		"K": classInfo.Order.DateEnd,
		"L": classInfo.CategoryStars,
		"M": classInfo.InfoAccredOrganization.AccredOrganizationShortName,
		"N": strconv.Itoa(roomCount),
		"O": strconv.Itoa(roomCount),
	}, nil
}
