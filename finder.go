package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "Google", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max Youtube results")
)

const developerKey = "DEVELOPER KEY"
const youtubeHome = "https://www.youtube.com/watch?v="

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printVideoInformation(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("%v\nURL: %v\n\n", title, youtubeHome+id)
	}
	fmt.Printf("\n\n")
}

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

	// Make the API call to Youtube.
	call := service.Search.List(searchKey).
		Q(*query).
		MaxResults(*maxResults)

	response, err := call.Do()
	handleError(err, "")

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		}
	}
	printVideoInformation("Videos", videos)
}
