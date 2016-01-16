package hostqueue

import (
	"testing"
)

func TestFoundValidResponder(t *testing.T) {
	group := createGroup()
	from := "test@test.com"

	isValid := isValidResponder(from, group)

	if !isValid {
		t.Error("Test failed")
	}
}

func TestFoundInvalidResponder(t *testing.T) {
	group := createGroup()
	from := "abcd@test.com"

	isValid := isValidResponder(from, group)

	if isValid {
		t.Error("Test failed")
	}
}

func createGroup() Group {
	var group Group
	var host1, host2 Host
	host1.Emails = "test@test.com, test2@test.com"
	host2.Emails = "test3@test.com, test4@test.com"
	group.Hosts = []Host{host1, host2}
	return group
}
