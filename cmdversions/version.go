package cmdversions

import (
	"encoding/json"
	"fmt"
	"github.com/william-olson/cmd-server/cmddb"
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdutils"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	inChannelResponse = "in_channel"
)

// Version is a response value returned when fetching versions
type Version struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Err       error  `json:"-"`
}

// GetDefaultOrErr fetches the version based on environment variables
func GetDefaultOrErr(deps *cmddeps.Deps) (cmdutils.SlackResponse, error) {

	// get client info
	rsp := cmdutils.SlackResponse{}
	config := deps.Get("config").(*cmdutils.Config)

	server := config.Get("VERSION_SERVER")
	host := config.Get("VERSION_HOST")
	path := config.Get("VERSION_ROUTE")

	// construct the fetch url and fetch the version
	url := fmt.Sprintf("http://%s.%s/%s", server, host, path)
	v := fetchVersion(url)

	if v.Err != nil {
		return rsp, v.Err
	}

	// TODO:
	//  format v.Timestamp for displaying proper build date
	//  capitalize cl.server string for SlackResponse.Text

	rsp.ResponseType = inChannelResponse
	rsp.Text = fmt.Sprintf("_*%s*_ is running version *%s*", title(server), v.Version)
	rsp.Attachments = []cmdutils.SlackAttachment{
		cmdutils.SlackAttachment{
			Title:     fmt.Sprintf("%s.%s", server, host),
			TitleLink: fmt.Sprintf("http://%s.%s", server, host),
			Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
		},
	}

	return rsp, nil

}

// GetVersionByUrlOrErr fetches a version from a given url
func GetVersionByUrlOrErr(url string) (cmdutils.SlackResponse, error) {

	rsp := cmdutils.SlackResponse{}
	rsp.ResponseType = inChannelResponse
	rsp.Attachments = []cmdutils.SlackAttachment{}
	v := fetchVersion(url)

	if v.Err != nil {
		return rsp, v.Err
	}

	rsp.Text = fmt.Sprintf("_*%s*_ is running version *%s*", url, v.Version)
	rsp.Attachments = append(rsp.Attachments, cmdutils.SlackAttachment{
		Title:     url,
		TitleLink: url,
		Text:      "Build Date: " + v.Timestamp,
	})

	return rsp, nil

}

// GetSlugVersionOrErr fetches a slack response for 1 slug
func GetSlugVersionOrErr(db *cmddb.DB, sc cmddb.SlackClient, slug string) (cmdutils.SlackResponse, error) {

	payload := cmdutils.SlackResponse{}
	payload.ResponseType = inChannelResponse
	payload.Attachments = []cmdutils.SlackAttachment{}

	// construct the fetch url and fetch the version
	url := fmt.Sprintf("http://%s.%s/%s", slug, sc.Host, sc.VersionPath)
	v := fetchVersion(url)

	if v.Err != nil {
		return payload, v.Err
	}

	// TODO:
	//  format v.Timestamp for displaying proper build date
	//  capitalize cl.server string for SlackResponse.Text

	payload.Text = fmt.Sprintf("_*%s*_ is running version *%s*", title(slug), v.Version)
	payload.Attachments = append(payload.Attachments, cmdutils.SlackAttachment{
		Title:     fmt.Sprintf("%s.%s", slug, sc.Host),
		TitleLink: fmt.Sprintf("http://%s.%s", slug, sc.Host),
		Text:      fmt.Sprintf("Build Date: %s", v.Timestamp),
	})

	// update slug
	err := updateSlugsOrErr(db, sc, slug)

	return payload, err

}

// GetMultiSlugVersionsOrErr fetches a version for each slug in slugs slice
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
			Title:     fmt.Sprintf("%s is running version %s", title(slg), v.Version),
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

	Fetches the version and timestamp from a given url

*/
func fetchVersion(url string) *Version {

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

	ch <- dat

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
		fmt.Printf("New slack_slugs entry: (%v) - %v\n", sc.Name, slug)
		err := db.CreateSlackSlugOrErr(newEntry)
		if err != nil {
			return err
		}
	}

	// TODO:
	//  update slug "updated_at" time if exists

	return nil

}

/*

	Titlize a string

*/
func title(s string) string {

	return strings.Title(s)

}
