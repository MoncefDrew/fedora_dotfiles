package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	WTTR_API_URL = "https://wttr.in"
	LOCATION     = "Batna,Algeria"
)

// Icons by condition text
var weatherIcons = map[string]string{
	"Clear":             "󰖙",
	"Sunny":             "󰖙",
	"PartlyCloudy":      "󰖕",
	"Cloudy":            "󰖐",
	"Overcast":          "󰖐",
	"Mist":              "󰖑",
	"Fog":               "󰖑",
	"LightRain":         "󰖗",
	"HeavyRain":         "󰖖",
	"Rain":              "󰖗",
	"LightSnow":         "󰖘",
	"HeavySnow":         "󰖘",
	"Snow":              "󰖘",
	"Thunderstorm":      "󰖓",
	"ThunderyShowers":   "󰖓",
	"ThunderySnow":      "󰖓",
	"Drizzle":           "󰖗",
	"LightShowers":      "󰖗",
	"HeavyShowers":      "󰖖",
	"Sleet":             "󰖘",
	"ClearNight":        "󰖔",
	"PartlyCloudyNight": "󰼱",
}

// Icons by weather code (more reliable)
var weatherCodeIcons = map[string]string{
	"113": "󰖙", // Sunny
	"116": "󰖕",
	"119": "󰖐",
	"122": "󰖐",
	"143": "󰖑",
	"176": "󰖗",
	"179": "󰖘",
	"182": "󰖘",
	"185": "󰖘",
	"200": "󰖓",
	"227": "󰖘",
	"230": "󰖘",
	"248": "󰖑",
	"260": "󰖑",
	"263": "󰖗",
	"266": "󰖗",
	"281": "󰖘",
	"284": "󰖘",
	"293": "󰖗",
	"296": "󰖗",
	"299": "󰖗",
	"302": "󰖗",
	"305": "󰖖",
	"308": "󰖖",
	"311": "󰖘",
	"314": "󰖘",
	"317": "󰖘",
	"320": "󰖘",
	"323": "󰖘",
	"326": "󰖘",
	"329": "󰖘",
	"332": "󰖘",
	"335": "󰖘",
	"338": "󰖘",
	"350": "󰖘",
	"353": "󰖗",
	"356": "󰖖",
	"359": "󰖖",
	"362": "󰖘",
	"365": "󰖘",
	"368": "󰖘",
	"371": "󰖘",
	"374": "󰖘",
	"377": "󰖘",
	"386": "󰖓",
	"389": "󰖓",
	"392": "󰖓",
	"395": "󰖓",
}

// API response
type WttrResponse struct {
	CurrentCondition []struct {
		TempC       string `json:"temp_C"`
		FeelsLikeC  string `json:"FeelsLikeC"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		WeatherCode   string `json:"weatherCode"`
		Humidity      string `json:"humidity"`
		WindspeedKmph string `json:"windspeedKmph"`
		WindDir16     string `json:"winddir16Point"`
		PrecipMM      string `json:"precipMM"`
	} `json:"current_condition"`

	NearestArea []struct {
		AreaName []struct {
			Value string `json:"value"`
		} `json:"areaName"`
		Country []struct {
			Value string `json:"value"`
		} `json:"country"`
	} `json:"nearest_area"`
}

// Waybar JSON output
type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Alt     string `json:"alt"`
	Class   string `json:"class"`
}

func getIcon(condition, code string) string {
	clean := strings.ReplaceAll(condition, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")

	if icon, ok := weatherIcons[clean]; ok {
		return icon
	}
	if icon, ok := weatherCodeIcons[code]; ok {
		return icon
	}
	return "󰖐"
}

func getWeather() (WaybarOutput, error) {
	url := fmt.Sprintf("%s/%s?format=j1&lang=en", WTTR_API_URL, LOCATION)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return WaybarOutput{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WaybarOutput{}, err
	}

	var data WttrResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return WaybarOutput{}, err
	}

	if len(data.CurrentCondition) == 0 {
		return WaybarOutput{}, fmt.Errorf("no weather data")
	}

	current := data.CurrentCondition[0]
	condition := current.WeatherDesc[0].Value
	icon := getIcon(condition, current.WeatherCode)

	text := fmt.Sprintf("%s %s°C", icon, current.TempC)

	tooltip := fmt.Sprintf(
		"Batna, Algeria\n%s\nTemperature: %s°C\nFeels like: %s°C\nHumidity: %s%%\nWind: %s km/h %s\nPrecipitation: %s mm\nUpdated: %s",
		condition,
		current.TempC,
		current.FeelsLikeC,
		current.Humidity,
		current.WindspeedKmph,
		current.WindDir16,
		current.PrecipMM,
		time.Now().Format("15:04:05"),
	)

	class := "normal"
	lc := strings.ToLower(condition)
	switch {
	case strings.Contains(lc, "rain"):
		class = "rain"
	case strings.Contains(lc, "snow"):
		class = "snow"
	case strings.Contains(lc, "clear"):
		class = "clear"
	case strings.Contains(lc, "cloud"):
		class = "cloud"
	case strings.Contains(lc, "thunder"):
		class = "thunder"
	case strings.Contains(lc, "fog"), strings.Contains(lc, "mist"):
		class = "fog"
	}

	return WaybarOutput{
		Text:    text,
		Tooltip: tooltip,
		Alt:     condition,
		Class:   class,
	}, nil
}

func main() {
	out, err := getWeather()
	if err != nil {
		fail := WaybarOutput{
			Text:    "󰅧 N/A",
			Tooltip: err.Error(),
			Alt:     "error",
			Class:   "error",
		}
		b, _ := json.Marshal(fail)
		fmt.Println(string(b))
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println(string(b))
}

