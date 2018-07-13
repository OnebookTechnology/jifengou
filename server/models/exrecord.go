package models

type ExchangeRecord struct {
	RecordId    int      `json:"record_id"`
	PhoneNumber int      `json:"phone_number"`
	BCodes      string   `json:"b_codes"`
	BCodeArray  []string `json:"b_code_array"`
	PCode       string   `json:"p_code"`
	ExTime      string   `json:"ex_time"`
	PId         int      `json:"p_id"`
	Name        string   `json:"name,omitempty"`
}
