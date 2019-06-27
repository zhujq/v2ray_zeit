package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
	"bytes"
	"compress/gzip"
)
const zhost string = `https://v2ray.14065567.now.sh/`

func Handler(w http.ResponseWriter, r *http.Request) {
	var (
		url =``
		realhost =``
	)

	switch r.URL.Path{
		case `/`:
			fmt.Fprintf(w, "Welcome you,your info:\r\n")
			fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
			fmt.Fprintf(w, "URL:\r\n")
			fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
			fmt.Fprintf(w, "SCHEME:"+r.URL.Scheme+"\r\n")
			fmt.Fprintf(w, "HOST:"+r.URL.Host+"URL-End\r\n")
    		fmt.Fprintf(w, "Proto:"+r.Proto+"\r\n")
			fmt.Fprintf(w, "HOST:"+r.Host+"\r\n")
			fmt.Fprintf(w, "RequestUrl:"+r.RequestURI+"\r\n")
			for k,_ := range r.Header {
			    fmt.Fprintf(w, k+""+r.Header.Get(k)+"\r\n")
			}
			return

		case `/google/`:    //google入口
			url = `https://www.google.com`
			realhost = `www.google.com`
			  
    	case `/youtube/`:   //youtube入口
			url = `https://www.youtube.com`
			realhost = `www.youtube.com`
		
		case `/search`:     //google search入口，由于暂时无法带上真实主机名导致
            url = `http://www.google.com` + r.URL.String() 
			realhost = `www.google.com`	
		
		
		case `/watch`:   //youtube入口
			url = `https://www.youtube.com`+ r.URL.String() 
			realhost = `www.youtube.com`

        default:    //  经google、youtube入口后重新返回的网址的处理，分离出真实主机名称 
    	 	var str string
			str = r.URL.String()
			str = strings.TrimLeft(str,`/`)
			realhost = string([]byte(str)[0:strings.Index(str,`/`)])  //去掉首位的/后截取host
			if realhost == `xjs` {             //google的xjs目录暂时无法带上www.google.com
				realhost = `www.google.com`
				str = realhost + `/` + str
			}

			if realhost == `youtubei` || realhost == `yts` {             //youtube的youtubei yts目录暂时无法带上www.youtube.com
				realhost = `www.youtube.com`
				str = realhost + `/` + str
			}
		
    	 	if realhost == ``{
				fmt.Fprintf(w, `Failed to handle RequestUrl:`+str+`\r\n`)
				return
			}
			if r.URL.Scheme == ``{
				r.URL.Scheme = `http`
			}
			url = r.URL.Scheme + `://` + str
        	
	}

	if toredirect(realhost){             //判断如果是国内域名，则指示重定向
		fmt.Println(r.Method,` URL:`+url,` LocalRealHost:`,realhost)	
		http.Redirect(w, r, `http://`+ realhost, 307)
		return
	}
	
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url, nil)

	req.Header = r.Header     //删除请求头压缩选项，否则无法对返回的文本的链接内容进行处理,20190625 调用compress/gzip进行压缩和解压缩,且只用gzip
	if  strings.Contains(string(req.Header.Get(`Accept-Encoding`)),`gzip`){
		req.Header.Set(`Accept-Encoding`,`gzip`)  
	}else {
		req.Header.Del(`Accept-Encoding`)   
	}
	if err != nil {
        panic(err)
    }

	resp, err := client.Do(req)
	
	fmt.Println(r.Method,` URL:`+url,`resp len,status,type,Enc:`,strconv.FormatInt(resp.ContentLength,10),resp.Status,resp.Header.Get(`content-type`),resp.Header.Get(`Content-Encoding`))	//记录访问记录

    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
        	
    for k, _ := range resp.Header{
   		w.Header().Set(k,resp.Header.Get(k))
	}
	
	body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
        panic(err)
	}
	
	if resp.StatusCode == 200 && strings.Contains(string(resp.Header.Get(`content-type`)),`text/html`){   //只有当返回200和文本类型时进行链接处理
		if len(body) == 0 {
			fmt.Println(`resp is empty`)
			return
		}
		
		if resp.Header.Get(`Content-Encoding`) == `gzip`{
			body,err = gzipdecode(body)
			if err != nil {
				panic(err)
			}
		}

		
		matching := false
		modifiedrsp := []byte{}
		tomodifystr := ``
		for _,v := range body {
			if string(v) == `<` {
				matching = true
				tomodifystr += string(v)
			}else if string(v) == `>` {
				matching = false
				tomodifystr += string(v)
				tomodifystr = modifylink(tomodifystr,realhost)
				for _,vv := range tomodifystr {
					modifiedrsp = append(modifiedrsp,byte(vv))
				}
				tomodifystr = ``
			}else{
				if matching == false {
					modifiedrsp = append(modifiedrsp,byte(v))
				}else{
					tomodifystr += string(v)
				}
			}
		}
		body = bytes.ReplaceAll(modifiedrsp,[]byte(`url(https://`),[]byte(`url(` + zhost ))
		body = bytes.ReplaceAll(body,[]byte(`url('https://`),[]byte(`url('` +zhost ))
		body = bytes.ReplaceAll(body,[]byte(`url(/`),[]byte(`url(` +zhost + realhost + `/`))
		body = bytes.ReplaceAll(body,[]byte(`s='/images`),[]byte(`s='` +zhost + realhost + `/images`))
		body = bytes.ReplaceAll(body,[]byte(`https:\/\/`),[]byte(`https:\/\/'` +zhost + `\/`))
		
		
		if resp.Header.Get("Content-Encoding") == "gzip" {    //如果resp指示压缩，还需要对解开的处理后的内容重新压缩
			body,err = gzipencode(body)
			if err != nil {
				panic(err)
			}

		}

	}

	w.Write([]byte(body))        
}

