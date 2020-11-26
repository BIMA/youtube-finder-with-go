# youtube-finder-with-go
Make Go script to custom youtube videos by using Youtube API

## How to setup client_secret key and API key
1. Enable the Youtube API Data V3 in your project
2. Create credentials > choose API Key
3. Follow the steps and voila! you get the key

## How to Download
```
git clone github.com/kaabima/youtube-finder-with-go.git
```

## How to use it
```
go run finder.go -query "your keyword" -max-results 10
```

after you run the code above, you will get the result like the following below:

```
Videos:
The Top 10 Most Promising Medical Technologies UPDATED - The Medical Futurist
URL: https://www.youtube.com/watch?v=dy5Hh_MTVRM
```
