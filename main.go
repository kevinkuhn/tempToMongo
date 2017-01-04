package main

import(
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "time"
    "github.com/wcharczuk/go-chart"
)

const (
    TEMP_JSON = "https://api.mlab.com/api/1/databases/arduinotest/collections/temperatureData/?apiKey=apiKey"
)

type TempData []struct {
    ID struct {
        Oid string `json:"$oid"`
    } `json:"_id"`
    Date struct {
        Date time.Time `json:"$date"`
    } `json:"date"`
    Temperature float64 `json:"temperature"`
    V int `json:"__v"`
}

type API struct {

}

func main() {
    http.HandleFunc("/", drawChart)
    http.ListenAndServe(":8080", nil)
}

func (r API) GetData() (*TempData,error) {

    body, err := MakeRequest(TEMP_JSON)
    if (err != nil) {
        return nil, err
    }
    s, err := ParseJSON(body)

    return  s, err
}

func MakeRequest(url string) ([]byte, error) {
    
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    
    return []byte(body), err
}

func ParseJSON(body []byte) (*TempData, error) {
    var s = new(TempData)

    err := json.Unmarshal(body, &s)
    if(err != nil){
        fmt.Println("error:", err)
    }

    return s, err
}

func drawChart(res http.ResponseWriter, req *http.Request) {

    API := new(API)
    data, err := API.GetData()
    if (err != nil) {
        fmt.Println(err)
    }

    var temps []float64
    var dates []time.Time

    for _, d := range *data {
        temps = append(temps, d.Temperature)
        dates = append(dates, d.Date.Date)
    }
    fmt.Println(temps)
    fmt.Println(dates)
    //
    graph := chart.Chart{
        XAxis: chart.XAxis{
            Name:      "Time",
            NameStyle: chart.StyleShow(),
            Style: chart.Style{
                Show: true,
                StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
                FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
            },
            ValueFormatter: chart.TimeHourValueFormatter,
        },
        YAxis: chart.YAxis{
            Name:      "Temperature",
            NameStyle: chart.StyleShow(),
            Style: chart.Style{
                Show: true, //enables / displays the y-axis
                StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
                FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
            },
        },
        Series: []chart.Series{
            chart.TimeSeries{
                XValues: dates,
                YValues: temps,
            },
        },
    }

    res.Header().Set("Content-Type", "image/png")
    graph.Render(chart.PNG, res)
}