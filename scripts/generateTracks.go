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
	"github.com/mattribution/api/pkg/api"
	wr "github.com/mroth/weightedrand"
)

// newStrPtr used to create a new pointer to a string
func newStrPtr(val string) *string {
	return &val
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	// Amount of users to mock
	userCount := 100
	// Mock owner id
	var ownerID int64 = 1
	trackLoopMax := 3
	baseURL := "https://mattribution.com"
	bar := pb.StartNew(userCount)

	convertTrack := api.Track{
		OwnerID:  ownerID,
		PageURL:  newStrPtr(baseURL + "/signup"),
		PagePath: newStrPtr("/signup"),
		Event:    newStrPtr("signup"),
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
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL),
			PagePath:        newStrPtr("/"),
			PageTitle:       newStrPtr("Home"),
			PageReferrer:    newStrPtr(randGoogleAdsReferrer.Pick().(string)),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr("AdWords"),
			CampaignMedium:  newStrPtr("banner"),
			CampaignName:    newStrPtr("Paid Ads"),
			CampaignContent: newStrPtr("Image of marketing connected by us"),
		}, Weight: 1},
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL + "/solutions"),
			PagePath:        newStrPtr("/solutions"),
			PageTitle:       newStrPtr("Solutions"),
			PageReferrer:    newStrPtr("https://google.com"),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr("AdWords"),
			CampaignMedium:  newStrPtr("paid search"),
			CampaignName:    newStrPtr("Paid Search"),
			CampaignContent: newStrPtr("Link to our home page"),
		}, Weight: 2},
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL + "/get-started"),
			PagePath:        newStrPtr("/get-started"),
			PageTitle:       newStrPtr("Get Started"),
			PageReferrer:    newStrPtr(randBlogPostReferrer.Pick().(string)),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr("Blog"),
			CampaignMedium:  newStrPtr("blog"),
			CampaignName:    newStrPtr("Content Marketing"),
			CampaignContent: newStrPtr("Blog posts about our technology"),
		}, Weight: 5},
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL),
			PagePath:        newStrPtr("/"),
			PageTitle:       newStrPtr("Home"),
			PageReferrer:    newStrPtr("https://twitter.com"),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr("Twitter"),
			CampaignMedium:  newStrPtr("twitter"),
			CampaignName:    newStrPtr("Social Media Presence"),
			CampaignContent: newStrPtr("Tweets about our software"),
		}, Weight: 2},
	)

	// Random generator for generic track
	randGeneralTrack := wr.NewChooser(
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL + "/about"),
			PagePath:        newStrPtr("/about"),
			PageTitle:       newStrPtr("About"),
			PageReferrer:    newStrPtr(""),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr(""),
			CampaignMedium:  newStrPtr(""),
			CampaignName:    newStrPtr(""),
			CampaignContent: newStrPtr(""),
		}, Weight: 2},
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL + "/how-it-works"),
			PagePath:        newStrPtr("/how-it-works"),
			PageTitle:       newStrPtr("How It Works"),
			PageReferrer:    newStrPtr(""),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr(""),
			CampaignMedium:  newStrPtr(""),
			CampaignName:    newStrPtr(""),
			CampaignContent: newStrPtr(""),
		}, Weight: 2},
		wr.Choice{Item: api.Track{
			OwnerID:         ownerID,
			PageURL:         newStrPtr(baseURL + "/pricing"),
			PagePath:        newStrPtr("/pricing"),
			PageTitle:       newStrPtr("Pricing"),
			PageReferrer:    newStrPtr(""),
			Event:           newStrPtr("key_page_view"),
			CampaignSource:  newStrPtr(""),
			CampaignMedium:  newStrPtr(""),
			CampaignName:    newStrPtr(""),
			CampaignContent: newStrPtr(""),
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

			fTrack := randFunnelTrack.Pick().(api.Track)
			fTrack.AnonymousID = &anonymousID
			fTrack.SentAt = &date
			storeTrack(fTrack)

			// Loop to generate "natural" web navigation
			for i := 0; i < rand.Intn(20); i++ {
				// Add to date (minutes) after every track to simulate user reading
				date = date.Add(time.Minute * time.Duration(rand.Intn(5)+1))
				if date.Unix() > time.Now().Unix() {
					break
				}

				genTrack := randGeneralTrack.Pick().(api.Track)
				genTrack.AnonymousID = &anonymousID
				genTrack.SentAt = &date
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
		convertTrack.AnonymousID = &anonymousID
		convertTrack.UserID = &userID
		convertTrack.SentAt = &date
		storeTrack(convertTrack)

		bar.Increment()
	}

	bar.Finish()
}

func storeTrack(t api.Track) {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	b64Str := b64.StdEncoding.EncodeToString(jsonBytes)
	// Create a Resty Client
	uri := fmt.Sprintf("http://localhost:3001/v1/pixel/track?data=%v", b64Str)
	client := resty.New()
	_, err = client.R().
		Get(uri)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}
