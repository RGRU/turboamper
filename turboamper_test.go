package turboamper

import (
	"fmt"
	"testing"
)

func TestFbToTurbo(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fstcnk%2Fposts%2F3384458724928901&width=500" width="500" height="498" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allow="encrypted-media"></iframe>`,
			`<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fstcnk%2Fposts%2F3384458724928901&width=500" width="500" height="498" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allow="encrypted-media"></iframe>`,
		},
		{
			`<iframe src="https://www.facebook.com/plugins/video.php?href=https%3A%2F%2Fwww.facebook.com%2Fnasaearth%2Fvideos%2F456540998570328%2F&show_text=0&width=560" width="560" height="373" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>`,
			`<iframe src="https://www.facebook.com/plugins/video.php?href=https%3A%2F%2Fwww.facebook.com%2Fnasaearth%2Fvideos%2F456540998570328%2F&show_text=0&width=560" width="560" height="373" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>`,
		},
		{
			`<iframe src="https://www.facebook.com/plugins/video.php?href=https%3A%2F%2Fwww.facebook.com%2Fbarsuksergey%2Fvideos%2F2720743767989363%2F&show_text=0&width=560" width="560" height="308" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>`,
			`<iframe src="https://www.facebook.com/plugins/video.php?href=https%3A%2F%2Fwww.facebook.com%2Fbarsuksergey%2Fvideos%2F2720743767989363%2F&show_text=0&width=560" width="560" height="308" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>`,
		},
		{
				//error
			`<iframe src="" width="560" height="308" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allowFullScreen="true"></iframe>`,
			`no src in the url`,
		},
	}

	for i, test := range tests {
		got, err := FbToTurbo([]byte(test.input))
		if err != nil {
			if fmt.Sprint(err) != test.want {
				t.Errorf("\n[%d]IframeToTurbo() = %q,\nwant ERR    %q\n", i+1, err, test.want)
			}
			continue
		}

		if string(got) != test.want {
			t.Errorf("\n[%d]IframeToTurbo() = %q,\nwant        %q\n", i+1, got, test.want)
		}
	}
}

