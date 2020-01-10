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
		post.IsVideo = true
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
	//`<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" data-instgrm-version="12" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;"> View this post on Instagram</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div></a> <p style=" margin:8px 0 0 0; padding:0 4px;"> <a href="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#000; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none; word-wrap:break-word;" target="_blank">Председатель общероссийской общественной организации защиты семьи «Родительское Всероссийское Сопротивление» (РВС) Мария Мамиконян написала открытое письмо депутату Государственной думы @opushkina Оксане Пушкиной. ⠀ На пресс-конференции, посвященной законопроекту о профилактике семейно-бытового насилия (СБН), Пушкина заявила, что ей придется «оправдываться» в Страсбурге за непринятие закона. В письме, опубликованном в ИА «Регнум», председатель РВС напомнила, что Пушкина занимает пост спецпредставителя Госдумы во Всемирном банке по женскому предпринимательству. ⠀ Мария Мамиконян порекомендовала Пушкиной рассказать о мерах, которые уже применяются в России для профилактики насилия, в том числе семейного. Председатель РВС напомнила о положительной практике перевода «побоев», которые причинены впервые (за совершенные повторно в РФ предусмотрено уголовное наказание - прим. РВС), из разряда уголовных преступлений в административные нарушения. ⠀ Эта мера обеспечила неотвратимость наказания и снизила латентность этого нарушения. «Так что этой мерой Вам можно отчитываться как достижением, а не требовать её отмены и возврата всех побоев в УК!» — считает председатель РВС. ⠀ Она также посоветовала Пушкиной не оправдываться, а обратить внимание на то, что уровень насилия в России явно и сильно снижается. «То есть для самой постановки вопроса о чрезвычайных мерах в России нет почвы», — говорится в письме. ⠀ Мамиконян порекомендовала Пушкиной, как члену ПАСЕ, предложить коллегам за рубежом изучить передовой и эффективный российский опыт. Она отметила, что нормы, предлагаемые в скандальном законе о СБН, дискриминационны, коррупциогенны, несовместимы с презумпцией невиновности и попросту не имеют доказанную эффективность. ⠀ «И коль скоро вас так тяготит необходимость „оправдываться“ в Страсбурге за суверенные решения Российской Федерации, то, быть может, вам стоит освободиться от этих обременительных обязательств? Быть может, Россию в ПАСЕ лучше представлять людям, которые не будут оправдываться, но станут защищать интересы нашей страны на международной арене, а не наоборот?»— заключает Мария Мамиконян. ⠀ #СемейноБытовоеНасилие #ДомашнееНасилие #ОксанаПушкина #ЯНеХотелаУмирать</a></p> <p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;">A post shared by <a href="https://www.instagram.com/rvs.news/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px;" target="_blank"> РВС - защита семьи 👨‍👩‍👧‍👦</a> (@rvs.news) on <time style=" font-family:Arial,sans-serif; font-size:14px; line-height:17px;" datetime="2019-12-28T09:32:04+00:00">Dec 28, 2019 at 1:32am PST</time></p></div></blockquote> <script async src="//www.instagram.com/embed.js"></script>`,
	//`<amp-instagram layout="responsive" data-shortcode="B6nHZAHl7JZ"></amp-instagram>`
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

// IframeToAMP convertes given instagram embeddable html to AMP
func IframeToAMP(htmlText []byte) ([]byte, error) {
	converted := make([]byte, 0, len(htmlText))
	return converted, nil
}
