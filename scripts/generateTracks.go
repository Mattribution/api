package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	b64 "encoding/base64"

	"github.com/cheggaaa/pb/v3"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mattribution/api/internal/app"
	wr "github.com/mroth/weightedrand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	// Amount of users to mock
	userCount := 100
	// Mock owner id
	var ownerID string = "0a1a69ea-7557-4716-8733-7daaeca91a54"
	trackLoopMax := 3
	baseURL := "https://mattribution.com"
	bar := pb.StartNew(userCount)

	convertTrack := app.Track{
		OwnerID:  ownerID,
		PageURL:  baseURL + "/signup",
		PagePath: "/signup",
		Event:    "signup",
	}

	// Random google ads referrer
	randGoogleAdsReferrer := wr.NewChooser(
		wr.Choice{Item: "https://google.com", Weight: 2},
		wr.Choice{Item: "https://youtube.com", Weight: 4},
	)

	// Random google ads referrer
	randBlogPostReferrer := wr.NewChooser(
		wr.Choice{Item: "https://mattribution.com/blog/posts/1", Weight: 1},
		wr.Choice{Item: "https://mattribution.com/blog/posts/2", Weight: 1},
		wr.Choice{Item: "https://mattribution.com/blog/posts/3", Weight: 2},
		wr.Choice{Item: "https://mattribution.com/blog/posts/4", Weight: 4},
		wr.Choice{Item: "https://mattribution.com/blog/posts/5", Weight: 2},
	)

	// Random generator for funnel
	randFunnelTrack := wr.NewChooser(
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL,
			PagePath:        "/",
			PageTitle:       "Home",
			PageReferrer:    randGoogleAdsReferrer.Pick().(string),
			Event:           "key_page_view",
			CampaignSource:  "AdWords",
			CampaignMedium:  "banner",
			CampaignName:    "Paid Ads",
			CampaignContent: "Image of marketing connected by us",
		}, Weight: 1},
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL + "/solutions",
			PagePath:        "/solutions",
			PageTitle:       "Solutions",
			PageReferrer:    "https://google.com",
			Event:           "key_page_view",
			CampaignSource:  "AdWords",
			CampaignMedium:  "paid search",
			CampaignName:    "Paid Search",
			CampaignContent: "Link to our home page",
		}, Weight: 2},
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL + "/get-started",
			PagePath:        "/get-started",
			PageTitle:       "Get Started",
			PageReferrer:    randBlogPostReferrer.Pick().(string),
			Event:           "key_page_view",
			CampaignSource:  "Blog",
			CampaignMedium:  "blog",
			CampaignName:    "Content Marketing",
			CampaignContent: "Blog posts about our technology",
		}, Weight: 5},
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL,
			PagePath:        "/",
			PageTitle:       "Home",
			PageReferrer:    "https://twitter.com",
			Event:           "key_page_view",
			CampaignSource:  "Twitter",
			CampaignMedium:  "twitter",
			CampaignName:    "Social Media Presence",
			CampaignContent: "Tweets about our software",
		}, Weight: 2},
	)

	// Random generator for generic track
	randGeneralTrack := wr.NewChooser(
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL + "/about",
			PagePath:        "/about",
			PageTitle:       "About",
			PageReferrer:    "",
			Event:           "key_page_view",
			CampaignSource:  "",
			CampaignMedium:  "",
			CampaignName:    "",
			CampaignContent: "",
		}, Weight: 2},
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL + "/how-it-works",
			PagePath:        "/how-it-works",
			PageTitle:       "How It Works",
			PageReferrer:    "",
			Event:           "key_page_view",
			CampaignSource:  "",
			CampaignMedium:  "",
			CampaignName:    "",
			CampaignContent: "",
		}, Weight: 2},
		wr.Choice{Item: app.Track{
			OwnerID:         ownerID,
			PageURL:         baseURL + "/pricing",
			PagePath:        "/pricing",
			PageTitle:       "Pricing",
			PageReferrer:    "",
			Event:           "key_page_view",
			CampaignSource:  "",
			CampaignMedium:  "",
			CampaignName:    "",
			CampaignContent: "",
		}, Weight: 2},
	)

	// Loop over fake user IDs and create paths for them
	for userID := 0; userID < userCount; userID++ {
		anonymousID := uuid.New().String()
		// Random day in the last 30 days
		date := time.Now().Add(-time.Hour * time.Duration(24*rand.Intn(30)))
		// Track loop will loop over creating a funnel track, then
		// natural page navigation
		for tLoopIndex := 0; tLoopIndex < trackLoopMax; tLoopIndex++ {
			// Add to date (days) if is next iterration
			if tLoopIndex != 0 {
				date = date.Add(time.Hour * time.Duration(24*rand.Intn(5)))
				if date.Unix() > time.Now().Unix() {
					break
				}
			}

			fTrack := randFunnelTrack.Pick().(app.Track)
			fTrack.AnonymousID = anonymousID
			fTrack.SentAt = date
			storeTrack(fTrack)

			// Loop to generate "natural" web navigation
			for i := 0; i < rand.Intn(20); i++ {
				// Add to date (minutes) after every track to simulate user reading
				date = date.Add(time.Minute * time.Duration(rand.Intn(5)+1))
				if date.Unix() > time.Now().Unix() {
					break
				}

				genTrack := randGeneralTrack.Pick().(app.Track)
				genTrack.AnonymousID = anonymousID
				genTrack.SentAt = date
				storeTrack(genTrack)
			}

			// Chance to convert (even chance during each action stage,
			// +1 to account for converting outside track loop)
			r := rand.Float32()
			chanceToConvert := 1 / (float32(trackLoopMax) + 1)
			if r <= chanceToConvert {
				break
			}
		}

		// Add to date (minutes) after every track to simulate user reading
		date = date.Add(time.Minute * time.Duration(rand.Intn(5)))
		if date.Unix() > time.Now().Unix() {
			continue
		}

		// If we got here we haven't converted, so convert
		userID := strconv.Itoa(userID)
		convertTrack.AnonymousID = anonymousID
		convertTrack.UserID = userID
		convertTrack.SentAt = date
		convertTrack.OwnerID = ownerID
		storeTrack(convertTrack)

		bar.Increment()
	}

	bar.Finish()
}

func storeTrack(t app.Track) {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	b64Str := b64.StdEncoding.EncodeToString(jsonBytes)
	// Create a Resty Client
	uri := fmt.Sprintf("http://localhost:3001/tracks/new?data=%v", b64Str)
	client := resty.New()
	resp, err := client.R().Get(uri)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
	if resp.StatusCode() != 200 {
		log.Println(resp.Status())
	}
}