func TestIframeToTurbo(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/nopolitics/video/706313-popokatepetl-stolb-pepel-3-km/video/5e1852ef02e8bd3b731db837" frameborder="0" allowfullscreen/></iframe></div>`,
			`<iframe allowfullscreen="true" frameborder="0" src="https://russian.rt.com/nopolitics/video/706313-popokatepetl-stolb-pepel-3-km/video/5e1852ef02e8bd3b731db837"></iframe>`,
		},
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb" frameborder="0"/></iframe></div>`,
			`<iframe frameborder="0" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb"></iframe>`,
		},
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb" frameborder="2"/></iframe></div>`,
			`<iframe frameborder="2" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb"></iframe>`,
		},
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb" frameborder="0" allowfullscreen/></iframe></div>`,
			`<iframe allowfullscreen="true" frameborder="0" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb"></iframe>`,
		},
		{
			//error
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="" frameborder="0" allowfullscreen/></iframe></div>`,
			`no src in the url`,
		},
	}

	for i, test := range tests {
		got, err := IframeToTurbo([]byte(test.input))
		if err != nil {
			if fmt.Sprint(err) != test.want {
				t.Errorf("\n[%d]IframeToTurbo() = %q,\nwant ERR    %q\n", i+1, err, test.want)
			}
			continue
		}

		if string(got) != test.want {
			t.Errorf("\n[%d]IframeToTurbo() = %q,\nwant        %q\n", i+1, got, test.want)
		}
	}
}

func TestYoutubeToTurbo(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<iframe width="560" height="315" src="https://www.youtube.com/embed/05klG-PTKqo" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<iframe width="560" height="315" allowfullscreen="true" frameborder="0" src="https://www.youtube.com/embed/05klG-PTKqo"></iframe>`,
		},
		{
			`<iframe src="https://www.youtube.com/embed/TVakXOkE2G4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<iframe allowfullscreen="true" frameborder="0" src="https://www.youtube.com/embed/TVakXOkE2G4"></iframe>`,
		},
		{
			`<iframe height="315" src="https://www.youtube.com/embed/TVakXOkE2G4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<iframe height="315" allowfullscreen="true" frameborder="0" src="https://www.youtube.com/embed/TVakXOkE2G4"></iframe>`,
		},
		{ //error
			`<iframe width="560" height="315" src="https://www.youtube.com/embed/" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`youtube url is malformed`,
		},
	}

	for i, test := range tests {
		got, err := YoutubeToTurbo([]byte(test.input))
		if err != nil {
			if fmt.Sprint(err) != test.want {
				t.Errorf("\n[%d]YoutubeToTurbo() = %q,\nwant ERR    %q\n", i+1, err, test.want)
			}
			continue
		}

		if string(got) != test.want {
			t.Errorf("\n[%d]YoutubeToTurbo() = %q,\nwant        %q\n", i+1, got, test.want)
		}
	}
}

func TestIframeToAMP(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/nopolitics/video/706313-popokatepetl-stolb-pepel-3-km/video/5e1852ef02e8bd3b731db837" frameborder="0" allowfullscreen/></iframe></div>`,
			`<amp-iframe width="480" height="315" sandbox="allow-scripts allow-same-origin" layout="responsive" frameborder="0" allowfullscreen src="https://russian.rt.com/nopolitics/video/706313-popokatepetl-stolb-pepel-3-km/video/5e1852ef02e8bd3b731db837"`,
		},
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb" frameborder="0"/></iframe></div>`,
			`<amp-iframe width="480" height="315" sandbox="allow-scripts allow-same-origin" layout="responsive" frameborder="0" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb"`,
		},
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb" frameborder="2"/></iframe></div>`,
			`<amp-iframe width="480" height="315" sandbox="allow-scripts allow-same-origin" layout="responsive" frameborder="2" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb"`,
		},
		{
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb" frameborder="0" allowfullscreen/></iframe></div>`,
			`<amp-iframe width="480" height="315" sandbox="allow-scripts allow-same-origin" layout="responsive" frameborder="0" allowfullscreen src="https://russian.rt.com/world/video/706283-posol-iran-oon-ssha-suleimani/video/5e184bbf02e8bd3f073eebeb"`,
		},
		{
			//error
			`<div style="position: relative;padding-bottom: 56.25%; padding-top: 25px; height: 0;"><iframe style="position: absolute;top: 0;left: 0;width: 100%;height: 100%;" src="" frameborder="0" allowfullscreen/></iframe></div>`,
			`no src in the url`,
		},
	}

	for i, test := range tests {
		got, err := IframeToAMP([]byte(test.input))
		if err != nil {
			if fmt.Sprint(err) != test.want {
				t.Errorf("\n[%d]IframeToAMP() = %q,\nwant ERR    %q\n", i+1, err, test.want)
			}
			continue
		}

		if string(got) != test.want {
			t.Errorf("\nIframeToAMP() = %q,\nwant        %q\n", got, test.want)
		}
	}
}

func TestYoutubeToAMP(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<iframe width="560" height="315" src="https://www.youtube.com/embed/05klG-PTKqo" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<amp-youtube layout="responsive" width="560" height="315" data-videoid="05klG-PTKqo"></amp-youtube>`,
		},
		{
			`<iframe src="https://www.youtube.com/embed/TVakXOkE2G4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<amp-youtube layout="responsive" data-videoid="TVakXOkE2G4"></amp-youtube>`,
		},
		{
			`<iframe height="315" src="https://www.youtube.com/embed/TVakXOkE2G4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<amp-youtube layout="responsive" height="315" data-videoid="TVakXOkE2G4"></amp-youtube>`,
		},
		{
			`<iframe width="560" height="315" src="https://www.youtube.com/embed/TVakXOkE2G4" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`<amp-youtube layout="responsive" width="560" height="315" data-videoid="TVakXOkE2G4"></amp-youtube>`,
		},
		{ //error
			`<iframe width="560" height="315" src="https://www.youtube.com/embed/" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
			`youtube url is malformed`,
		},
	}

	for i, test := range tests {
		got, err := YoutubeToAMP([]byte(test.input))
		if err != nil {
			if fmt.Sprint(err) != test.want {
				t.Errorf("\n[%d]YoutubeToAMP() = %q,\nwant ERR    %q\n", i+1, err, test.want)
			}
			continue
		}

		if string(got) != test.want {
			t.Errorf("\nYoutubeToAMP() = %q,\nwant        %q\n", got, test.want)
			t.Errorf("\nERROR: %q", err)
		}
	}
}

