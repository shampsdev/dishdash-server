package twogis

import "C"
import (
	"encoding/json"
	"fmt"
	"io"
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

func getParamsMap(tags []string, lon, lat float64, radiusOptional ...int) map[string]string {
	radius := 4000
	if len(radiusOptional) > 0 && radiusOptional[0] > 0 {
		radius = radiusOptional[0]
	}

	return map[string]string{
		"q":      joinTags(tags),
		"point":  fmt.Sprintf("%f,%f", lon, lat),
		"fields": "items.point,items.name,items.description,items.external_content,items.rubrics,items.reviews,items.attribute_groups,items.schedule",
		"radius": strconv.Itoa(radius),
		"key":    ApiKey,
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
		return "", err
	}

	query := reqUrl.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	reqUrl.RawQuery = query.Encode()

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ParseApiResponse(responseBody string) ([]domain.TwoGisPlace, error) {
	var response ApiResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return nil, err
	}

	var twoGisPlaces []domain.TwoGisPlace
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
						return nil, err
					}
					averagePrice = number
				}
			}
		}

		twoGisPlaces = append(twoGisPlaces, domain.TwoGisPlace{
			Name:         item.Name,
			Address:      item.AddressName,
			Lat:          item.Point.Lat,
			Lon:          item.Point.Lon,
			PhotoURL:     photoURL,
			ReviewRating: item.Reviews.GeneralRating,
			ReviewCount:  item.Reviews.GeneralReviewCount,
			Rubrics:      rubrics,
			AveragePrice: averagePrice,
		})
	}

	return twoGisPlaces, nil
}

// TODO tags!!!!
func FetchPlacesForLobbyFromAPI(lobby *domain.Lobby) ([]*domain.Place, error) {
	params := getParamsMap(lobby.TagNames(), lobby.Location.Lon, lobby.Location.Lat)

	apiResponse, err := GetPlacesFromApi(params)
	if err != nil {
		return nil, err
	}

	apiPlaces, _ := ParseApiResponse(apiResponse)

	places := make([]*domain.Place, len(apiPlaces))

	for i, apiPlace := range apiPlaces {
		place := apiPlace.ToPlace()

		for _, rubric := range apiPlace.Rubrics {
			fmt.Println(rubric) // TODO тут как-то привязать тег
		}
		places[i] = place
	}
	return places, nil
}
