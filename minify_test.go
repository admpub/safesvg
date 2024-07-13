package safesvg

import "testing"

func TestMinify(t *testing.T) {

	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><style>.cls-1{isolation:isolate;}.cls-2{fill:#735462;}.cls-3{fill:#6a4c5b;}.cls-3,.cls-4,.cls-5,.cls-6,.cls-7,.cls-8{fill-rule:evenodd;}.cls-4{fill:#ec7752;}.cls-5{fill:#e36b43;}.cls-6{fill:url(#GradientFill_1);}.cls-7{fill:#e77e49;}.cls-10,.cls-11,.cls-8{fill:#fff;}.cls-9{opacity:0.5;mix-blend-mode:multiply;}.cls-11{font-size:42px;font-family:Impact;}</style><path fill="none" d="M0 0h24v24H0V0z"/><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"/></svg>`)

	minified, err := Minify(svg)
	if err != nil {
		t.Errorf("Unexptected error %v", err)
	}
	t.Logf("%s\nbefore:%d\nafter:%d", string(minified), len(svg), len(minified))
}
