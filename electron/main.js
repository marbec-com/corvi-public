'use strict';

var spawn = require('child_process').spawn;
var path = require('path');
var backend = null;

const electron = require('electron');
const app = electron.app;
const BrowserWindow = electron.BrowserWindow;

let mainWindow;

function createWindow() {
	mainWindow = new BrowserWindow({ width: 800, height: 600, minWidth: 400, minHeight: 500, 'title': 'Corvi', 'titleBarStyle': 'hidden-inset' });
	mainWindow.loadURL('http://localhost:8080/app/');
	mainWindow.on('closed', function () {
		mainWindow = null;
	});
}

function startBackend() {

	if (process.platform === 'darwin') {
		backend = spawn('./corvi-backend', [], { env: { USER_DATA: app.getPath("userData") }, cwd: __dirname });
	} else if (process.platform === 'win32')  {
		backend = spawn('./corvi-backend.exe', [], { env: { USER_DATA: app.getPath("userData") }, cwd: __dirname });
	} else if (process.platform === 'linux') {
		backend = spawn('./corvi-backend', [], { env: { USER_DATA: app.getPath("userData") }, cwd: __dirname });
	} else {
		throw new Error("Incompatible operating system: ", process.platform);
	}

	backend.stdout.on('data', function (data) {
		console.log('Go : ' + data);
	});
	backend.stderr.on('data', function (data) {
		console.log('Go Error: ' + data);
	});
	backend.on('close', function (code) {
		console.log('Backend exited: ' + code);
	});

}

app.on('ready', function () {
	startBackend();
	createWindow();
});

app.on('window-all-closed', function () {
	if (process.platform !== 'darwin') {
		app.quit();
	}
});

app.on('will-quit', function () {
	backend.kill();
});

app.on('activate', function () {
	if (mainWindow === null) {
		createWindow();
	}
});
