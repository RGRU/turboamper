package turboamper

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

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
	post.VideoId = submatch[1]

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
