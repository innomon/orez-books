package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"orez-books/internal/database"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

// MCPServer wraps the Orez Books database manager and exposes it as an MCP server.
type MCPServer struct {
	dbManager *database.Manager
	server    *mcp.Server
}

// NewMCPServer creates a new MCPServer that interacts with the given database manager.
func NewMCPServer(dbManager *database.Manager) *MCPServer {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "orez-books-mcp",
		Version: "0.1.0",
	}, nil)

	m := &MCPServer{
		dbManager: dbManager,
		server:    server,
	}

	m.registerTools()

	return m
}

// registerTools registers all accounting tools on the server.
func (m *MCPServer) registerTools() {
	// Phase 2 Tools (Placeholders for now to verify structure)
	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "list_parties",
		Description: "Get a list of customers or suppliers from the accounting system",
	}, m.listParties)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "list_invoices",
		Description: "Get a list of Sales or Purchase invoices",
	}, m.listInvoices)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "get_invoice_details",
		Description: "Get full details of a specific invoice including items",
	}, m.getInvoiceDetails)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "get_account_balances",
		Description: "Get current balances for ledger accounts",
	}, m.getAccountBalances)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "get_profit_and_loss",
		Description: "Calculate Profit and Loss for a given date range",
	}, m.getPnL)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "get_balance_sheet",
		Description: "Get the Balance Sheet as of a specific date",
	}, m.getBalanceSheet)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "list_accounts",
		Description: "Export the full Chart of Accounts",
	}, m.listAccounts)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "create_account",
		Description: "Add a new account to the Chart of Accounts (Import)",
	}, m.createAccount)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "update_account",
		Description: "Modify an existing account in the Chart of Accounts",
	}, m.updateAccount)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "create_party",
		Description: "Create a new Customer or Supplier",
	}, m.createParty)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "create_item",
		Description: "Create a new product or service item",
	}, m.createItem)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "create_invoice",
		Description: "Create a new Sales or Purchase invoice with line items",
	}, m.createInvoice)

	mcp.AddTool(m.server, &mcp.Tool{
		Name:        "create_journal_entry",
		Description: "Create a manual or automated Journal Entry (Bank, Cash, Contra, etc.)",
	}, m.createJournalEntry)
}

// listParties handles the list_parties tool call.
func (m *MCPServer) listParties(ctx context.Context, req *mcp.CallToolRequest, args ListPartiesArgs) (*mcp.CallToolResult, ListPartiesResult, error) {
	db := m.dbManager.GetDB()
	var parties []PartyItem

	query := db.Table("Party")
	if args.Type != "" && args.Type != "Both" {
		query = query.Where("role = ?", args.Type)
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 20
	}

	err := query.Limit(limit).Find(&parties).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching parties: %v", err)}},
			IsError: true,
		}, ListPartiesResult{}, nil
	}

	var text string
	for _, p := range parties {
		text += fmt.Sprintf("- %s (%s): %s, %s\n", p.Name, p.Role, p.Email, p.Phone)
	}

	if len(parties) == 0 {
		text = "No parties found matching the criteria."
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, ListPartiesResult{Parties: parties, Count: len(parties)}, nil
}

// listInvoices handles the list_invoices tool call.
func (m *MCPServer) listInvoices(ctx context.Context, req *mcp.CallToolRequest, args ListInvoicesArgs) (*mcp.CallToolResult, ListInvoicesResult, error) {
	db := m.dbManager.GetDB()
	var invoices []InvoiceItem

	tableName := "SalesInvoice"
	if args.Type == "Purchase" {
		tableName = "PurchaseInvoice"
	}

	query := db.Table(tableName).Select("name, party, date, \"grandTotal\" as total, \"outstandingAmount\" as outstanding, submitted, cancelled")

	if args.Status == "Unpaid" {
		query = query.Where("\"outstandingAmount\" > 0 AND submitted = 1")
	} else if args.Status == "Paid" {
		query = query.Where("\"outstandingAmount\" <= 0 AND submitted = 1")
	}

	limit := args.Limit
	if limit <= 0 {
		limit = 20
	}

	rows, err := query.Limit(limit).Rows()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching invoices: %v", err)}},
			IsError: true,
		}, ListInvoicesResult{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var name, party, date string
		var total, outstanding float64
		var submitted, cancelled bool
		err := rows.Scan(&name, &party, &date, &total, &outstanding, &submitted, &cancelled)
		if err != nil {
			continue
		}

		status := "Draft"
		if cancelled {
			status = "Cancelled"
		} else if submitted {
			if outstanding > 0 {
				status = "Unpaid"
			} else {
				status = "Paid"
			}
		}

		invoices = append(invoices, InvoiceItem{
			Name:        name,
			Party:       party,
			Date:        date,
			Total:       total,
			Outstanding: outstanding,
			Status:      status,
		})
	}

	var text string
	for _, inv := range invoices {
		text += fmt.Sprintf("- %s: %s, %s, Total: %.2f, Status: %s\n", inv.Name, inv.Party, inv.Date, inv.Total, inv.Status)
	}

	if len(invoices) == 0 {
		text = "No invoices found matching the criteria."
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, ListInvoicesResult{Invoices: invoices, Count: len(invoices)}, nil
}

