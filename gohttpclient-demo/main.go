package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var addr = ":8999"
func get()  {
	resp,err := http.Get("http://localhost"+addr+"/?a=b&c=d")
	if err !=nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}

func post()  {
	resp,err := http.Post("http://localhost"+addr, "application/x-www-form-urlencoded",
		strings.NewReader("a=b&c=d"))
	if err !=nil {
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}

func postform()  {
	resp,err := http.PostForm("http://localhost"+addr, url.Values{"a": {"b"}, "c": {"d"}})
	if err !=nil {
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}

func fileupload()  {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("uploadfile", "test.txt") //第一个字段名，第二个是参数名
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
	}
	srcFile, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("%Open source file failed: s\n", err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
	}
	writer.Close()
	resp,err := http.Post("http://localhost"+addr+"/file", writer.FormDataContentType(), buf)
	if err !=nil {
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}

func postjson()  {
	jsonStr :=[]byte(`{{"a":"b"},{"c":"d"}}`)
	req, err := http.NewRequest("POST", "http://localhost"+addr+"/json", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}

func main()  {
	go func() {
		http.HandleFunc("/", func(write http.ResponseWriter, r* http.Request){
			//r.ParseForm()

			file, header, err :=r.FormFile("uploadfile")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			nameParts := strings.Split(header.Filename, ".")
			ext := nameParts[1]
			savedPath := nameParts[0] + "_up."+ext
			f, err := os.OpenFile(savedPath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				panic(err)
			}
			fmt.Printf("method:%s,a:%s,c:%s,file:%s",r.Method,r.FormValue("a"),r.FormValue("c"),header.Filename)
		})
		http.HandleFunc("/json", func(write http.ResponseWriter, r* http.Request){
			b,_:= ioutil.ReadAll(r.Body)
			fmt.Println(string(b))
		})
		http.HandleFunc("/file", func(write http.ResponseWriter, r* http.Request){
			file, header, err :=r.FormFile("uploadfile")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			nameParts := strings.Split(header.Filename, ".")
			ext := nameParts[1]
			savedPath := nameParts[0] + "_up."+ext
			f, err := os.OpenFile(savedPath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				panic(err)
			}

		})
		http.ListenAndServe(addr,nil)
	}()
	time.Sleep(1e9)
/*	get()
	post()
	postform()
	postjson()
	fileupload()*/
	time.Sleep(100e9)
}
