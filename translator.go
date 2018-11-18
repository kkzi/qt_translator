package main

import (
	"flag"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"strings"
	"os/exec"
	"encoding/xml"
)

type Arguments struct {
	Input string
	Output string
	Dict string
	Qt string
}

type TS struct {
	Language string `xml:"language,attr"`
	Version string `xml:"version,attr"`
	Contexts []Context `xml:"context"`
}

type Context struct {
	XMLName xml.Name `xml:"context"`
	Name string `xml:"name"`
	Messages []Message `xml:"message"`
}

type Message struct {
	Source string `xml:"source"`
	Trans Translation `xml:"translation"`
}

type Translation struct {
	Type string `xml:"type,attr"`
	Text string `xml:",innerxml"`
}


var args = Arguments{}
var dict = map[string]string{}
var todoList []string
var tsFile = "zh_cn.ts"
var tsDone = TS{}

func main() {
	log.Println(os.Args)
	help := flag.Bool("help", false, "help")
	flag.StringVar(&args.Input, "input", "../../Src" ,"input source path")
	flag.StringVar(&args.Output,"output", "../../Src/Client/translations/zh_cn.qm" ,"output qm file path")
	flag.StringVar(&args.Dict, "dict", "zh_dict.json" ,"translation dict file")
	flag.StringVar(&args.Qt, "qt", "D:/local/Qt/Qt5.9.4/5.9.4/msvc2015_64/bin/" ,"qt home path")
	flag.Parse()
	log.Println(args)

	if *help {
		flag.Usage()
		return
	}

	_, err := os.Stat(args.Qt)
	if err != nil {
		log.Fatal("qt path ", args.Qt, " not exists")
	}
	_, err = os.Stat(args.Dict)
	if err != nil {
		log.Fatal("dict file ", args.Dict, " not exists ", err)
	}
	_, err = os.Stat(args.Input)
	if err != nil {
		log.Fatal("input path ", args.Input, " not exists")
	}

	loadDict()
	runCommand()

	log.Println("translation completed")
}

func loadDict() {
	text, err := ioutil.ReadFile(args.Dict)
	if err != nil {
		log.Fatal("read dict file ", args.Dict, " failed:", err)
	}
	json.Unmarshal(text, &dict)
}

func runCommand() {
	generateTsFile()
	createTodoFile()
	createQmFile()
}

func generateTsFile() {
	cmd := exec.Command(filepath.Clean(args.Qt + "/lupdate.exe"), args.Input, "-silent", "-recursive", "-ts", tsFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	translateTsFile()
}

func translateTsFile() {
	text, err := ioutil.ReadFile(tsFile)
	if err != nil {
		log.Fatal(err)
	}
	ctx := TS{}
	xml.Unmarshal(text, &ctx)
	for _, c := range ctx.Contexts {
		for i, m := range c.Messages {
			key := strings.ToLower(m.Source)
			if value, ok := dict[key]; ok {
				value = strings.Replace(value, "<", "&lt;", -1)
				value = strings.Replace(value, ">", "&gt;", -1)
				m.Trans.Type = "finished"
				m.Trans.Text = value
				c.Messages[i] = m
			} else {
				todoList = append(todoList, key)
			}
		}
		tsDone.Contexts = append(tsDone.Contexts, c)
	}
}

func createTodoFile() {
	var lines []string
	for _, it := range todoList {
		lines = append(lines, `"` + it + `":""`)
	}
	text := strings.Join(lines, ",\n")
	if err:=ioutil.WriteFile("todo.txt", []byte(text), os.ModePerm); err != nil {
		log.Println("create todo.txt failed: ", err)
	}
}

func createQmFile() {
	tsDone.Language = "zh_CN"
	tsDone.Version = "2.1"

	text, err := xml.MarshalIndent(tsDone, "", "    ")
	if err != nil {
		log.Fatal("create ts file failed: ", err)
	}
	if err := ioutil.WriteFile("done.ts", []byte(string(text)), os.ModePerm); err != nil {
		log.Fatal("create ts file failed: ", err)
	}

	cmd := exec.Command(filepath.Clean(args.Qt + "/lrelease.exe"), "done.ts")
	output, err := cmd.CombinedOutput()
	log.Println(string(output))

	os.Rename("done.qm", args.Output)
}
