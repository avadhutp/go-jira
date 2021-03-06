package jira

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAcquireSessionCookie_Fail(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/auth/1/session", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/rest/auth/1/session")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error in read body: %s", err)
		}
		if bytes.Index(b, []byte(`"username":"foo"`)) < 0 {
			t.Error("No username found")
		}
		if bytes.Index(b, []byte(`"password":"bar"`)) < 0 {
			t.Error("No password found")
		}

		// Emulate error
		w.WriteHeader(http.StatusInternalServerError)
	})

	res, err := testClient.Authentication.AcquireSessionCookie("foo", "bar")
	if err == nil {
		t.Errorf("Expected error, but no error given")
	}
	if res == true {
		t.Error("Expected error, but result was true")
	}
}

func TestAcquireSessionCookie_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/auth/1/session", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/rest/auth/1/session")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error in read body: %s", err)
		}
		if bytes.Index(b, []byte(`"username":"foo"`)) < 0 {
			t.Error("No username found")
		}
		if bytes.Index(b, []byte(`"password":"bar"`)) < 0 {
			t.Error("No password found")
		}

		fmt.Fprint(w, `{"session":{"name":"JSESSIONID","value":"12345678901234567890"},"loginInfo":{"failedLoginCount":10,"loginCount":127,"lastFailedLoginTime":"2016-03-16T04:22:35.386+0000","previousLoginTime":"2016-03-16T04:22:35.386+0000"}}`)
	})

	res, err := testClient.Authentication.AcquireSessionCookie("foo", "bar")
	if err != nil {
		t.Errorf("No error expected. Got %s", err)
	}
	if res == false {
		t.Error("Expected result was true. Got false")
	}
}
