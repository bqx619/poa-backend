package biz

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func testAsync() {
	dc := make(chan int, 10)
	wg := sync.WaitGroup{}
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(dc chan <- int, j int) {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(j%10))
			dc <- j
			//wg.Done()
		}(dc, i)
	}
	go func() {
		c := 0
		for {
			select {
			case op := <-dc:
				c ++
				fmt.Println(op, c)
			}
		}
	}()
	wg.Wait()

}



func TestAsync(t *testing.T) {
	//testAsync()
	//a := time.Now().UnixMicro()
	//fmt.Println(RandomIPv6(10))
	//fmt.Println(time.UnixMicro(a))
	response, err := http.Get("https://www.coder.work/article/25431")
	if err != nil {
		return
	}
	defer response.Body.Close()
	reader, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}
	fmt.Println(reader.Find("title").Text())
	cnBody := GetStrCn(reader.Find("body").Text())
	fmt.Println(cnBody)
}
