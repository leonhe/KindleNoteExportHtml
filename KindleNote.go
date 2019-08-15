package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"net/http"
)

type Page struct{
	Title string	
	Body []byte
}
type Note struct{
	Title string
	CreateTime string
	Content string
}

type Notes struct{
	Id string
	data Note;
}

var data map[string][]Note;

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main(){
	content,err := ioutil.ReadFile("/Users/heyuanfei/Documents/MyClippings.txt")
	if err !=nil {
		log.Fatal(err)
	}
	//匹配笔记内容
	re := regexp.MustCompile(`[\s\S]*?\r\n={10}\r\n`)
	titmeRe := regexp.MustCompile(`.*\r\n`)
	dateTimeRe:=regexp.MustCompile(` 添加于 (.*)\r\n`)
	contentRe:=regexp.MustCompile(`\r\n(.*)\r\n={10}`)
	contentReplaceRe:=regexp.MustCompile(`\r\n={10}`)
	a:=re.FindAll(content,-1)
	count:=len(a)
	if count>0 {
		data= make(map[string][]Note)
	}
	for _,sourceContent:= range a {
		//		fmt.Printf("d---%s",content) 
		title:=titmeRe.Find(sourceContent)
		//fmt.Printf("Title:%s\n",title)
		dateTime:=dateTimeRe.Find(sourceContent)
		// fmt.Printf("Date:%s\n",dateTime)
		findContent:=contentRe.Find(sourceContent)
		content:=contentReplaceRe.ReplaceAll(findContent,[]byte{} )
		// fmt.Printf("Content:%s\n",content)
		
		notes:=data[string(title)]
		fmt.Println(notes)
		note := Note{string(title),string(dateTime),string(content)}
		data[string(title)]=append(notes,note)
	}
	
	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(":3001", nil))
}