func modifylink(s string,realhost string) string{
	if s == ""{
		return s
	}

	tempstr := s
	olds :=``
	news :=``
																
	olds = `href="https://`                                     //先改https，否则会重复改
	news = `href="` + zhost
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `href="http://`
	news = `href="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `href="//`                                          // href="//  后是绝大路径
	news = `href="` +zhost  
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `href="/`                                                     // href="/ 后是相对路径
	news = `href="` + zhost + realhost + "/"
	tempstr = strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}
	
	olds = `<a href="https://`                                     //先改https，否则会重复改
	news = `<a href="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `<a href="http://`                                     
	news = `<a href="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `<a href="//`                                     
	news = `<a href="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `<a href="/`
	news = `<a href="` + zhost + realhost + "/"
	tempstr = strings.Replace(tempstr,olds,news,-1)
	if len(tempstr) > len(s){
		return tempstr
	}

	olds = `src="https://`
	news = `src="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	

	olds = `src="http://`
	news = `src="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	

	olds = `src="//`
	news = `src="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	

	olds = `src="/`
	news = `src="` +zhost  + realhost+ "/"
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `srcset="https://`
	news = `srcset="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	

	olds = `srcset="http://`
	news = `srcset="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	

	olds = `srcset="//`
	news = `srcset="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `srcset="/`
	news = `srcset="` + zhost  + realhost+ "/"
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	

	olds = `<meta content="https://`                             //先改https，否则会重复改
	news = `<meta content="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `<meta content="http://`
	news = `<meta content="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)

	olds = `<meta content="//`
	news = `<meta content="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `<meta content="/`
	news = `<meta content="` + zhost  + realhost + "/"
	tempstr =  strings.Replace(tempstr,olds,news,-1)

	olds = `<iframe src="https://`                           //先改https，否则会重复改
	news = `<iframe src="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `<iframe src="http://`
	news = `<iframe src="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `itemtype="https://`                           //先改https，否则会重复改
	news = `itemtype="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
	
	olds = `itemtype="http://`
	news = `itemtype="` + zhost 
	tempstr =  strings.Replace(tempstr,olds,news,-1)
			
	return tempstr

}

func toredirect(s string) bool{
	if strings.HasSuffix(s, ".cn"){
		return true
	}
	localurls := []string{"baidu","taobao","sina","163.com","tmall","jd.com","sohu","qq.com","ifeng.com","qunae.com","toutiao.com","alipay.com","ctrip.com","weibo.com","zhihu"}
	for _, localurl := range localurls {
		if strings.Contains(s,localurl){
			return true
		}
	}
	
	return false
}

func gzipencode(in []byte) ([]byte, error) {
    var (
        buffer bytes.Buffer
        out    []byte
		err    error
       	)
        writer := gzip.NewWriter(&buffer)
        _, err = writer.Write(in)
        if err != nil {
            writer.Close()
            return out, err
        }
        err = writer.Close()
        if err != nil {
            return out, err
        }
   	     return buffer.Bytes(), nil
}
	
func gzipdecode(in []byte) ([]byte, error) {
        reader, err := gzip.NewReader(bytes.NewReader(in))
        if err != nil {
            var out []byte
            return out, err
        }
        defer reader.Close()
        return ioutil.ReadAll(reader)
}