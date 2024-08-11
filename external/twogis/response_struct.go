package twogis

import (
	"encoding/json"
)

type ApiResponse struct {
	Meta struct {
		APIVersion string `json:"api_version"`
		Code       int    `json:"code"`
		IssueDate  string `json:"issue_date"`
	} `json:"meta"`
	Result struct {
		Items []struct {
			AddressName     string `json:"address_name,omitempty"`
			AttributeGroups []struct {
				Attributes []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Tag  string `json:"tag"`
				} `json:"attributes"`
				IconURL   string   `json:"icon_url,omitempty"`
				IsContext bool     `json:"is_context"`
				IsPrimary bool     `json:"is_primary"`
				Name      string   `json:"name"`
				RubricIds []string `json:"rubric_ids"`
			} `json:"attribute_groups"`
			ExternalContent []struct {
				Count        int    `json:"count"`
				MainPhotoURL string `json:"main_photo_url"`
				Subtype      string `json:"subtype"`
				Type         string `json:"type"`
			} `json:"external_content"`
			ID    string `json:"id"`
			Name  string `json:"name"`
			Point struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"point"`
			Reviews struct {
				GeneralRating               float64 `json:"general_rating"`
				GeneralReviewCount          int     `json:"general_review_count"`
				GeneralReviewCountWithStars int     `json:"general_review_count_with_stars"`
				IsReviewable                bool    `json:"is_reviewable"`
				IsReviewableOnFlamp         bool    `json:"is_reviewable_on_flamp"`
				Items                       []struct {
					IsReviewable bool        `json:"is_reviewable"`
					Tag          string      `json:"tag"`
					Rating       json.Number `json:"rating,omitempty"`
					ReviewCount  int         `json:"review_count,omitempty"`
				} `json:"items"`
				OrgRating               float64 `json:"org_rating"`
				OrgReviewCount          int     `json:"org_review_count"`
				OrgReviewCountWithStars int     `json:"org_review_count_with_stars"`
				Rating                  float64 `json:"rating"`
				ReviewCount             int     `json:"review_count"`
			} `json:"reviews"`
			Rubrics []struct {
				Alias    string `json:"alias"`
				ID       string `json:"id"`
				Kind     string `json:"kind"`
				Name     string `json:"name"`
				ParentID string `json:"parent_id"`
				ShortID  int    `json:"short_id"`
			} `json:"rubrics"`
			Schedule struct {
				Fri struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"working_hours"`
				} `json:"Fri"`
				Mon struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Mon"`
				} `json:"Mon"`
				Sat struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Sat"`
				} `json:"Sat"`
				Sun struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Sun"`
				} `json:"Sun"`
				Thu struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Thu"`
				} `json:"Thu"`
				Tue struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Tue"`
				} `json:"Tue"`
				Wed struct {
					WorkingHours []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Wed"`
				} `json:"Wed"`
			} `json:"schedule"`
			Type           string `json:"type"`
			AddressComment string `json:"address_comment,omitempty"`
			FullName       string `json:"full_name,omitempty"`
			PurposeName    string `json:"purpose_name,omitempty"`
		} `json:"items"`
		Total int `json:"total"`
	} `json:"result"`
}
