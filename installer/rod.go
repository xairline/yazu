package installer

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

func GetNumberOfPages(page *rod.Page) (int, error) {
	res, err := strconv.Atoi(
		strings.TrimSpace(
			strings.Split(page.
				MustElementX("//li[contains(@class, 'ipsPagination_pageJump')]/a").
				MustText(), "of ")[1],
		),
	)
	if err != nil {
		logrus.New().Error(err)
		return 0, err
	}
	return res, nil
}

func GetIconBase64(iconElem *rod.Element) (string, error) {
	imageURL := iconElem.MustAttribute("data-src")

	if imageURL == nil {
		return "", errors.New("no url found in icon element")
	}
	// Step 1: Fetch the image content
	response, err := http.Get(*imageURL)
	if err != nil {
		logrus.New().Error(err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logrus.New().Errorf("failed to download image: %s", response.Status)
		return "", err
	}

	imageData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.New().Error(err)
		return "", err
	}

	// Step 2: Encode the image content to base64
	base64Encoding := base64.StdEncoding.EncodeToString(imageData)
	return base64Encoding, nil
}
