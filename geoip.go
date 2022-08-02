package geoip

import (
	"bufio"
	"fmt"
	"github.com/civet148/log"
	"os"
	"strings"
)

const (
	CZNET = "纯真网络"
	CZ88  = "CZ88.NET"
)

type GeoLocation struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
}

type GeoData struct {
	IPBegin       string `json:"ip_begin"`
	IPBeginUint32 uint32 `json:"ip_begin_uint32"`
	IPEnd         string `json:"ip_end"`
	IPEndUint32   uint32 `json:"ip_end_uint32"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Area          string `json:"area"`
}

type GeoIP struct {
	data []*GeoData
}

func NewGeoIP(strDatFile string) (*GeoIP, error) {
	data, err := loadGeoData(strDatFile)
	if err != nil {
		return nil, err
	}
	return &GeoIP{
		data: data,
	}, nil
}

func (g *GeoIP) Find(ip string) *GeoLocation {
	uip := IP2Uint(ip)
	for _, v := range g.data {
		if uip >= v.IPBeginUint32 && uip <= v.IPEndUint32 {
			return &GeoLocation{
				IP:       ip,
				Country:  v.Country,
				Province: v.Province,
				City:     v.City,
				Area:     v.Area,
			}
		}
	}
	return &GeoLocation{}
}

func loadGeoData(strDatFile string) (data []*GeoData, err error) {
	var fil *os.File
	fil, err = os.Open(strDatFile)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	log.Enter()
	defer log.Leave()

	reader := bufio.NewScanner(fil)
	for reader.Scan() {
		strLine := reader.Text()
		if strings.TrimSpace(strLine) == "" {
			break
		}
		var strIPBegin, strIPEnd, strCountryOrProvince, strArea string
		_, err = fmt.Sscanf(strLine, "%s %s %s %s", &strIPBegin, &strIPEnd, &strCountryOrProvince, &strArea)
		if err != nil {
			fmt.Sscanf(strLine, "%s %s %s", &strIPBegin, &strIPEnd, &strCountryOrProvince)
		}

		if strIPBegin == "" || strIPEnd == "" || strCountryOrProvince == "" {
			log.Errorf("line [%s] parse error", strLine)
			continue
		}
		if strCountryOrProvince == CZNET {
			continue
		}
		if strArea == CZ88 {
			strArea = ""
		}

		v := handleCountryProvinceCity(&GeoData{
			IPBegin:       strIPBegin,
			IPBeginUint32: IP2Uint(strIPBegin),
			IPEnd:         strIPEnd,
			IPEndUint32:   IP2Uint(strIPEnd),
			Country:       strCountryOrProvince,
			Province:      "",
			City:          "",
			Area:          strArea,
		})
		data = append(data, v)
		//log.Infof("begin [%s] end [%s] country [%s] province [%s] city [%s] area [%s]", v.IPBegin, v.IPEnd, v.Country, v.Province, v.City, v.Area)
	}
	return
}

func handleCountryProvinceCity(v *GeoData) *GeoData {
	var strCN = "中国"
	var individualProvinces = []string{"宁夏", "新疆", "广西", "西藏", "内蒙古"}
	var individualCities = []string{"北京市", "天津市", "重庆市", "上海市", "香港"}

	strCountry := v.Country
	if strings.Contains(strCountry, "省") {
		v.Country = strCN
		pc := strings.Split(strCountry, "省")
		if len(pc) == 1 {
			v.Province = pc[0]
		} else if len(pc) == 2 {
			v.Province = pc[0]
			v.City = pc[1]
		}
	} else {
		//特殊省份
		for _, p := range individualProvinces {
			if strings.Contains(strCountry, p) {
				v.Country = strCN
				v.Province = p
				v.City = strings.TrimPrefix(strCountry, p)
				break
			}
		}
		//直辖市
		for _, c := range individualCities {
			if strings.Contains(strCountry, c) {
				v.Country = strCN
				v.Province = c
				v.City = c
				break
			}
		}
	}
	//区或县
	strCity := v.City
	if strings.Contains(strCity, "区") {
		ss := strings.Split(strCity,"市")
		if len(ss) == 2 {
			v.City = ss[0]
			v.Area = ss[1]
		}
	}
	if strings.Contains(strCity, "县") {
		ss := strings.Split(strCity,"市")
		if len(ss) == 2 {
			v.City = ss[0]
			v.Area = ss[1]
		}
	}
	return v
}
