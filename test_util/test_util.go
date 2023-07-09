package test_util

import (
	"bytes"
	"colmanback/objects"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func CheckExists(t *testing.T, router *mux.Router, getURL string, expectExists bool) {
	req, err := http.NewRequest(http.MethodGet, getURL, nil)

	if err != nil {
		t.Errorf("An error has been reported when preparing the get request (CheckExists): %v\n", err)
	} else {
		resp := ExecuteRequest(router, req)
		if !expectExists {
			if resp.Code != http.StatusNotFound {
				t.Errorf("Object found when it was not expected. Response code: %d", resp.Code)
			}
		} else {
			if resp.Code != http.StatusOK {
				t.Errorf("Object not found when it was expected. Response code: %d", resp.Code)
			}
		}
	}
}

func CheckField(t *testing.T, fieldName string, expectedValue string, actualValue string) {
	if expectedValue != actualValue {
		t.Errorf("Field %s does not match. Expected %s but got %s.", fieldName, expectedValue, actualValue)
	}
}

func CheckPut(t *testing.T, router *mux.Router, jsonString string, putURL string) {
	req, err := http.NewRequest(http.MethodPut, putURL, bytes.NewBuffer([]byte(jsonString)))

	if err != nil {
		t.Errorf("An error has been reported when preparing the put request for %s. Error: %v\n", jsonString, err)
	} else {
		resp := ExecuteRequest(router, req)
		if resp.Code != http.StatusOK {
			t.Errorf("Status code not as expected after put.")
		}
	}
}

func CheckList[K objects.Object](t *testing.T, router *mux.Router, listURL string, objectInstList *[]K) {
	req, err := http.NewRequest(http.MethodGet, listURL, nil)
	if err != nil {
		t.Errorf("An error has been reported when preparing the get request for the full list. Error: %v\n", err)
	} else {
		resp := ExecuteRequest(router, req)
		if resp.Code != http.StatusOK {
			t.Errorf("Status code not as expected for the get of the full list. Code %d", resp.Code)
		} else {
			body := resp.Body.Bytes()
			errUnmarshall := json.Unmarshal(body, objectInstList)
			if errUnmarshall != nil {
				t.Errorf("An error has been reported when unmarshalling full list. Error: %v", errUnmarshall)
			} else {
				if len(*objectInstList) == 0 {
					t.Errorf("The full list retrieved from the API is not complete. Entries: %d", len(*objectInstList))
				}
			}
		}
	}
}

func CheckFields[K objects.Object](t *testing.T, router *mux.Router, expectedObjectInst K, getURL string, newObjectInstMem K, comparator func(t *testing.T, origObject K, newObject K)) {
	req, err := http.NewRequest(http.MethodGet, getURL, nil)

	if err != nil {
		t.Errorf("An error has been reported when preparing the get request (CheckFields): %v\n", err)
	} else {
		resp := ExecuteRequest(router, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Object found when it wasn't expected. Response code: %d", resp.Code)
		} else {
			unmarshallErr := json.Unmarshal(resp.Body.Bytes(), newObjectInstMem)
			if unmarshallErr != nil {
				t.Errorf("An error has occurred when unmarshalling single object instance. Error: %v", unmarshallErr)
			} else {
				comparator(t, expectedObjectInst, newObjectInstMem)
			}
		}
	}
}

func CheckDelete(t *testing.T, router *mux.Router, deleteURL string, expectOK bool) {
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)

	if err != nil {
		t.Errorf("An error has been reported when preparing delete request: %v\n", err)
	} else {
		resp := ExecuteRequest(router, req)
		if (expectOK && resp.Code != http.StatusOK) || (!expectOK && resp.Code != http.StatusNotFound) {
			t.Errorf("Status code not as expected.")
		}
	}
}

func ExecuteRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}
