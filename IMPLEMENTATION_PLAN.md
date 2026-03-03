# Implementation Plan: Frappe Books Re-architecture (Electron to Wails)

## Phase 1: Preparation & Setup
- [x] Research and Document Electron-specific API usage.
- [x] Initialize Wails project in the root directory.
  - `wails init -n orez-books -t vue`
- [x] Configure `wails.json` and project structure.
- [x] Integrate existing Vue 3 frontend (from source repo) into Wails frontend directory (`/frontend`).
- [x] Verify frontend build pipeline works with Wails.

## Phase 2: Backend Logic Migration (Go)
- [ ] Implement Go-based SQLite integration.
  - [x] Choose SQLite driver (`modernc.org/sqlite` chosen for Cgo-free experience).
  - [x] Implement database connection and management logic in Go.
- [x] Port business logic from Electron's main process to Go methods.
  - [x] Identify and port file system operations.
  - [x] Identify and port data processing tasks.
- [x] Set up GORM or equivalent Go-based query builder.
- [x] Implement database migrations in Go (Dynamic schema creation from JSON).

## Phase 3: Communication Bridge (IPC to Bindings)
- [x] Define Go methods to be exposed to the frontend.
- [x] Use `wails build` or `wails dev` to auto-generate the JS client bindings.
- [x] Replace `ipcRenderer` calls in the frontend with Wails bindings (via new `ipc` abstraction).
- [x] Verify data flow between Go backend and Vue frontend (Verified via build & Go tests).

## Phase 4: Native Features & OS Integration
- [x] Implement application menus using Wails APIs.
- [ ] Implement tray icons (Optional/Future).
- [x] Port native dialogs (Open, Save, Alerts) to Wails runtime.
- [x] Replace `electron-store` with Go-based persistent configuration (e.g., using `os.UserConfigDir`).
- [x] Implement window management logic (maximize, minimize, close) using Wails runtime.

## Phase 5: Testing & Optimization
- [x] Conduct unit tests for Go backend logic (Config and Database).
- [ ] Perform integration tests for frontend-backend communication.
- [ ] Verify SQLite database compatibility and migration success.
- [ ] Optimize Go binary size and runtime memory footprint.
- [ ] Test cross-platform builds (macOS, Windows, Linux).

## Phase 6: Finalization & Distribution
- [ ] Configure `wails build` for all target platforms.
- [ ] Set up automated CI/CD for building releases.
- [ ] Verify auto-update mechanisms (if applicable).
- [ ] Create installation packages (DMG, EXE, DEB/RPM).
- [ ] Document the new build and development processes.
