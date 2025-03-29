package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Follow struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Login       string    `json:"login"`
	Avatar      string    `json:"avatar"`
	FollowedAt  time.Time `json:"followedAt"`
	IsLive      bool      `json:"isLive"`
}

type Mod struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Login       string    `json:"login"`
	Avatar      string    `json:"avatar"`
	GrantedAt   time.Time `json:"granted_at"`
	Banned      bool      `json:"banned"`
}

type Vip struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Login       string    `json:"login"`
	Avatar      string    `json:"avatar"`
	GrantedAt   time.Time `json:"granted_at"`
	Banned      bool      `json:"banned"`
}

type Founders struct {
	ID           string    `json:"id"`
	DisplayName  string    `json:"displayName"`
	Login        string    `json:"login"`
	FirstMonth   time.Time `json:"firstMonth"`
	IsSubscribed bool      `json:"isSubscribed"`
	Avatar       string    `json:"avatar"`
	Banned       bool      `json:"banned"`
}

// Fetcher handles HTTP requests to retrieve Twitch channel data.
type Fetcher struct {
	client *http.Client
}

// NewFetcher creates a new Fetcher instance with a configured HTTP client.
//
// Returns:
//
//	A pointer to a new Fetcher instance
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchFollows retrieves the list of users followed by the specified Twitch user.
//
// Parameters:
//
//	ctx - Context for controlling request cancellation
//	username - Twitch username to fetch follows for
//
// Returns:
//
//	A slice of Follow structs and an error if any
func (f *Fetcher) FetchFollows(ctx context.Context, username string) ([]Follow, error) {
	url := fmt.Sprintf("https://tools.2807.eu/api/getfollows/%s", username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch follows: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("the user does not follow any channel")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var follows []Follow
	if err := json.NewDecoder(resp.Body).Decode(&follows); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return follows, nil
}

// FetchMods retrieves the list of moderators for the specified Twitch channel.
//
// Parameters:
//
//	ctx - Context for controlling request cancellation
//	username - Twitch username to fetch moderators for
//
// Returns:
//
//	A slice of Mod structs and an error if any
func (f *Fetcher) FetchMods(ctx context.Context, username string) ([]Mod, error) {
	url := fmt.Sprintf("https://tools.2807.eu/api/getmods/%s", username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mods: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("the user does not have any moderators on their channel")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var mods []Mod
	if err := json.NewDecoder(resp.Body).Decode(&mods); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return mods, nil
}

// FetchVips retrieves the list of VIPs for the specified Twitch channel.
//
// Parameters:
//
//	ctx - Context for controlling request cancellation
//	username - Twitch username to fetch VIPs for
//
// Returns:
//
//	A slice of Vip structs and an error if any
func (f *Fetcher) FetchVips(ctx context.Context, username string) ([]Vip, error) {
	url := fmt.Sprintf("https://tools.2807.eu/api/getvips/%s", username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vips: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("the user does not have any VIPs on their channel")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vips []Vip
	if err := json.NewDecoder(resp.Body).Decode(&vips); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return vips, nil
}

// FetchFounders retrieves the list of founders for the specified Twitch channel.
//
// Parameters:
//
//	ctx - Context for controlling request cancellation
//	username - Twitch username to fetch founders for
//
// Returns:
//
//	A slice of Founders structs and an error if any
func (f *Fetcher) FetchFounders(ctx context.Context, username string) ([]Founders, error) {
	url := fmt.Sprintf("https://tools.2807.eu/api/getfounders/%s", username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch founders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("the user does not have any founders on their channel")
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var founders []Founders
	if err := json.NewDecoder(resp.Body).Decode(&founders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return founders, nil
}
