// Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package turboamper provides support for making some HTML structured texts validateable by Yandex Turbo and Google AMP services.
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

// AMP gives you amp-representation of html and its type
// If it cannot recognize your html, it returns simple error.
func AMP(htmlText []byte) ([]byte, string, error) {
	got, err := VkToAMP(htmlText)
	if err == nil {
		return got, `vkontakte`, nil
	}
	got, err = FbToAMP(htmlText)
	if err == nil {
		return got, `facebook`, nil
	}
	got, err = InstaToAMP(htmlText)
	if err == nil {
		return got, `instagram`, nil
	}
	got, err = TwitToAMP(htmlText)
	if err == nil {
		return got, `twitter`, nil
	}
	got, err = YoutubeToAMP(htmlText)
	if err == nil {
		return got, `youtube`, nil
	}
	got, err = IframeToAMP(htmlText)
	if err == nil {
		return got, `iframe`, nil
	}
	got, err = PlaybuzzToAMP(htmlText)
	if err == nil {
		return got, `playbuzz`, nil
	}

	return nil, ``, fmt.Errorf("unknown embed")
}

type iframePost struct {
	Width       int64
	Height      int64
	AllowFS     bool
	Frameborder int64
	Src         string
}

// printAMP returns ready to handle AMP with given parameters
func (ifrPost *iframePost) printAMP() []byte {
	var attributes string
	if ifrPost.AllowFS {
		attributes += ` allowfullscreen`
	}

	template := `<amp-iframe width="480" height="315" sandbox="allow-scripts allow-same-origin" layout="responsive" frameborder="%d"%s src="%s"`

	amp := fmt.Sprintf(template, ifrPost.Frameborder, attributes, ifrPost.Src)

	return []byte(amp)
}

type youtubePost struct {
	AllowFS     bool
	Frameborder int64
	Width       int64
	Height      int64
	VideoID     string
	Src         string
}

// printAMP returns ready to handle AMP with given parameters
func (ypost *youtubePost) printAMP() []byte {
	var attributes string
	if ypost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, ypost.Width)
	}
	if ypost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, ypost.Height)
	}

	template := `<amp-youtube layout="responsive"%s data-videoid="%s"></amp-youtube>`

	amp := fmt.Sprintf(template, attributes, ypost.VideoID)

	return []byte(amp)
}

// tweetPost contents instagram data
type tweetPost struct {
	ID     string
	Width  int64
	Height int64
	Src    string
}

func (post *tweetPost) printAMP() []byte {
	if post.Width == 0 {
		post.Width = 380
	}
	if post.Height == 0 {
		post.Height = 480
	}
	//template := `<amp-instagram layout="responsive"%s data-shortcode="%s"></amp-instagram>`
	template := `<amp-twitter layout="responsive" height="%d" width="%d" data-tweetid="%s"></amp-twitter>`

	amp := fmt.Sprintf(template, post.Height, post.Width, post.ID)

	return []byte(amp)
}

// instaPost contents instagram data
type instaPost struct {
	IsCaptioned bool
	Shortcode   string
	Width       int64
	Height      int64
	Src         string
}

func (post *instaPost) printAMP() []byte {
	attributes := ""
	if post.Width == 0 {
		post.Width = 400
	}
	if post.Height == 0 {
		post.Height = 400
	}
	if post.IsCaptioned {
		attributes += ` data-captioned`
	}
	template := `<amp-instagram layout="responsive" height="%d" width="%d"%s data-shortcode="%s"></amp-instagram>`

	amp := fmt.Sprintf(template, post.Height, post.Width, attributes, post.Shortcode)

	return []byte(amp)
}

// playbuzzPost contents playbuzz data
type playbuzzPost struct {
	DataItem string
	Height   int64
	Src      string
}

func (post *playbuzzPost) printAMP() []byte {
	if post.Height == 0 {
		post.Height = 500
	}
	template := `<amp-playbuzz layout="responsive" height="%d" data-item="%s"></amp-playbuzz>`

	amp := fmt.Sprintf(template, post.Height, post.DataItem)

	return []byte(amp)
}

