package spoonacular

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client mendefinisikan klien untuk Spoonacular API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient membuat instance baru Client.
func NewClient(apiKey, baseURL string, timeout time.Duration) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GenerateMealPlan memanggil endpoint /mealplanner/generate.
func (c *Client) GenerateMealPlan(ctx context.Context, targetCalories int, timeFrame string) (*MealPlanResponse, error) {
	url := fmt.Sprintf("%s/mealplanner/generate?timeFrame=%s&targetCalories=%d&apiKey=%s", c.baseURL, timeFrame, targetCalories, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("generate meal plan: status code tidak terduga %d", resp.StatusCode)
	}

	var mealPlan MealPlanResponse
	if err := json.NewDecoder(resp.Body).Decode(&mealPlan); err != nil {
		return nil, err
	}
	return &mealPlan, nil
}

// GetRecipeInformation memanggil endpoint /recipes/{id}/information dengan parameter includeNutrition.
func (c *Client) GetRecipeInformation(ctx context.Context, recipeID int, includeNutrition bool) (*RecipeInformation, error) {
	url := fmt.Sprintf("%s/recipes/%d/information?includeNutrition=%t&apiKey=%s", c.baseURL, recipeID, includeNutrition, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get recipe information: status code tidak terduga %d", resp.StatusCode)
	}

	var recipeInfo RecipeInformation
	if err := json.NewDecoder(resp.Body).Decode(&recipeInfo); err != nil {
		return nil, err
	}
	return &recipeInfo, nil
}