// getInvoiceDetails handles the get_invoice_details tool call.
func (m *MCPServer) getInvoiceDetails(ctx context.Context, req *mcp.CallToolRequest, args GetInvoiceDetailsArgs) (*mcp.CallToolResult, GetInvoiceDetailsResult, error) {
	db := m.dbManager.GetDB()
	
	// Check both Sales and Purchase invoices
	var invoice map[string]interface{}
	err := db.Table("SalesInvoice").Where("name = ?", args.Name).First(&invoice).Error
	if err != nil {
		// Try Purchase
		err = db.Table("PurchaseInvoice").Where("name = ?", args.Name).First(&invoice).Error
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invoice %s not found", args.Name)}},
				IsError: true,
			}, GetInvoiceDetailsResult{}, nil
		}
	}

	// Fetch items
	itemTable := "SalesInvoiceItem"
	if _, ok := invoice["customer"]; !ok { // Simple check to see if it's purchase
		itemTable = "PurchaseInvoiceItem"
	}
	
	var items []map[string]interface{}
	db.Table(itemTable).Where("parent = ?", args.Name).Find(&items)
	invoice["items"] = items

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Details for invoice %s retrieved", args.Name)}},
	}, GetInvoiceDetailsResult{Invoice: invoice}, nil
}

// getAccountBalances handles the get_account_balances tool call.
func (m *MCPServer) getAccountBalances(ctx context.Context, req *mcp.CallToolRequest, args GetAccountBalancesArgs) (*mcp.CallToolResult, GetAccountBalancesResult, error) {
	db := m.dbManager.GetDB()
	var balances []AccountBalanceItem

	query := db.Table("AccountingLedgerEntry as ale").
		Select("ale.account, SUM(ale.debit) - SUM(ale.credit) as balance").
		Joins("JOIN Account as a ON a.name = ale.account").
		Group("ale.account")

	if args.AccountType != "" {
		query = query.Where("a.accountType = ?", args.AccountType)
	}

	err := query.Scan(&balances).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching balances: %v", err)}},
			IsError: true,
		}, GetAccountBalancesResult{}, nil
	}

	var text string
	for _, b := range balances {
		text += fmt.Sprintf("- %s: %.2f\n", b.Account, b.Balance)
	}

	if len(balances) == 0 {
		text = "No account balances found."
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, GetAccountBalancesResult{Balances: balances}, nil
}

// getPnL handles the get_profit_and_loss tool call.
func (m *MCPServer) getPnL(ctx context.Context, req *mcp.CallToolRequest, args GetPnLArgs) (*mcp.CallToolResult, PnLResult, error) {
	db := m.dbManager.GetDB()

	var income, expense float64

	// Total Income (Credit - Debit for Income accounts)
	db.Table("AccountingLedgerEntry as ale").
		Select("SUM(ale.credit) - SUM(ale.debit)").
		Joins("JOIN Account as a ON a.name = ale.account").
		Where("a.rootType = ? AND ale.date BETWEEN ? AND ?", "Income", args.StartDate, args.EndDate).
		Scan(&income)

	// Total Expense (Debit - Credit for Expense accounts)
	db.Table("AccountingLedgerEntry as ale").
		Select("SUM(ale.debit) - SUM(ale.credit)").
		Joins("JOIN Account as a ON a.name = ale.account").
		Where("a.rootType = ? AND ale.date BETWEEN ? AND ?", "Expense", args.StartDate, args.EndDate).
		Scan(&expense)

	netProfit := income - expense

	text := fmt.Sprintf("Profit and Loss from %s to %s:\n", args.StartDate, args.EndDate)
	text += fmt.Sprintf("- Total Income: %.2f\n", income)
	text += fmt.Sprintf("- Total Expense: %.2f\n", expense)
	text += fmt.Sprintf("- Net Profit: %.2f\n", netProfit)

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, PnLResult{
		TotalIncome:  income,
		TotalExpense: expense,
		NetProfit:    netProfit,
	}, nil
}

