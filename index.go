package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
		case "/google/":
        	resp, err := http.Get("http://www.google.com")
        	if err != nil {
            	panic(err)
        	}

        	defer resp.Body.Close()
        	
        	body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
         		panic(err)
    		}
   		 	w.Header().Set("content-type", "text/html;charset=utf-8")
            w.Write([]byte(body))

          
    		//	fmt.Fprintf(w,string(body))    
    	case "/images/":
    	 	var str string
    	 	str = r.URL.Path
    	 	fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
    	 	if strings.HasSuffix(str, "png"){
    	 		str = "http://www.google.com" + str
    	 		resp, err := http.Get(str)
    	 		if err != nil {
            	panic(err)
        		}

        		defer resp.Body.Close()
        	
        		body, err := ioutil.ReadAll(resp.Body)
    			if err != nil {
         			panic(err)
    			}
   		 		w.Header().Set("content-type", "image/png")
            	w.Write([]byte(body))

			}

        default:  // 
        	

	}
}
