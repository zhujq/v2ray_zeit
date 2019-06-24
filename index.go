package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var url, realhost string

	switch r.URL.Path{
		case "/":
			fmt.Fprintf(w, "Welcome you,your info:\r\n")
			fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
			fmt.Fprintf(w, "URL:\r\n")
			fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
			fmt.Fprintf(w, "SCHEME:"+r.URL.Scheme+"\r\n")
			fmt.Fprintf(w, "HOST:"+r.URL.Host+"URL-End\r\n")
    		fmt.Fprintf(w, "Proto:"+r.Proto+"\r\n")
			fmt.Fprintf(w, "HOST:"+r.Host+"\r\n")
			fmt.Fprintf(w, "RequestUrl:"+r.RequestURI+"\r\n")
			return

		case "/google/":    //google入口
			url = "https://www.google.com"
			realhost = "www.google.com"
			  
    	case "/youtube/":   //youtube入口
			url = "https://www.youtube.com"
			realhost = "www.youtube.com"
		
		case "/search":     //google search入口，由于暂时无法带上真实主机名导致
            url = "http://www.google.com" + r.URL.String() 
            realhost = "www.google.com"

		// case "/url":
        // 		url = "http://www.google.com" + r.URL.String() 
			    	

        default:    //  经google、youtube入口后重新返回的网址的处理，分离出真实主机名称 
    	 	var str string
			str = r.URL.String()
			str = strings.TrimLeft(str,"/")
			realhost = string([]byte(str)[0:strings.Index(str,"/")])  //去掉首位的/后截取host
			fmt.Println(realhost)
    	 	if realhost == ""{
				fmt.Fprintf(w, "Failed to handle RequestUrl:"+str+"\r\n")
				return
			}
			if r.URL.Scheme == ""{
				r.URL.Scheme = "http"
			}
			url = r.URL.Scheme + "://" + str
        	
	}

	

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url, nil)
	req.Header = r.Header
	if err != nil {
        panic(err)
    }

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
        	
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
   		 //	w.Header().Set("content-type", "text/html;charset=utf-8")
   	for k, _ := range resp.Header{
   		w.Header().Set(k,resp.Header.Get(k))
   	}
	
	if strings.Contains(string(resp.Header.Get("content-type")),"text/html"){  

		fmt.Println("start to match")
		tempstr := string(body)
		olds := "<a href=\"/"
		news := "<a href=\"https://v2ray.14065567.now.sh/" + realhost + "/"
		tempstr = strings.ReplaceAll(string(body),olds,news)
		
		if len(tempstr) != len(string(body)){
			fmt.Println("matched")
		}

		olds = "src=\"/"
		news = "src=\"https://v2ray.14065567.now.sh/" + realhost+ "/"
		tempstr =  strings.ReplaceAll(tempstr,olds,news)

		olds = "href=\"http://"
		news = "href=\"https://v2ray.14065567.now.sh/" 
		tempstr =  strings.ReplaceAll(tempstr,olds,news)

		olds = "href=\"https://"
		news = "href=\"https://v2ray.14065567.now.sh/" 
		tempstr =  strings.ReplaceAll(tempstr,olds,news)

		olds = "<meta content=\"https://"
		news = "<meta content=\"https://v2ray.14065567.now.sh/" 
		tempstr =  strings.ReplaceAll(tempstr,olds,news)

		olds = "<meta content=\"/"
		news = "<meta content=\"https://v2ray.14065567.now.sh/" + realhost + "/"
		tempstr =  strings.ReplaceAll(tempstr,olds,news)
		
		body = []byte(tempstr)
	}
	fmt.Println(r.Method," URL:"+url," RealHost:",realhost,resp.Header.Get("content-type"))		
			
    w.Write([]byte(body))

          
    //	fmt.Fprintf(w,string(body))  
}
