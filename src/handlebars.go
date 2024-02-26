package src

import (
	"github.com/flowchartsman/handlebars/v3"
	"strconv"
	"strings"
)

type WinReg5Hbs struct {
	RegKey    string
	ExePath   string
	CtxMOText string
}

func GenerateRegFile(ctx WinReg5Hbs) (string, error) {
	fpath, err := Pathify("./templates/win.reg5.hbs")
	if err != nil {
		return "", err
	}
	template, err := handlebars.ParseFile(fpath)
	if err != nil {
		return "", err
	}
	template.RegisterHelper("escape", func(str string) string {
		return strings.Trim(strconv.Quote(str), "\"")
	})
	result, err := template.Exec(ctx)
	if err != nil {
		return "", err
	}
	return result, nil
}

func GenerateRegDeleteFile(ctx WinReg5Hbs) (string, error) {
	fpath, err := Pathify("./templates/win.reg5.del.hbs")
	if err != nil {
		return "", err
	}
	template, err := handlebars.ParseFile(fpath)
	if err != nil {
		return "", err
	}
	result, err := template.Exec(ctx)
	if err != nil {
		return "", err
	}
	return result, nil
}
