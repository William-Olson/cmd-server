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

// GetDefaultOrErr fetches the version based on environment variables
func GetDefaultOrErr(deps *cmddeps.Deps) (cmdutils.SlackResponse, error) {

	// get client info
	rsp := cmdutils.SlackResponse{}
	config := deps.Get("config").(*cmdutils.Config)
	cl := getEnvClientInfo(config)

	// construct the fetch url and fetch the version
	url := fmt.Sprintf("http://%s.%s/%s", cl.server, cl.host, cl.path)
	v := fetchVersion(url)

	if v.Err != nil {
		return rsp, v.Err
	}

	// TODO:
	//  format v.Timestamp for displaying proper build date
	//  capitalize cl.server string for SlackResponse.Text

	rsp.ResponseType = inChannelResponse
	rsp.Text = fmt.Sprintf("_*%s*_ is running version *%s*", cl.server, v.Version)
	rsp.Attachments = []cmdutils.SlackAttachment{
		cmdutils.SlackAttachment{
			Title:     fmt.Sprintf("%s.%s", cl.server, cl.host),
			TitleLink: fmt.Sprintf("http://%s.%s", cl.server, cl.host),
			Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
		},
	}

	return rsp, nil

}

/*

	Fetches the version and timestamp from a given url

*/
func fetchVersion(url string) *Version {

	fmt.Printf("fetching version from: %v\n", url)

	v := Version{}
	resp, err := http.Get(url)

	if v.Err = err; v.Err != nil {
		return &v
	}

	bt, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if v.Err = err; v.Err != nil {
		return &v
	}

	v.Err = json.Unmarshal(bt, &v)
	return &v

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

	payload := cmdutils.SlackResponse{}
	payload.ResponseType = inChannelResponse
	payload.Attachments = []cmdutils.SlackAttachment{}

	// construct the fetch url and fetch the version
	url := fmt.Sprintf("http://%s.%s/%s", cl.server, cl.host, cl.path)
	v := fetchVersion(url)

	if v.Err != nil {
		return payload, v.Err
	}

	// TODO:
	//  format v.Timestamp for displaying proper build date
	//  capitalize cl.server string for SlackResponse.Text

	payload.Text = fmt.Sprintf("_*%s*_ is running version *%s*", cl.server, v.Version)
	payload.Attachments = append(payload.Attachments, cmdutils.SlackAttachment{
		Title:     fmt.Sprintf("%s.%s", cl.server, cl.host),
		TitleLink: fmt.Sprintf("http://%s.%s", cl.server, cl.host),
		Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
	})

	// update slug
	err := updateSlugsOrErr(db, slackClient, slug)

	return payload, err

}

// GetMultiSlugVersionsOrErr
func GetMultiSlugVersionsOrErr(db *cmddb.DB, sc cmddb.SlackClient, slugs []string) (cmdutils.SlackResponse, error) {

	ch := make(chan Version)

	// spin off a go routine for each slug fetch
	for _, slg := range slugs {
		go chFetch(fmt.Sprintf("http://%s.%s/%s", slg, sc.Host, sc.VersionPath), ch)
	}

	// create a response struct
	payload := cmdutils.SlackResponse{}
	payload.ResponseType = inChannelResponse
	payload.Attachments = []cmdutils.SlackAttachment{}
	payload.Text = fmt.Sprintf("Listing Versions")

	// build up results under Attachments field
	for _, slg := range slugs {
		v := <-ch

		if v.Err != nil {
			return payload, v.Err
		}

		// TODO:
		//  format v.Timestamp for displaying proper build date
		//  capitalize cl.server string for SlackResponse.Text

		payload.Attachments = append(payload.Attachments, cmdutils.SlackAttachment{
			Title:     fmt.Sprintf("%s is running version %s", slg, v.Version),
			TitleLink: fmt.Sprintf("http://%s.%s", slg, sc.Host),
			Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
		})

		// update slug for client
		if err := updateSlugsOrErr(db, sc, slg); err != nil {
			return payload, err
		}
	}

	return payload, nil

}

/*

	Attempt updating a slug in db or return error

*/
func updateSlugsOrErr(db *cmddb.DB, sc cmddb.SlackClient, slug string) error {

	// check if slug is already in db
	slugExists := false
	for _, slg := range sc.GetSlugs() {
		if slg.Name == slug {
			slugExists = true
		}
	}

	if slugExists == false {
		// new slug
		newEntry := cmddb.SlackSlug{
			SlackClientID: sc.ID,
			Name:          slug,
		}

		// insert the new slug in db
		fmt.Printf("Creating new slack_slugs entry: %v\n", newEntry)
		err := db.CreateSlackSlugOrErr(newEntry)
		if err != nil {
			return err
		}
	}

	// TODO:
	//  update slug "updated_at" time if exists

	return nil

}
