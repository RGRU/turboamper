package turboamper

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Turbo gives you YandexTurbo-representation of html and its type
// If it cannot recognize your html, it returns simple error.
func Turbo(htmlText []byte) ([]byte, string, error) {
	got, err := VkToTurbo(htmlText)
	if err == nil {
		return got, `vkontakte`, nil
	}
	got, err = FbToTurbo(htmlText)
	if err == nil {
		return got, `facebook`, nil
	}
	got, err = InstaToTurbo(htmlText)
	if err == nil {
		return got, `instagram`, nil
	}
	got, err = TwitToTurbo(htmlText)
	if err == nil {
		return got, `twitter`, nil
	}
	got, err = YoutubeToTurbo(htmlText)
	if err == nil {
		return got, `youtube`, nil
	}
	got, err = IframeToTurbo(htmlText)
	if err == nil {
		return got, `iframe`, nil
	}

	return nil, ``, fmt.Errorf("unknown embed")
}

// printTurbo returns ready to handle Turbo with given parameters
func (ifrPost *iframePost) printTurbo() []byte {
	var attributes string
	if ifrPost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, ifrPost.Width)
	}
	if ifrPost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, ifrPost.Height)
	}
	if ifrPost.AllowFS {
		attributes += ` allowfullscreen="true"`
	}

	template := `<iframe%s frameborder="%d" src="%s"></iframe>`

	turbo := fmt.Sprintf(template, attributes, ifrPost.Frameborder, ifrPost.Src)

	return []byte(turbo)
}

// printTurbo returns ready to handle Turbo with given parameters
func (ypost *youtubePost) printTurbo() []byte {
	var attributes string
	if ypost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, ypost.Width)
	}
	if ypost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, ypost.Height)
	}
	if ypost.AllowFS {
		attributes += ` allowfullscreen="true"`
	}

	template := `<iframe%s frameborder="%d" src="https://www.youtube.com/embed/%s"></iframe>`

	amp := fmt.Sprintf(template, attributes, ypost.Frameborder, ypost.VideoID)

	return []byte(amp)
}

// VkToTurbo validates given vkontakte widget post for Yandex Turbo
// What is that? Look https://vk.com/dev/widget_post
func VkToTurbo(htmlText []byte) ([]byte, error) {
	if !bytes.Contains(htmlText, []byte(`VK.Widgets.Post`)) {
		return nil, fmt.Errorf("given string is not a VK widget post")
	}

	re := regexp.MustCompile(`VK.Widgets.Post\("vk_post_(-?\d+)_(-?\d+)", (-?\d+), (-?\d+), '(\S+?)'`)
	widgetParsed := re.FindSubmatch(htmlText)
	if widgetParsed == nil {
		return nil, fmt.Errorf("cannot parse vk widget")
	}

	// 1st and 3rd, 2nd and 4th should match
	if string(widgetParsed[1]) != string(widgetParsed[3]) || string(widgetParsed[2]) != string(widgetParsed[4]) {
		return nil, fmt.Errorf("parsed string does not match Vk widget post format")
	}

	_, err := strconv.ParseInt(string(widgetParsed[1]), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse owner id")
	}

	_, err = strconv.ParseInt(string(widgetParsed[2]), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse post id")
	}

	return htmlText, nil
}

// TwitToTurbo convertes given twitter embeddable html for Yandex Turbo
func TwitToTurbo(htmlText []byte) ([]byte, error) {
	pointerNode, err := html.Parse(bytes.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("cannot parse twitter html")
	}
	var post tweetPost

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					re := regexp.MustCompile(`https://twitter.com/[a-zA-Z_]{1,15}/status/(\d+)`)
					submatch := re.FindStringSubmatch(a.Val)
					if submatch == nil {
						continue
					}
					post.ID = submatch[1]
					post.Src = a.Val
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(pointerNode)

	if len(post.Src) < 1 {
		return nil, fmt.Errorf("no twitter ID in the url")
	}

	return htmlText, nil
}

