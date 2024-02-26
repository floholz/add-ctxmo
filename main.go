package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	. "github.com/floholz/add-ctxmo/src"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var outPath string

var entryRegKey *widget.Entry
var entryExePath *widget.Entry
var entryCtxMOText *widget.Entry

var labelSuccess *canvas.Text
var labelFailed *canvas.Text

var fileDialog *dialog.FileDialog
var fileButton *widget.Button

var checkDontExec *widget.Check
var buttonCreate *widget.Button

func main() {
	initOutPath()

	a := app.New()
	w := a.NewWindow("Add Contextmenu Options")
	w.Resize(fyne.NewSize(675, 745))
	w.SetIcon(theme.SettingsIcon())

	if IsAdmin() {
		w.SetTitle(w.Title() + " [Admin]")

		entryRegKey = widget.NewEntry()
		entryRegKey.OnChanged = onInputChanged
		entryRegKey.Validator = NoEmptyStringValidator

		entryExePath = widget.NewEntry()
		entryExePath.OnChanged = onInputChanged
		entryExePath.Validator = validation.NewAllStrings(NoEmptyStringValidator, ValidPathValidator)

		entryCtxMOText = widget.NewEntry()
		entryCtxMOText.OnChanged = onInputChanged
		entryCtxMOText.Validator = NoEmptyStringValidator

		fileDialog = dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			fpath, _ := Pathify(closer.URI().Path())
			fmt.Println(fpath)
			entryExePath.SetText(fpath)
		}, w)
		home, err := os.UserHomeDir()
		if err == nil {
			uri := storage.NewFileURI(home + "/Downloads")
			luri, err := storage.ListerForURI(uri)
			if err == nil {
				fileDialog.SetLocation(luri)
			}
		}
		fileButton = widget.NewButtonWithIcon("", theme.FileApplicationIcon(), func() {
			fileDialog.Show()
		})

		checkDontExec = widget.NewCheck("Don't execute. Only create registry file.", onCheckExecChanged)

		labelSuccess = canvas.NewText("Successfully created Contextmenu Option!", theme.SuccessColor())
		labelSuccess.Hide()
		labelFailed = canvas.NewText("Failed to create Contextmenu Option!", theme.ErrorColor())
		labelFailed.Hide()

		buttonCreate = widget.NewButtonWithIcon("Create Contextmenu-Option", theme.ContentAddIcon(), onButtonClick)

		w.SetContent(container.NewVBox(
			widget.NewForm(
				widget.NewFormItem("Registry Key", entryRegKey),
				widget.NewFormItem("EXE Filepath", container.NewBorder(nil, nil, nil, fileButton, entryExePath)),
				widget.NewFormItem("Ctxmenu Text", entryCtxMOText),
			),
			checkDontExec,
			buttonCreate,
			layout.NewSpacer(),
			labelSuccess,
			labelFailed,
		))
	} else {
		w.Resize(fyne.NewSize(800, 245))
		w.SetContent(container.NewGridWithRows(4,
			layout.NewSpacer(),
			container.NewCenter(canvas.NewText("This application must be run with elevated permissions. Please relaunch the application as administrator.", theme.WarningColor())),
			container.NewGridWithColumns(3,
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Relaunch as Admin", theme.WarningIcon(), RelaunchWithElevatedPerms),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		))
	}

	// Launch window
	w.ShowAndRun()
}

func onInputChanged(input string) {
	resetLabels()
}

func onCheckExecChanged(checked bool) {
	if checked {
		buttonCreate.SetText("Create Contextmenu-Option registry file")
	} else {
		buttonCreate.SetText("Create Contextmenu-Option")
	}
}

func onButtonClick() {
	if !validateInputs() {
		fail(nil)
		return
	}

	ctx := WinReg5Hbs{
		RegKey:    entryRegKey.Text,
		ExePath:   entryExePath.Text,
		CtxMOText: entryCtxMOText.Text,
	}

	/* [start] generate registry file  */
	res, err := GenerateRegFile(ctx)
	if err != nil {
		fail(err)
		return
	}

	regFile, err := Pathify(outPath + "/" + entryRegKey.Text + ".reg")
	if err != nil {
		fail(err)
		return
	}

	err = SaveToDisk([]byte(res), regFile)
	if err != nil {
		fail(err)
		return
	}
	/* [end] generate registry file */

	/* [start] generate registry delete file  */
	regDeleteString, err := GenerateRegDeleteFile(ctx)
	if err != nil {
		fail(err)
		return
	}

	regDeleteFile, err := Pathify(outPath + "/" + entryRegKey.Text + ".del.reg")
	if err != nil {
		fail(err)
		return
	}

	err = SaveToDisk([]byte(regDeleteString), regDeleteFile)
	if err != nil {
		fail(err)
		return
	}
	/* [end] generate registry delete file */

	if !checkDontExec.Checked {
		cmd := exec.Command("regedit.exe", "/s", regFile)
		err = cmd.Run()
		if err != nil {
			fail(err)
			return
		}
		fmt.Println("Context Menu Option created")
	}
	labelSuccess.Show()
}

func fail(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
	labelSuccess.Hide()
	labelFailed.Show()
}

func resetLabels() {
	labelSuccess.Hide()
	labelFailed.Hide()
}

func validateInputs() bool {
	if entryRegKey.Text == "" {
		fmt.Println("Invalid input [Registry Key]")
		return false
	}
	if entryExePath.Text == "" || ValidPathValidator(entryExePath.Text) != nil {
		fmt.Println("Invalid input [Path to EXE]")
		return false
	}
	if entryCtxMOText.Text == "" {
		fmt.Println("Invalid input [CtxMO Text]")
		return false
	}
	return true
}

func initOutPath() {
	envPath := os.Getenv("ADDCTXMO_PATH")
	if ValidPathValidator(envPath) == nil {
		outPath = envPath
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	outPath = filepath.Join(home, ".addctxmo")
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
