# Notices

## Frappe Books
This project is a re-architecture of **Frappe Books**, originally developed by **Frappe Technologies Pvt. Ltd. and contributors**.

- **Original Project URL:** [https://github.com/frappe/books](https://github.com/frappe/books)
- **Original License:** GNU Affero General Public License v3.0 (AGPL-3.0)

## Attribution & Changes
This project, **Orez Books**, maintains the original Vue.js 3 frontend logic and SQLite database structure of Frappe Books while replacing the **Electron** desktop framework with **Wails**.

Significant changes include:
1.  **Backend Transition:** Migration of the main process from Node.js/Electron to Go/Wails.
2.  **IPC Replacement:** Replacement of Electron IPC with Wails Go-to-JS bindings.
3.  **Database Driver:** Transition from `better-sqlite3` (Node.js) to Go-native SQLite drivers.
4.  **Native Integration:** Utilization of Wails Runtime APIs for system menus, tray icons, and dialogs.

Copyright (c) 2024 Frappe Technologies Pvt. Ltd. and contributors.
Copyright (c) 2024 [Orez.Digital, Ashish Banerjee] and contributors.