// vkPost contents widget data
type vkPost struct {
	OwnerID int64
	PostID  int64
	Hash    string
	Width   int64
	Height  int64
}

// printAMP returns ready to handle AMP with given parameters
func (post *vkPost) printAMP() []byte {
	if post.Width == 0 {
		post.Width = 500
	}
	if post.Height == 0 {
		post.Height = 300
	}
	template := `<amp-vk height="%d" width="%d" data-embedtype="post" layout="responsive" data-owner-id="%d" data-post-id="%d" data-hash="%s"></amp-vk>`

	amp := fmt.Sprintf(template, post.Height, post.Width, post.OwnerID, post.PostID, post.Hash)

	return []byte(amp)
}

type fbPost struct {
	IsVideo bool
	Width   int64
	Height  int64
	Href    string
	Src     string
}

// printAMP returns ready to handle AMP with given parameters
func (post *fbPost) printAMP() []byte {
	var attributes string
	if post.Width == 0 {
		post.Width = 500
	}
	if post.Height == 0 {
		post.Height = 500
	}
	if post.IsVideo {
		attributes += ` data-embed-as="video"`
	}
	template := `<amp-facebook height="%d" width="%d" layout="responsive"%s data-href="%s"></amp-facebook>`

	amp := fmt.Sprintf(template, post.Height, post.Width, attributes, post.Href)

	return []byte(amp)
}

// FbToAMP convertes given facebook embeddable html to AMP
func FbToAMP(htmlText []byte) ([]byte, error) {
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
		post.IsVideo = true
	}

	post.Href = urlPtr.Query().Get("href")

	return post.printAMP(), nil
}

// VkToAMP convertes given vkontakte widget post to AMP
// What is that? Look https://vk.com/dev/widget_post
func VkToAMP(htmlText []byte) ([]byte, error) {
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

	ownerID, err := strconv.ParseInt(string(widgetParsed[1]), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse owner id")
	}

	postID, err := strconv.ParseInt(string(widgetParsed[2]), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse owner id")
	}

	data := &vkPost{OwnerID: ownerID, PostID: postID, Hash: string(widgetParsed[5])}

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

	if !(len(post.Src) > 0) {
		return nil, fmt.Errorf("no src in the url")
	}

	urlPtr, err := url.Parse(post.Src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url")
	}

	if !strings.Contains(urlPtr.Hostname(), "instagram.com") {
		return nil, fmt.Errorf("it is not instagram url")
	}

	if bytes.Contains(htmlText, []byte(` data-instgrm-captioned`)) {
		post.IsCaptioned = true
	}

	re := regexp.MustCompile(`p/(\S+?)/`)
	submatch := re.FindStringSubmatch(urlPtr.Path)
	if submatch == nil {
		return nil, fmt.Errorf("instagram url is malformed")
	}
	post.Shortcode = submatch[1]

	return post.printAMP(), nil
}

// TwitToAMP convertes given twitter embeddable html to AMP
func TwitToAMP(htmlText []byte) ([]byte, error) {
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

	if !(len(post.Src) > 0) {
		return nil, fmt.Errorf("no twitter ID in the url")
	}

	return post.printAMP(), nil
}

// YoutubeToAMP convertes given youtube embeddable html to AMP
func YoutubeToAMP(htmlText []byte) ([]byte, error) {
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

	return post.printAMP(), nil
}

// IframeToAMP convertes some custom iframe embeddable html to AMP
// Tested on Russia Today
func IframeToAMP(htmlText []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("amp supports only https iframe scheme")
	}

	return post.printAMP(), nil
}

// PlaybuzzToAMP convert playbuzz code
func PlaybuzzToAMP(htmlText []byte) ([]byte, error) {
	r := regexp.MustCompile(`<div(.+?)(class="playbuzz") data-id="(.+?)"(.+?)<\/div>`)
	data := r.FindSubmatch(htmlText)

	var post playbuzzPost
	if len(data) > 3 {
		post.DataItem = string(data[3])
	}

	return post.printAMP(), nil
}
