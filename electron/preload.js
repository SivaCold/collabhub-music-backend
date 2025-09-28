const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
  selectFolder: () => ipcRenderer.invoke('select-folder'),
  listVideos: (folderPath) => ipcRenderer.invoke('list-videos', folderPath),
  moveVideo: (data) => ipcRenderer.invoke('move-video', data)
});