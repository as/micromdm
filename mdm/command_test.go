package mdm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/groob/plist"
)

// Basic tests will attempt to marshal and unmarshal mdm command structures to identify any naming or tag errors.

// Make sure a command can be marshalled to json
func testMarshalJSON(t *testing.T, cmd interface{}) {
	jsonCmd, err := json.Marshal(cmd)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(jsonCmd))
}

// Make sure a command can be marshalled to plist
func testMarshalPlist(t *testing.T, cmd interface{}) {
	plistCmd, err := plist.MarshalIndent(cmd, "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(plistCmd))
}

func TestInstallProfile(t *testing.T) {
	cmd := InstallProfile{Payload: []byte{00}}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestRemoveProfile(t *testing.T) {
	cmd := RemoveProfile{Identifier: "io.micromdm.test.profile"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestInstallProvisioningProfile(t *testing.T) {
	cmd := InstallProvisioningProfile{ProvisioningProfile: []byte{00}}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestRemoveProvisioningProfile(t *testing.T) {
	cmd := RemoveProvisioningProfile{UUID: "1111-2222-3333"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestInstalledApplicationList(t *testing.T) {
	cmd := InstalledApplicationList{Identifiers: []string{"io.micromdm.application"}, ManagedAppsOnly: true}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestDeviceInformation(t *testing.T) {
	cmd := DeviceInformation{Queries: []string{"SerialNumber"}}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestDeviceLock(t *testing.T) {
	cmd := DeviceLock{PIN: "123456", Message: "Locked", PhoneNumber: "123-4567"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestClearPasscode(t *testing.T) {
	cmd := ClearPasscode{UnlockToken: []byte{00}}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestEraseDevice(t *testing.T) {
	cmd := EraseDevice{PIN: "123456"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestRequestMirroring(t *testing.T) {
	cmd := RequestMirroring{DestinationName: "Apple TV", ScanTime: "30", Password: "sekret"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

//func TestRestrictions(t *testing.T) {
//	cmd := Restrictions{}
//}

func TestSettings(t *testing.T) {
	cmd := Settings{
		Settings: []Setting{
			Setting{
				Item:       "DeviceName",
				DeviceName: stringPtr("foo"),
			},
		},
	}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func stringPtr(s string) *string {
	return &s
}

func TestDeleteUser(t *testing.T) {
	cmd := DeleteUser{UserName: "joe", ForceDeletion: false}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestEnableLostMode(t *testing.T) {
	cmd := EnableLostMode{Message: "Lost!", PhoneNumber: "123-4567", Footnote: "This is lost"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestInstallApplication(t *testing.T) {
	cmd := InstallApplication{ITunesStoreID: 1234567}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestApplyRedemptionCode(t *testing.T) {
	cmd := ApplyRedemptionCode{Identifier: "id", RedemptionCode: "abcdefg"}
	testMarshalJSON(t, cmd)
	testMarshalPlist(t, cmd)
}

func TestMarshalWithOverlappingKeys(t *testing.T) {
	rp := RemoveProfile{Identifier: "io.micromdm.test.profile"}

	cmd := Command{RequestType: "RemoveProfile", RemoveProfile: rp}
	out, err := plist.Marshal(&cmd)
	if err != nil {
		t.Fatal(err)
	}

	// check that output contains Identifier field
	if !bytes.Contains(out, []byte("Identifier")) {
		t.Errorf("expected %v plist to contain identifier key, got \n%s", cmd.RequestType, out)
	}

}

func TestUnmarshalOverlaplingKeys(t *testing.T) {
	data := []byte(`{
		"request_type":"RemoveProfile",
		"udid":"abcd",
		"identifier":"aaaa"
	}`)

	var req CommandRequest
	if err := json.Unmarshal(data, &req); err != nil {
		t.Fatal(err)
	}

	if req.RemoveProfile.Identifier != "aaaa" {
		t.Errorf("expected RemoveProfile struct with Identifier field")
	}
}

func Test_PayloadWithNoFields(t *testing.T) {
	data := []byte(`{
		"request_type":"ProfileList",
		"udid":"abcd"
	}`)

	var req CommandRequest
	if err := json.Unmarshal(data, &req); err != nil {
		t.Fatal(err)
	}

	payload, err := NewPayload(&req)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if err := plist.NewEncoder(buf).Encode(payload); err != nil {
		t.Fatal(err)
	}

	var have Payload
	if err := plist.NewDecoder(buf).Decode(&have); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, *payload) {
		t.Errorf("have %+v\n want %+v", have, *payload)
	}
}
