// Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package turboamper provides support for making some HTML structured texts validateable by Yandex Turbo and Google AMP services.
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

// vkPost contents widget data
type vkPost struct {
	OwnerId int64
	PostId  int64
	Hash    string
	Width   int64
	Height  int64
}

// printAMP returns ready to handle AMP with given parameters
func (vkpost *vkPost) printAMP() []byte {
	attributes := ""
	if vkpost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, vkpost.Width)
	}
	if vkpost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, vkpost.Height)
	}
	template := `<amp-vk%s data-embedtype="post" layout="responsive" data-owner-id="%d" data-post-id="%d" data-hash="%s"></amp-vk>`

	amp := fmt.Sprintf(template, attributes, vkpost.OwnerId, vkpost.PostId, vkpost.Hash)

	return []byte(amp)
}

type fbPost struct {
	isVideo bool
	Width   int64
	Height  int64
	Href    string
	Src     string
}

// printAMP returns ready to handle AMP with given parameters
func (fbpost *fbPost) printAMP() []byte {
	var attributes string
	if fbpost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, fbpost.Width)
	}
	if fbpost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, fbpost.Height)
	}
	if fbpost.isVideo {
		attributes += ` data-embed-as="video"`
	}
	template := `<amp-facebook layout="responsive"%s data-href="%s"></amp-facebook>`

	amp := fmt.Sprintf(template, attributes, fbpost.Href)

	return []byte(amp)
}

// FbToAMP convertes given facebook embeddable html to AMP
func FbToAMP(htmlText []byte) ([]byte, error) {
	//`<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fstcnk%2Fposts%2F3384458724928901&width=500" width="500" height="498" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allow="encrypted-media"></iframe>`,
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
		return nil, fmt.Errorf("cannot parse fb url")
	}

	if !strings.Contains(urlPtr.Hostname(), "facebook.com") {
		return nil, fmt.Errorf("it is not facebook url")
	}

	if strings.Contains(urlPtr.Path, "video.php") {
		post.isVideo = true
	}

	post.Href = urlPtr.Query().Get("href")

	return post.printAMP(), nil
}

// VkToAMP convertes given vkontakte widget post to AMP
// What is that? Look https://vk.com/dev/widget_post
func VkToAMP(htmlText []byte) ([]byte, error) {
	// VK.Widgets.Post("vk_post_1_45616", 1, 45616, 'ZMk4b98xpQZMJJRXVsL1ig', {width: 500})
	//VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM')
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

	ownerId, err := strconv.ParseInt(string(widgetParsed[1]), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse owner id")
	}

	postId, err := strconv.ParseInt(string(widgetParsed[2]), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse owner id")
	}

	data := &vkPost{OwnerId: ownerId, PostId: postId, Hash: string(widgetParsed[5])}

	// let's extract width
	whRe := regexp.MustCompile(`VK.Widgets.Post\(.+?{width: (\d+)(?:, height: (\d+))?}\)`)
	widthHeight := whRe.FindSubmatch(htmlText)
	if widthHeight != nil {
		w, err := strconv.ParseInt(string(widthHeight[1]), 10, 0)
		if err == nil {
			data.Width = w
		}

		h, err := strconv.ParseInt(string(widthHeight[2]), 10, 0)
		if err == nil {
			data.Height = h
		}
	}

	return data.printAMP(), nil
}

// InstaToAMP convertes given instagram embeddable html to AMP
func InstaToAMP(htmlText []byte) ([]byte, error) {
	converted := make([]byte, 0, len(htmlText))
	return converted, nil
}

// TwitToAMP convertes given instagram embeddable html to AMP
func TwitToAMP(htmlText []byte) ([]byte, error) {
	converted := make([]byte, 0, len(htmlText))
	return converted, nil
}

// IframeToAMP convertes given instagram embeddable html to AMP
func IframeToAMP(htmlText []byte) ([]byte, error) {
	converted := make([]byte, 0, len(htmlText))
	return converted, nil
}
