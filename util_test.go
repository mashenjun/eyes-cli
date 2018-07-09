package main

import (
	"fmt"
	"qiniupkg.com/x/log.v7"
	"os/exec"
	"testing"
	"gopkg.in/cheggaaa/pb.v1"
	"time"
	"os"

	"strings"
	"path"
)

func TestTouchFile(t *testing.T) {
	name := ".track"
	_, err := TouchFile(name)
	log.Println(err)
}

func TestNewUuid(t *testing.T) {

	uid, err := Gen(12)
	log.Println(err)
	log.Println(len(uid))
}

func TestNewErrorLog(t *testing.T) {
	file := NewErrorLog("error_qiniu_testGroup")

	// don't forget to close it
	defer file.Close()

	// assign it to the standard logger
	ErrorLog.Println("{zzzz}")
	ErrorLog.Output(10, `{"xxx":123,"yy":43}`)

}

func TestExecCommand(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//commandDir := filepath.Dir(dir)
	//log.Println(commandDir)
	//cmd := exec.Command("tail", "-1","error_qiniu_testGroup")
	//cmd.Dir = commandDir
	//out, err := cmd.Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	commandDir := path.Join(dir, "new")
	wc := exec.Command("ls", )
	wc.Dir = commandDir
	out, err := wc.Output()
	if err != nil {
		log.Fatal(err)
	}
	outStr := strings.TrimRight(string(out), " \n")
	rlt:= strings.Split(outStr, "\n")
	f, err := os.OpenFile("metadata.jsonl", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	for _ ,item := range rlt{
		tokens := strings.Fields(item)
		name := strings.Split(tokens[1],".")[0]
		idCard := tokens[0]
		text := fmt.Sprintf(`{"name":"%v", "sex":0, "id_card_num":"%v", "images":["%v"]}`,name,idCard,"new/"+item)
		if _, err = f.WriteString(text+"\n"); err != nil {
			panic(err)
		}
	}
}

func TestNewfileUploadRequest(t *testing.T) {
	_, err := newfileUploadRequest("http://localhost:8080", map[string]string{}, "file","new/0658 葛利亚.jpg")
	log.Println(err)

}


func TestLoadingBarA(t *testing.T){
	count := 10000
	bar := pb.New(count)
	bar.SetWidth(100)
	bar.ShowTimeLeft = false
	bar.Start()
	bar.Width = 100
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")
}


func TestDotEffect(t *testing.T) {
	f, err:= os.Stat("new/1295 应开翔.JPG")
	log.Println(err)
	log.Println(f.Size())
}

func TestTotalLine(t *testing.T) {
	lines,err := TotalLine("metadata.jsonl")
	log.Println(err)
	log.Println(lines)
}
