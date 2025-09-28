const { app, BrowserWindow, ipcMain, dialog } = require('electron');
const path = require('path');
const fs = require('fs');

let mainWindow;

app.on('ready', () => {
  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, './preload.js'),
      webSecurity: false 
    }
  });

  mainWindow.loadFile(path.join(__dirname, '../build/index.html'));
});

// Ouvrir un répertoire
ipcMain.handle('select-folder', async () => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openDirectory']
  });
  return result.filePaths[0] || null;
});


// Lister les vidéos du répertoire
ipcMain.handle('list-videos', async (event, folderPath) => {
  if (!folderPath) return [];
  const files = fs.readdirSync(folderPath);
  return files.filter(f => f.match(/\.(mp4|mov|avi|mkv)$/i)).map(f => path.join(folderPath, f));
});

// Déplacer une vidéo (exemple : swipe gauche = corbeille)
ipcMain.handle('move-video', async (event, { filePath, targetDir }) => {
  const fileName = path.basename(filePath);
  const destPath = path.join(targetDir, fileName);
  fs.renameSync(filePath, destPath);
  return destPath;
});
