package utils

import (
	_ "bytes"
	"encoding/json"
	"io"
	"net/http"
)

// FetchData fetchData effectue une requête HTTP et renvoie la réponse sous forme d'objet ou de tableau
func FetchData(url, method string) (interface{}, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		//log
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//log
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//log
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		//log
		return nil, err
	}

	return data, nil
}
