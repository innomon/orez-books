package main

// ListPartiesArgs contains the parameters for listing customers or suppliers.
type ListPartiesArgs struct {
	Type  string `json:"type" jsonschema:"optional,Filter by Customer or Supplier (Both, Customer, Supplier)"`
	Limit int    `json:"limit" jsonschema:"optional,Maximum number of results to return (default: 20)"`
}

// PartyItem represents a single party (customer/supplier) in the result list.
type PartyItem struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// ListPartiesResult is the response returned when listing parties.
type ListPartiesResult struct {
	Parties []PartyItem `json:"parties"`
	Count   int         `json:"count"`
}

// ListInvoicesArgs contains parameters for listing invoices.
type ListInvoicesArgs struct {
	Type   string `json:"type" jsonschema:"optional,Sales or Purchase"`
	Status string `json:"status" jsonschema:"optional,Filter by status (e.g. Unpaid, Paid)"`
	Limit  int    `json:"limit" jsonschema:"optional,Maximum number of results to return"`
}

// InvoiceItem represents a single invoice in the result list.
type InvoiceItem struct {
	Name            string  `json:"name"`
	Party           string  `json:"party"`
	Date            string  `json:"date"`
	Total           float64 `json:"total"`
	Outstanding     float64 `json:"outstanding_amount"`
	Status          string  `json:"status"`
}

// ListInvoicesResult is the response returned when listing invoices.
type ListInvoicesResult struct {
	Invoices []InvoiceItem `json:"invoices"`
	Count    int           `json:"count"`
}

// GetInvoiceDetailsArgs contains parameters for fetching full invoice details.
type GetInvoiceDetailsArgs struct {
	Name string `json:"name" jsonschema:"required,The unique name/ID of the invoice"`
}

// GetInvoiceDetailsResult is the response returned with full invoice details.
type GetInvoiceDetailsResult struct {
	Invoice map[string]interface{} `json:"invoice"`
}

// GetAccountBalancesArgs contains parameters for fetching account balances.
type GetAccountBalancesArgs struct {
	AccountType string `json:"account_type" jsonschema:"optional,Filter by account type (e.g. Bank, Cash, Receivable, Payable)"`
}

// AccountBalanceItem represents a single account balance.
type AccountBalanceItem struct {
	Account string  `json:"account"`
	Balance float64 `json:"balance"`
}

// GetAccountBalancesResult is the response returned with account balances.
type GetAccountBalancesResult struct {
	Balances []AccountBalanceItem `json:"balances"`
}

// GetPnLArgs contains parameters for calculating Profit and Loss.
type GetPnLArgs struct {
	StartDate string `json:"start_date" jsonschema:"required,Start date (YYYY-MM-DD)"`
	EndDate   string `json:"end_date" jsonschema:"required,End date (YYYY-MM-DD)"`
}

// PnLResult is the response returned with Profit and Loss data.
type PnLResult struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpense  float64 `json:"total_expense"`
	NetProfit     float64 `json:"net_profit"`
}

// CreatePartyArgs contains parameters for creating a new party.
type CreatePartyArgs struct {
	Name  string `json:"name" jsonschema:"required,Full Name"`
	Role  string `json:"role" jsonschema:"required,Role (Customer, Supplier, or Both)"`
	Email string `json:"email" jsonschema:"optional,Email address"`
	Phone string `json:"phone" jsonschema:"optional,Phone number"`
}

// CreatePartyResult is the response returned after creating a party.
type CreatePartyResult struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
}

// CreateItemArgs contains parameters for creating a new item.
type CreateItemArgs struct {
	Name           string  `json:"name" jsonschema:"required,Item Name"`
	ItemType       string  `json:"item_type" jsonschema:"optional,Product or Service (default: Product)"`
	Rate           float64 `json:"rate" jsonschema:"optional,Standard rate"`
	IncomeAccount  string  `json:"income_account" jsonschema:"required,Sales account for this item"`
	ExpenseAccount string  `json:"expense_account" jsonschema:"required,Purchase account for this item"`
}

// CreateItemResult is the response returned after creating an item.
type CreateItemResult struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
}


