package cmdversions

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// GetAppVersionOrErr fetches the apps current version from file
func GetAppVersionOrErr() (Version, error) {

	var v Version
	curdir, err := os.Getwd()

	if err != nil {
		return v, err
	}
	fullpath := curdir + "/version.json"

	v, err = readVersionFromFileOrErr(fullpath)
	return v, err

}

func readVersionFromFileOrErr(filename string) (Version, error) {

	v := Version{}
	raw, err := ioutil.ReadFile(filename)

	if err != nil {
		return v, err
	}

	err = json.Unmarshal(raw, &v)
	return v, err

}
