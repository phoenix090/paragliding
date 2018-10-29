# paragliding
### Assignment 2 in IMT2681 Cloud Technologies.
Similar to assignment 1, igcinfo but with more api- endpoints and webhooks implemented.
This api uses mlab (mongodb) to restore all the tracks.

All the below return 404 or 501 (invalid req method) when an error occours or the track is not found.

### Heroku
https://apiparagliding.herokuapp.com/paragliding/api/
### Openstack
http://10.212.136.91:8080/paragliding/api/

### GET: /paragliding/api:
Gives basic info about the api.
returns: 
{ "uptime": <uptime>,  "info": "Service for Paragliding tracks.",  "version": "v1" }

### POST: /api/track:
returns: {  "id": "<id>" } when given correct igc file url.

### GET: /api/track:
Same as assignment1, returns:
[<id1>, <id2>, ...]
  
### GET: /api/track/<id> 
 returns:
{
"H_date": <date from File Header, H-record>,
"pilot": <pilot>,
"glider": <glider>,
"glider_id": <glider_id>,
"track_length": <calculated total track length>,
"track_src_url": <the original URL used to upload the track, ie. the URL used with POST>
}

### GET: /api/track/<id>/<field> 
  
 returns:
<pilot> for pilot
<glider> for glider
<glider_id> for glider_id
<calculated total track length> for track_length
<H_date> for H_date
<track_src_url> for track_src_url
  
### GET: /api/ticker/latest 
returns latest timestamp.

### GET /api/ticker/
returns: 
{
"t_latest": <latest added timestamp>,
"t_start": <the first timestamp of the added track>, this will be the oldest track recorded
"t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
"tracks": [<id1>, <id2>, ...],
"processing": <time in ms of how long it took to process the request>
}
If you want to set the cap, use get- variable name "cap" in the url. example:
/api/ticker/?cap=4

### GET: /api/ticker/<timestamp>
returns the same as the one above with a little modification. Starts with the track's timestamp HIGHER then the timestamp given. I use timestamps returned by time.Now(), so remember to send a time- format accourdingly to that when sending with timestamp. E.g: /api/ticker/2018-10-23T08:55:04.332Z
  
### POST: /api/webhook/new_track/
Request should be in the format:
{
    "webhookURL": "http://remoteUrl:8080/randomWebhookPath",
    "minTriggerValue": 2
}
returns the id of the webhook registered.
when posted again with the same "webhookURL", a webhook is triggered and sent to that url.

### GET: /api/webhook/new_track/<webhook_id>
if id exists, returns the webhook in format: 
{
    "webhookURL": "http://remoteUrl:8080/randomWebhookPath",
    "minTriggerValue": 2
}

### GET: /admin/api/tracks_count
returns the current count of all tracks in the DB
i maid a little auth so everyone should not have access to admin- endpoints.
Use admincode var, and set it to 12345 in the url. E.g /admin/api/tracks_count?admincode=12345
##### otherwise it will deny and return 404. (Use this on all admin endpoints)

### DELETE: /admin/api/tracks
Deletes all the tracks from the db and returns "OK!" if everything went well.

### OTHERS
The clock trigger is implemented onyl in Heroku, didn't see the point of doing it in two places and the assignment didn't specify where.


 
