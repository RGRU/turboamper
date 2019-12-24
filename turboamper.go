// Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package turboamper provides support for making some HTML structured texts validateable by Yandex Turbo and Google AMP services.
package turboamper

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

// VkPost contents widget data
type VkPost struct {
	OwnerId int64
	PostId  int64
	Hash    string
	Width   int64
	Height  int64
}

// printAMP returns ready to handle AMP with given parameters
func (vkpost *VkPost) printAMP() []byte {
	attributes := ""
	if vkpost.Width > 0 {
		attributes +=  fmt.Sprintf(` width="%d"`, vkpost.Width)
	}
	if vkpost.Height > 0 {
		attributes +=  fmt.Sprintf(` height="%d"`, vkpost.Height)
	}
	template := `<amp-vk%s data-embedtype="post" layout="responsive" data-owner-id="%d" data-post-id="%d" data-hash="%s"></amp-vk>`

	amp := fmt.Sprintf(template, attributes, vkpost.OwnerId, vkpost.PostId, vkpost.Hash)

	return []byte(amp)
}

// FbToAMP convertes given facebook embeddable html to AMP
func FbToAMP(htmlText []byte) ([]byte, error) {
	converted := make([]byte, 0, len(htmlText))
	return converted, nil
}

// VkToAMP convertes given vkontakte embeddable html to AMP
func VkToAMP(htmlText []byte) ([]byte, error) {
	// VK.Widgets.Post("vk_post_1_45616", 1, 45616, 'ZMk4b98xpQZMJJRXVsL1ig', {width: 500})
	//VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM')
	if !bytes.Contains(htmlText, []byte(`VK.Widgets.Post`)) {
		return nil, fmt.Errorf("given string is not a VK widget post")
	}

	re := regexp.MustCompile(`VK.Widgets.Post\("vk_post_(.+?)_(.+?)", (-?\d+?), (-?\d+?), '(\S+?)'`)
	widgetParsed := re.FindSubmatch(htmlText)

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

	data := &VkPost{OwnerId: ownerId, PostId: postId, Hash: string(widgetParsed[5])}

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
