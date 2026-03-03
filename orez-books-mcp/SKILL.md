---
name: orez-books-mcp
description: Advanced accounting and bookkeeping capabilities using the Orez Books SQLite backend. Use this skill to manage customers, vendors, items, invoices, journal entries, and generate financial reports like Profit & Loss and Balance Sheets.
license: AGPL-3.0
metadata:
  version: 1.0.0
  author: Innomon
allowed-tools:
  - list_parties
  - create_party
  - list_invoices
  - create_invoice
  - get_invoice_details
  - create_item
  - list_accounts
  - create_account
  - update_account
  - get_profit_and_loss
  - get_balance_sheet
  - create_journal_entry
---

# Orez Books Accounting Skill

This skill enables an agent to perform professional bookkeeping and financial analysis by interacting with the Orez Books local database via its MCP server.

## When to Use
- **Automated Invoicing:** When a user wants to create sales or purchase invoices from text descriptions or uploaded data.
- **Master Data Management:** When adding new customers, suppliers, or items to the system.
- **Financial Reporting:** When asked for Profit & Loss statements, Balance Sheets, or account balances.
- **Journal Entries:** For manual adjustments, bank transfers, or complex multi-ledger transactions.
- **Chart of Accounts (COA):** For exporting, importing, or restructuring the accounting hierarchy.

## Step-by-Step Instructions

### 1. Initial Discovery
Always start by listing the Chart of Accounts using `list_accounts` to understand the existing ledger structure before attempting to create invoices or journal entries. This ensures you use the correct account names.

### 2. Entity Verification
Before creating an invoice for a customer or supplier, use `list_parties` to check if they already exist. If not, use `create_party`.

### 3. Creating Invoices
When creating an invoice:
- Ensure all items exist (`list_items` or `create_item`).
- Provide a valid `party` name.
- Specify the main `account` (e.g., "Debtors" for sales, "Creditors" for purchases).
- The `create_invoice` tool will automatically handle the transaction safety.

### 4. Recording Journal Entries
For transactions that don't fit the invoice model (like internal transfers):
- Use `create_journal_entry`.
- Ensure the sum of `debit` amounts equals the sum of `credit` amounts.
- The tool will reject unbalanced entries.

### 5. Financial Analysis
Use `get_profit_and_loss` for performance over time and `get_balance_sheet` for a snapshot of the company's financial position at a specific date.

## Examples

### Example 1: Creating a Customer and an Invoice
**User:** "Add a new customer named 'Acme Corp' and send them an invoice for 5 'Consulting Services' at $100 each today."
**Agent Workflow:**
1. Call `create_party` with `name: "Acme Corp"`, `role: "Customer"`.
2. Call `create_invoice` with:
   - `type: "Sales"`
   - `party: "Acme Corp"`
   - `date: "2024-03-03"`
   - `items: [{ "item": "Consulting Services", "quantity": 5, "rate": 100 }]`
   - `account: "Debtors"`

### Example 2: Financial Health Check
**User:** "How much profit did we make last month?"
**Agent Workflow:**
1. Call `get_profit_and_loss` with `start_date: "2024-02-01"`, `end_date: "2024-02-29"`.
2. Summarize the `net_profit` to the user.

## Safety & Integrity
- **Transactional Consistency:** Do not worry about partial writes; the server uses SQL transactions for complex tools like `create_invoice`.
- **Validation:** Hierarchical accounts require a parent. Always ensure the parent account is a "Group" account before nesting.
