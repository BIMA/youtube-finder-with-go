package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// Variables declaration
var (
	query      = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max Youtube results")
)

// Constant declaration
const developerKey = "DEVELOPER_KEY"
const youtubeHome = "https://www.youtube.com/watch?v="

// Function declaration
func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func videoCount(id string, key []string, s *youtube.Service) uint64 {
	var count uint64
	call := s.Videos.List(key).
		Id(id).
		MaxResults(*maxResults)
	response, err := call.Do()
	handleError(err, "")

	for _, item := range response.Items {
		switch item.Kind {
		case "youtube#video":
			count = item.Statistics.ViewCount
		}
	}
	return count
}

func subscriberCount(id string, key []string, s *youtube.Service) uint64 {
	var count uint64
	call := s.Channels.List(key).
		Id(id).
		MaxResults(*maxResults)

	response, err := call.Do()
	handleError(err, "")

	for _, item := range response.Items {
		switch item.Kind {
		case "youtube#channel":
			count = item.Statistics.SubscriberCount
		}
	}
	return count
}

func printVideoInformation(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for key, value := range matches {
		switch true {
		case strings.Contains(key, "view_count"):
			fmt.Printf("View count: %v\n", value)
		case strings.Contains(key, "subs_count"):
			fmt.Printf("Subscriber count: %v\n", value)
		default:
			fmt.Printf("\n\n%v\nURL: %v\n", value, youtubeHome+key)
		}
	}
	fmt.Printf("\n\n")
}

// Another calculation function
// still progress

func main() {
	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new Youtube client: %v", err)
	}

	searchKey := []string{"id", "snippet"}
	statsKey := []string{"statistics"}

	// Make the API call to Youtube.
	searchCall := service.Search.List(searchKey).
		Q(*query).
		MaxResults(*maxResults)

	searchResponse, err := searchCall.Do()
	handleError(err, "")

	// Group video in separate lists
	videos := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range searchResponse.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
			videos["view_count_"+item.Id.VideoId] = fmt.Sprint(
				videoCount(item.Id.VideoId, statsKey, service),
			)
			videos["subs_count_"+item.Id.VideoId] = fmt.Sprint(
				subscriberCount(item.Id.VideoId, statsKey, service),
			)
		}
	}
	printVideoInformation("Videos", videos)

}
