package models

// Result описывает результат, возвращаемый из Python-скрипта
type Result struct {
	AccountTransactions map[string][]map[string]interface{} `json:"account_transactions"`
	FirstPageText       string                              `json:"first_page_text"`
	StatementType       string                              `json:"statement_type"`
}
