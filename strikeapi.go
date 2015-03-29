package strikeapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIEndpoint represents the APIEnpoint
var APIEndpoint = "https://getstrike.net/api/v2"

// Custom errors
var (
	ErrEmptyHashes = errors.New("empty hash array given")
)

// Categories
const (
	Anime        = "Anime"
	Applications = "Applications"
	Books        = "Books"
	Games        = "Games"
	Movies       = "Movies"
	Music        = "Music"
	Other        = "Other"
	TV           = "TV"
	XXX          = "XXX"
)

// SubCategories
const (
	HighresMovies     = "Highres Movies"
	Hentai            = "Hentai"
	HDVideo           = "HD Video"
	Handheld          = "Handheld"
	Fiction           = "Fiction"
	EnglishTranslated = "English-translated"
	Ebooks            = "Ebooks"
	DubbedMovies      = "Dubbed Movies"
	Documentary       = "Documentary"
	Concerts          = "Concerts"
	Comics            = "Comics"
	Bollywood         = "Bollywood"
	AudioBooks        = "Audio books"
	Asian             = "Asian"
	AnimeMusicVideo   = "Anime Music Video"
	Animation         = "Animation"
	Android           = "Android"
	Academic          = "Academic"
	AAC               = "AAC"
	Movies3D          = "3D Movies"
	XBOX360           = "XBOX360"
	Windows           = "Windows"
	Wii               = "Wii"
	Wallpapers        = "Wallpapers"
	Video             = "Video"
	Unsorted          = "Unsorted"
	UNIX              = "UNIX"
	UltraHD           = "UltraHD"
	Tutorials         = "Tutorials"
	Transcode         = "Transcode"
	Trailer           = "Trailer"
	Textbooks         = "Textbooks"
	Subtitles         = "Subtitles"
	Soundtrack        = "Soundtrack"
	SoundClips        = "Sound clips"
	RadioShows        = "Radio Shows"
	PSP               = "PSP"
	PS3               = "PS3"
	PS2               = "PS2"
	Poetry            = "Poetry"
	Pictures          = "Pictures"
	PC                = "PC"
	OtherXXX          = "Other XXX"
	OtherTV           = "Other TV"
	OtherMusic        = "Other Music"
	OtherMovies       = "Other Movies"
	OtherGames        = "Other Games"
	OtherBooks        = "Other Books"
	OtherApplications = "Other Applications"
	OtherAnime        = "Other Anime"
	NonFiction        = "Non-fiction"
	Newspapers        = "Newspapers"
	MusicVideos       = "Music videos"
	Mp3               = "Mp3"
	MovieClips        = "Movie clips"
	Magazines         = "Magazines"
	Mac               = "Mac"
	Lossless          = "Lossless"
	Linux             = "Linux"
	Karaoke           = "Karaoke"
	iOS               = "iOS"
)

// Response represents the information given by the API
type Response struct {
	ResultSize   int       `json:"results"`
	Status       int       `json:"statuscode"`
	ResponseTime float64   `json:"responsetime"`
	Torrents     []Torrent `json:"torrents"`
}

// ResponseStatus represents a basic response information given by the API
type ResponseStatus struct {
	Status  int    `json:"statuscode"`
	Message string `json:"message"`
}

// ResponseStatusInt represents a basic response information given by the API
// where message is a int
type ResponseStatusInt struct {
	Status  int `json:"statuscode"`
	Message int `json:"message"`
}

// FilesInfo represents an array of FilesInfo returned by the API
type FilesInfo struct {
	FileInfo []FileInfo
}

// FileInfo reprensents information about a File returned by the API
type FileInfo struct {
	FileName string
	FileSize float64
}

// Torrent represents a torrent
type Torrent struct {
	Title            string     `json:"torrent_title"`
	Hash             string     `json:"torrent_hash"`
	Category         string     `json:"torrent_category"`
	SubCategory      string     `json:"sub_category"`
	Seeds            int        `json:"seeds"`
	Leeches          int        `json:"leeches"`
	FileCount        int        `json:"file_count,omitempty"`
	DownloadCount    int        `json:"download_count,omitempty"`
	Page             string     `json:"page,omitempty"`
	RSSFeed          string     `json:"rss_feed,omitempty"`
	Size             float64    `json:"size"`
	UploadDate       string     `json:"upload_date"`
	UploaderUsername string     `json:"uploader_username"`
	MagnetURI        string     `json:"magnet_uri"`
	FilesInfo        *FilesInfo `json:"file_info"`
}

// Special date struct to unmarshall properly
type torrentDate struct{ *time.Time }

// UnmarshalJSON is a custom unmarshal function to handle FileInfo struct
func (f *FilesInfo) UnmarshalJSON(data []byte) error {
	dataBytes := bytes.NewReader(data)
	var tempFilesInfo struct {
		FileNames   []string  `json:"file_names"`
		FileLenghts []float64 `json:"file_lengths"`
	}

	files := []FileInfo{}

	// Decode json into the aux struct
	if err := json.NewDecoder(dataBytes).Decode(&tempFilesInfo); err != nil {
		return err
	}

	// Set the FileInfos
	for i, name := range tempFilesInfo.FileNames {
		fileInfo := FileInfo{FileName: name, FileSize: tempFilesInfo.FileLenghts[i]}
		files = append(files, fileInfo)
	}
	f.FileInfo = files

	return nil
}

