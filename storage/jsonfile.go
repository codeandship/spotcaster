
package storage

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/iwittkau/spotcaster"
)

func WriteToken(t spotcaster.Token, name string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, os.ModePerm)
}

func ReadToken(name string) (spotcaster.Token, error) {
	t := spotcaster.Token{}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(data, &t)
	return t, err
}