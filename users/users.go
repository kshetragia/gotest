// +build windows

package users

import (
	"gotest/winapi"

	"github.com/pkg/errors"

	"golang.org/x/sys/windows"
)

type User struct {
	Name             string        `json:"Name"`
	Domain           string        `json:"Domain"`
	SID              string        `json:"SID"`
	SessionID        uint32        `json:"SessionID"`
	LastSuccessLogon string        `json:"LastSuccessLogon"`
	AuthenticationID *windows.LUID `json:"LUID"`
}

type Users map[int64]User

// Init creates new Users map
func Init() Users {
	return Users(make(map[int64]User))
}

// Add takes windows LUID structure, collects info
// about user related to mentioned LUID and saves collected data into Users map
func (u Users) Add(luid *windows.LUID) error {
	key, err := u.Key(luid)
	if err != nil {
		return errors.Wrap(err, "get Users map Key")
	}

	// Record already exists
	if _, ok := u[key]; ok {
		return nil
	}

	// Getting LSA Logon info
	var sessionData *winapi.SecurityLogonSessionData
	err = winapi.LsaGetLogonSessionData(luid, &sessionData)
	if err != nil {
		return errors.Wrap(err, "get logon session data")
	}

	var user = User{
		Name:             winapi.LsaUnicodeToString(sessionData.UserName),
		Domain:           winapi.LsaUnicodeToString(sessionData.LogonDomain),
		SID:              sessionData.Sid.String(),
		LastSuccessLogon: winapi.WinToUnixTime(sessionData.LogonTime).String(),
		SessionID:        sessionData.Session,
		AuthenticationID: luid,
	}

	u[key] = user

	return nil
}

func (u Users) Get(luid *windows.LUID) (*User, error) {
	key, err := u.Key(luid)
	if err != nil {
		return nil, errors.Wrap(err, "get Users map Key")
	}
	if val, ok := u[key]; ok {
		return &val, nil
	}
	return nil, nil
}

// Key takes windows.LUID structure and converts it into int64 luid number
func (u Users) Key(luid *windows.LUID) (int64, error) {
	if luid == nil {
		return 0, errors.New("got empty LUID pointer")
	}
	key := int64(int64(luid.HighPart<<32) + int64(luid.LowPart))
	return key, nil
}
