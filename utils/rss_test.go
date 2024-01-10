package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestFetchRssFeedValidURL tests fetching an RSS feed with a valid URL
func TestFetchRssFeedValidURL(t *testing.T) {
	// You should replace this with a URL you know is valid and stable for testing
	validURL := "https://skymatixva.com/tfiles/feed.xml"

	rss := NewRss(validURL, logrus.New())
	assert.NotNil(t, *rss)

	fullInstallItems := rss.GetFullInstallItems()
	assert.NotEqual(t, 0, len(*fullInstallItems))

	patchInstallItems := rss.GetPatchInstallItems()
	assert.NotEqual(t, 0, len(*patchInstallItems))

	latestVersion := rss.GetLatestVersion()
	assert.NotEqual(t, "", latestVersion)

	for _, fullInstallItem := range *fullInstallItems {
		assert.NotNil(t, fullInstallItem)
		assert.NotNil(t, fullInstallItem.Version)
	}
}

// TestFetchRssFeedInvalidURL tests fetching an RSS feed with an invalid URL
func TestFetchRssFeedInvalidURL(t *testing.T) {
	invalidURL := "http://thisisnotarealurl"

	rss := NewRss(invalidURL, logrus.New())
	assert.Nil(t, rss)
}

// TestGetFullInstallItems
