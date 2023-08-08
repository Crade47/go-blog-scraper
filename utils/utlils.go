package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CheckMediumURL(urlString string) error {
	if !strings.Contains(urlString, "medium") {
		fmt.Println("Need a valid medium.com link")
		return errors.New("Need a valid medium.com link")
	}
	if !strings.Contains(urlString, "https://") {
		urlString = "https://" + urlString
	}
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		fmt.Println("Error in parsing the url / Invalid URL")
		return err
	}

	resp, err := http.Get(urlString)
	if err != nil {
		fmt.Println("Error in performing GET request on URL")
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("The page does not exist.")
		return errors.New("page does not exist")
	} else {
		errorString := "Received status code:" + strconv.Itoa(resp.StatusCode)
		fmt.Println(errorString)
		return errors.New(errorString)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
