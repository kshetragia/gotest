package privilege

import (
	"testing"
)

func TestPrivilege(t *testing.T) {
	name := "SeShutdownPrivilege"

	status := IsEnabled(name)

	// Let's try to change permissions
	if !Set("SeShutdownPrivilege", !status) {
		t.Errorf("Change status '%v' privilege", name)
	}

	if IsEnabled(name) == status {
		t.Errorf("Privilege '%v' was not changed", name)
	}

	// Revert perms back
	if !Set("SeShutdownPrivilege", status) {
		t.Errorf("Change status back '%v' privilege", name)
	}

	if IsEnabled(name) != status {
		t.Errorf("Privilege '%v' was not changed back", name)
	}
}