// getBalanceSheet handles the get_balance_sheet tool call.
func (m *MCPServer) getBalanceSheet(ctx context.Context, req *mcp.CallToolRequest, args GetBalanceSheetArgs) (*mcp.CallToolResult, GetBalanceSheetResult, error) {
	db := m.dbManager.GetDB()
	var results []BalanceSheetItem

	// Assets, Liabilities, Equity (Debit - Credit)
	query := db.Table("AccountingLedgerEntry as ale").
		Select("a.rootType as root_type, SUM(ale.debit) - SUM(ale.credit) as balance").
		Joins("JOIN Account as a ON a.name = ale.account").
		Where("a.rootType IN (?) AND ale.date <= ?", []string{"Asset", "Liability", "Equity"}, args.AsOfDate).
		Group("a.rootType")

	err := query.Scan(&results).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching balance sheet: %v", err)}},
			IsError: true,
		}, GetBalanceSheetResult{}, nil
	}

	var text string
	for _, r := range results {
		text += fmt.Sprintf("- %s: %.2f\n", r.RootType, r.Balance)
	}

	if len(results) == 0 {
		text = "No balance sheet data found."
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, GetBalanceSheetResult{Items: results}, nil
}

// listAccounts handles the list_accounts tool call.
func (m *MCPServer) listAccounts(ctx context.Context, req *mcp.CallToolRequest, args ListAccountsArgs) (*mcp.CallToolResult, ListAccountsResult, error) {
	db := m.dbManager.GetDB()
	var accounts []AccountItem

	query := db.Table("Account").Select("name, rootType as root_type, accountType as account_type, parentAccount as parent_account, isGroup as is_group")
	if args.RootType != "" {
		query = query.Where("rootType = ?", args.RootType)
	}

	err := query.Find(&accounts).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error fetching accounts: %v", err)}},
			IsError: true,
		}, ListAccountsResult{}, nil
	}

	var text string
	for _, a := range accounts {
		groupSuffix := ""
		if a.IsGroup {
			groupSuffix = " (Group)"
		}
		text += fmt.Sprintf("- %s: %s / %s%s\n", a.Name, a.RootType, a.AccountType, groupSuffix)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, ListAccountsResult{Accounts: accounts}, nil
}

// createAccount handles the create_account tool call.
func (m *MCPServer) createAccount(ctx context.Context, req *mcp.CallToolRequest, args CreateAccountArgs) (*mcp.CallToolResult, map[string]interface{}, error) {
	db := m.dbManager.GetDB()
	now := time.Now().Format("2006-01-02 15:04:05")

	if args.ParentAccount != "" {
		if err := m.validateParentAccount(args.ParentAccount); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
				IsError: true,
			}, nil, nil
		}
	}

	account := map[string]interface{}{
		"name":          args.Name,
		"rootType":      args.RootType,
		"accountType":   args.AccountType,
		"parentAccount": args.ParentAccount,
		"isGroup":       args.IsGroup,
		"created":       now,
		"modified":      now,
		"createdBy":     "mcp-server",
		"modifiedBy":    "mcp-server",
	}

	err := db.Table("Account").Create(&account).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error creating account: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Account %s created successfully", args.Name)}},
	}, map[string]interface{}{"success": true, "name": args.Name}, nil
}

