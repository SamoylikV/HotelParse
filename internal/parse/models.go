package parse

type Content struct {
	ID string `json:"_id"`
}

type Response struct {
	Content []Content `json:"content"`
}

type ClassificationInformation struct {
	Order                  Order              `json:"order"`
	CategoryStars          string             `json:"categoryStars"`
	InfoAccredOrganization AccredOrganization `json:"InfoAccredOrganization"`
}

type Order struct {
	LicenseNumber string `json:"licenseNumber"`
	DateIssued    string `json:"licenseDateIssued"`
	DateEnd       string `json:"dateEnd"`
}

type AccredOrganization struct {
	AccredOrganizationShortName string `json:"accredOrganizationShortName"`
}

type Address struct {
	FullAddress string `json:"fullAddress"`
}

type Region struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type XsdData struct {
	View               string                    `json:"View"`
	Email              string                    `json:"Email"`
	Phone              string                    `json:"Phone"`
	City               string                    `json:"City"`
	SiteUrl            string                    `json:"SiteUrl"`
	ClassificationInfo ClassificationInformation `json:"ClassificationInformation"`
	InformationRooms   []InformationRooms        `json:"InformationRooms"`
}

type InformationRoomsBlock struct {
	NumberRooms string `json:"numberRooms"`
}

type InformationRooms struct {
	InformationRoomsBlock InformationRoomsBlock `json:"InformationRoomsBlock"`
}

type Content2 struct {
	Objects []Objects `json:"objects"`
}

type Objects struct {
	HotelName string  `json:"name"`
	XsdData   XsdData `json:"xsdData"`
	Address   Address `json:"address"`
	Region    Region  `json:"region"`
}

type Response2 struct {
	Content []Content2 `json:"content"`
}
