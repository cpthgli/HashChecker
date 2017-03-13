const http = require('http');
const server = http.createServer((req, res) => {});
let port;
server.listen(0);
server.on('listening', () => {
  port = server.address().port;
})

const electron = require('electron');
const app = electron.app;
const Menu = electron.Menu;
const menu = new Menu()
const BrowserWindow = electron.BrowserWindow;
let mainWindow;

function createMainWindow() {
  server.close();
  mainWindow = new BrowserWindow({
    webPreferences: {
      nodeIntegration: false
    },
    width: 590,
    height: 238,
    resizable: false,
    show: false
  });
  mainWindow.setMenu(null);

  function load() {
    var mainAddr = 'http://localhost:' + port;
    mainWindow.loadURL(mainAddr);
  }

  server.close();
  var go = require('child_process').execFile('./server', ['-path', process.argv.slice(1).join(' '), '-port', port]);
  load();

  mainWindow.once('ready-to-show', () => {
    mainWindow.show()
  })

  mainWindow.on('closed', function () {
    mainWindow = null;
    go.kill('SIGINT');
  });
}
app.on('ready', createMainWindow);

app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') {
    app.quit();
  }
})

app.on('activate', function () {
  if (mainWindow === null) {
    createMainWindow();
  }
})
