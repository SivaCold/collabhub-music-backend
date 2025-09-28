const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');
const ffmpegPath = require('ffmpeg-static'); // binaire intégré
const ffmpeg = require('fluent-ffmpeg');

let mainWindow;

app.on('ready', () => {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js'),
      webSecurity: false 
    }
  });

  mainWindow.loadFile(path.join(__dirname, '../build/index.html'));
});

// Configurer fluent-ffmpeg avec le bon binaire
ffmpeg.setFfmpegPath(ffmpegPath);

// IPC : découper une vidéo
ipcMain.handle('cut-video', async (event, filePath) => {
  return new Promise((resolve, reject) => {
    const output = path.join(app.getPath("desktop"), "output.mp4");

    ffmpeg(filePath)
      .setStartTime("00:00:05")
      .setDuration(10)
      .output(output)
      .on('end', () => resolve(`Vidéo générée : ${output}`))
      .on('error', (err) => reject(err.message))
      .run();
  });
});