// updateAccount handles the update_account tool call.
func (m *MCPServer) updateAccount(ctx context.Context, req *mcp.CallToolRequest, args UpdateAccountArgs) (*mcp.CallToolResult, map[string]interface{}, error) {
	db := m.dbManager.GetDB()
	now := time.Now().Format("2006-01-02 15:04:05")

	if parent, ok := args.Updates["parent_account"].(string); ok && parent != "" {
		if err := m.validateParentAccount(parent); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
				IsError: true,
			}, nil, nil
		}
	}

	args.Updates["modified"] = now
	args.Updates["modifiedBy"] = "mcp-server"

	err := db.Table("Account").Where("name = ?", args.Name).Updates(args.Updates).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error updating account: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Account %s updated successfully", args.Name)}},
	}, map[string]interface{}{"success": true}, nil
}

func (m *MCPServer) validateParentAccount(parentName string) error {
	db := m.dbManager.GetDB()
	var parent map[string]interface{}
	err := db.Table("Account").Where("name = ?", parentName).First(&parent).Error
	if err != nil {
		return fmt.Errorf("parent account %s not found", parentName)
	}
	
	isGroup, _ := parent["isGroup"].(int64)
	if isGroup == 0 {
		// Try boolean check if driver returns it differently
		isGroupBool, _ := parent["isGroup"].(bool)
		if !isGroupBool {
			return fmt.Errorf("account %s is not a Group account and cannot be a parent", parentName)
		}
	}
	return nil
}

// createParty handles the create_party tool call.
func (m *MCPServer) createParty(ctx context.Context, req *mcp.CallToolRequest, args CreatePartyArgs) (*mcp.CallToolResult, CreatePartyResult, error) {
	db := m.dbManager.GetDB()
	now := time.Now().Format("2006-01-02 15:04:05")

	party := map[string]interface{}{
		"name":       args.Name,
		"role":       args.Role,
		"email":      args.Email,
		"phone":      args.Phone,
		"created":    now,
		"modified":   now,
		"createdBy":  "mcp-server",
		"modifiedBy": "mcp-server",
	}

	err := db.Table("Party").Create(&party).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error creating party: %v", err)}},
			IsError: true,
		}, CreatePartyResult{Success: false}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Party %s created successfully", args.Name)}},
	}, CreatePartyResult{Success: true, Name: args.Name}, nil
}

// createItem handles the create_item tool call.
func (m *MCPServer) createItem(ctx context.Context, req *mcp.CallToolRequest, args CreateItemArgs) (*mcp.CallToolResult, CreateItemResult, error) {
	db := m.dbManager.GetDB()
	now := time.Now().Format("2006-01-02 15:04:05")

	itemType := args.ItemType
	if itemType == "" {
		itemType = "Product"
	}

	item := map[string]interface{}{
		"name":           args.Name,
		"itemType":       itemType,
		"rate":           args.Rate,
		"incomeAccount":  args.IncomeAccount,
		"expenseAccount": args.ExpenseAccount,
		"created":        now,
		"modified":       now,
		"createdBy":      "mcp-server",
		"modifiedBy":     "mcp-server",
		"for":            "Both",
	}

	err := db.Table("Item").Create(&item).Error
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error creating item: %v", err)}},
			IsError: true,
		}, CreateItemResult{Success: false}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Item %s created successfully", args.Name)}},
	}, CreateItemResult{Success: true, Name: args.Name}, nil
}

