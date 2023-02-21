package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/neurosnap/sentences"
)

var bookName string
var SentenceNum int = 1
var exeName string = "MMTrimV0.2.4"

func main() {
	createGUI()
}

func start() {
	locked = true
	defer func() { locked = false }()
	getArgumentes(currentPath)
	transformFile(currentPath)
	locked = false
	output(bookName + "转换结束。\r\n\r\n")
}

func getArgumentes(path string) {
	temp := strings.Split(path, "\\")
	filepath := temp[len(temp)-1]
	bookName = strings.Split(filepath, ".")[0]
	output("文件名/书名认定为: " + bookName)
}

func transformFile(filepath string) {
	output("开始读取文件……")
	content, err := readPdf(filepath) // Read local pdf file
	if err != nil {
		output("文件读取失败 " + err.Error())
		return
	}
	output("开始输出为json文件……")
	outputJson(content)

}

func outputJson(content string) {
	filePath := "./" + bookName + ".json"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		output("文件打开失败 " + err.Error())
		return
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
}

func readPdf(path string) (string, error) {
	_, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()
	output("文件共有" + strconv.Itoa(totalPage) + "页")
	// progress.MaxValue = totalPage
	var textBuilder bytes.Buffer
	result := ""

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		// output("正在处理第" + strconv.Itoa(pageIndex) + "页")
		// progress.Value = pageIndex
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		result, _ = p.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		// result = trimStr(result)
		// result = "???"

		textBuilder.WriteString(result)
	}

	output("文件处理尾声……")
	result, err = toFinalJSONString(trimStr(textBuilder.String()))
	if err != nil {
		return "", err
	}
	output := "[" + strings.TrimRight(result, ",") + "]"

	return output, err
}

func trimStr(str string) []tempEntry {
	result := ""
	b := []byte(TrainedData)
	training, _ := sentences.LoadTraining(b)
	tokenizer := sentences.NewSentenceTokenizer(training)
	sentences := tokenizer.Tokenize(str)
	SentenceNum = 0
	var entrySlice []tempEntry
	for _, s := range sentences {
		SentenceNum++
		result = strings.TrimSpace(s.Text)
		result = strings.ReplaceAll(result, "\n", "")
		entry := tempEntry{SentenceNum: SentenceNum, str: result}
		entrySlice = append(entrySlice, entry)
	}

	return entrySlice
}

func toFinalJSONString(te []tempEntry) (string, error) {
	var textBuilder bytes.Buffer
	var err error
	// textBuilder.WriteString(string("["))
	for _, s := range te {

		key := exeName + " " + bookName + " Sentence" + strconv.Itoa(s.SentenceNum)
		entry := Entry{Key: key, Original: s.str, Translation: ""}

		b, err := json.Marshal(entry)
		if err != nil {
			return "", err
		}

		textBuilder.WriteString(string(b))
		textBuilder.WriteString(",")

	}
	// textBuilder.WriteString("]")
	return textBuilder.String(), err
}
