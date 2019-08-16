package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"net/http"
	"net/url"
	"html/template"
)

type Page struct{
	Title string	
	Body template.HTML
}
type Note struct{
	Title string
	CreateTime string
	Content string
	Pos string
}

type Notes struct{
	Id string
	data Note;
}

var data map[string][]Note;

func handler(w http.ResponseWriter, r *http.Request) {
	page_content:="<h2>Books</h2><ul>"
	for k,_:=range data {
		title,_:=url.QueryUnescape(k)
		page_content+=("<li><a href='/note/"+title+"'>"+title+"</a></li>")
		
	}
	page_content+="</ul>"
	//fmt.Fprintf(w, page_content, r.URL.Path[1:])
	p:=&Page{Title:"Kindle Note Home",Body:template.HTML(page_content)}
	t,_:=template.ParseFiles("template/index.html")
	t.Execute(w,p)	
}

func bookHandler(w http.ResponseWriter,r *http.Request){
	params:=(r.URL.Path[len("/note/"):])
	notes,ok:=data[params]
	content:= "<h3>"+params+"</h3><p><a href='/'>Back Home</a></p>"
	if !ok {
		fmt.Fprintf(w,content+"<p>Not Found:%q</p>",params)
		return
	}
	for _,v:=range notes{
		content+="<div class=\"card mb-3\"><div class=\"card-body\">"
                content+="<h5 class=\"card-title\">"+v.Pos+"</h5>"
    content+="<p class=\"card-text\">"+v.Content+"</p>"
content+="<p class=\"card-text\"><small class=\"text-muted\">"+v.CreateTime+"</small></p>"
content+=" </div></div>"
		// content+="<p>"+v.CreateTime+"</p><p>"+v.Content+"</p>"
	}
	// fmt.Fprintf(w,content)
	p:=&Page{Title:params,Body:template.HTML(content)}
	t,_:=template.ParseFiles("template/index.html")
	t.Execute(w,p)
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
	posRe:=regexp.MustCompile(`\s#[\s\S]*\s\|`)
	posReplace:=regexp.MustCompile(`\|`)
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
		pos:=posRe.Find(sourceContent)
		pos = posReplace.ReplaceAll(pos,[]byte{})
		// fmt.Printf("Date:%s\n",dateTime)
		findContent:=contentRe.Find(sourceContent)
		content:=contentReplaceRe.ReplaceAll(findContent,[]byte{} )
		// fmt.Printf("Content:%s\n",content)

		var key string=string(title)//url.QueryEscape(string(title))
		_,ok:=data[key]
		if ok==false {
			data[key] = []Note{}
		}
		note := Note{string(title),string(dateTime),string(content),string(pos)}
		 data[key]=append(data[key],note)
	}
	
	 http.HandleFunc("/", handler)
	 http.HandleFunc("/note/", bookHandler)
	 log.Fatal(http.ListenAndServe(":3002", nil))
}
