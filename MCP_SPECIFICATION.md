# MCP Specification: Orez Books Accounting Server

## 1. Overview
The Orez Books MCP Server provides a standardized Model Context Protocol (MCP) interface, enabling AI agents (like Claude, Gemini, Cursor) to directly interact with the Orez Books local SQLite accounting database. This allows agents to perform automated bookkeeping, financial analysis, and cross-application workflows seamlessly.

## 2. Architecture
- **Language:** Go (v1.21+)
- **SDK:** `github.com/modelcontextprotocol/go-sdk/mcp`
- **Database Layer:** Interacts with the `internal/database/manager.go` using GORM and Cgo-free SQLite (`github.com/glebarez/sqlite`).
- **Transport:** Supports both `stdio` (for local desktop integration) and `sse`/`streamable-http` (for remote agent access).

## 3. Capabilities and Tools exposed to MCP
The server will expose a set of tools categorized by function. Unlike the RAG template which handles unstructured documents, this server handles structured financial data.

### 3.1 Entity Management (CRUD)
Tools to Add, Modify, and Delete core accounting entities.

- **`create_party`**: Create a new Customer or Supplier.
  - Args: `name` (string), `type` (Customer/Supplier), `email` (string), `phone` (string).
- **`update_party`**: Update existing Customer/Supplier details.
  - Args: `id` (string), updates mapping.
- **`create_item`**: Create a new product or service item.
  - Args: `name` (string), `description` (string), `default_price` (float).
- **`create_invoice`**: Create a Sales Invoice or Purchase Invoice.
  - Args: `party_id` (string), `type` (Sales/Purchase), `items` (array of item IDs and quantities), `date` (string).
- **`create_payment`**: Record a payment against an invoice.
  - Args: `invoice_id` (string), `amount` (float), `date` (string), `account` (string).

### 3.2 Data Retrieval (Read)
Tools for agents to query lists of entities and specific transaction details.

- **`list_parties`**: Get a list of customers/suppliers.
  - Args: `type` (optional, filter by Customer/Supplier), `limit` (int).
- **`list_invoices`**: Get unpaid or all invoices.
  - Args: `status` (optional, e.g., Unpaid, Paid), `limit` (int).
- **`get_invoice_details`**: Get full line items and details for a specific invoice.
  - Args: `invoice_id` (string).

### 3.3 Financial Reporting (Read)
Tools specifically designed for financial analysis.

- **`get_profit_and_loss`**: Retrieve P&L data for a specific date range.
  - Args: `start_date` (string), `end_date` (string).
- **`get_balance_sheet`**: Retrieve the Balance Sheet as of a specific date.
  - Args: `as_of_date` (string).
- **`get_account_balances`**: Get the current balance of specific ledger accounts (e.g., Bank, Cash).
  - Args: `account_type` (string).

## 4. Security & Authentication
- Since Orez Books is a local-first application, the primary deployment of this MCP server will be via `stdio` running on the user's local machine, inheriting the user's local file permissions.
- The server will locate the active SQLite database using the `App.GetDbDefaultPath()` logic established in the Wails backend.
- Write operations (Create/Update) will rigorously enforce database schemas and foreign key constraints to prevent data corruption by the AI agent.
