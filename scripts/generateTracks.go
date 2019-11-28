package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattribution/api/pkg/api"
	wr "github.com/mroth/weightedrand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	var host, username, password, dbName = "localhost", "postgres", "password", "mattribution"
	port := 5432
	// Amount of users to mock
	userCount := 100
	// Mock owner id
	var ownerID int64 = 1
	trackLoopMax := 3
	baseURL := "https://mattribution.com"
	bar := pb.StartNew(userCount)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	log.Println(psqlInfo)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	convertTrack := api.Track{
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
		wr.Choice{Item: "https://mattribution.com/blog/posts/1", Weight: 2},
		wr.Choice{Item: "https://mattribution.com/blog/posts/2", Weight: 4},
		wr.Choice{Item: "https://mattribution.com/blog/posts/3", Weight: 8},
		wr.Choice{Item: "https://mattribution.com/blog/posts/4", Weight: 10},
		wr.Choice{Item: "https://mattribution.com/blog/posts/5", Weight: 24},
	)
	// Random generator for funnel
	randFunnelTrack := wr.NewChooser(
		wr.Choice{Item: api.Track{
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
		}, Weight: 10},
		wr.Choice{Item: api.Track{
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
		}, Weight: 5},
		wr.Choice{Item: api.Track{
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
		}, Weight: 20},
		wr.Choice{Item: api.Track{
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
		wr.Choice{Item: api.Track{
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
		}, Weight: 8},
	)

	// Loop over fake user IDs and create paths for them
	for userID := 0; userID < userCount; userID++ {
		// Random day in the last 30 days
		date := time.Now().Add(-time.Hour * time.Duration(24*rand.Intn(30)))
		// Track loop will loop over creating a funnel track, then
		// natural page navigation
		for tLoopIndex := 0; tLoopIndex < trackLoopMax; tLoopIndex++ {
			// Add to date (days) if is next iterration
			if tLoopIndex != 0 {
				date.Add(time.Hour * time.Duration(24*rand.Intn(5)))
				if date.Unix() > time.Now().Unix() {
					break
				}
			}

			fTrack := randFunnelTrack.Pick().(api.Track)
			fTrack.UserID = strconv.Itoa(userID)
			fTrack.SentAt = date
			fTrack.ReceivedAt = date
			storeTrack(db, fTrack)

			// Add to date (minutes) for every track
			date.Add(time.Minute * time.Duration(rand.Intn(5)))
			if date.Unix() > time.Now().Unix() {
				break
			}

			// Loop to generate "natural" web navigation
			for i := 0; i < rand.Intn(20); i++ {
				genTrack := randGeneralTrack.Pick().(api.Track)
				genTrack.UserID = strconv.Itoa(userID)
				storeTrack(db, genTrack)

				// Add to date (minutes) after every track to simulate user reading
				date.Add(time.Minute * time.Duration(rand.Intn(5)))
				if date.Unix() > time.Now().Unix() {
					break
				}
			}

			// Chance to convert (even chance during each action stage,
			// +1 to account for converting outside track loop)
			r := rand.Float32()
			var chanceToConvert float32 = 1 / (float32(trackLoopMax) + 1)
			if r <= chanceToConvert {
				break
			}
		}

		// If we got here we haven't converted, so convert
		storeTrack(db, convertTrack)

		bar.Increment()
	}

	bar.Finish()
}

func storeTrack(db *sqlx.DB, t api.Track) {
	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	// Set default json value (so postgres doesn't get mad)
	if t.Extra == "" {
		t.Extra = "{}"
	}

	_, err := db.Exec(sqlStatement, t.OwnerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt.Format(time.RFC3339), t.ReceivedAt.Format(time.RFC3339), t.Extra)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}
