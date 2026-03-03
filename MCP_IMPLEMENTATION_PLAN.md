# Implementation Plan: Orez Books MCP Server

This plan outlines the steps required to build and integrate the Model Context Protocol (MCP) server for Orez Books, allowing AI agents to interact with the local accounting database.

## Phase 1: Setup and Project Structure
- [x] Initialize the MCP server directory within the project (e.g., `cmd/mcp-server`).
- [x] Add the required Go MCP SDK dependency: `go get github.com/modelcontextprotocol/go-sdk/mcp`.
- [x] Create the core `mcp_server.go` file outlining the `MCPServer` struct.
- [x] Implement the `Run` method supporting `stdio` transport.
- [x] Integrate the existing `internal/database` manager into the `MCPServer` struct so it can connect to the local SQLite database.

## Phase 2: Implementation of Data Retrieval Tools (Read)
- [x] Define argument structs and result structs for Read tools (e.g., `ListPartiesArgs`, `ListPartiesResult`).
- [x] Implement `list_parties` tool: Query the `Party` table using GORM and return formatted JSON results.
- [x] Implement `list_invoices` tool: Query the `SalesInvoice` and `PurchaseInvoice` tables, allowing filtering by status.
- [x] Implement `get_invoice_details` tool: Join `SalesInvoice` with `SalesInvoiceItem` to return the full document.
- [x] Register these tools in the `MCPServer.registerTools()` method.

## Phase 3: Implementation of Financial Reporting Tools (Read)
- [x] Define structs for financial reports (e.g., `GetPnLArgs`, `ReportResult`).
- [x] Implement `get_account_balances` tool: Aggregate totals from the `AccountingLedgerEntry` table based on account types.
- [x] Implement `get_profit_and_loss` tool: Calculate revenue minus expenses over a given date range.
- [ ] Implement `get_balance_sheet` tool: Calculate assets, liabilities, and equity as of a specific date.
- [x] Register these reporting tools.

## Phase 4: Implementation of Entity Management Tools (Write)
- [x] Define argument structs for Write tools with strict validation (e.g., `CreateInvoiceArgs`).
- [x] Implement `create_party` tool: Validate input and insert a new record into the `Party` table.
- [x] Implement `create_item` tool: Validate input and insert into the `Item` table.
- [ ] Implement `create_invoice` tool:
  - Validate the `party_id` exists.
  - Validate all `items` exist.
  - Insert the main invoice record.
  - Insert associated line items into the child table.
  - Handle this entire operation within a single GORM transaction to ensure data integrity.
- [x] Register these CRUD tools.

## Phase 5: Testing and Integration
- [ ] Write Go unit tests for each MCP tool handler, mocking the database or using an in-memory SQLite instance.
- [ ] Compile the MCP server into a standalone executable.
- [ ] Test the server locally using the MCP Inspector or Claude Desktop.
- [ ] (Optional) Add `sse` / `streamable-http` transport support for remote agent access.
