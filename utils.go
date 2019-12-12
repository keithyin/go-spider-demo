package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

type JsonFmt struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type JsFmt struct {
	Html string `json:"html"`
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func Today() time.Time {
	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		panic(fmt.Sprintf("time parse error %+v", err))
	}
	year, month, day := today.Date()
	today = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return today
}

func ProcessBodyStringInfo(body_string string) (res string) {
	//TODO: 将 body_string 转成 io.Reader
	// reader := io.NewReader(body_string)
	// doc, err := goquery.NewDocumentFromReader(reader)

	today := Today()

	reader := bytes.NewReader([]byte(body_string))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Println("NewDocumentFromReader error=", err)
		panic("")
	}
	repost_reg := regexp.MustCompile(`转发微博`)
	doc.Find("div.WB_feed_detail.clearfix").Each(func(i int, s *goquery.Selection) {
		s.Find("div.WB_from.S_txt2 a").First().Each(func(i int, sub_s *goquery.Selection) {
			// "a:first" 这种标识不可用，可以尝试 .First()然后再 .Each()
			val, exist := sub_s.Attr("title")
			if exist {
				post_day, err := time.Parse("2006-01-02 15:04", val)

				if err != nil {
					panic(fmt.Sprintf("time parse error: %+v", err))
				}
				if post_day.After(today) {
					post := s.Find(".WB_text.W_f14").Text()
					post = strings.Trim(post, "\r\n \t")
					find_repost := repost_reg.FindString(post)
					// 如果是转发微博，读取转发微博的内容
					if find_repost != "" {
						div_text_node := s.Find(".WB_feed_expand .WB_expand.S_bg1 .WB_text")
						post = div_text_node.Text()
						post = strings.Trim(post, "\r\n \t")
					}

					img_list := make([]string, 0, 5)
					if find_repost != "" {
						s.Find(".WB_feed_expand .media_box img").Each(func(i int, img_s *goquery.Selection) {
							img_src, _ := img_s.Attr("src")
							img_list = append(img_list, img_src)
						})
					} else {
						s.Find(".media_box img").Each(func(i int, img_s *goquery.Selection) {
							img_src, _ := img_s.Attr("src")
							img_list = append(img_list, img_src)
						})
					}
					res = fmt.Sprintf("%s\r\n %v\r\n ------------\r\n", post, img_list)

				}
			}
		})
	})
	return
}

func ExtractUrlFromJs(str string) string {
	reg := regexp.MustCompile(`location.replace\((.*)\);`)
	res := reg.FindAllStringSubmatch(str, -1)
	if res == nil {
		return ""
	}
	return strings.Trim(res[0][1], "\"")
}

// must call defer response.Body.Close() after call this function
func RequestAndResponse(client *http.Client, url string, cookie string, page_idx int) *http.Response {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("request error: page=%d, error message=%s \n", page_idx, err)
		fmt.Println("error url = ", url)
		panic("byebye")
	}
	request.Header.Set("Cookie", cookie)
	response, err2 := client.Do(request)
	if err2 != nil {
		fmt.Printf("response error: page=%d, error message=%s \n", page_idx, err2)
		panic("byebye")
	}
	return response
}

// Append something to the json respond!!!!
func GetResponseBody(response *http.Response) string {
	buffer := make([]byte, 4*1024)
	var body_string string
	var n int
	for {
		n, _ = response.Body.Read(buffer)
		if n == 0 {
			break
		}
		body_string += string(buffer[:n])
	}
	content_type := response.Header["Content-Type"][0]
	content_type = strings.Split(content_type, ";")[0]
	// fmt.Printf(":%s:\r\n", content_type)
	if content_type == "application/json" {
		var jsonfmt JsonFmt
		err := json.Unmarshal([]byte(body_string), &jsonfmt)
		if err != nil {
			fmt.Println("convert to json error : ", err)
			panic("")
		}
		body_string = jsonfmt.Data
	} else {
		// WB use javascipt DOM technology to contruct html page,
		// this function extract html code from js code
		body_string = ProcessFirstPage(body_string)
	}

	// do something for the first page !!!!!!!!!!!!!!!
	return body_string
}

// true, the title is ok, false, need to do a redirect
func CheckTitle(body_string string) bool {
	reg := regexp.MustCompile(`<title>(.*)</title>`)
	matches := reg.FindAllStringSubmatch(body_string, -1)
	if matches == nil {
		return true
	}

	title := matches[0][1]
	title = strings.Replace(title, " ", "", -1)
	title = ConvertToString(title, "gbk", "utf-8")

	if title == "新浪通行证" {
		return false
	}
	return true
}

/*
<script>FM.View({"html": "<div></div>"}) ----> <div></div>
*/
func ProcessFirstPage(body_string string) string {
	ret_string := body_string
	reg := regexp.MustCompile(`<script>FM.view\((.*)\)</script>`)
	matched := reg.FindAllStringSubmatch(body_string, -1)
	for _, v := range matched {
		var tmp JsFmt
		// fmt.Println(v[1])
		// if i > 0 {
		// 	return ret_string
		// }
		err := json.Unmarshal([]byte(v[1]), &tmp)
		if err != nil {
			panic(fmt.Sprintln("json to JsFmt error: ", err))
		}

		if tmp.Html != "" {
			ret_string = strings.Replace(ret_string, v[0], tmp.Html, 1)
		}
	}
	return ret_string
}
