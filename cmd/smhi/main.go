package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/tomyl/smhi"
)

func printForecast(forecast *smhi.Forecast) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintf(w, "Time\tWeather\tTemperature\tMax precipitation\tWind speed\n")

	for _, item := range forecast.TimeSeries {
		ts := item.ValidTime.Local().Format("Mon 15:04")
		weather := item.WeatherSymbol()
		fmt.Fprintf(w, "%s\t%s %s\t%.1f°C\t%.1f mm/h\t%.1f m/s\n", ts, weather.FixedWidth(), weather.Meaning, item.Temperature(), item.MaxPrecipitation(), item.WindSpeed())
	}

	w.Flush()
}

func run() error {
	lon := flag.Float64("lon", 0, "Longitude")
	lat := flag.Float64("lat", 0, "Latitude")
	name := flag.String("file", "", "Read data from file")
	flag.Parse()

	if *name != "" {
		buf, err := os.ReadFile(*name)
		if err != nil {
			return err
		}
		var forecast smhi.Forecast
		if err := json.Unmarshal(buf, &forecast); err != nil {
			return err
		}
		printForecast(&forecast)
		return nil
	}

	forecast, err := smhi.GetForecast(*lon, *lat)
	if err != nil {
		return err
	}

	printForecast(forecast)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
