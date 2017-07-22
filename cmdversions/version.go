package cmdversions

import (
	"encoding/json"
	"fmt"
	"github.com/william-olson/cmd-server/cmddb"
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdutils"
	"io/ioutil"
	"net/http"
)

type client struct {
	server string
	host   string
	path   string
}

type Version struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Err       error  `json:"-"`
}

const (
	inChannelResponse = "in_channel"
)

// GetDefault fetches the version based on environment variables
func GetDefault(deps *cmddeps.Deps) cmdutils.SlackResponse {

	// get client info
	config := deps.Get("config").(*cmdutils.Config)
	cl := getEnvClientInfo(config)

	// construct the fetch url and fetch the version
	url := fmt.Sprintf("http://%s.%s/%s", cl.server, cl.host, cl.path)
	v := fetchVersion(url)

	// TODO:
	//  format v.timestamp for displaying proper build date
	//  capitalize cl.server string for SlackResponse.Text

	return cmdutils.SlackResponse{
		// show to everyone in channel
		ResponseType: inChannelResponse,
		Text:         fmt.Sprintf("_*%s*_ is running version *%s*", cl.server, v.Version),

		// display build info as attachment
		Attachments: []cmdutils.SlackAttachment{
			cmdutils.SlackAttachment{
				Title:     fmt.Sprintf("%s.%s", cl.server, cl.host),
				TitleLink: fmt.Sprintf("http://%s.%s", cl.server, cl.host),
				Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
			},
		},
	}

}

/*

	Fetches the version and timestamp from a given url

*/
func fetchVersion(url string) Version {

	fmt.Printf("fetching version from: %v\n", url)

	// just return dummy version for now
	return Version{
		Version:   "3.0",
		Timestamp: "017-06-16T04:57:40.439Z",
		Err:       nil,
	}

}

/*

	Fetch the version and return it via channel

*/
func chFetch(url string, ch chan<- Version) {

	dat := Version{}
	resp, err := http.Get(url)

	if err != nil {
		dat.Err = err
		ch <- dat
		return
	}

	bt, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		dat.Err = err
		ch <- dat
		return
	}

	if err = json.Unmarshal(bt, &dat); err != nil {
		dat.Err = err
		ch <- dat
		return
	}

	fmt.Println("resp", dat)

	ch <- dat

}

/*

	Get a client based on environment variables

*/
func getEnvClientInfo(config *cmdutils.Config) client {

	return client{
		server: config.Get("VERSION_SERVER"),
		host:   config.Get("VERSION_HOST"),
		path:   config.Get("VERSION_ROUTE"),
	}

}

// GetSlugVersionOrErr fetches a slack response for 1 slug
func GetSlugVersionOrErr(db *cmddb.DB, slackClient cmddb.SlackClient, slug string) (cmdutils.SlackResponse, error) {

	cl := client{
		host:   slackClient.Host,
		path:   slackClient.VersionPath,
		server: slug,
	}

	// construct the fetch url and fetch the version
	url := fmt.Sprintf("http://%s.%s/%s", cl.server, cl.host, cl.path)
	v := fetchVersion(url)

	// TODO:
	//  format v.timestamp for displaying proper build date
	//  capitalize cl.server string for SlackResponse.Text

	payload := cmdutils.SlackResponse{
		ResponseType: inChannelResponse,
		Text:         fmt.Sprintf("_*%s*_ is running version *%s*", cl.server, v.Version),
		Attachments: []cmdutils.SlackAttachment{
			cmdutils.SlackAttachment{
				Title:     fmt.Sprintf("%s.%s", cl.server, cl.host),
				TitleLink: fmt.Sprintf("http://%s.%s", cl.server, cl.host),
				Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
			},
		},
	}

	// check if slug is already in db
	slugExists := false
	for _, slg := range slackClient.GetSlugs() {
		if slg.Name == slug {
			slugExists = true
		}
	}

	if slugExists == false {
		// new slug
		newEntry := cmddb.SlackSlug{
			SlackClientID: slackClient.ID,
			Name:          slug,
		}

		// insert the new slug in db
		fmt.Printf("Creating new slack_slugs entry: %v\n", newEntry)
		err := db.CreateSlackSlugOrErr(newEntry)
		if err != nil {
			return payload, err
		}
	}

	// TODO:
	//  update slug "updated_at" time if exists

	return payload, nil

}
