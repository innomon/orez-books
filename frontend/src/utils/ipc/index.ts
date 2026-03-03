import * as App from '../../../wailsjs/go/main/App';
import { EventsOn } from '../../../wailsjs/runtime';
import { IPC_CHANNELS } from '../messages';

const ipc = {
  desktop: true,

  reloadWindow() {
    window.location.reload();
  },

  minimizeWindow() {
    App.MinimizeMainWindow();
  },

  toggleMaximize() {
    App.MaximizeMainWindow();
  },

  isMaximized() {
    // Wails doesn't have a direct isMaximized binding in JS runtime by default
    // We might need to implement this in Go if really needed
    return Promise.resolve(false);
  },

  isFullscreen() {
    return Promise.resolve(false);
  },

  closeWindow() {
    App.CloseMainWindow();
  },

  async getCreds() {
    // Placeholder for now
    return { url: '', token: '' };
  },

  async getLanguageMap(code: string) {
    // Placeholder - we should port this logic to Go
    return {
      languageMap: {},
      success: true,
      message: '',
    };
  },

  async getTemplates(posTemplateWidth?: number) {
    // Placeholder
    return [];
  },

  async initScheduler(time: string) {
    // Placeholder
  },

  async selectFile(options: any) {
    const path = await App.GetOpenFilePath(options.title || 'Select File');
    return {
      filePath: path,
      success: !!path,
      canceled: !path,
      name: path ? path.split('/').pop() : '',
      data: new Uint8Array(), // We might need to read the file in Go
    };
  },

  async getSaveFilePath(options: any) {
    const path = await App.GetSaveFilePath(options.title || 'Save File', options.defaultPath || '');
    return {
      filePath: path,
      canceled: !path,
    };
  },

  async getOpenFilePath(options: any) {
    const path = await App.GetOpenFilePath(options.title || 'Open File');
    return {
      filePaths: path ? [path] : [],
      canceled: !path,
    };
  },

  async checkDbAccess(filePath: string) {
    return await App.CheckDbAccess(filePath);
  },

  async checkForUpdates() {
    // Placeholder
  },

  openLink(link: string) {
    App.OpenExternal(link);
  },

  async deleteFile(filePath: string) {
    try {
      await App.DeleteFile(filePath);
      return { success: true };
    } catch (err) {
      return { success: false, error: String(err) };
    }
  },

  async saveData(data: string, savePath: string) {
    await App.SaveData(data, savePath);
  },

  showItemInFolder(filePath: string) {
    // Placeholder - Wails runtime doesn't have showItemInFolder directly
  },

  async makePDF(html: string, savePath: string, width: number, height: number) {
    // Placeholder
    return false;
  },

  async printDocument(html: string, width: number, height: number) {
    // Placeholder
    return false;
  },

  async getDbList() {
    // Placeholder - needs Go implementation
    return [];
  },

  async getDbDefaultPath(companyName: string) {
    return await App.GetDbDefaultPath(companyName);
  },

  async getEnv() {
    return await App.GetEnv() as any;
  },

  openExternalUrl(url: string) {
    App.OpenExternal(url);
  },

  async showError(title: string, content: string) {
    await App.ShowError(title, content);
  },

  async sendError(body: string) {
    // Placeholder
  },

  async sendAPIRequest(endpoint: string, options: RequestInit | undefined) {
    // We can use native fetch or port to Go
    const response = await fetch(endpoint, options);
    return await response.json();
  },

  registerMainProcessErrorListener(listener: any) {
    EventsOn(IPC_CHANNELS.LOG_MAIN_PROCESS_ERROR, listener);
  },

  registerTriggerFrontendActionListener(listener: any) {
    EventsOn(IPC_CHANNELS.TRIGGER_ERPNEXT_SYNC, listener);
  },

  registerConsoleLogListener(listener: any) {
    EventsOn(IPC_CHANNELS.CONSOLE_LOG, listener);
  },

  db: {
    async getSchema() {
      // NEEDS GO IMPLEMENTATION
      return { success: false, error: 'Not implemented' };
    },

    async create(dbPath: string, countryCode?: string) {
      try {
        await App.DbCreate(dbPath, countryCode || '');
        return { success: true };
      } catch (err) {
        return { success: false, error: String(err) };
      }
    },

    async connect(dbPath: string, countryCode?: string) {
      try {
        await App.DbConnect(dbPath);
        return { success: true };
      } catch (err) {
        return { success: false, error: String(err) };
      }
    },

    async call(method: string, ...args: unknown[]) {
      // NEEDS GO IMPLEMENTATION
      return { success: false, error: 'Not implemented' };
    },

    async bespoke(method: string, ...args: unknown[]) {
      // NEEDS GO IMPLEMENTATION
      return { success: false, error: 'Not implemented' };
    },
  },

  store: {
    get(key: string) {
      return App.ConfigGet(key);
    },

    set(key: string, value: any) {
      return App.ConfigSet(key, value);
    },

    delete(key: string) {
      return App.ConfigDelete(key);
    },
  },
};

export default ipc;
export type IPC = typeof ipc;
