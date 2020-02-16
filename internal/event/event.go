package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/denisbrodbeck/machineid"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/version"
	"github.com/spf13/viper"
)

var (
	proactionURI = "https://oss.proaction.io"
	userAgent    = fmt.Sprintf("Proaction/%s", version.Version())
)

type Event struct {
	Name string `json:"name"`
}

func Init(v *viper.Viper) error {
	if v.GetBool("no-track") {
		return nil
	}

	machineID, err := machineid.ProtectedID("proaction")
	if err != nil {
		return errors.Wrap(err, "failed to get machine id")
	}

	e := Event{
		Name: "scan",
	}
	b, err := json.Marshal(e)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/event", proactionURI), bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Authorization", machineID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}
	defer resp.Body.Close()

	return nil
}
