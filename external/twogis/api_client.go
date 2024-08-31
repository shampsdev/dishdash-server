package twogis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"dishdash.ru/internal/domain"

	"dishdash.ru/cmd/server/config"
)

var (
	ApiKey = config.C.TwoGisApi.Key
	ApiUrl = config.C.TwoGisApi.Url
)

func joinTags(tags []string) string {
	return strings.Join(tags, ",")
}

func getParamsMap(tags []string, lon, lat float64, page, pageSize int, radiusOptional ...int) map[string]string {
	radius := 4000
	if len(radiusOptional) > 0 && radiusOptional[0] > 0 {
		radius = radiusOptional[0]
	}

	return map[string]string{
		"q":         joinTags(tags),
		"point":     fmt.Sprintf("%f,%f", lon, lat),
		"fields":    "items.point,items.name,items.description,items.external_content,items.rubrics,items.reviews,items.attribute_groups,items.schedule",
		"radius":    strconv.Itoa(radius),
		"page":      strconv.Itoa(page),
		"page_size": strconv.Itoa(pageSize),
		"key":       ApiKey,
	}
}

func extractNumber(s string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("no number found in the string")
	}
	number, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func GetPlacesFromApi(params map[string]string) (string, error) {
	reqUrl, err := url.Parse(ApiUrl)
	if err != nil {
		log.Printf("Error parsing API URL: %v", err)
		return "", err
	}

	query := reqUrl.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	reqUrl.RawQuery = query.Encode()

	safeUrlString := reqUrl.String()
    if apiKey := query.Get("key"); apiKey != "" {
        safeUrlString = strings.Replace(safeUrlString, apiKey, "******", 1)
    }
    log.Printf("Sending request to API URL: %s", safeUrlString)

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		log.Printf("Error making HTTP request to API: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	return string(body), nil
}

func ParseApiResponse(responseBody string) ([]*domain.TwoGisPlace, error) {
	var response ApiResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, err
	}

	var twoGisPlaces []*domain.TwoGisPlace
	for _, item := range response.Result.Items {
		var photoURL string
		if len(item.ExternalContent) > 0 {
			photoURL = item.ExternalContent[0].MainPhotoURL
		}

		var rubrics []string
		for _, rubric := range item.Rubrics {
			rubrics = append(rubrics, rubric.Name)
		}

		var averagePrice int
		for _, group := range item.AttributeGroups {
			for _, attribute := range group.Attributes {
				if attribute.Tag == "food_service_avg_price" {
					number, err := extractNumber(attribute.Name)
					if err != nil {
						log.Printf("Error extracting average price: %v", err)
						return nil, err
					}
					averagePrice = number
				}
			}
		}

		twoGisPlace := &domain.TwoGisPlace{
			Name:         item.Name,
			Address:      item.AddressName,
			Lat:          item.Point.Lat,
			Lon:          item.Point.Lon,
			PhotoURL:     photoURL,
			ReviewRating: item.Reviews.GeneralRating,
			ReviewCount:  item.Reviews.GeneralReviewCount,
			Rubrics:      rubrics,
			AveragePrice: averagePrice,
		}
		log.Printf("Processed place: %s, Address: %s", twoGisPlace.Name, twoGisPlace.Address)
		twoGisPlaces = append(twoGisPlaces, twoGisPlace)
	}

	return twoGisPlaces, nil
}

func FetchPlacesForLobbyFromAPI(lobby *domain.Lobby) ([]*domain.TwoGisPlace, error) {
	var allApiPlaces []*domain.TwoGisPlace
	page := 1
	pageSize := 10

	for {
		log.Printf("Fetching places from API for lobby: %s, Page: %d", lobby.ID, page)
		params := getParamsMap(lobby.TagNames(), lobby.Location.Lon, lobby.Location.Lat, page, pageSize)

		apiResponse, err := GetPlacesFromApi(params)
		if err != nil {
			log.Printf("Error fetching places from API: %v", err)
			return nil, err
		}

		apiPlaces, err := ParseApiResponse(apiResponse)
		if err != nil {
			log.Printf("Error parsing API response: %v", err)
			return nil, err
		}

		if len(apiPlaces) == 0 {
			log.Printf("No more places found, stopping fetch.")
			break
		}

		allApiPlaces = append(allApiPlaces, apiPlaces...)

		if len(apiPlaces) < pageSize {
			log.Printf("Less places returned than requested, likely end of data.")
			break
		}

		page++
	}

	log.Printf("Total places fetched for lobby: %s, Count: %d", lobby.ID, len(allApiPlaces))
	return allApiPlaces, nil
}