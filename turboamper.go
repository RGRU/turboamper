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

type youtubePost struct {
	Width   int64
	Height  int64
	VideoId    string
	Src     string
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

	amp := fmt.Sprintf(template, attributes, ypost.VideoId)

	return []byte(amp)
}

// tweetPost contents instagram data
type tweetPost struct {
	ID     string
	Width  int64
	Height int64
	Src    string
}

func (tpost *tweetPost) printAMP() []byte {
	attributes := ""
	if tpost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, tpost.Width)
	}
	if tpost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, tpost.Height)
	}
	//template := `<amp-instagram layout="responsive"%s data-shortcode="%s"></amp-instagram>`
	template := `<amp-twitter layout="responsive"%s data-tweetid="%s"></amp-twitter>`

	amp := fmt.Sprintf(template, attributes, tpost.ID)

	return []byte(amp) //
}

// instaPost contents instagram data
type instaPost struct {
	IsCaptioned bool
	Shortcode   string
	Width       int64
	Height      int64
	Src         string
}

func (ipost *instaPost) printAMP() []byte {
	attributes := ""
	if ipost.Width > 0 {
		attributes += fmt.Sprintf(` width="%d"`, ipost.Width)
	}
	if ipost.Height > 0 {
		attributes += fmt.Sprintf(` height="%d"`, ipost.Height)
	}
	if ipost.IsCaptioned {
		attributes += ` data-captioned`
	}
	template := `<amp-instagram layout="responsive"%s data-shortcode="%s"></amp-instagram>`

	amp := fmt.Sprintf(template, attributes, ipost.Shortcode)

	return []byte(amp) //
}

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
	IsVideo bool
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
	if fbpost.IsVideo {
		attributes += ` data-embed-as="video"`
	}
	template := `<amp-facebook layout="responsive"%s data-href="%s"></amp-facebook>`

	amp := fmt.Sprintf(template, attributes, fbpost.Href)

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
	//<iframe width="560" height="315" src="https://www.youtube.com/embed/TVakXOkE2G4" frameborder="0"
	//  allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
	//`<amp-youtube layout="responsive" width="560" height="315" data-videoid="TVakXOkE2G4"></amp-youtube>`,
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
	post.VideoId = submatch[1]

	return post.printAMP(), nil
}

// IframeToAMP convertes given instagram embeddable html to AMP
func IframeToAMP(htmlText []byte) ([]byte, error) {
	converted := make([]byte, 0, len(htmlText))
	return converted, nil
}