func TestTwitToAMP(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<blockquote class="twitter-tweet"><p lang="en" dir="ltr"><a href="https://twitter.com/hashtag/TodayInHistory?src=hash&amp;ref_src=twsrc%5Etfw">#TodayInHistory</a> | 1999 – Prime minister Vladimir Putin becomes acting president <a href="https://t.co/pHbyk12C0s">pic.twitter.com/pHbyk12C0s</a></p>&mdash; WION (@WIONews) <a href="https://twitter.com/WIONews/status/1211912897590202368?ref_src=twsrc%5Etfw">December 31, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`,
			`<amp-twitter layout="responsive" data-tweetid="1211912897590202368"></amp-twitter>`,
		},
		{
			`<blockquote class="twitter-tweet"><p lang="ru" dir="ltr">Андрей Сошенко. Когда рванет второй Чернобыль? <br>Рано или поздно, но на Украине обязательно сотворят глобальную катастрофу <a href="https://t.co/EQGPtpvxVF">https://t.co/EQGPtpvxVF</a> <a href="https://t.co/WBIrRCAvZq">pic.twitter.com/WBIrRCAvZq</a></p>&mdash; газета Завтра (@ZavtraRu) <a href="https://twitter.com/ZavtraRu/status/1215336058755436547?ref_src=twsrc%5Etfw">January 9, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`,
			`<amp-twitter layout="responsive" data-tweetid="1215336058755436547"></amp-twitter>`,
		},
		{
			`<blockquote class="twitter-tweet"><p lang="ru" dir="ltr">Андрей Сошенко. Когда рванет второй Чернобыль? <br>Рано или поздно, но на Украине обязательно сотворят глобальную катастрофу <a href="https://t.co/EQGPtpvxVF">https://t.co/EQGPtpvxVF</a> <a href="https://t.co/WBIrRCAvZq">pic.twitter.com/WBIrRCAvZq</a></p>&mdash; газета Завтра (@ZavtraRu) <a href="https://twitter.com/ZavtraRu/status/?ref_src=twsrc%5Etfw">January 9, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`,
			`no twitter ID in the url`,
		},
	}

	for i, test := range tests {
		got, err := TwitToAMP([]byte(test.input))
		if err != nil {
			if fmt.Sprint(err) != test.want {
				t.Errorf("\n[%d]TwitToAMP() = %q,\nwant ERR    %q\n", i+1, err, test.want)
			}
			continue
		}

		if string(got) != test.want {
			t.Errorf("\nTwitToAMP() = %q,\nwant        %q\n", got, test.want)
			t.Errorf("\nERROR: %q", err)
		}
	}
}

