package extensions

//Finally, we can make a GET call to a generic URL and query the HTML response using go-query syntax:
//{{ $webpage := fetch "https://www.google.it/" | html }}{{ dump $webpage}}

// func FetchHTML(url string) ([]*goquery.Selection, error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []*goquery.Selection{doc.Selection}, nil
// }

// func FilterHTML(selector string, selection *goquery.Selection) (*goquery.Selection, error) {
// 	selection.Find(selector).Each(func(i int, s *goquery.Selection) {
// 		// For each item found, get the title
// 		title := s.Find("a").Text()
// 		fmt.Printf("Review %d: %s\n", i, title)
// 	})
// }

// func HTML(selection *goquery.Selection) (string, error) {
// 	if selection != nil {
// 		return selection.Text(), nil
// 	}
// 	return "", nil
// }
