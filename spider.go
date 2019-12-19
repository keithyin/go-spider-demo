package main

/**  不知道为啥 tabsize = 8
1. 获取 URL
2. 爬数据
3. 处理数据，碰到 URL，放到之前的 URL List里面
4. 如何处理数据
*/

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func dododo() {
	var i int = -1
	cookie := GetCookie()
	client := &http.Client{}
	out_file, err := os.OpenFile("post.txt", os.O_WRONLY|os.O_CREATE, 666)

	if err != nil {
		panic(fmt.Sprintln("open file error: ", err))
	}
	defer out_file.Close()
	writer := bufio.NewWriter(out_file)

	for {
		i++
		url, page_idx := GenerateUrl(i)
		response := RequestAndResponse(client, url, cookie, page_idx)
		body_string := GetResponseBody(response)
		response.Body.Close()

		need_redirect := !CheckTitle(body_string)
		if need_redirect {
			redirect_url := ExtractUrlFromJs(body_string)
			response = RequestAndResponse(client, redirect_url, cookie, page_idx)
			body_string = GetResponseBody(response)
			response.Body.Close()
		}

		file_handle, _ := os.Create(fmt.Sprintf("file-%d.html", page_idx))
		file_handle.WriteString(body_string)
		post_info := ProcessBodyStringInfo(body_string)
		_, err := writer.WriteString(post_info + "\r\n")
		if err != nil {
			panic(fmt.Sprintln("writeString error: ", err))
		}
		if i > 0 {
			break
		}

	}
	writer.Flush()
}

func main() {
	dododo()
}
