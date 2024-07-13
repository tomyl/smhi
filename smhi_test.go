package smhi_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tomyl/smhi"
)

func TestParseForecast(t *testing.T) {
	buf, err := os.ReadFile("testdata/data.json")
	require.Nil(t, err)

	var forecast smhi.Forecast
	require.Nil(t, json.Unmarshal(buf, &forecast))

	item := forecast.TimeSeries[10]
	require.Equal(t, 18.6, item.Temperature())
	require.Equal(t, 2.6, item.MaxPrecipitation())
	require.Equal(t, 5.6, item.WindSpeed())

	symbol := item.WeatherSymbol()
	require.Equal(t, 19, symbol.Value)
	require.Equal(t, "Moderate rain", symbol.Meaning)
	require.Equal(t, "🌧 ", symbol.FixedWidth())
}
