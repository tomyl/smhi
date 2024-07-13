package smhi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// https://opendata.smhi.se/apidocs/metfcst/parameters.html
var ParameterDescriptions = map[string]ParameterDescription{
	"msl": {
		Name:        "msl",
		LevelType:   "hmsl",
		Level:       0,
		Unit:        "hPa",
		Description: "Air pressure",
		ValueRange:  "Decimal number, one decimal",
	},
	"t": {
		Name:        "t",
		LevelType:   "hl",
		Level:       2,
		Unit:        "C",
		Description: "Air temperature",
		ValueRange:  "Decimal number, one decimal",
	},
	"vis": {
		Name:        "vis",
		LevelType:   "hl",
		Level:       2,
		Unit:        "km",
		Description: "Horizontal visibility",
		ValueRange:  "Decimal number, one decimal",
	},
	"wd": {
		Name:        "wd",
		LevelType:   "hl",
		Level:       10,
		Unit:        "degree",
		Description: "Wind direction",
		ValueRange:  "Integer",
	},
	"ws": {
		Name:        "ws",
		LevelType:   "hl",
		Level:       10,
		Unit:        "m/s",
		Description: "Wind speed",
		ValueRange:  "Decimal number, one decimal",
	},
	"r": {
		Name:        "r",
		LevelType:   "hl",
		Level:       2,
		Unit:        "%",
		Description: "Relative humidity",
		ValueRange:  "Integer, 0-100",
	},
	"tstm": {
		Name:        "tstm",
		LevelType:   "hl",
		Level:       0,
		Unit:        "%",
		Description: "Thunder probability",
		ValueRange:  "Integer, 0-100",
	},
	"tcc_mean": {
		Name:        "tcc_mean",
		LevelType:   "hl",
		Level:       0,
		Unit:        "octas",
		Description: "Mean value of total cloud cover",
		ValueRange:  "Integer, 0-8",
	},
	"lcc_mean": {
		Name:        "lcc_mean",
		LevelType:   "hl",
		Level:       0,
		Unit:        "octas",
		Description: "Mean value of low level cloud cover",
		ValueRange:  "Integer, 0-8",
	},
	"mcc_mean": {
		Name:        "mcc_mean",
		LevelType:   "hl",
		Level:       0,
		Unit:        "octas",
		Description: "Mean value of medium level cloud cover",
		ValueRange:  "Integer, 0-8",
	},
	"hcc_mean": {
		Name:        "hcc_mean",
		LevelType:   "hl",
		Level:       0,
		Unit:        "octas",
		Description: "Mean value of high level cloud cover",
		ValueRange:  "Integer, 0-8",
	},
	"gust": {
		Name:        "gust",
		LevelType:   "hl",
		Level:       10,
		Unit:        "m/s",
		Description: "Wind gust speed",
		ValueRange:  "Decimal number, one decimal",
	},
	"pmin": {
		Name:        "pmin",
		LevelType:   "hl",
		Level:       0,
		Unit:        "mm/h",
		Description: "Minimum precipitation intensity",
		ValueRange:  "Decimal number, one decimal",
	},
	"pmax": {
		Name:        "pmax",
		LevelType:   "hl",
		Level:       0,
		Unit:        "mm/h",
		Description: "Maximum precipitation intensity",
		ValueRange:  "Decimal number, one decimal",
	},
	"spp": {
		Name:        "spp",
		LevelType:   "hl",
		Level:       0,
		Unit:        "%",
		Description: "Percent of precipitation in frozen form",
		ValueRange:  "Integer, -9 or 0-100",
	},
	"pcat": {
		Name:        "pcat",
		LevelType:   "hl",
		Level:       0,
		Unit:        "category",
		Description: "Precipitation category",
		ValueRange:  "Integer, 0-6",
	},
	"pmean": {
		Name:        "pmean",
		LevelType:   "hl",
		Level:       0,
		Unit:        "mm/h",
		Description: "Mean precipitation intensity",
		ValueRange:  "Decimal number, one decimal",
	},
	"pmedian": {
		Name:        "pmedian",
		LevelType:   "hl",
		Level:       0,
		Unit:        "mm/h",
		Description: "Median precipitation intensity",
		ValueRange:  "Decimal number, one decimal",
	},
	"wsymb2": {
		Name:        "wsymb2",
		LevelType:   "hl",
		Level:       0,
		Unit:        "code",
		Description: "Weather symbol",
		ValueRange:  "Integer, 1-27",
	},
}

