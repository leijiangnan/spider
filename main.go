package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

// <img alt="外媒：美国法院正式撤销孟晚舟“银行欺诈”等指控" src="https://imagecloud.thepaper.cn/thepaper/image/228/138/686.jpg?x-oss-process=image/resize,w_332" width="318" height="182">

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)
	if err != nil {
		fmt.Printf("read content failed:%v\n", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("read content failed:%v", err)
	}

	doc.Find(`img[width="318"][height="182"]`).Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("alt")
		fmt.Printf("Review %d: %s\n", i, val)
	})

}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%d\n", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)

}

func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		fmt.Printf("fetch error:%v\n", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
