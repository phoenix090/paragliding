package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// for restoring all webhooks
var allhooks map[int]WebhookInfo

// Keeps the current nr of tracks in memory
var Count int

type Webhook struct {
	Content string `json:"content"`
}

type WebhookInfo struct {
	WebhookURL      string `json:"webhookURL"`
	MinTriggerValue int    `json:"minTriggerValue"`
}

// WebHookResponse When webhook is invoked POST request
type WebHookResponse struct {
	Tlatest    time.Time `json:"t_latest"`
	Tracks     []int     `json:"tracks"`
	Processing float64   `json:"processing"`
}

// POST /api/webhook/new_track/
// RegisterWebhook handles registration of new webhook and
// Calling other functions when handling GET and Delete api
func RegisterWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case "POST":
		var newHook WebhookInfo
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&newHook); err != nil {
			fmt.Println(err)
			http.Error(w, "Error with the request! The api only accepts JSON", 404)
			return
		}

		// Setting mintrigger to default if not omitted
		if newHook.MinTriggerValue <= 0 {
			newHook.MinTriggerValue = 1
		}

		if allhooks == nil {
			allhooks = make(map[int]WebhookInfo)
		}
		id := len(allhooks) + 1

		// Handling Invoking a registered webhook
		for _, webhook := range allhooks {
			if webhook.WebhookURL == newHook.WebhookURL {
				before := time.Now()
				PopulateTickerInfo(newHook.MinTriggerValue, false, "")
				after := time.Now()
				tot := after.Sub(before).Seconds() * 1000
				resObj := WebHookResponse{Tlatest: ticker.TLatest, Tracks: ticker.Tracks, Processing: tot}
				err := WebhookToDiscord(resObj, webhook.WebhookURL)
				if err != nil {
					http.Error(w, "something went wrong", 404)
				}
				return
			}
		}
		json.NewEncoder(w).Encode(id)
		allhooks[id] = newHook
	case "GET":
		getWebhookByID(w, r)
	case "DELETE":
		DeleteWbHook(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

//GetWebhookByID ...
func getWebhookByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	//fmt.Println("len", len(parts), parts, parts[4])
	if len(parts) == 5 && parts[4] != "" {
		ID, conErr := strconv.Atoi(parts[4])
		if conErr != nil {
			http.Error(w, "Invalid id provided.", 404)
			return
		}
		for i, hook := range allhooks {
			if ID == i {
				json.NewEncoder(w).Encode(hook)
				return
			}
		}
		http.Error(w, "No webhook with that track is registered", http.StatusNotFound)
	} else {
		http.Error(w, "Empty or not a number was provided", http.StatusNotFound)
	}
}

//DeleteWbHook deletes a spesific webhook by its id
func DeleteWbHook(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	//fmt.Println("len", len(parts), parts, parts[4])
	if len(parts) == 5 && parts[4] != "" {
		ID, conErr := strconv.Atoi(parts[4])
		if conErr != nil {
			http.Error(w, "Invalid id provided.", 404)
			return
		}
		if ID > len(allhooks)+1 {
			http.Error(w, "No webhook with that track is registered", http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(allhooks[ID])
			delete(allhooks, ID)
		}
	}
}

//WebhookToDiscord Sends webhook message to discord
func WebhookToDiscord(payload WebHookResponse, url string) error {
	proc := strconv.Itoa(int(payload.Processing))
	totTracks := strconv.Itoa(len(payload.Tracks))
	var ids string
	var message string
	if len(payload.Tracks) == 0 {
		message = "There are total of " + totTracks +
			" tracks in the record\n. Time to process: " + proc
	} else {
		for _, eachID := range payload.Tracks {
			ids += strconv.Itoa(eachID) + ", "
		}
		message = "Latest timestamp: " + payload.Tlatest.String() + ".\nThere are " +
			totTracks + " new tracks in the record. \nIDs: " + ids + "\nTime to process: " + proc
	}
	raw, err := json.Marshal(Webhook{Content: message + "\n"})
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(raw))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 206 {
		return fmt.Errorf("Error with responsecode %v", resp.StatusCode)
	}
	return nil
}

// TriggerWebhook sends webook when triggervalue is met
// Updates triggervalue when new track is registered
func TriggerWebhook(count int) error {
	// Make extra field that counts the trigger value down?
	before := time.Now()
	for _, hook := range allhooks {
		if count%hook.MinTriggerValue == 0 {
			PopulateTickerInfo(hook.MinTriggerValue, false, "")
			after := time.Now()
			tot := after.Sub(before).Seconds() * 1000
			resObj := WebHookResponse{Tlatest: ticker.TLatest, Tracks: ticker.Tracks, Processing: tot}
			err := WebhookToDiscord(resObj, hook.WebhookURL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SendLogToDiscord Sends webhook to a predefined webhook.
// Sends webhook only if there are new tracks in the record.
func SendLogToDiscord(url string, count *int) (int, error) {
	before := time.Now()
	all, ok := GlobalDB.GetAll()
	if !ok {
		return 0, fmt.Errorf("error")
	}
	diff := len(all) - *count
	if diff != 0 {
		PopulateTickerInfo(diff, false, "")
		after := time.Now()
		tot := after.Sub(before).Seconds() * 1000
		resObj := WebHookResponse{Tlatest: ticker.TLatest, Tracks: ticker.Tracks, Processing: tot}
		err := WebhookToDiscord(resObj, url)
		if err != nil {
			return 0, err
		}
	}

	*count = len(all)
	return *count, nil
}
