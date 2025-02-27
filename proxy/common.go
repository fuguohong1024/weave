package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/fuguohong1024/weave/common"
	weavedocker "github.com/fuguohong1024/weave/common/docker"
)

var (
	containerIDRegexp   = regexp.MustCompile(`^(/v[0-9.]*)?/containers/([^/]*)/.*`)
	weaveWaitEntrypoint = []string{"/w/w"}

	Log = common.Log
)

func unmarshalRequestBody(r *http.Request, target interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	Log.Debugf("->requestBody: %s", body)
	if err := r.Body.Close(); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber() // don't want large numbers in scientific format
	return d.Decode(&target)
}

func marshalRequestBody(r *http.Request, body interface{}) error {
	newBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	Log.Debugf("<-requestBody: %s", newBody)
	r.Body = ioutil.NopCloser(bytes.NewReader(newBody))
	r.ContentLength = int64(len(newBody))
	return nil
}

func unmarshalResponseBody(r *http.Response, target interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	Log.Debugf("->responseBody: %s", body)
	if err := r.Body.Close(); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber() // don't want large numbers in scientific format
	return d.Decode(&target)
}

func marshalResponseBody(r *http.Response, body interface{}) error {
	newBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	Log.Debugf("<-responseBody: %s", newBody)
	r.Body = ioutil.NopCloser(bytes.NewReader(newBody))
	r.ContentLength = int64(len(newBody))
	// Stop it being chunked, because that hangs
	r.TransferEncoding = nil
	return nil
}

func inspectContainerInPath(client *weavedocker.Client, path string) (*docker.Container, error) {
	subs := containerIDRegexp.FindStringSubmatch(path)
	if subs == nil {
		err := fmt.Errorf("No container id found in request with path %s", path)
		Log.Warningln(err)
		return nil, err
	}
	containerID := subs[2]

	container, err := client.InspectContainer(containerID)
	if err != nil {
		Log.Warningf("Error inspecting container %s: %v", containerID, err)
	}
	return container, err
}

func addVolume(hostConfig jsonObject, source, target, mode string) error {
	configBinds, err := hostConfig.StringArray("Binds")
	if err != nil {
		return err
	}

	var binds []string
	for _, bind := range configBinds {
		s := strings.Split(bind, ":")
		if len(s) >= 2 && s[1] == target {
			continue
		}
		binds = append(binds, bind)
	}
	bind := source + ":" + target
	if mode != "" {
		bind += ":" + mode
	}
	hostConfig["Binds"] = append(binds, bind)
	return nil
}