// InstaToTurbo validates given instagram embeddable html for Yandex Turbo
func InstaToTurbo(htmlText []byte) ([]byte, error) {
	pointerNode, err := html.Parse(bytes.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("cannot parse insta html")
	}
	var post instaPost

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.Blockquote {
			for _, bq := range n.Attr {
				switch bq.Key {
				case "data-instgrm-permalink":
					post.Src = bq.Val
				case "width":
					w, err := strconv.ParseInt(bq.Val, 10, 0)
					if err == nil {
						post.Width = w
					}
				case "height":
					h, err := strconv.ParseInt(bq.Val, 10, 0)
					if err == nil {
						post.Height = h
					}
				}
			}
			if len(post.Src) > 0 {
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(pointerNode)

	if len(post.Src) < 1 {
		return nil, fmt.Errorf("no src in the url")
	}

	urlPtr, err := url.Parse(post.Src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url")
	}

	if !strings.Contains(urlPtr.Hostname(), "instagram.com") {
		return nil, fmt.Errorf("it is not instagram url")
	}

	re := regexp.MustCompile(`p/(\S+?)/`)
	submatch := re.FindStringSubmatch(urlPtr.Path)
	if submatch == nil {
		return nil, fmt.Errorf("instagram url is malformed")
	}

	return htmlText, nil
}

// FbToTurbo validates Facebook html for Yandex Turbo
func FbToTurbo(htmlText []byte) ([]byte, error) {
	pointerNode, err := html.Parse(bytes.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("cannot parse fb iframe")
	}
	var post fbPost

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.Iframe {
			for _, iframe := range n.Attr {
				switch iframe.Key {
				case "src":
					post.Src = iframe.Val
				case "width":
					w, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Width = w
					}
				case "height":
					h, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Height = h
					}
				}
			}
			if len(post.Src) > 0 {
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(pointerNode)

	if len(post.Src) < 1 {
		return nil, fmt.Errorf("no src in the url")
	}

	urlPtr, err := url.Parse(post.Src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse fb url")
	}

	if !strings.Contains(urlPtr.Hostname(), "facebook.com") {
		return nil, fmt.Errorf("it is not facebook url")
	}

	return htmlText, nil
}

// YoutubeToTurbo convertes Youtube embeddable html to Yandex Turbo
func YoutubeToTurbo(htmlText []byte) ([]byte, error) {
	pointerNode, err := html.Parse(bytes.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("cannot parse youtube iframe")
	}
	var post youtubePost

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.Iframe {
			for _, iframe := range n.Attr {
				switch iframe.Key {
				case "src":
					post.Src = iframe.Val
				case "width":
					w, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Width = w
					}
				case "height":
					h, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Height = h
					}
				case "allowfullscreen":
					post.AllowFS = true
				}
			}
			if len(post.Src) > 0 {
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(pointerNode)

	if !(len(post.Src) > 0) {
		return nil, fmt.Errorf("no src in the url")
	}

	urlPtr, err := url.Parse(post.Src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse youtube url")
	}

	if !strings.Contains(urlPtr.Hostname(), "youtube.com") {
		return nil, fmt.Errorf("it is not youtube url")
	}

	re := regexp.MustCompile(`embed/([A-Za-z0-9_-]{11})`)
	submatch := re.FindStringSubmatch(urlPtr.Path)
	if submatch == nil {
		return nil, fmt.Errorf("youtube url is malformed")
	}
	post.VideoID = submatch[1]

	return post.printTurbo(), nil
}

// IframeToTurbo convertes some custom iframe embeddable html to Yandex Turbo
func IframeToTurbo(htmlText []byte) ([]byte, error) {
	pointerNode, err := html.Parse(bytes.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("cannot parse iframe")
	}
	var post iframePost

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.DataAtom == atom.Iframe {
			for _, iframe := range n.Attr {
				switch iframe.Key {
				case "src":
					post.Src = iframe.Val
				case "width":
					w, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Width = w
					}
				case "height":
					h, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Height = h
					}
				case "allowfullscreen":
					post.AllowFS = true
				case "frameborder":
					fb, err := strconv.ParseInt(iframe.Val, 10, 0)
					if err == nil {
						post.Frameborder = fb
					}
				}
			}
			if len(post.Src) > 0 {
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(pointerNode)

	if len(post.Src) < 1 {
		return nil, fmt.Errorf("no src in the url")
	}

	urlPtr, err := url.Parse(post.Src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse iframe url")
	}

	if urlPtr.Scheme != `https` {
		return nil, fmt.Errorf("yandex Turbo supports only https iframe scheme")
	}

	return post.printTurbo(), nil
}
