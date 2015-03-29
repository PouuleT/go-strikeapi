package strikeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestEmptyMessage tests if the message empty
func TestEmptyHashArray(t *testing.T) {
	emptyHashes := []string{}

	torrentList, err := GetTorrentsInfos(emptyHashes)
	if torrentList != nil {
		t.Errorf("Shouldn't get a torrent list")
	}
	if err != ErrEmptyHashes {
		t.Errorf("Should get an ErrEmptyHashes")
	}
}

func TestGetTorrentInfos(t *testing.T) {
	rawHTMLResponse := `{"results":1,"statuscode":200,"responsetime":0.0031,"torrents":[{"torrent_hash":"B425907E5755031BDA4A8D1B6DCCACA97DA14C04","torrent_title":"Arch Linux 2015.01.01 (x86/x64)","torrent_category":"Applications","sub_category":"","seeds":645,"leeches":13,"file_count":1,"size":615514112,"upload_date":"Jan  6, 2015","uploader_username":"The_Doctor-","file_info":{"file_names":["archlinux-2015.01.01-dual.iso"],"file_lengths":[615514112]},"magnet_uri":"magnet:?xt=urn:btih:B425907E5755031BDA4A8D1B6DCCACA97DA14C04&dn=Arch+Linux+2015.01.01+%28x86%2Fx64%29&tr=udp://open.demonii.com:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://exodus.desync.com:6969"}]}`

	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()
	APIEndpoint = ts.URL

	torrent, err := GetTorrentInfos("B425907E5755031BDA4A8D1B6DCCACA97DA14C04")
	if err != nil {
		t.Errorf("Error getting a torrent infos")
	}

	// Expected result
	expectedFilesInfo := &FilesInfo{
		FileInfo: []FileInfo{
			{
				FileName: "archlinux-2015.01.01-dual.iso",
				FileSize: 615514112,
			},
		},
	}
	expectedTorrents := &Torrent{
		Title:            "Arch Linux 2015.01.01 (x86/x64)",
		Hash:             "B425907E5755031BDA4A8D1B6DCCACA97DA14C04",
		Category:         "Applications",
		SubCategory:      "",
		Seeds:            645,
		Leeches:          13,
		FileCount:        1,
		Size:             615514112,
		UploadDate:       "Jan  6, 2015",
		UploaderUsername: "The_Doctor-",
		MagnetURI:        "magnet:?xt=urn:btih:B425907E5755031BDA4A8D1B6DCCACA97DA14C04&dn=Arch+Linux+2015.01.01+%28x86%2Fx64%29&tr=udp://open.demonii.com:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://exodus.desync.com:6969",
		FilesInfo:        expectedFilesInfo,
	}

	if reflect.DeepEqual(torrent, expectedTorrents) == false {
		t.Errorf("Response not properly set")
	}
}

func TestGetDescription(t *testing.T) {
	rawHTMLResponse := `{"statuscode":200,"message":"VGhpcyB0b3JyZW50IGhhcyBubyBkZXNjcmlwdGlvbg=="}`
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()

	APIEndpoint = ts.URL

	desc, err := GetDescription("B425907E5755031BDA4A8D1B6DCCACA97DA14C04")
	if err != nil {
		t.Errorf("Error getting description from hash")
	}
	if desc != "This torrent has no description" {
		t.Errorf("Bad description")
	}
}

func TestGetDownloadLink(t *testing.T) {
	rawHTMLResponse := `{"statuscode":200,"message":"https://getstrike.net/torrents/api/download/0EB6605E041F1846B84BAA63346012A82706A95D.torrent"}`
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()

	APIEndpoint = ts.URL

	link, err := GetDownloadLink("B425907E5755031BDA4A8D1B6DCCACA97DA14C04")
	if err != nil {
		t.Errorf("Error getting download link from hash")
	}
	if link != "https://getstrike.net/torrents/api/download/0EB6605E041F1846B84BAA63346012A82706A95D.torrent" {
		t.Errorf("Bad download link")
	}
}

func TestSearch(t *testing.T) {
	rawHTMLResponse := `{"results":1,"statuscode":200,"responsetime":0.4725,"torrents":[{"torrent_hash":"156B69B8643BD11849A5D8F2122E13FBB61BD041","torrent_title":"Slackware 14.1 x86_64 DVD ISO","torrent_category":"Applications","sub_category":"","seeds":192,"leeches":9,"file_count":4,"size":2437393940.48,"download_count":40,"upload_date":"Feb 24, 2014","uploader_username":"Nusantara","page":"https://getstrike.net/torrents/156B69B8643BD11849A5D8F2122E13FBB61BD041","rss_feed":"https://getstrike.net/torrents/156B69B8643BD11849A5D8F2122E13FBB61BD041?rss=1","magnet_uri":"magnet:?xt=urn:btih:156B69B8643BD11849A5D8F2122E13FBB61BD041&dn=Slackware+14.1+x86_64+DVD+ISO&tr=udp://open.demonii.com:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://exodus.desync.com:6969"}]}`

	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()

	APIEndpoint = ts.URL

	torrentList, err := SearchWithCategory("Slackware 14.1 x86_64 DVD ISO", Applications)
	if err != nil {
		t.Errorf("Error searching for torrent")
	}
	// Expected result
	expectedTorrents := []Torrent{
		{
			Title:            "Slackware 14.1 x86_64 DVD ISO",
			Hash:             "156B69B8643BD11849A5D8F2122E13FBB61BD041",
			Category:         "Applications",
			SubCategory:      "",
			Seeds:            192,
			Leeches:          9,
			FileCount:        4,
			DownloadCount:    40,
			Size:             2437393940.48,
			UploadDate:       "Feb 24, 2014",
			UploaderUsername: "Nusantara",
			Page:             "https://getstrike.net/torrents/156B69B8643BD11849A5D8F2122E13FBB61BD041",
			RSSFeed:          "https://getstrike.net/torrents/156B69B8643BD11849A5D8F2122E13FBB61BD041?rss=1",
			MagnetURI:        "magnet:?xt=urn:btih:156B69B8643BD11849A5D8F2122E13FBB61BD041&dn=Slackware+14.1+x86_64+DVD+ISO&tr=udp://open.demonii.com:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://exodus.desync.com:6969",
		},
	}

	if reflect.DeepEqual(torrentList, expectedTorrents) == false {
		t.Errorf("Torrent result not properly set")
	}
}

func TestCountTorrents(t *testing.T) {
	rawHTMLResponse := `{"statuscode":200,"message":6355272}`
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()

	APIEndpoint = ts.URL

	count, err := CountTorrents()
	if err != nil {
		t.Errorf("Error counting torrents")
	}
	// Expected result
	if count != 6355272 {
		t.Errorf("Bad count response not properly set")
	}
}
