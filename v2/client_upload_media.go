package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UploadMediaResp struct {
	MediaID          string      `json:"media_id"`
	MediaKey         string      `json:"media_key"`
	MediaIDString    string      `json:"media_id_string"`
	Size             int         `json:"size"`
	ExpiresAfterSecs int         `json:"expires_after_secs"`
	Image            ImageInfoV1 `json:"image"`

	RateLimit *RateLimit `json:"-"`
}

type ImageInfoV1 struct {
	ImageType string `json:"image_type"`
	W         int    `json:"w"`
	H         int    `json:"h"`
}

// SimpleUpload will upload the media and return the media_id
func (c *Client) SimpleUpload(ctx context.Context, filename string, file io.Reader, mediaData, mediaCategory, additionalOwners string) (interface{}, error) {
	if file != nil && mediaData != "" {
		return nil, fmt.Errorf("either provide file or mediaData, not both")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if file != nil {
		part, err := writer.CreateFormFile("media", filename)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}

	} else {
		_ = writer.WriteField("media_data", mediaData)
	}

	if mediaCategory != "" {
		_ = writer.WriteField("media_category", mediaCategory)
	}

	if additionalOwners != "" {
		_ = writer.WriteField("additional_owners", additionalOwners)
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", `https://upload.twitter.com/1.1/media/upload.json`, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")

	if c.Authorizer != nil {
		// Assuming Authorizer has an Add method to add authorization headers
		c.Authorizer.Add(req)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var raw interface{}
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		b, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return nil, fmt.Errorf("decode response error: %w", err)
		}
		return nil, fmt.Errorf("decode response: %w, body: %v", err, string(b))
	}

	return raw, nil
}
