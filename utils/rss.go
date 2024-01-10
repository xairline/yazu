package utils

import (
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Define the structures to hold the RSS data
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
	Url     string
	log     *logrus.Logger
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Version     string
	//PubDate     string `xml:"pubDate"`
}

func GetRss(url string, log *logrus.Logger) *Rss {
	// Fetch the RSS feed
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error fetching the RSS feed: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading the response body: %v", err)
		return nil
	}

	// Parse the XML
	var rss Rss
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Errorf("Error parsing XML: %v", err)
		return nil
	}
	rss.Url = url
	rss.log = log
	return &rss
}

func NewRss(url string, log *logrus.Logger) *Rss {
	return GetRss(url, log)
}

func (r *Rss) GetFullInstallItems() *[]Item {
	var items []Item
	versionRegEx, _ := regexp.Compile(`XP12.*\.zip\.torrent`)
	for _, item := range r.Channel.Items {
		if strings.Contains(item.Title, "XP12") &&
			strings.Contains(item.Title, "_full") {
			matchedVersion := versionRegEx.FindString(item.Link)
			matchedVersion = strings.Replace(matchedVersion, "XP12_", "", 1)
			matchedVersion = strings.Replace(matchedVersion, "_full.zip.torrent", "", 1)
			matchedVersion = strings.Replace(matchedVersion, "_", ".", 2)
			item.Version = matchedVersion
			items = append(items, item)
		}
	}
	return &items
}

func (r *Rss) GetPatchInstallItems() *[]Item {
	var items []Item
	versionRegEx, _ := regexp.Compile(`XP12.*\.zip\.torrent`)
	for _, item := range r.Channel.Items {
		if strings.Contains(item.Title, "XP12") &&
			!strings.Contains(item.Title, "_full") {
			matchedVersion := versionRegEx.FindString(item.Link)
			matchedVersion = strings.Replace(matchedVersion, "XP12_", "", 1)
			matchedVersion = strings.Replace(matchedVersion, ".zip.torrent", "", 1)
			matchedVersion = strings.Replace(matchedVersion, "_", ".", 2)
			item.Version = matchedVersion
			items = append(items, item)
		}
	}
	return &items
}

func (r *Rss) GetLatestVersion() string {
	patchItems := *r.GetPatchInstallItems()
	if len(patchItems) > 0 {
		return patchItems[len(patchItems)-1].Version
	} else {
		fullItems := *r.GetFullInstallItems()
		return fullItems[len(fullItems)-1].Version
	}
}
