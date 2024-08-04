package twogis

import (
	"dishdash.ru/cmd/server/config"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var ApiKey = config.C.TwoGisApi.Key

func GetParamsMap(tags string, twogisApi string, lon, lat float64, radius int) map[string]string {
	return map[string]string{
		"q":      tags,
		"point":  strconv.FormatFloat(lon, 'f', -1, 64) + "," + strconv.FormatFloat(lat, 'f', -1, 64),
		"fields": "items.point,items.name,items.description,items.external_content,items.rubrics,items.reviews,items.attribute_groups,items.schedule",
		"radius": strconv.Itoa(radius),
		"key":    twogisApi,
	}
}

func GetPlacesFromApi(apiUrl string, params map[string]string) string {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	return string(body)
}
