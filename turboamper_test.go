package turboamper

import (
	"fmt"
	"testing"
)

func TestFBToAMPTable(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fstcnk%2Fposts%2F3384458724928901&width=500" width="500" height="498" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allow="encrypted-media"></iframe>`,
			`<amp-facebook layout="responsive" width="500" height="498" data-href="https://www.facebook.com/stcnk/posts/3384458724928901"></amp-facebook>`,
		},
		{
			`<iframe src="https://www.facebook.com/plugins/video.php?href=https%3A%2F%2Fwww.facebook.com%2Fnasaearth%2Fvideos%2F456540998570328%2F&show_text=0&width=560" width="560" height="373" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>`,
			`<amp-facebook layout="responsive" width="560" height="373" data-embed-as="video" data-href="https://www.facebook.com/nasaearth/videos/456540998570328/"></amp-facebook>`,
		},
		{
			`<iframe src="https://www.facebook.com/plugins/video.php?href=https%3A%2F%2Fwww.facebook.com%2Fbarsuksergey%2Fvideos%2F2720743767989363%2F&show_text=0&width=560" width="560" height="308" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>
`,
			`<amp-facebook layout="responsive" width="560" height="308" data-embed-as="video" data-href="https://www.facebook.com/barsuksergey/videos/2720743767989363/"></amp-facebook>`,
		},
	}

	for _, test := range tests {
		if got, _ := FbToAMP([]byte(test.input)); string(got) != test.want {
			t.Errorf("\nFbToAMP() = %q,\nwant        %q", got, test.want)
		}
	}

}

func TestVkToAMPTable(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{`<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM');
	}());
	</script>`, `<amp-vk data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`,
		},
		{
			`<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM', {width: 500});
	}());
	</script>`, `<amp-vk width="500" data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`,
		},
		{
			`<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM', {width: 500, height: 300});
	}());
	</script>`, `<amp-vk width="500" height="300" data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`,
		},
	}

	for _, test := range tests {
		if got, _ := VkToAMP([]byte(test.input)); string(got) != test.want {
			t.Errorf("VkToAMP() = %q, want %q", got, test.want)
		}
	}

}

func TestVkToAMPWrong(t *testing.T) {
	html := `<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", 175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM', {width: 500, height: 300});
	}());
	</script>`

	if got, err := VkToAMP([]byte(html)); err == nil {
		t.Errorf("VkToAMP() = %q; want err, got result", got)
	}
}

func ExampleVkToAMP() {
	html := `<div id="vk_post_-165546713_21078"></div>
<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
<script type="text/javascript">
  (function() {
    VK.Widgets.Post( "vk_post_-165546713_21078", -165546713, 21078, 'UDKmSYMw9-_LHr7Lcgz8oAVE3Xg', {width: 600});
  }());
</script>`
	amp, err := VkToAMP([]byte(html))
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	fmt.Printf("AMPfied: %s", amp)
}
