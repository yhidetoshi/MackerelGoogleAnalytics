package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/mackerelio/mackerel-client-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/option"
)

var (
	mkrKey = os.Getenv("MKRKEY")
	client = mackerel.NewClient(mkrKey)

	json   = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	viewID = os.Getenv("VIEW_ID")
)

const (
	// GA
	startDate = "today"
	endDate   = "today"

	metricsUsers     = "users"
	metricsPVs       = "pagePath!=/preview"
	metricsPageViews = "pageviews"

	dimensionsTitle = "pageTitle"
	dimensionsPath  = "pagePath"

	// Mackerel
	serviceName = "GoogleAnalytics"

	// Lambda
	timezone = "Asia/Tokyo"
	offset   = 9 * 60 * 60
)

func main() {
	lambda.Start(Handler)
}

// func main (){

// Handler Lambda
func Handler() {
	jst := time.FixedZone(timezone, offset)
	nowTime := time.Now().In(jst)

	ctx := context.Background()
	jwtConfig, err := google.JWTConfigFromJSON([]byte(json), analytics.AnalyticsReadonlyScope)
	if err != nil {
		fmt.Println(err)
	}

	ts := jwtConfig.TokenSource(ctx)
	client, err := analytics.NewService(ctx, option.WithTokenSource(ts))

	// Get Users
	resUsers, err := client.Data.Ga.Get("ga:"+viewID, startDate, endDate, "ga:"+metricsUsers).Do()
	if err != nil {
		fmt.Println(err)
	}

	// Get PVs
	resPVs, err := client.Data.Ga.Get(
		"ga:"+viewID, startDate, endDate, "ga:"+metricsPageViews).Dimensions(
		"ga:" + dimensionsTitle + "," + "ga:" + dimensionsPath).Filters("ga:" + metricsPVs).Do()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(nowTime)
	fmt.Println(resUsers.TotalsForAllResults["ga:"+metricsUsers])
	fmt.Println(resPVs.TotalsForAllResults["ga:"+metricsPageViews])

	intResultUsers, _ := strconv.Atoi(resUsers.TotalsForAllResults["ga:"+metricsUsers])
	intResultPVs, _ := strconv.Atoi(resPVs.TotalsForAllResults["ga:"+metricsPageViews])

	mkrErr := PostValuesToMackerel(intResultPVs, intResultUsers, nowTime)
	if err != nil {
		fmt.Println(mkrErr)
	}
}

// PostValuesToMackerel Post Metrics to Mackerel
func PostValuesToMackerel(intResultUsers int, intResultPVs int, nowTime time.Time) error {
	// Post users
	errUser := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Users.users",
			Time:  nowTime.Unix(),
			Value: intResultUsers,
		},
	})
	if errUser != nil {
		fmt.Println(errUser)
	}

	// Post PV
	errPV := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "PV.PVs",
			Time:  nowTime.Unix(),
			Value: intResultPVs,
		},
	})
	if errPV != nil {
		fmt.Println(errPV)
	}

	return nil
}
