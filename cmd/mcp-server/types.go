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

// InvoiceItemInput represents a line item for invoice creation.
type InvoiceItemInput struct {
	Item     string  `json:"item" jsonschema:"required,Item ID/Name"`
	Quantity float64 `json:"quantity" jsonschema:"required,Quantity"`
	Rate     float64 `json:"rate" jsonschema:"optional,Rate (if not provided, standard item rate is used)"`
	Account  string  `json:"account" jsonschema:"optional,Override income/expense account"`
}

// CreateInvoiceArgs contains parameters for creating a new invoice.
type CreateInvoiceArgs struct {
	Type    string             `json:"type" jsonschema:"required,Sales or Purchase"`
	Party   string             `json:"party" jsonschema:"required,Customer or Supplier ID/Name"`
	Date    string             `json:"date" jsonschema:"required,Invoice date (YYYY-MM-DD)"`
	Items   []InvoiceItemInput `json:"items" jsonschema:"required,Line items"`
	Account string             `json:"account" jsonschema:"required,Main ledger account (e.g., Debtors or Creditors)"`
}

// CreateInvoiceResult is the response returned after creating an invoice.
type CreateInvoiceResult struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
}

// GetBalanceSheetArgs contains parameters for fetching the balance sheet.
type GetBalanceSheetArgs struct {
	AsOfDate string `json:"as_of_date" jsonschema:"required,Report date (YYYY-MM-DD)"`
}

// BalanceSheetItem represents a root type total in the balance sheet.
type BalanceSheetItem struct {
	RootType string  `json:"root_type"`
	Balance  float64 `json:"balance"`
}

// GetBalanceSheetResult is the response returned with balance sheet data.
type GetBalanceSheetResult struct {
	Items []BalanceSheetItem `json:"items"`
}

// ListAccountsArgs contains parameters for listing accounts.
type ListAccountsArgs struct {
	RootType string `json:"root_type" jsonschema:"optional,Filter by Asset, Liability, Equity, Income, Expense"`
}

// AccountItem represents a single account in the COA.
type AccountItem struct {
	Name          string `json:"name"`
	RootType      string `json:"root_type"`
	AccountType   string `json:"account_type"`
	ParentAccount string `json:"parent_account"`
	IsGroup       bool   `json:"is_group"`
}

// ListAccountsResult is the response returned when listing accounts.
type ListAccountsResult struct {
	Accounts []AccountItem `json:"accounts"`
}

// CreateAccountArgs contains parameters for creating a new account.
type CreateAccountArgs struct {
	Name          string `json:"name" jsonschema:"required,Account Name"`
	RootType      string `json:"root_type" jsonschema:"required,Asset, Liability, Equity, Income, or Expense"`
	AccountType   string `json:"account_type" jsonschema:"optional,e.g. Bank, Cash, Receivable, Payable"`
	ParentAccount string `json:"parent_account" jsonschema:"optional,Name of the parent group account"`
	IsGroup       bool   `json:"is_group" jsonschema:"optional,True if this is a folder for other accounts"`
}

// UpdateAccountArgs contains parameters for modifying an account.
type UpdateAccountArgs struct {
	Name    string                 `json:"name" jsonschema:"required,Current name of the account to update"`
	Updates map[string]interface{} `json:"updates" jsonschema:"required,Map of fields to update (e.g. parent_account, accountType)"`
}

// JournalAccountInput represents a single account line in a journal entry.
type JournalAccountInput struct {
	Account string  `json:"account" jsonschema:"required,Account Name"`
	Debit   float64 `json:"debit" jsonschema:"optional,Debit amount"`
	Credit  float64 `json:"credit" jsonschema:"optional,Credit amount"`
}

// CreateJournalEntryArgs contains parameters for creating a journal entry.
type CreateJournalEntryArgs struct {
	EntryType string                `json:"entryType" jsonschema:"required,e.g. Journal Entry, Bank Entry, Cash Entry"`
	Date      string                `json:"date" jsonschema:"required,YYYY-MM-DD"`
	Accounts  []JournalAccountInput `json:"accounts" jsonschema:"required,List of account entries (debits and credits)"`
	Remark    string                `json:"userRemark" jsonschema:"optional,Internal notes"`
}





