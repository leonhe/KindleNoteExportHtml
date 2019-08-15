package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)
func main(){
	content,err := ioutil.ReadFile("/Users/heyuanfei/Documents/MyClippings.txt")
	if err !=nil {
		log.Fatal(err)
	}
	// a:="Hll"
	//	fmt.Printf("Hello Word: %s",content)
	// re := regexp.MustCompile(`\\ufeff([\s\S]*)\r\n={10}`)
	//匹配笔记内容
	re := regexp.MustCompile(`[\s\S]*?\r\n={10}\r\n`)
	titmeRe := regexp.MustCompile(`.*\r\n`)
	dateTimeRe:=regexp.MustCompile(` 添加于 (.*)\r\n`)
	contentRe:=regexp.MustCompile(`\r\n(.*)\r\n={10}`)
	contentReplaceRe:=regexp.MustCompile(`\r\n={10}`)
	a:=re.FindAll(content,-1)
	for _,sourceContent:= range a {
		//		fmt.Printf("d---%s",content) 
		title:=titmeRe.Find(sourceContent)
		fmt.Printf("Title:%s\n",title)
		dateTime:=dateTimeRe.Find(sourceContent)
		fmt.Printf("Date:%s\n",dateTime)
		findContent:=contentRe.Find(sourceContent)
		content:=contentReplaceRe.ReplaceAll(findContent,[]byte{} )
		fmt.Printf("Content:%s\n",content)
		

	}
	

}