type ParameterDescription struct {
	Name        string
	LevelType   string
	Level       int
	Unit        string
	Description string
	ValueRange  string
}

var WeatherSymbols = []WeatherSymbol{
	{0, "No weather", "?", 1},
	{1, "Clear sky", "\u2600", 1},                   // ☀
	{2, "Nearly clear sky", "\u26c5", 2},            // ⛅
	{3, "Variable cloudiness", "\u26c5", 2},         // ⛅
	{4, "Halfclear sky", "\u26c5", 2},               // ⛅
	{5, "Cloudy sky", "\u2601", 1},                  // ☁
	{6, "Overcast", "\u2601", 1},                    // ☁
	{7, "Fog", "\U0001f32B", 1},                     // 🌫
	{8, "Light rain showers", "\U0001f326", 1},      // 🌦
	{9, "Moderate rain showers", "\U0001f326", 1},   // 🌦
	{10, "Heavy rain showers", "\U0001f327", 1},     // 🌧
	{11, "Thunderstorm", "\u26a1", 2},               // ⚡
	{12, "Light sleet showers", "\U0001f328", 1},    // 🌨
	{13, "Moderate sleet showers", "\U0001f328", 1}, // 🌨
	{14, "Heavy sleet showers", "\U0001f328", 1},    // 🌨
	{15, "Light snow showers", "\U0001f328", 1},     // 🌨
	{16, "Moderate snow showers", "\U0001f328", 1},  // 🌨
	{17, "Heavy snow showers", "\U0001f328", 1},     // 🌨
	{18, "Light rain", "\U0001f327", 1},             // 🌧
	{19, "Moderate rain", "\U0001f327", 1},          // 🌧
	{20, "Heavy rain", "\U0001f327", 1},             // 🌧
	{21, "Thunder", "\u26a1", 2},                    // ⚡
	{22, "Light sleet", "\U0001f328", 1},            // 🌨
	{23, "Moderate sleet", "\U0001f328", 1},         // 🌨
	{24, "Heavy sleet", "\U0001f328", 1},            // 🌨
	{25, "Light snowfall", "\U0001f328", 1},         // 🌨
	{26, "Moderate snowfall", "\U0001f328", 1},      // 🌨
	{27, "Heavy snowfall", "\U0001f328", 1},         // 🌨
}

type WeatherSymbol struct {
	Value        int
	Meaning      string
	Unicode      string
	UnicodeWidth int
}

func (s WeatherSymbol) FixedWidth() string {
	if s.UnicodeWidth == 1 {
		return s.Unicode + " "
	}
	return s.Unicode + "\u200b"
}

type Forecast struct {
	ApprovedTime  time.Time
	ReferenceTime time.Time
	Geometry      Geometry
	TimeSeries    []TimeSeriesItem
}

type Geometry struct {
	Type        string
	Coordinates []Point
}

type Point [2]float64

type TimeSeriesItem struct {
	ValidTime  time.Time
	Parameters []Parameter
}

func (i TimeSeriesItem) Float64(name string) float64 {
	for _, p := range i.Parameters {
		if p.Name == name {
			return p.Values[0]
		}
	}
	return 0
}

func (i TimeSeriesItem) Int(name string) int {
	for _, p := range i.Parameters {
		if p.Name == name {
			return int(p.Values[0])
		}
	}
	return 0
}

func (i TimeSeriesItem) Temperature() float64 {
	return i.Float64("t")
}

func (i TimeSeriesItem) MaxPrecipitation() float64 {
	return i.Float64("pmax")
}

func (i TimeSeriesItem) WindSpeed() float64 {
	return i.Float64("ws")
}

func (i TimeSeriesItem) WeatherSymbol() WeatherSymbol {
	idx := i.Int("Wsymb2")
	if idx >= 1 && idx < len(WeatherSymbols) {
		return WeatherSymbols[idx]
	}
	return WeatherSymbol{}
}

type Parameter struct {
	Name      string
	LevelType string
	Level     int
	Unit      string
	Values    []float64
}

func GetForecast(lon, lat float64) (*Forecast, error) {
	resp, err := http.Get(fmt.Sprintf("https://opendata-download-metfcst.smhi.se/api/category/pmp3g/version/2/geotype/point/lon/%f/lat/%f/data.json", lon, lat))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status is not ok: %s", buf)
	}

	var forecast Forecast
	if err := json.Unmarshal(buf, &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}
