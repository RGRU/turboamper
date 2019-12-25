package turboamper

import "testing"

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
