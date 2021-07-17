package crawl_test

import (
	"GoCrawl/crawl"
	"testing"
)

func TestGetUrlDocument(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{name: "invalidSiteUrl", url: "hello_world", wantErr: true},
		{name: "validSiteUrlNonStatusOk", url: "https://disney.com.tw/doesnotexist", wantErr: true},
		{name: "validSiteUrlStatusOk", url: "https://www.google.com", wantErr: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := crawl.GetUrlDocument(tc.url)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetUrlDocument() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
