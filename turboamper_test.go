package turboamper

import "testing"

func TestVkToAMP(t *testing.T) {
	html := `<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM');
	}());
	</script>`

	//want := []byte(`<amp-vk width="500" height="300" data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`)
	want := []byte(`<amp-vk data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`)

	if got, _ := VkToAMP([]byte(html)); string(got) != string(want) {
		t.Errorf("VkToAMP() = %q, want %q", got, want)
	}
}

func TestVkToAMPWidth(t *testing.T) {
	html := `<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM', {width: 500});
	}());
	</script>`

	want := []byte(`<amp-vk width="500" data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`)

	if got, _ := VkToAMP([]byte(html)); string(got) != string(want) {
		t.Errorf("VkToAMP() = %q, want %q", got, want)
	}
}

func TestVkToAMPWidthHeight(t *testing.T) {
	html := `<div id="vk_post_-175249128_1156"></div>
	<script type="text/javascript" src="https://vk.com/js/api/openapi.js?162"></script>
	<script type="text/javascript">
		(function() {
		VK.Widgets.Post("vk_post_-175249128_1156", -175249128, 1156, 'HmCFKRSM81NEzJ8mY9gzgXOlEFM', {width: 500, height: 300});
	}());
	</script>`

	want := []byte(`<amp-vk width="500" height="300" data-embedtype="post" layout="responsive" data-owner-id="-175249128" data-post-id="1156" data-hash="HmCFKRSM81NEzJ8mY9gzgXOlEFM"></amp-vk>`)

	if got, _ := VkToAMP([]byte(html)); string(got) != string(want) {
		t.Errorf("VkToAMP() = %q, want %q", got, want)
	}
}
