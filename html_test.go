package linkchecker

import (
	"reflect"
	"strings"
	"testing"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []Node
	}{
		{
			name: "only <a> tag",
			body: `
<html>
<head></head>
<body>
	<a href="https://example.com/1">example1</a>
	<a href="https://example.com/2">example2</a>
	<a href="https://example.com/3">example3</a>
</body>
</html>`,
			want: []Node{
				&AnchorNode{
					Href: "https://example.com/1",
					Text: "example1",
				},
				&AnchorNode{
					Href: "https://example.com/2",
					Text: "example2",
				},
				&AnchorNode{
					Href: "https://example.com/3",
					Text: "example3",
				},
			},
		},
		{
			name: "only <img> tag",
			body: `
<html>
<head></head>
<body>
	<img src="https://example.com/1" alt="example1">
	<img src="https://example.com/2" alt="example2">
	<img src="https://example.com/3" alt="example3">
</body>
</html>`,
			want: []Node{
				&ImgNode{
					Src: "https://example.com/1",
					Alt: "example1",
				},
				&ImgNode{
					Src: "https://example.com/2",
					Alt: "example2",
				},
				&ImgNode{
					Src: "https://example.com/3",
					Alt: "example3",
				},
			},
		},
		{
			name: "<a> and <img> tag",
			body: `
<html>
<head></head>
<body>
	<a href="https://example.com/1">example1</a>
	<a href="https://example.com/2">example2</a>
	<img src="https://example.com/1" alt="example1">
	<img src="https://example.com/2" alt="example2">
</body>
</html>`,
			want: []Node{
				&AnchorNode{
					Href: "https://example.com/1",
					Text: "example1",
				},
				&AnchorNode{
					Href: "https://example.com/2",
					Text: "example2",
				},
				&ImgNode{
					Src: "https://example.com/1",
					Alt: "example1",
				},
				&ImgNode{
					Src: "https://example.com/2",
					Alt: "example2",
				},
			},
		},
		{
			name: "Nested <img> tag",
			body: `
<html>
<head></head>
<body>
	<a href="https://example.com/1"><img src="https://image.example.com/1" alt="example_image1"></a>
	<a href="https://example.com/2"><img src="https://image.example.com/2" alt="example_image2"></a>
</body>
</html>`,
			want: []Node{
				&AnchorNode{
					Href: "https://example.com/1",
					Text: "",
				},
				&ImgNode{
					Src: "https://image.example.com/1",
					Alt: "example_image1",
				},
				&AnchorNode{
					Href: "https://example.com/2",
					Text: "",
				},
				&ImgNode{
					Src: "https://image.example.com/2",
					Alt: "example_image2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.body)
			got, err := Parse(r)
			if err != nil {
				t.Errorf("html.Parse() error=%w", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("html.Parse() got=%+v, want=%+v\n", got, tt.want)
			}
		})
	}
}
