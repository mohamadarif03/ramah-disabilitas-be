package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// ExtractVideoID mengambil ID video dari URL Youtube
func ExtractVideoID(url string) string {
	// Regex untuk berbagai format URL Youtube
	re := regexp.MustCompile(`(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GetYoutubeTranscript mengambil transkrip/caption dari video Youtube
func GetYoutubeTranscript(videoID string) (string, error) {
	if videoID == "" {
		return "", errors.New("video ID kosong")
	}

	// 1. Get Video Page
	resp, err := http.Get("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	htmlContent := string(bodyBytes)

	// 2. Find ytInitialPlayerResponse
	// Mencari JSON object di dalam script
	re := regexp.MustCompile(`var ytInitialPlayerResponse = (\{.*?\});`)
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) < 2 {
		return "", errors.New("gagal mengambil data player (ytInitialPlayerResponse tidak ditemukan)")
	}
	jsonStr := matches[1]

	var playerResponse struct {
		Captions struct {
			PlayerCaptionsTracklistRenderer struct {
				CaptionTracks []struct {
					BaseUrl string `json:"baseUrl"`
					Name    struct {
						SimpleText string `json:"simpleText"`
					} `json:"name"`
					LanguageCode string `json:"languageCode"`
				} `json:"captionTracks"`
			} `json:"playerCaptionsTracklistRenderer"`
		} `json:"captions"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &playerResponse); err != nil {
		return "", errors.New("gagal parsing data player Youtube")
	}

	tracks := playerResponse.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks
	if len(tracks) == 0 {
		return "", errors.New("video ini tidak memiliki caption/transkrip otomatis")
	}

	// 3. Select Track (Prioritas: Indonesia -> Inggris -> Lainnya)
	var selectedURL string
	for _, track := range tracks {
		if strings.HasPrefix(track.LanguageCode, "id") { // Indonesia
			selectedURL = track.BaseUrl
			break
		}
	}
	if selectedURL == "" {
		for _, track := range tracks {
			if strings.HasPrefix(track.LanguageCode, "en") { // Inggris
				selectedURL = track.BaseUrl
				break
			}
		}
	}
	if selectedURL == "" {
		selectedURL = tracks[0].BaseUrl // Fallback ke yang pertama
	}

	// 4. Get Transcript XML
	respTrans, err := http.Get(selectedURL)
	if err != nil {
		return "", err
	}
	defer respTrans.Body.Close()

	bodyTrans, err := io.ReadAll(respTrans.Body)
	if err != nil {
		return "", err
	}

	// 5. Parse XML
	// Format: <transcript><text start="0" dur="2">Hello</text>...</transcript>
	type Text struct {
		Content string `xml:",chardata"`
	}
	type Transcript struct {
		Texts []Text `xml:"text"`
	}

	var t Transcript
	if err := xml.Unmarshal(bodyTrans, &t); err != nil {
		return "", errors.New("gagal parsing XML transkrip")
	}

	var fullText strings.Builder
	for _, item := range t.Texts {
		decoded := html.UnescapeString(item.Content)
		fullText.WriteString(decoded)
		fullText.WriteString(" ")
	}

	return fullText.String(), nil
}
