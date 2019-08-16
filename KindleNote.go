package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"net/http"
	"net/url"
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
	page_content:=""
	for k,_:=range data {
		title,_:=url.QueryUnescape(k)
		page_content+="<h1><a href='/note/"+title+"'>"+title+"</a></h1>"
		
	}
    fmt.Fprintf(w, page_content, r.URL.Path[1:])
}

func bookHandler(w http.ResponseWriter,r *http.Request){
	params:=(r.URL.Path[len("/note/"):])
	notes,ok:=data[params]
	content:= "<p><a href='/'>Back Home</a></p>"
	if !ok {
		fmt.Fprintf(w,content+"<p>Not Found:%q</p>",params)
		return
	}
	for _,v:=range notes{
		content+="<p>"+v.CreateTime+"</p><p>"+v.Content+"</p>"
	}
	fmt.Fprintf(w,content)
}

func main(){
	content,err := ioutil.ReadFile("/Users/heyuanfei/Documents/MyClipping.txt")
	if err !=nil {
		log.Fatal(err)
	}
	//匹配笔记内容
	re := regexp.MustCompile(`[\s\S]*?\r\n={10}\r\n`)
	titmeRe := regexp.MustCompile(`.*\r\n`)
	dateTimeRe:=regexp.MustCompile(` 添加于 (.*)\r\n`)
	contentRe:=regexp.MustCompile(`\r\n(.*)\r\n={10}`)
	contentReplaceRe:=regexp.MustCompile(`\r\n={10}`)
	titleAndAuthor:=regexp.MustCompile(`[\s\S]*\s\(`)
	titleReplaceWorld:=regexp.MustCompile(`[\s\(].*`)
	a:=re.FindAll(content,-1)
	count:=len(a)
	if count>0 {
		data= make(map[string][]Note)
	}
	for _,sourceContent:= range a {
		//		fmt.Printf("d---%s",content) 
		title:=titmeRe.Find(sourceContent)
		title=titleAndAuthor.Find(title)
		title=titleReplaceWorld.ReplaceAll(title, []byte{})
		// fmt.Printf("Title:%s\n",title)
		dateTime:=dateTimeRe.Find(sourceContent)
		// fmt.Printf("Date:%s\n",dateTime)
		findContent:=contentRe.Find(sourceContent)
		content:=contentReplaceRe.ReplaceAll(findContent,[]byte{} )
		// fmt.Printf("Content:%s\n",content)

		var key string=string(title)//url.QueryEscape(string(title))
		_,ok:=data[key]
		if ok==false {
			data[key] = []Note{}
		}
		note := Note{string(title),string(dateTime),string(content)}
		 data[key]=append(data[key],note)
	}
	
	 http.HandleFunc("/", handler)
	 http.HandleFunc("/note/", bookHandler)
	 log.Fatal(http.ListenAndServe(":3002", nil))
}
