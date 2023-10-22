package rss

import (
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Поток 1",
			args: args{
				url: "https://habr.com/ru/rss/best/daily/?fl=ru",
			},
		},
		{
			name: "Поток 2",
			args: args{
				url: "https://cprss.s3.amazonaws.com/golangweekly.com.xml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := Parse(tt.args.url)
			if err != nil {
				t.Fatal(err)
			}
			if len(stream) == 0 {
				t.Fatal("Данные не получены")
			}
			t.Logf("получено %d публикаций", len(stream))
		})
	}
}
