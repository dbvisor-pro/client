package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func CreatePostRequest(data []byte, url string, token *string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("something wrong with POST request data: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	var errMsg string = ""

	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
		errMsg = "Your token has expired. Use the login command to update it"
	} else {
		errMsg = "Invalid credentials"
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("bad status: %s. %s", resp.Status, errMsg)
		} else {
			return nil, fmt.Errorf("bad status: %s", resp.Status)
		}
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return result, nil
}

func CreateGetRequest(url string, token *string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("something wrong with GET request data: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("bad status: %s. Your token has expired. Use the login command to update it", resp.Status)
		} else {
			return nil, fmt.Errorf("bad status: %s", resp.Status)
		}
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return result, nil
}
