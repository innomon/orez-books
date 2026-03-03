# Orez Books

Orez Books is a re-architecture of the original **Frappe Books**, migrating the desktop application from **Electron** to **Wails**. It maintains the powerful Vue.js 3 frontend logic while replacing the Node.js main process with a high-performance **Go** backend and a Cgo-free SQLite implementation.

## 🚀 Why Wails?

- **Binary Size:** Drastically smaller distribution size (typically < 20MB vs 100MB+ for Electron).
- **Resource Efficiency:** Reduced memory footprint by replacing the V8 backend with a compiled Go binary.
- **Native Experience:** Direct access to native OS APIs through the Wails runtime without the overhead of heavy IPC.
- **Cgo-Free:** Uses `modernc.org/sqlite` via GORM for a seamless cross-platform build experience without needing a C compiler.

## 🛠 Tech Stack

- **Frontend:** Vue.js 3, Vite, Tailwind CSS, TypeScript.
- **Backend:** Go (v1.21+), Wails (v2).
- **Database:** SQLite managed via **GORM** with dynamic schema generation from JSON.
- **Communication:** Wails auto-generated JS/TS bindings.

## 🤖 Model Context Protocol (MCP) Server

Orez Books includes a built-in **MCP Server**, allowing AI agents (like Claude Desktop, Gemini, or Cursor) to interact directly with your accounting data.

### MCP Tools included:
- **Master Data:** List and Create Customers, Suppliers, Items, and Ledger Accounts.
- **Transactions:** Create and List Invoices (Sales/Purchase) and multi-line Journal Entries.
- **Reporting:** Retrieve real-time Profit & Loss and Balance Sheet data.
- **Chart of Accounts:** Export, Import, and Modify the COA with hierarchical validation.

### Running the MCP Server:
```bash
# Build the MCP server
go build -o build/bin/orez-books-mcp ./cmd/mcp-server

# Run with your database
./build/bin/orez-books-mcp -db ~/Documents/Frappe\ Books/my_company.books.db
```

## 💻 Development

### Prerequisites
- [Go](https://golang.org/dl/) (v1.21+)
- [Node.js](https://nodejs.org/) & npm
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Live Development
To start the app in development mode with hot-reloading:
```bash
wails dev
```

### Building for Production
To generate a production-ready native binary:
```bash
wails build
```

## 📁 Project Structure

- `frontend/`: Vue.js 3 frontend source code.
- `internal/`: Core Go packages for database management and configuration.
- `cmd/mcp-server/`: Source code for the Model Context Protocol server.
- `app.go`: Main Wails application logic and backend bindings.
- `main.go`: Entry point for the Wails desktop application.

## 📜 Credits & License

Orez Books is built upon the incredible foundation of [Frappe Books](https://github.com/frappe/books) by **Frappe Technologies Pvt. Ltd.**

- **Original License:** GNU Affero General Public License v3.0 (AGPL-3.0)
- **Orez Books License:** GNU Affero General Public License v3.0 (AGPL-3.0)

See `NOTICE.md` for full attribution details.
