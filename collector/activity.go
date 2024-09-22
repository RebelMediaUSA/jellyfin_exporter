// Copyright 2010 Rebel Media
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !noactivity
// +build !noactivity

package collector

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rebelmediausa/jellyfin_exporter/collector/utils"
	"github.com/rebelmediausa/jellyfin_exporter/config"
)

var (
	jellyfinReportDays = kingpin.Flag("collector.activity.days", "Jellyfin Playback Reporting search in days (Default to 100 Years).").Default("36525").String()
	activityCount      float64
)

type activityCollector struct {
	activityReport *prometheus.Desc
	logger         log.Logger
}

func init() {
	registerCollector("activity", defaultDisabled, NewActivityCollector)
}

func NewActivityCollector(logger log.Logger) (Collector, error) {
	const subsystem = "activity"
	activityReport := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "count",
		),
		"Playback Reporting activity.",
		[]string{"user_id", "username", "last_seen", "total_time"}, nil,
	)
	return &activityCollector{
		activityReport: activityReport,
		logger:         logger,
	}, nil
}

func (c *activityCollector) Update(ch chan<- prometheus.Metric) error {
	jellyfinURL, jellyfinToken, nil := config.JellyfinInfo(c.logger)

	jellyfinAPIURL := fmt.Sprintf("%s/user_usage_stats/user_activity?days=%s", jellyfinURL, *jellyfinReportDays)
	rawData := utils.GetHTTP(jellyfinAPIURL, jellyfinToken)
	data := rawData.([]interface{})

	for _, item := range data {
		activityMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		activityCount = activityMap["total_count"].(float64)
		level.Debug(c.logger).Log("msg", "Jellyfin Playback Reporting for", "User", activityMap["user_name"].(string))
		ch <- prometheus.MustNewConstMetric(c.activityReport,
			prometheus.CounterValue,
			activityCount,
			activityMap["user_id"].(string),
			activityMap["user_name"].(string),
			strings.TrimSpace(activityMap["last_seen"].(string)),
			strings.TrimSpace(activityMap["total_play_time"].(string)),
		)
	}

	return nil
}