// createInvoice handles the create_invoice tool call.
func (m *MCPServer) createInvoice(ctx context.Context, req *mcp.CallToolRequest, args CreateInvoiceArgs) (*mcp.CallToolResult, CreateInvoiceResult, error) {
	db := m.dbManager.GetDB()
	now := time.Now().Format("2006-01-02 15:04:05")

	tableName := "SalesInvoice"
	itemTableName := "SalesInvoiceItem"
	if args.Type == "Purchase" {
		tableName = "PurchaseInvoice"
		itemTableName = "PurchaseInvoiceItem"
	}

	invoiceName := m.generateName(args.Type)
	var netTotal, grandTotal float64

	// Start Transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Process items and calculate totals
		for _, itemInput := range args.Items {
			rate := itemInput.Rate
			if rate == 0 {
				// Fetch standard rate from Item table
				tx.Table("Item").Where("name = ?", itemInput.Item).Select("rate").Scan(&rate)
			}

			incomeAccount := itemInput.Account
			if incomeAccount == "" {
				// Fetch standard account
				accountField := "incomeAccount"
				if args.Type == "Purchase" {
					accountField = "expenseAccount"
				}
				tx.Table("Item").Where("name = ?", itemInput.Item).Select(accountField).Scan(&incomeAccount)
			}

			amount := itemInput.Quantity * rate
			netTotal += amount

			itemRow := map[string]interface{}{
				"name":        m.generateName("Item"),
				"parent":      invoiceName,
				"parenttype":  tableName,
				"parentfield": "items",
				"item":        itemInput.Item,
				"quantity":    itemInput.Quantity,
				"rate":        rate,
				"amount":      amount,
				"account":     incomeAccount,
				"created":     now,
				"modified":    now,
			}

			if err := tx.Table(itemTableName).Create(&itemRow).Error; err != nil {
				return err
			}
		}

		grandTotal = netTotal // Simplified: no taxes for now

		invoiceRow := map[string]interface{}{
			"name":              invoiceName,
			"party":             args.Party,
			"date":              args.Date,
			"account":           args.Account,
			"netTotal":          netTotal,
			"grandTotal":        grandTotal,
			"outstandingAmount": grandTotal,
			"submitted":         0,
			"cancelled":         0,
			"created":           now,
			"modified":          now,
			"createdBy":         "mcp-server",
			"modifiedBy":        "mcp-server",
		}

		if err := tx.Table(tableName).Create(&invoiceRow).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error creating invoice: %v", err)}},
			IsError: true,
		}, CreateInvoiceResult{Success: false}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("%s %s created successfully", args.Type, invoiceName)}},
	}, CreateInvoiceResult{Success: true, Name: invoiceName}, nil
}

// createJournalEntry handles the create_journal_entry tool call.
func (m *MCPServer) createJournalEntry(ctx context.Context, req *mcp.CallToolRequest, args CreateJournalEntryArgs) (*mcp.CallToolResult, map[string]interface{}, error) {
	db := m.dbManager.GetDB()
	now := time.Now().Format("2006-01-02 15:04:05")

	var totalDebit, totalCredit float64
	for _, acc := range args.Accounts {
		totalDebit += acc.Debit
		totalCredit += acc.Credit
	}

	if totalDebit != totalCredit {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Journal Entry is not balanced. Total Debit: %.2f, Total Credit: %.2f", totalDebit, totalCredit)}},
			IsError: true,
		}, nil, nil
	}

	jvName := m.generateName("JV")

	err := db.Transaction(func(tx *gorm.DB) error {
		// Create Child Entries
		for _, acc := range args.Accounts {
			itemRow := map[string]interface{}{
				"name":        m.generateName("JV-Item"),
				"parent":      jvName,
				"parenttype":  "JournalEntry",
				"parentfield": "accounts",
				"account":     acc.Account,
				"debit":       acc.Debit,
				"credit":      acc.Credit,
				"created":     now,
				"modified":    now,
			}
			if err := tx.Table("JournalEntryAccount").Create(&itemRow).Error; err != nil {
				return err
			}
		}

		// Create Main Entry
		jvRow := map[string]interface{}{
			"name":       jvName,
			"entryType":  args.EntryType,
			"date":       args.Date,
			"userRemark": args.Remark,
			"submitted":  0,
			"cancelled":  0,
			"created":    now,
			"modified":   now,
			"createdBy":  "mcp-server",
			"modifiedBy": "mcp-server",
		}
		if err := tx.Table("JournalEntry").Create(&jvRow).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error creating journal entry: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Journal Entry %s created successfully", jvName)}},
	}, map[string]interface{}{"success": true, "name": jvName}, nil
}

func (m *MCPServer) generateName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// Run starts the MCP server using the specified transport ("stdio" or "sse").
func (m *MCPServer) Run(ctx context.Context, transport, addr string) error {
	switch transport {
	case "stdio", "":
		return m.server.Run(ctx, &mcp.StdioTransport{})
	case "sse":
		handler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
			return m.server
		}, nil) // Pass nil for default options
		fmt.Fprintf(os.Stderr, "MCP SSE server listening on %s\n", addr)
		srv := &http.Server{Addr: addr, Handler: handler}
		go func() {
			<-ctx.Done()
			srv.Close()
		}()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported transport: %q (use stdio or sse)", transport)
	}
}