func TestInstaToAMP(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{
			`<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" data-instgrm-version="12" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;"> View this post on Instagram</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div></a> <p style=" margin:8px 0 0 0; padding:0 4px;"> <a href="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#000; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none; word-wrap:break-word;" target="_blank">Председатель общероссийской общественной организации защиты семьи «Родительское Всероссийское Сопротивление» (РВС) Мария Мамиконян написала открытое письмо депутату Государственной думы @opushkina Оксане Пушкиной. ⠀ На пресс-конференции, посвященной законопроекту о профилактике семейно-бытового насилия (СБН), Пушкина заявила, что ей придется «оправдываться» в Страсбурге за непринятие закона. В письме, опубликованном в ИА «Регнум», председатель РВС напомнила, что Пушкина занимает пост спецпредставителя Госдумы во Всемирном банке по женскому предпринимательству. ⠀ Мария Мамиконян порекомендовала Пушкиной рассказать о мерах, которые уже применяются в России для профилактики насилия, в том числе семейного. Председатель РВС напомнила о положительной практике перевода «побоев», которые причинены впервые (за совершенные повторно в РФ предусмотрено уголовное наказание - прим. РВС), из разряда уголовных преступлений в административные нарушения. ⠀ Эта мера обеспечила неотвратимость наказания и снизила латентность этого нарушения. «Так что этой мерой Вам можно отчитываться как достижением, а не требовать её отмены и возврата всех побоев в УК!» — считает председатель РВС. ⠀ Она также посоветовала Пушкиной не оправдываться, а обратить внимание на то, что уровень насилия в России явно и сильно снижается. «То есть для самой постановки вопроса о чрезвычайных мерах в России нет почвы», — говорится в письме. ⠀ Мамиконян порекомендовала Пушкиной, как члену ПАСЕ, предложить коллегам за рубежом изучить передовой и эффективный российский опыт. Она отметила, что нормы, предлагаемые в скандальном законе о СБН, дискриминационны, коррупциогенны, несовместимы с презумпцией невиновности и попросту не имеют доказанную эффективность. ⠀ «И коль скоро вас так тяготит необходимость „оправдываться“ в Страсбурге за суверенные решения Российской Федерации, то, быть может, вам стоит освободиться от этих обременительных обязательств? Быть может, Россию в ПАСЕ лучше представлять людям, которые не будут оправдываться, но станут защищать интересы нашей страны на международной арене, а не наоборот?»— заключает Мария Мамиконян. ⠀ #СемейноБытовоеНасилие #ДомашнееНасилие #ОксанаПушкина #ЯНеХотелаУмирать</a></p> <p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;">A post shared by <a href="https://www.instagram.com/rvs.news/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px;" target="_blank"> РВС - защита семьи 👨‍👩‍👧‍👦</a> (@rvs.news) on <time style=" font-family:Arial,sans-serif; font-size:14px; line-height:17px;" datetime="2019-12-28T09:32:04+00:00">Dec 28, 2019 at 1:32am PST</time></p></div></blockquote> <script async src="//www.instagram.com/embed.js"></script>`,
			`<amp-instagram layout="responsive" data-shortcode="B6nHZAHl7JZ"></amp-instagram>`,
		},
		{
			`<blockquote class="instagram-media" data-instgrm-captioned data-instgrm-permalink="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" data-instgrm-version="12" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;"> View this post on Instagram</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div></a> <p style=" margin:8px 0 0 0; padding:0 4px;"> <a href="https://www.instagram.com/p/B6nHZAHl7JZ/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#000; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none; word-wrap:break-word;" target="_blank">Председатель общероссийской общественной организации защиты семьи «Родительское Всероссийское Сопротивление» (РВС) Мария Мамиконян написала открытое письмо депутату Государственной думы @opushkina Оксане Пушкиной. ⠀ На пресс-конференции, посвященной законопроекту о профилактике семейно-бытового насилия (СБН), Пушкина заявила, что ей придется «оправдываться» в Страсбурге за непринятие закона. В письме, опубликованном в ИА «Регнум», председатель РВС напомнила, что Пушкина занимает пост спецпредставителя Госдумы во Всемирном банке по женскому предпринимательству. ⠀ Мария Мамиконян порекомендовала Пушкиной рассказать о мерах, которые уже применяются в России для профилактики насилия, в том числе семейного. Председатель РВС напомнила о положительной практике перевода «побоев», которые причинены впервые (за совершенные повторно в РФ предусмотрено уголовное наказание - прим. РВС), из разряда уголовных преступлений в административные нарушения. ⠀ Эта мера обеспечила неотвратимость наказания и снизила латентность этого нарушения. «Так что этой мерой Вам можно отчитываться как достижением, а не требовать её отмены и возврата всех побоев в УК!» — считает председатель РВС. ⠀ Она также посоветовала Пушкиной не оправдываться, а обратить внимание на то, что уровень насилия в России явно и сильно снижается. «То есть для самой постановки вопроса о чрезвычайных мерах в России нет почвы», — говорится в письме. ⠀ Мамиконян порекомендовала Пушкиной, как члену ПАСЕ, предложить коллегам за рубежом изучить передовой и эффективный российский опыт. Она отметила, что нормы, предлагаемые в скандальном законе о СБН, дискриминационны, коррупциогенны, несовместимы с презумпцией невиновности и попросту не имеют доказанную эффективность. ⠀ «И коль скоро вас так тяготит необходимость „оправдываться“ в Страсбурге за суверенные решения Российской Федерации, то, быть может, вам стоит освободиться от этих обременительных обязательств? Быть может, Россию в ПАСЕ лучше представлять людям, которые не будут оправдываться, но станут защищать интересы нашей страны на международной арене, а не наоборот?»— заключает Мария Мамиконян. ⠀ #СемейноБытовоеНасилие #ДомашнееНасилие #ОксанаПушкина #ЯНеХотелаУмирать</a></p> <p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;">A post shared by <a href="https://www.instagram.com/rvs.news/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px;" target="_blank"> РВС - защита семьи 👨‍👩‍👧‍👦</a> (@rvs.news) on <time style=" font-family:Arial,sans-serif; font-size:14px; line-height:17px;" datetime="2019-12-28T09:32:04+00:00">Dec 28, 2019 at 1:32am PST</time></p></div></blockquote> <script async src="//www.instagram.com/embed.js"></script>`,
			`<amp-instagram layout="responsive" data-captioned data-shortcode="B6nHZAHl7JZ"></amp-instagram>`,
		},
		{
			`<blockquote class="instagram-media" data-instgrm-captioned data-instgrm-permalink="https://www.instagram.com/p/B6n1kfKoLmr/?utm_source=ig_embed&amp;utm_campaign=loading" data-instgrm-version="12" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/B6n1kfKoLmr/?utm_source=ig_embed&amp;utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;"> View this post on Instagram</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div></a> <p style=" margin:8px 0 0 0; padding:0 4px;"> <a href="https://www.instagram.com/p/B6n1kfKoLmr/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#000; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none; word-wrap:break-word;" target="_blank">Во Франции власти неоднократно пытались повысить пенсионный возраст. Несмотря на более мягкие условия повышения пенсионного возраста по сравнению с Россией, в декабре этого года противостояние достигло предельного накала. К забастовкам во Франции присоединился спецназ полиции, а количество протестующих перевалило за миллион. Но в России, по мнению сотрудника Федерального научно-исследовательского социологического центра РАН Анны Мытиль, повторение французского опыта с приостановкой повышения пенсионного возраста невозможно по нескольким причинам. Первая заключается в том, что в России, повысили возраст выхода на пенсию только для „гражданских лиц“. Для „льготных категорий“ работников, к которым относятся военнослужащие, сотрудники правоохранительной системы, в том числе, судов и прокуратуры, порядок выхода на пенсию остался прежним. Одного недовольства для появления протестов, аналогичных Французским, недостаточно.  Для этого и нужны лидеры, организации — для объединения, выработки стратегии, воодушевления. И беда граждан, что одни из этих лидеров и организаций превратились в системных бюрократов, а другим еще необходимо набираться политического опыта. И, конечно, если бы, на митинги (не разовые) вышел миллион граждан, которые подписали петицию движения „Суть времени“, то, вероятно, что власти проявили бы гораздо больше готовности к диалогу. Профсоюзы, которые сразу после объявления о повышении пенсионного возраста, инициировали манифестации, петиции, как-то быстро свернули эту деятельность. Усилия КПРФ и других партий, в программах которых обозначена защита интересов „простых“ людей, вообще можно назвать имитационными. Французы добились отмены повышения пенсионного возраста до 64 лет для тех, кто родился до 1975, а также отставки идеолога реформы Жан-Поля Делевое. Но протесты не закончились и протестующие намерены полностью отменить реформу. #рвс #rvs #пенсионнаяреформа2018</a></p> <p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;">A post shared by <a href="https://www.instagram.com/rvs.news/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px;" target="_blank"> РВС - защита семьи 👨‍👩‍👧‍👦</a> (@rvs.news) on <time style=" font-family:Arial,sans-serif; font-size:14px; line-height:17px;" datetime="2019-12-28T16:15:35+00:00">Dec 28, 2019 at 8:15am PST</time></p></div></blockquote> <script async src="//www.instagram.com/embed.js"></script>`,
			`<amp-instagram layout="responsive" data-captioned data-shortcode="B6n1kfKoLmr"></amp-instagram>`,
		},
		{
			`<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/B6n1kfKoLmr/?utm_source=ig_embed&amp;utm_campaign=loading" data-instgrm-version="12" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/B6n1kfKoLmr/?utm_source=ig_embed&amp;utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;"> View this post on Instagram</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center; margin-bottom: 24px;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 224px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 144px;"></div></div></a><p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;"><a href="https://www.instagram.com/p/B6n1kfKoLmr/?utm_source=ig_embed&amp;utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none;" target="_blank">A post shared by РВС - защита семьи 👨‍👩‍👧‍👦 (@rvs.news)</a> on <time style=" font-family:Arial,sans-serif; font-size:14px; line-height:17px;" datetime="2019-12-28T16:15:35+00:00">Dec 28, 2019 at 8:15am PST</time></p></div></blockquote> <script async src="//www.instagram.com/embed.js"></script>`,
			`<amp-instagram layout="responsive" data-shortcode="B6n1kfKoLmr"></amp-instagram>`,
		},
	}

	for _, test := range tests {
		if got, err := InstaToAMP([]byte(test.input)); string(got) != test.want {
			t.Errorf("\nInstaToAMP() = %q,\nwant        %q\n", got, test.want)
			t.Errorf("\nERROR: %q", err)
		}
	}

}

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

func ExampleFbToAMP() {
	html := `<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fstcnk%2Fposts%2F3384458724928901&width=500" width="500" height="498" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowTransparency="true" allow="encrypted-media"></iframe>`
	amp, err := FbToAMP([]byte(html))
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	fmt.Printf("AMPfied: %s", amp)
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
