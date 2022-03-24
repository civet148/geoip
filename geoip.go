package client

import (
	"net/url"

	"github.com/civet148/httpc"
	"github.com/civet148/log"
)

type Language string

const (
	Language_CN Language = "zh-CN"
	Language_EN Language = "en-US"
)

const (
	baseUrl = "http://www.geoplugin.net/json.gp"
)

type IpMsg struct {
	Request                string  `json:"geoplugin_request"`
	Status                 int     `json:"geoplugin_status"`
	Delay                  string  `json:"geoplugin_delay"`
	Credit                 string  `json:"geoplugin_credit"`
	City                   string  `json:"geoplugin_city"`
	Region                 string  `json:"geoplugin_region"`
	RegionCode             string  `json:"geoplugin_regionCode"`
	RegionName             string  `json:"geoplugin_regionName"`
	AreaCode               string  `json:"geoplugin_areaCode"`
	DmaCode                string  `json:"geoplugin_dmaCode"`
	CountryCode            string  `json:"geoplugin_countryCode"`
	CountryName            string  `json:"geoplugin_countryName"`
	InEU                   int     `json:"geoplugin_inEU"`
	EuVATrate              bool    `json:"geoplugin_euVATrate"`
	ContinentCode          string  `json:"geoplugin_continentCode"`
	ContinentName          string  `json:"geoplugin_continentName"`
	Latitude               string  `json:"geoplugin_latitude"`
	Longitude              string  `json:"geoplugin_longitude"`
	LocationAccuracyRadius string  `json:"geoplugin_locationAccuracyRadius"`
	Timezone               string  `json:"geoplugin_timezone"`
	CurrencyCode           string  `json:"geoplugin_currencyCode"`
	CurrencySymbol         string  `json:"geoplugin_currencySymbol"`
	CurrencySymbolUTF8     string  `json:"geoplugin_currencySymbol_UTF8"`
	CurrencyConverter      float64 `json:"geoplugin_currencyConverter"`
}

func GetIpMsg(lang Language, ip string) (*IpMsg, error) {
	var c = httpc.NewHttpClient(5)
	var msg IpMsg

	value := url.Values{
		"lang": []string{string(lang)},
		"ip":   []string{ip},
	}
	r, err := c.Get(baseUrl, value)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = r.Unmarshal(&msg)
	return &msg, err
}

func GetLocation(lang Language, ip string) (lng, lat string, err error) {
	var msg *IpMsg
	if msg, err = GetIpMsg(lang, ip); err != nil {
		log.Error(err.Error())
		return
	}
	return msg.Longitude, msg.Latitude, nil
}

func GetFullAddr(lang Language, ip string) (string, error) {
	var err error
	var msg *IpMsg
	if msg, err = GetIpMsg(lang, ip); err != nil {
		log.Error(err.Error())
		return "", err
	}
	strAddr := msg.CountryName + msg.RegionName + msg.City
	return strAddr, nil
}

func GetProvinceCity(lang Language, ip string) (string, error) {
	var err error
	var msg *IpMsg
	if msg, err = GetIpMsg(lang, ip); err != nil {
		log.Error(err.Error())
		return "", err
	}
	return msg.RegionName + msg.City, nil
}

func GetProvince(lang Language, ip string) (string, error) {
	var err error
	var msg *IpMsg
	if msg, err = GetIpMsg(lang, ip); err != nil {
		log.Error(err.Error())
		return "", err
	}
	return msg.RegionName, nil
}

func GetCity(lang Language, ip string) (string, error) {
	var err error
	var msg *IpMsg
	if msg, err = GetIpMsg(lang, ip); err != nil {
		log.Error(err.Error())
		return "", err
	}
	return msg.City, nil
}
