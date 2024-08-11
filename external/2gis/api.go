package twogis

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"dishdash.ru/cmd/server/config"
)

var ApiKey = config.C.TwoGisApi.Key

func GetParamsMap(tags string, lon, lat float64, radiusOptional ...int) map[string]string {
	radius := 2000
	if len(radiusOptional) > 0 && radiusOptional[0] > 0 {
		radius = radiusOptional[0]
	}

	return map[string]string{
		"q":      tags,
		"point":  strconv.FormatFloat(lon, 'f', -1, 64) + "," + strconv.FormatFloat(lat, 'f', -1, 64),
		"fields": "items.point,items.name,items.description,items.external_content,items.rubrics,items.reviews,items.attribute_groups,items.schedule",
		"radius": strconv.Itoa(radius),
		"key":    ApiKey,
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	return string(body)
}
