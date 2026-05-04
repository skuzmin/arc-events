package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"arc-events/cache"
	"arc-events/models"
)

const URL = "https://arcraiders.com/map-conditions"
const ICONS_URL = "http://localhost/icons"

var (
	reMapNames   = regexp.MustCompile(`"mapNames":(\[[^\]]*\])`)
	reConditions = regexp.MustCompile(`"conditionItems":(\[[^\]]*\])`)
	reServerNow  = regexp.MustCompile(`"serverNow":(\d+)`)
	reEvents     = regexp.MustCompile(`"entries":(\[[^\]]*\])`)
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func getArcData() []byte {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("RSC", "1")
	req.Header.Set("Accept", "*/*")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error GET arc data: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error READ body: %v\n", err)
		return nil
	}
	return body
}

func extractArcEventsData(body []byte) ([]models.ArcEvent, models.ArcMetadata) {
	var (
		metadata   models.ArcMetadata
		mapNames   []string
		conditions []models.ConditionDetails
		events     []models.ArcEvent
		serverTime int64
		wg         sync.WaitGroup
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		extractJSON(body, reMapNames, "mapNames", &mapNames)
	}()

	go func() {
		defer wg.Done()
		extractJSON(body, reConditions, "conditionItems", &conditions)
	}()

	go func() {
		defer wg.Done()
		extractJSON(body, reEvents, "entries", &events)
	}()

	go func() {
		defer wg.Done()
		t := reServerNow.FindSubmatch(body)
		if t == nil {
			log.Println("serverNow not found")
			metadata.ServerTime = 0
			return
		}
		n, err := strconv.ParseInt(string(t[1]), 10, 64)
		if err != nil {
			log.Printf("parse serverNow: %v", err)
			metadata.ServerTime = 0
			return
		}
		serverTime = n
	}()

	wg.Wait()

	metadata.Maps = make(map[int]string, len(mapNames))
	for i, m := range mapNames {
		metadata.Maps[i] = m
	}

	metadata.Conditions = make(map[int]models.ConditionDetails, len(conditions))
	for i, c := range conditions {
		c.Icon = strings.Join(strings.Split(strings.ToLower(c.Name), " "), "-") + ".svg"
		metadata.Conditions[i] = c
	}

	metadata.ServerTime = serverTime
	metadata.IconsUrl = ICONS_URL

	return events, metadata
}

func createArcDataResponse(events []models.ArcEvent, metadata models.ArcMetadata) models.ArcDataResponse {
	var result models.ArcDataResponse
	result.Meta = metadata
	result.Data = make([]models.ArcEventResponse, 0, len(events))

	reverseMaps := make(map[string]int, len(metadata.Maps))
	for id, m := range metadata.Maps {
		reverseMaps[m] = id
	}

	reverseConditions := make(map[string]int, len(metadata.Conditions))
	for id, c := range metadata.Conditions {
		reverseConditions[c.Name] = id
	}

	for _, event := range events {
		result.Data = append(result.Data, models.ArcEventResponse{
			MapId:       reverseMaps[event.MapName],
			ConditionId: reverseConditions[event.ConditionName],
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
		})
	}

	return result
}

func extractJSON(body []byte, re *regexp.Regexp, name string, dst any) {
	m := re.FindSubmatch(body)
	if m == nil {
		log.Printf("%s: not found", name)
		return
	}
	if err := json.Unmarshal(m[1], dst); err != nil {
		log.Printf("parse %s: %v", name, err)
	}
}

func getArcEvents() {
	arcData := getArcData()
	events, metadata := extractArcEventsData(arcData)
	result := createArcDataResponse(events, metadata)
	cache.GetInstance().Set(result)
}

func startCacheUpdate() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			getArcEvents()
		}
	}()
}

func ArcEventsHandler() {
	getArcEvents()
	startCacheUpdate()
}
