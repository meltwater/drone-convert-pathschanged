package main

import (
	"fmt"
	"testing"
)

func TestValidateSecretMissing(t *testing.T) {
	s := &spec{
		Secret: "",
	}

	got := validate(s)

	want := "missing secret key"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestValidateProviderMissing(t *testing.T) {
	s := &spec{
		Provider: "",
		Secret:   "abcdefg",
	}

	got := validate(s)

	want := "missing provider"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestValidateProviderUnsupported(t *testing.T) {
	s := &spec{
		Provider: "unsupported",
		Secret:   "abcdefg",
	}

	got := validate(s)

	want := "unsupported provider"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestValidateTokenMissing(t *testing.T) {
	// bitbucket-server/stash and github use tokens for authentication
	providers := []string{
		"bitbucket-server",
		"github",
		"stash",
	}
	for _, provider := range providers {
		s := &spec{
			Provider: provider,
			Secret:   "abcdefg",
			Token:    "",
		}

		got := validate(s)

		want := "missing token"

		if got.Error() != want {
			t.Errorf("wanted %s, got %s", want, got)
		}
	}
}

func TestValidateBitbucketUserMissing(t *testing.T) {
	s := &spec{
		BitBucketUser: "",
		Provider:      "bitbucket",
		Secret:        "abcdefg",
	}

	got := validate(s)

	want := "missing bitbucket user"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestValidateBitbucketPasswordMissing(t *testing.T) {
	s := &spec{
		BitBucketUser:     "centauri",
		BitBucketPassword: "",
		Provider:          "bitbucket",
		Secret:            "abcdefg",
	}

	got := validate(s)

	want := "missing bitbucket password"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestValidateBitbucketServerAddressMissing(t *testing.T) {
	s := &spec{
		BitBucketAddress: "",
		Provider:         "bitbucket-server",
		Secret:           "abcdefg",
		Token:            "abcdefg",
	}

	got := validate(s)

	want := "missing bitbucket server address"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

// this tests backwards compatibility with bitbucket-server for stash
func TestValidateBitbucketServerStashCompatibility(t *testing.T) {
	s := &spec{
		BitBucketAddress: "example.com",
		Provider:         "bitbucket-server",
		Secret:           "abcdefg",
		Token:            "abcdefg",
	}

	err := validate(s)
	if err != nil {
		t.Error(err)
	}

	got := fmt.Sprintf("%s %s", s.Provider, s.StashServer)
	want := fmt.Sprintf("stash %s", s.BitBucketAddress)

	if got != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestValidateStashServerMissing(t *testing.T) {
	s := &spec{
		Provider:    "stash",
		Secret:      "abcdefg",
		StashServer: "",
		Token:       "abcdefg",
	}

	got := validate(s)

	want := "missing stash server"
	if got.Error() != want {
		t.Errorf("wanted %s, got %s", want, got)
	}
}
