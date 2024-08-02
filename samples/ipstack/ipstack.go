// SPDX-License-Identifier: MIT
//go:build samples

package ipstack

import (
	"encoding/json"
	"errors"
	"github.com/spandigital/with"
	"net/http"
	"net/url"
	"time"
)

type options struct {
	httpClient *http.Client
	apiKey     string
	baseURL    *url.URL
}

func WithHttpClient(httpClient *http.Client) with.Func[options] {
	return func(options *options) (err error) {
		options.httpClient = httpClient
		return
	}
}

func WithAPIKey(apiKey string) with.Func[options] {
	return func(options *options) (err error) {
		options.apiKey = apiKey
		return
	}
}

func WithBaseURL(baseURL *url.URL) with.Func[options] {
	return func(options *options) (err error) {
		options.baseURL = baseURL
		return
	}
}

func WithRawBaseURL(rawBaseURL string) with.Func[options] {
	return func(options *options) (err error) {
		options.baseURL, err = url.Parse(rawBaseURL)
		return
	}
}

type IPLocation struct {
	Ip            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      struct {
		GeonameId int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEu                    bool   `json:"is_eu"`
	} `json:"location"`
	TimeZone struct {
		Id               string    `json:"id"`
		CurrentTime      time.Time `json:"current_time"`
		GmtOffset        int       `json:"gmt_offset"`
		Code             string    `json:"code"`
		IsDaylightSaving bool      `json:"is_daylight_saving"`
	} `json:"time_zone"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Plural       string `json:"plural"`
		Symbol       string `json:"symbol"`
		SymbolNative string `json:"symbol_native"`
	} `json:"currency"`
	Connection struct {
		Asn int    `json:"asn"`
		Isp string `json:"isp"`
	} `json:"connection"`
	Security struct {
		IsProxy     bool        `json:"is_proxy"`
		ProxyType   interface{} `json:"proxy_type"`
		IsCrawler   bool        `json:"is_crawler"`
		CrawlerName interface{} `json:"crawler_name"`
		CrawlerType interface{} `json:"crawler_type"`
		IsTor       bool        `json:"is_tor"`
		ThreatLevel string      `json:"threat_level"`
		ThreatTypes interface{} `json:"threat_types"`
	} `json:"security"`
}

type IPStack interface {
	GetIPLocation() (ipLocation *IPLocation, err error)
}

type ipstack struct {
	httpClient *http.Client
	url        string
}

func NewIPStack(withOptions ...with.Func[options]) (newIpStack *ipstack, err error) {
	var newOptions *options
	if newOptions, err = with.Build(&options{
		httpClient: http.DefaultClient,
		apiKey:     "",
		baseURL:    must(url.Parse("https://api.ipstack.com")),
	}, func(options *options) (err error) {
		if options.apiKey == "" {
			return errors.New("apiKey is required")
		}
		return
	}, withOptions...); err == nil {
		url := *newOptions.baseURL
		url.Query().Set("access_key", newOptions.apiKey)
		newIpStack = &ipstack{
			httpClient: newOptions.httpClient,
			url:        url.RawQuery,
		}
	}
	return
}

func (i *ipstack) GetIPLocation() (ipLocation *IPLocation, err error) {
	var response *http.Response
	if response, err = i.httpClient.Get(i.url); err == nil {
		defer response.Body.Close()
		var ipLocation IPLocation
		err = json.NewDecoder(response.Body).Decode(&ipLocation)
	}
	return
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
