package dto

type Cert struct {
	ID               string  `json:"id"`
	Labname          string  `json:"labname"`
	MedicineName     string  `json:"medicinename"`
	Contry           string  `json:"country"`
	NoOfParticipants string  `json:"noOfParticipants"`
	Category         string  `json:"category"`
	AgeGroup         string  `json:"age"`
	Placebo          int     `json:"placebo"`
	CurrencyType     string  `json:"currency"`
	PremiumAmount    float64 `json:"premiumamount"`
	Signature        
}

type ConversionRequest struct {
	Amount      float64 `json:"amount"`
	CountryCode string  `json:"country_code"`
}

type CurCountry struct {
	Code string
	Name string
}

// Parse the currencyconverter response
type CurrencyResponse struct {
	ConvertedAmount float64 `json:"converted_amount"`
	CurrencyCode    string  `json:"currency_code"`
	CurrencyName    string  `json:"currency_name"`
}

type CurrencyRequest struct {
	Amount      float64 `json:"amount"`
	CountryCode string  `json:"country_code"`
}