// GetTorrentsInfos will get all the infos from a list of Torrent
func GetTorrentsInfos(hashes []string) ([]Torrent, error) {
	// Check arguments
	if len(hashes) == 0 {
		return nil, ErrEmptyHashes
	}
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/torrents/info/", APIEndpoint))
	if err != nil {
		return nil, err
	}
	// Add parameters
	urlValues := &url.Values{}
	urlValues.Add("hashes", strings.Join(hashes, ","))
	u.RawQuery = urlValues.Encode()

	// Make the request
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Counldn't make the GET ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	response, err := newResponse(resp)
	if err != nil {
		return nil, err
	}

	return response.Torrents, nil
}

// GetTorrentInfos will get all the infos of a Torrent from a hash
func GetTorrentInfos(hash string) (*Torrent, error) {
	torrentList, err := GetTorrentsInfos([]string{hash})
	if err != nil {
		return nil, err
	}
	if len(torrentList) == 0 {
		return nil, nil
	}
	return &torrentList[0], nil
}

// CountTorrents will return the number of torrents
func CountTorrents() (int, error) {
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/torrents/count/", APIEndpoint))
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Counldn't make the GET ", err)
		return 0, err
	}
	defer resp.Body.Close()

	// Parse the response
	response, err := newResponseStatusInt(resp)
	if err != nil {
		return 0, err
	}

	if response.Status != 200 {
		return 0, fmt.Errorf("Error while getting the number of torrents")
	}

	return response.Message, nil
}

// GetDescription will get the description of a hash torrent
func GetDescription(hash string) (string, error) {
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/torrents/descriptions/", APIEndpoint))
	if err != nil {
		return "", err
	}
	urlValues := &url.Values{}
	urlValues.Add("hash", hash)

	u.RawQuery = urlValues.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Counldn't make the GET ", err)
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	response, err := newResponseStatus(resp)
	if err != nil {
		return "", err
	}
	if response.Status != 200 {
		return "", fmt.Errorf("Error while getting the number of torrents")
	}

	description, err := base64.StdEncoding.DecodeString(response.Message)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(description), nil
}

// GetDescription will get a download link of a Torrent
func (t *Torrent) GetDescription() (string, error) {
	return GetDescription(t.Hash)
}

// Search will search for torrents
func Search(phrase string) ([]Torrent, error) {
	return SearchWithCategoryAndSubCategory(phrase, "", "")
}

// SearchWithCategory will search for torrents with category
func SearchWithCategory(phrase, category string) ([]Torrent, error) {
	return SearchWithCategoryAndSubCategory(phrase, category, "")
}

// SearchWithCategoryAndSubCategory will search with category and subcategory
func SearchWithCategoryAndSubCategory(phrase, category, subCategory string) ([]Torrent, error) {
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/torrents/search/", APIEndpoint))
	if err != nil {
		return nil, err
	}
	urlValues := &url.Values{}
	urlValues.Add("phrase", phrase)
	if category != "" {
		urlValues.Add("category", category)
	}
	if subCategory != "" {
		urlValues.Add("subcategory", subCategory)
	}
	u.RawQuery = urlValues.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Counldn't make the GET ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	response, err := newResponse(resp)
	if err != nil {
		return nil, err
	}

	return response.Torrents, nil
}

// GetDownloadLink will get a download link of a Torrent
func (t *Torrent) GetDownloadLink() (string, error) {
	return GetDownloadLink(t.Hash)
}

// GetDownloadLink will get a download link of a Torrent from a hash
func GetDownloadLink(hash string) (string, error) {
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/torrents/download/", APIEndpoint))
	if err != nil {
		return "", err
	}
	urlValues := &url.Values{}
	urlValues.Add("hash", hash)

	u.RawQuery = urlValues.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Counldn't make the GET ", err)
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	response, err := newResponseStatus(resp)
	if err != nil {
		return "", err
	}
	return response.Message, nil
}

// newResponse will parse a Response struct
func newResponse(resp *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read response body", err)
		return nil, err
	}

	response := &Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Couln't unmarshall result : ", err)
		return nil, err
	}
	return response, nil
}

// newResponse will parse a ResponseStatus struct
func newResponseStatus(resp *http.Response) (*ResponseStatus, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read response body", err)
		return nil, err
	}

	response := &ResponseStatus{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Couln't unmarshall result : ", err)
		return nil, err
	}
	return response, nil
}

// newResponse will parse a ResponseStatusInt struct
func newResponseStatusInt(resp *http.Response) (*ResponseStatusInt, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read response body", err)
		return nil, err
	}

	response := &ResponseStatusInt{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Couln't unmarshall result : ", err)
		return nil, err
	}
	return response, nil
}
