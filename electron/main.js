'use strict';

var spawn = require('child_process').spawn;
var path = require('path');
var backend = null;

const electron = require('electron');
// Module to control application life.
const app = electron.app;
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow;

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow;

function createWindow () {
  // Create the browser window.
  mainWindow = new BrowserWindow({width: 800, height: 600, minWidth: 400, minHeight: 500, 'title': 'Corvi', 'titleBarStyle': 'hidden-inset'});

  // and load the index.html of the app.
  mainWindow.loadURL('http://localhost:8080/app/');

  // Open the DevTools.
  //mainWindow.webContents.openDevTools();

  // Emitted when the window is closed.
  mainWindow.on('closed', function() {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null;
  });
}

function startBackend() {
	
	if (process.platform === 'darwin') {
		backend = spawn('./darwin_amd64_corvi-backend', [], {cwd: __dirname});
	}else if (process.platform === 'win32') {
		backend = spawn('./windows_amd64_corvi-backend.exe', [], {cwd: __dirname});
	}else if (process.platform === 'linux') {
		backend = spawn('./linux_amd64_corvi-backend', [], {cwd: __dirname});
	}else{
		throw new Error("Incompatible operating system: ", process.platform);
	}
	
	backend.stdout.on('data', function(data) {
		console.log('Go : ' + data);
	});
	backend.stderr.on('data', function(data) {
		console.error('Go Error: ' + data);
	});
	backend.on('close', function(code) {
		console.log('Backend exited: ' + code);
	});

}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
app.on('ready', function() {
	startBackend();
	createWindow();
});

// Quit when all windows are closed.
app.on('window-all-closed', function () {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('will-quit', function() {
	backend.kill();
});

app.on('activate', function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow();
  }
});
