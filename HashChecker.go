package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"errors"

	"github.com/andlabs/ui"
)

type mainWindow struct {
	window     *ui.Window
	configPage *ui.Box
	resultPage *ui.Box
	config     *Config
	fileName   string
	fileBuf    []byte
}

// page function
func loadConfigPage(mwp *mainWindow) {
	config := &mwp.config.AutoCheck

	box := ui.NewVerticalBox()

	lConfigFilePath := ui.NewLabel(os.Args[0] + "\n")

	cAutoCheck := ui.NewCheckbox("Auto Check")
	cAutoCheck.SetChecked(config.Enable)
	cAutoCheck.OnToggled(func(*ui.Checkbox) {
		config.Enable = !config.Enable
		fmt.Printf("%v\nSuccess toggle config.Enable\n", config.Enable)
	})

	cMd5 := ui.NewCheckbox("MD5")
	cMd5.SetChecked(config.Md5)
	cMd5.OnToggled(func(*ui.Checkbox) {
		config.Md5 = !config.Md5
		fmt.Printf("%v\nSuccess toggle config.Md5\n", config.Md5)
	})

	cSha1 := ui.NewCheckbox("SHA1")
	cSha1.SetChecked(config.Sha1)
	cSha1.OnToggled(func(*ui.Checkbox) {
		config.Sha1 = !config.Sha1
		fmt.Printf("%v\nSuccess toggle config.Sha1\n", config.Sha1)
	})

	cSha256 := ui.NewCheckbox("SHA256")
	cSha256.SetChecked(config.Sha256)
	cSha256.OnToggled(func(*ui.Checkbox) {
		config.Sha256 = !config.Sha256
		fmt.Printf("%v\nSuccess toggle config.Sha256\n", config.Sha256)
	})

	bSave := ui.NewButton("Save")
	bSave.OnClicked(func(*ui.Button) {
		save(mwp.config)
	})

	bCheck := ui.NewButton("Check")
	bCheck.OnClicked(func(*ui.Button) {
		loadResultPage(mwp)
	})

	checkboxBox := ui.NewVerticalBox()
	checkboxBox.Append(cAutoCheck, false)
	checkboxBox.Append(cMd5, false)
	checkboxBox.Append(cSha1, false)
	checkboxBox.Append(cSha256, false)

	buttonBox := ui.NewHorizontalBox()
	buttonBox.Append(bSave, false)
	buttonBox.Append(bCheck, false)

	box.Append(lConfigFilePath, false)
	box.Append(checkboxBox, false)
	box.Append(buttonBox, false)
	mwp.configPage = box
	mwp.window.SetChild(mwp.configPage)
	mwp.window.SetMargined(true)
}
func loadResultPage(mwp *mainWindow) {
	config := &mwp.config.AutoCheck

	box := ui.NewVerticalBox()

	lFileName := ui.NewLabel(mwp.fileName + "\n")

	lMd5 := ui.NewLabel("MD5: ")
	if config.Md5 {
		vMd5 := md5.Sum(mwp.fileBuf)
		lMd5.SetText("MD5: " + hex.EncodeToString(vMd5[:]))
	}

	lSha1 := ui.NewLabel("SHA1: ")
	if config.Sha1 {
		vSha1 := sha1.Sum(mwp.fileBuf)
		lSha1.SetText("SHA1: " + hex.EncodeToString(vSha1[:]))
	}

	lSha256 := ui.NewLabel("SHA256: ")
	if config.Sha256 {
		vSha256 := sha256.Sum256(mwp.fileBuf)
		lSha256.SetText("SHA256: " + hex.EncodeToString(vSha256[:]))
	}

	bConfig := ui.NewButton("Config")
	bConfig.OnClicked(func(*ui.Button) {
		loadConfigPage(mwp)
	})

	box.Append(lFileName, false)
	box.Append(lMd5, false)
	box.Append(lSha1, false)
	box.Append(lSha256, false)
	box.Append(bConfig, false)
	mwp.resultPage = box
	mwp.window.SetChild(mwp.resultPage)
	mwp.window.SetMargined(true)
}

// target file function
func getTargetFilesPath() []string {
	if paths := os.Getenv("NAUTILUS_SCRIPT_SELECTED_FILE_PATHS"); len(paths) != 0 {
		fmt.Println(paths)
		fmt.Printf("%v\nSuccess get files path(getenv)\n", paths)
		return strings.Split(paths[:len(paths)-1], " ")
	}
	if len(os.Args) > 1 {
		fmt.Println(os.Args)
		fmt.Printf("%v\nSuccess get files path(os.args)\n", os.Args[1:])
		return os.Args[1:]
	}
	return nil
}
func loadTargetFile(path string) ([]byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var b []byte
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		b = append(b, buf[:n]...)
	}
	fmt.Println("Success load target file")
	return b, nil
}

// window function
func setMainWindow() error {
	mwp := new(mainWindow)

	paths := getTargetFilesPath()
	if paths == nil {
		return errors.New("Files path are not found")
	}
	var err error
	mwp.fileName = paths[0]
	mwp.fileBuf, err = loadTargetFile(mwp.fileName)
	if err != nil {
		return err
	}

	mwp.window = ui.NewWindow("HashChecker", 200, 1, false)
	mwp.window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	if mwp.config, err = load(); err != nil {
		return err
	}

	if mwp.config.AutoCheck.Enable {
		loadResultPage(mwp)
	} else {
		loadConfigPage(mwp)
	}

	mwp.window.Show()
	fmt.Println("Success show window")
	return nil
}
func setErrorWindow(err error) {
	window := ui.NewWindow("Error", 0, 0, false)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Message:"), false)
	box.Append(ui.NewLabel(err.Error()), false)
	window.SetChild(box)
	window.Show()
}

func main() {
	err := ui.Main(func() {
		if err := setMainWindow(); err != nil {
			setErrorWindow(err)
		}
	})
	if err != nil {
		panic(err)
	}
}
