package main

import (
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var locked bool = false
var outTE *walk.TextEdit
var currentPath string
var button PushButton = PushButton{
	Text:      "开始转换.jpg",
	OnClicked: guiStart,
}

// var lines int = 0

// var progress ProgressBar = ProgressBar{MinValue: 0}

func createGUI() {
	MainWindow{
		Title:   exeName,
		MinSize: Size{Width: 600, Height: 400},
		MaxSize: Size{Width: 600, Height: 400},
		Size:    Size{Width: 600, Height: 400},
		Layout:  VBox{},
		OnDropFiles: func(files []string) {
			currentPath = strings.Join(files, "\r\n")
			outTE.AppendText(strings.Join(files, "\r\n") + "\r\n")
		},
		Children: []Widget{
			TextEdit{AssignTo: &outTE, ReadOnly: true, Text: "" + exeName + " 作者:苦力Dora\r\n把文件拖放进来以开始。\r\n\r\n"},
			button,
			// progress,
		},
	}.Run()
}

func output(s string) {
	outTE.AppendText(s + "\r\n")
}

func guiStart() {
	if locked {
		return
	}
	if currentPath == "" {
		output("请选择文件再启动捏")
		return
	}
	outTE.SetText("")
	go start()
}
