package main

import (
	"context"
	"log"
	"time"
	"fmt"
	"strings"
	"github.com/chromedp/chromedp"
	"regexp"
)

func main(){
	const link = `https://online.carrefour.com.tw/zh/homepage`
	ctx, cancel := chromedp.NewExecAllocator(context.Background())
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var sel_1 = `#online-shop`
	var sel_2 = `#online-shop > ul`
	var All_Subject string
	
	if err := chromedp.Run(ctx,
		chromedp.Navigate(link),
		chromedp.WaitVisible(sel_1),
		chromedp.Click(sel_1),
		chromedp.WaitVisible(sel_2),
		chromedp.Text(`//*[@id="online-shop"]`,&All_Subject),
		chromedp.Sleep(5*time.Second),
	); err != nil {
		log.Println(err)
	} else {
		log.Println("抓到所有分類啦")
	}

	fmt.Println(All_Subject)
	split_string := strings.Split(All_Subject,"\n")
	url_table := get_all_subject_url(split_string)
	var content string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url_table[1]),
		chromedp.WaitVisible(`body > div.page > section > section > div > div.flright.list > div.hot-recommend.clearfix > div:nth-child(4) > div.desc-operation-wrapper`),
		chromedp.OuterHTML(`//*[@class="hot-recommend clearfix"]`,&content),
			); err != nil {
		log.Println(err)
	} else {
		log.Println("done")
	}

	content2 := strings.Split(content,"\n")
	url_img := filter_image_url(content2)
	content3 := strings.Split(content," ")
	product,price := filter_product_and_price(content3)
	for i:= range price{
		fmt.Println("Product: ",product[i])
		fmt.Println("Price: ",price[i])
		fmt.Println("IMG Url Source: ",url_img[i],"\n")
	}

}

func get_all_subject_url(content []string) []string{
	const base_url = `https://online.carrefour.com.tw/zh/`
	var url_table []string
	for i:= range content{
		string:= strings.TrimSpace(content[i])
		string = base_url + string
		url_table = append(url_table,string)
	}
	return url_table
}
func filter_product_and_price(content []string)([]string,[]string){
	var price []string
	var product []string
	var lens int
	reg := regexp.MustCompile(`[^\d]`)
	reg2 := regexp.MustCompile(`data-name="|"`)
	reg3 := regexp.MustCompile(`class="packageQty"|</?div>?|</span>|\s`)
	for j:= range content {
		if strings.Contains(content[j],"current-price"){
			price_string := reg.ReplaceAllString(content[j],"${1}")
			price = append(price,price_string)
		} else if strings.Contains(content[j],"data-name"){
			input_string := reg2.ReplaceAllString(content[j],"${1}")
			//input_string := strings.Trim(content[j],`data-name=`)
			if !IsContain(product,input_string){
				product = append(product,input_string)
			}
		}
		if strings.Contains(content[j],"packageQty"){
			lens = len(product)
			Qty := reg3.ReplaceAllString(content[j],"${1}")
			product[lens-1] = product[lens-1]+Qty
		}
	}
	return product,price
}

func filter_image_url(content []string)[]string{
	var clear_image_url []string
	reg := regexp.MustCompile(`src=`)
	for j:= range content {
		if strings.Contains(content[j],"m_lazyload"){
			new_strings := strings.Split(content[j]," ")
			for i:= range new_strings{
				if strings.Contains(new_strings[i],"src"){
					src_url := reg.ReplaceAllString(new_strings[i],"${1}")
					clear_image_url = append(clear_image_url,src_url)
				}
			}
		}
	}
	return clear_image_url
	
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}