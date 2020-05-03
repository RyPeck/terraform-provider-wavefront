package wavefront

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type DerivedMetric struct {
	ID                       *string  `json:"id,omitempty"`
	Name                     string   `json:"name,omitempty"`
	Query                    string   `json:"query,omitempty"`
	Minutes                  int      `json:"minutes,omitempty"`
	Tags                     WFTags   `json:"tags,omitempty"`
	Status                   []string `json:"status,omitempty"`
	InTrash                  bool     `json:"inTrash,omitempty"`
	QueryFailing             bool     `json:"queryFailing,omitempty"`
	LastFailedTime           int      `json:"lastFailedTime,omitempty"`
	LastErrorMessage         string   `json:"lastErrorMessage,omitempty"`
	AdditionalInformation    string   `json:"additionalInformation,omitempty"`
	HostsUsed                []string `json:"hostsUsed,omitempty"`
	UpdateUserId             string   `json:"updateUserId,omitempty"`
	CreateUserId             string   `json:"createUserId,omitempty"`
	LastProcessedMillis      int      `json:"lastProcessedMillis,omitempty"`
	ProcessRateMinutes       int      `json:"processRateMinutes,omitempty"`
	PointsScannedAtLastQuery int      `json:"pointsScannedAtLastQuery,omitempty"`
	IncludeObsoleteMetrics   bool     `json:"includeObsoleteMetrics,omitempty"`
	LastQueryTime            int      `json:"lastQueryTime,omitempty"`
	MetricsUsed              []string `json:"metricsUsed,omitempty"`
	QueryQBEnabled           bool     `json:"queryQBEnabled,omitempty"`
	UpdatedEpochMillis       int      `json:"updatedEpochMillis,omitempty"`
	CreatedEpochMillis       int      `json:"createdEpochMillis,omitempty"`
	Deleted                  bool     `json:"deleted,omitempty"`
}

type DerivedMetrics struct {
	client Wavefronter
}

type WFTags struct {
	CustomerTags []string `json:"customerTags"`
}

const baseDerivedMetricsPath = "/api/v2/derivedmetric"

func (c *Client) DerivedMetrics() *DerivedMetrics {
	return &DerivedMetrics{client: c}
}

// Get is used to retrieve an existing DerivedMetric by ID.
// The ID field must be specified
func (dm DerivedMetrics) Get(metric *DerivedMetric) error {
	if *metric.ID == "" {
		return fmt.Errorf("id must be specified")
	}
	return basicCrud(dm.client, "GET",
		fmt.Sprintf("%s/%s", baseDerivedMetricsPath, *metric.ID), metric, nil)
}

// Find returns all DerivedMetrics filtered by the given search conditions.
// If filter is nil, all DerivedMetrics are returned.
func (dm DerivedMetrics) Find(filter []*SearchCondition) ([]*DerivedMetric, error) {
	search := &Search{
		client: dm.client,
		Type:   "derivedmetric",
		Params: &SearchParams{
			Conditions: filter,
		},
	}

	var results []*DerivedMetric
	moreItems := true
	for moreItems == true {
		resp, err := search.Execute()
		if err != nil {
			return nil, err
		}
		var tmpres []*DerivedMetric
		err = json.Unmarshal(resp.Response.Items, &tmpres)
		if err != nil {
			return nil, err
		}
		results = append(results, tmpres...)
		moreItems = resp.Response.MoreItems
		search.Params.Offset = resp.NextOffset
	}

	return results, nil
}

// Create a DerivedMetric, name, query, and minutes are required
func (dm DerivedMetrics) Create(metric *DerivedMetric) error {
	if metric.Name == "" || metric.Query == "" || metric.Minutes == 0 {
		return fmt.Errorf("name, query, and minutes must be specified to create a derived metric")
	}

	return basicCrud(dm.client, "POST", baseDerivedMetricsPath, metric, nil)
}

// Update a DerivedMetric all fields are optional except for ID
func (dm DerivedMetrics) Update(metric *DerivedMetric) error {
	if *metric.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return basicCrud(dm.client, "PUT",
		fmt.Sprintf("%s/%s", baseDerivedMetricsPath, *metric.ID), metric, nil)
}

// Delete a DerivedMetric all fields are optional except for ID
func (dm DerivedMetrics) Delete(metric *DerivedMetric, skipTrash bool) error {
	if *metric.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	params := map[string]string{
		"skipTrash": strconv.FormatBool(skipTrash),
	}

	err := basicCrud(dm.client, "DELETE",
		fmt.Sprintf("%s/%s", baseDerivedMetricsPath, *metric.ID), metric, &params)
	if err != nil {
		return err
	}
	*metric.ID = ""
	return nil
}
