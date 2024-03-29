// Code generated by "requestgen -method GET -responseType .APIResponse -responseDataField Result -url /v5/market/kline -type GetKLinesRequest -responseDataType .KLinesResponse"; DO NOT EDIT.

package bybitapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func (g *GetKLinesRequest) Category(category Category) *GetKLinesRequest {
	g.category = category
	return g
}

func (g *GetKLinesRequest) Symbol(symbol string) *GetKLinesRequest {
	g.symbol = symbol
	return g
}

func (g *GetKLinesRequest) Interval(interval string) *GetKLinesRequest {
	g.interval = interval
	return g
}

func (g *GetKLinesRequest) StartTime(startTime time.Time) *GetKLinesRequest {
	g.startTime = &startTime
	return g
}

func (g *GetKLinesRequest) EndTime(endTime time.Time) *GetKLinesRequest {
	g.endTime = &endTime
	return g
}

func (g *GetKLinesRequest) Limit(limit uint64) *GetKLinesRequest {
	g.limit = &limit
	return g
}

// GetQueryParameters builds and checks the query parameters and returns url.Values
func (g *GetKLinesRequest) GetQueryParameters() (url.Values, error) {
	var params = map[string]interface{}{}
	// check category field -> json key category
	category := g.category

	// TEMPLATE check-valid-values
	switch category {
	case "spot":
		params["category"] = category

	default:
		return nil, fmt.Errorf("category value %v is invalid", category)

	}
	// END TEMPLATE check-valid-values

	// assign parameter of category
	params["category"] = category
	// check symbol field -> json key symbol
	symbol := g.symbol

	// assign parameter of symbol
	params["symbol"] = symbol
	// check interval field -> json key interval
	interval := g.interval

	// TEMPLATE check-valid-values
	switch interval {
	case "1", "3", "5", "15", "30", "60", "120", "240", "360", "720", "D", "W", "M":
		params["interval"] = interval

	default:
		return nil, fmt.Errorf("interval value %v is invalid", interval)

	}
	// END TEMPLATE check-valid-values

	// assign parameter of interval
	params["interval"] = interval
	// check startTime field -> json key start
	if g.startTime != nil {
		startTime := *g.startTime

		// assign parameter of startTime
		// convert time.Time to milliseconds time stamp
		params["start"] = strconv.FormatInt(startTime.UnixNano()/int64(time.Millisecond), 10)
	} else {
	}
	// check endTime field -> json key end
	if g.endTime != nil {
		endTime := *g.endTime

		// assign parameter of endTime
		// convert time.Time to milliseconds time stamp
		params["end"] = strconv.FormatInt(endTime.UnixNano()/int64(time.Millisecond), 10)
	} else {
	}
	// check limit field -> json key limit
	if g.limit != nil {
		limit := *g.limit

		// assign parameter of limit
		params["limit"] = limit
	} else {
	}

	query := url.Values{}
	for _k, _v := range params {
		query.Add(_k, fmt.Sprintf("%v", _v))
	}

	return query, nil
}

// GetParameters builds and checks the parameters and return the result in a map object
func (g *GetKLinesRequest) GetParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}

	return params, nil
}

// GetParametersQuery converts the parameters from GetParameters into the url.Values format
func (g *GetKLinesRequest) GetParametersQuery() (url.Values, error) {
	query := url.Values{}

	params, err := g.GetParameters()
	if err != nil {
		return query, err
	}

	for _k, _v := range params {
		if g.isVarSlice(_v) {
			g.iterateSlice(_v, func(it interface{}) {
				query.Add(_k+"[]", fmt.Sprintf("%v", it))
			})
		} else {
			query.Add(_k, fmt.Sprintf("%v", _v))
		}
	}

	return query, nil
}

// GetParametersJSON converts the parameters from GetParameters into the JSON format
func (g *GetKLinesRequest) GetParametersJSON() ([]byte, error) {
	params, err := g.GetParameters()
	if err != nil {
		return nil, err
	}

	return json.Marshal(params)
}

// GetSlugParameters builds and checks the slug parameters and return the result in a map object
func (g *GetKLinesRequest) GetSlugParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}

	return params, nil
}

func (g *GetKLinesRequest) applySlugsToUrl(url string, slugs map[string]string) string {
	for _k, _v := range slugs {
		needleRE := regexp.MustCompile(":" + _k + "\\b")
		url = needleRE.ReplaceAllString(url, _v)
	}

	return url
}

func (g *GetKLinesRequest) iterateSlice(slice interface{}, _f func(it interface{})) {
	sliceValue := reflect.ValueOf(slice)
	for _i := 0; _i < sliceValue.Len(); _i++ {
		it := sliceValue.Index(_i).Interface()
		_f(it)
	}
}

func (g *GetKLinesRequest) isVarSlice(_v interface{}) bool {
	rt := reflect.TypeOf(_v)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	}
	return false
}

func (g *GetKLinesRequest) GetSlugsMap() (map[string]string, error) {
	slugs := map[string]string{}
	params, err := g.GetSlugParameters()
	if err != nil {
		return slugs, nil
	}

	for _k, _v := range params {
		slugs[_k] = fmt.Sprintf("%v", _v)
	}

	return slugs, nil
}

// GetPath returns the request path of the API
func (g *GetKLinesRequest) GetPath() string {
	return "/v5/market/kline"
}

// Do generates the request object and send the request object to the API endpoint
func (g *GetKLinesRequest) Do(ctx context.Context) (*KLinesResponse, error) {

	// no body params
	var params interface{}
	query, err := g.GetQueryParameters()
	if err != nil {
		return nil, err
	}

	var apiURL string

	apiURL = g.GetPath()

	req, err := g.client.NewRequest(ctx, "GET", apiURL, query, params)
	if err != nil {
		return nil, err
	}

	response, err := g.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}

	type responseValidator interface {
		Validate() error
	}
	validator, ok := interface{}(apiResponse).(responseValidator)
	if ok {
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}
	var data KLinesResponse
	if err := json.Unmarshal(apiResponse.Result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
