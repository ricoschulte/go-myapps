package sysclient_test

import (
	"testing"

	"github.com/ricoschulte/go-myapps/sysclient"
	"gotest.tools/assert"
)

func TestMarshalIdentifyMessage(t *testing.T) {
	tests := []struct {
		Name string

		Id           string
		Product      string
		Version      string
		FwBuild      string
		BcBuild      string
		Major        string
		Fw           string
		Bc           string
		Mini         bool
		PbxActive    bool
		Other        bool
		Platform     sysclient.Platform
		Digest       string
		EthIfs       []sysclient.EthIf
		expectedJson string
	}{
		{
			Name:      "all booleans false",
			Id:        "f09033480af9",
			Product:   "IP222",
			Version:   "13r2 dvl [13.4250/131286/1300]",
			FwBuild:   "134250",
			BcBuild:   "131286",
			Major:     "13r2",
			Fw:        "ip222.bin",
			Bc:        "boot222.bin",
			Mini:      false,
			PbxActive: false,
			Other:     false,
			Platform: sysclient.Platform{
				Type: "PHONE",
				Fxs:  false,
			},
			Digest: "f3d37205d8eb93a4770035853c12984",
			EthIfs: []sysclient.EthIf{
				{
					If:   "ETH0",
					Ipv4: "172.16.4.141",
					Ipv6: "2002:91fd:9d07:0:290:33ff:fe46:af2",
				},
			},
			expectedJson: `{"id":"f09033480af9","product":"IP222","version":"13r2 dvl [13.4250/131286/1300]","fwBuild":"134250","bcBuild":"131286","major":"13r2","fw":"ip222.bin","bc":"boot222.bin","mini":false,"pbxActive":false,"other":false,"platform":{"type":"PHONE","fxs":false},"digest":"f3d37205d8eb93a4770035853c12984","ethIfs":[{"if":"ETH0","ipv4":"172.16.4.141","ipv6":"2002:91fd:9d07:0:290:33ff:fe46:af2"}]}`,
		},
		{
			Name:      "no EthIfs",
			Id:        "f09033480af9",
			Product:   "IP222",
			Version:   "13r2 dvl [13.4250/131286/1300]",
			FwBuild:   "134250",
			BcBuild:   "131286",
			Major:     "13r2",
			Fw:        "ip222.bin",
			Bc:        "boot222.bin",
			Mini:      false,
			PbxActive: false,
			Other:     false,
			Platform: sysclient.Platform{
				Type: "PHONE",
				Fxs:  false,
			},
			Digest: "f3d37205d8eb93a4770035853c12984",
			// EthIfs: []sysclient.EthIf{
			// 	{
			// 		If:   "ETH0",
			// 		Ipv4: "172.16.4.141",
			// 		Ipv6: "2002:91fd:9d07:0:290:33ff:fe46:af2",
			// 	},
			// },
			expectedJson: `{"id":"f09033480af9","product":"IP222","version":"13r2 dvl [13.4250/131286/1300]","fwBuild":"134250","bcBuild":"131286","major":"13r2","fw":"ip222.bin","bc":"boot222.bin","mini":false,"pbxActive":false,"other":false,"platform":{"type":"PHONE","fxs":false},"digest":"f3d37205d8eb93a4770035853c12984"}`,
		},
		{
			Name:      "two EthIfs",
			Id:        "f09033480af9",
			Product:   "IP222",
			Version:   "13r2 dvl [13.4250/131286/1300]",
			FwBuild:   "134250",
			BcBuild:   "131286",
			Major:     "13r2",
			Fw:        "ip222.bin",
			Bc:        "boot222.bin",
			Mini:      false,
			PbxActive: false,
			Other:     false,
			Platform: sysclient.Platform{
				Type: "PHONE",
				Fxs:  false,
			},
			//Digest: "f3d37205d8eb93a4770035853c12984",
			EthIfs: []sysclient.EthIf{
				{
					If:   "ETH0",
					Ipv4: "172.16.4.141",
					Ipv6: "2002:91fd:9d07:0:290:33ff:fe46:af2",
				},
				{
					If:   "ETH1",
					Ipv4: "172.16.2.141",
					Ipv6: "2002:91fd:9d07:1:290:33ff:fe46:af2",
				},
			},
			expectedJson: `{"id":"f09033480af9","product":"IP222","version":"13r2 dvl [13.4250/131286/1300]","fwBuild":"134250","bcBuild":"131286","major":"13r2","fw":"ip222.bin","bc":"boot222.bin","mini":false,"pbxActive":false,"other":false,"platform":{"type":"PHONE","fxs":false},"ethIfs":[{"if":"ETH0","ipv4":"172.16.4.141","ipv6":"2002:91fd:9d07:0:290:33ff:fe46:af2"},{"if":"ETH1","ipv4":"172.16.2.141","ipv6":"2002:91fd:9d07:1:290:33ff:fe46:af2"}]}`,
		},
		{
			Name:      "no digest",
			Id:        "f09033480af9",
			Product:   "IP222",
			Version:   "13r2 dvl [13.4250/131286/1300]",
			FwBuild:   "134250",
			BcBuild:   "131286",
			Major:     "13r2",
			Fw:        "ip222.bin",
			Bc:        "boot222.bin",
			Mini:      false,
			PbxActive: false,
			Other:     false,
			Platform: sysclient.Platform{
				Type: "PHONE",
				Fxs:  false,
			},
			//Digest: "f3d37205d8eb93a4770035853c12984",
			EthIfs: []sysclient.EthIf{
				{
					If:   "ETH0",
					Ipv4: "172.16.4.141",
					Ipv6: "2002:91fd:9d07:0:290:33ff:fe46:af2",
				},
			},
			expectedJson: `{"id":"f09033480af9","product":"IP222","version":"13r2 dvl [13.4250/131286/1300]","fwBuild":"134250","bcBuild":"131286","major":"13r2","fw":"ip222.bin","bc":"boot222.bin","mini":false,"pbxActive":false,"other":false,"platform":{"type":"PHONE","fxs":false},"ethIfs":[{"if":"ETH0","ipv4":"172.16.4.141","ipv6":"2002:91fd:9d07:0:290:33ff:fe46:af2"}]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			identityA := sysclient.Identity{
				Id:        test.Id,
				Product:   test.Product,
				Version:   test.Version,
				FwBuild:   test.FwBuild,
				BcBuild:   test.BcBuild,
				Major:     test.Major,
				Fw:        test.Fw,
				Bc:        test.Bc,
				Mini:      test.Mini,
				PbxActive: test.PbxActive,
				Other:     test.Other,
				Platform: sysclient.Platform{
					Type: test.Platform.Type,
					Fxs:  test.Platform.Fxs,
				},
				Digest: test.Digest,
				EthIfs: test.EthIfs,
			}

			// convert to Json
			json_bytes, err := identityA.ToBytes()
			assert.NilError(t, err, "converting to bytes failed with error")
			assert.DeepEqual(t, test.expectedJson, string(json_bytes))

			// convert from Bytes to Instance
			identity, err_parsing := sysclient.NewIdentity(json_bytes)
			assert.NilError(t, err_parsing, "error parsing from json")

			assert.Equal(t, test.Id, identity.Id, "Id is not the expected one")
			assert.Equal(t, test.Product, identity.Product, "Product is not the expected one")
			assert.Equal(t, test.Version, identity.Version, "Version is not the expected one")
			assert.Equal(t, test.FwBuild, identity.FwBuild, "FwBuild is not the expected one")
			assert.Equal(t, test.BcBuild, identity.BcBuild, "BcBuild is not the expected one")
			assert.Equal(t, test.Major, identity.Major, "Major is not the expected one")
			assert.Equal(t, test.Fw, identity.Fw, "Fw is not the expected one")
			assert.Equal(t, test.Bc, identity.Bc, "Bc is not the expected one")
			assert.Equal(t, test.Mini, identity.Mini, "Mini is not the expected one")
			assert.Equal(t, test.PbxActive, identity.PbxActive, "PbxActive is not the expected one")
			assert.Equal(t, test.Other, identity.Other, "Other is not the expected one")
			assert.Equal(t, test.Platform, identity.Platform, "Platform is not the expected one")
			assert.Equal(t, test.Digest, identity.Digest, "Digest is not the expected one")
			assert.Equal(t, test.Platform, identity.Platform, "Platform is not the expected one")
			for i, eth := range test.EthIfs {
				assert.Equal(t, test.EthIfs[i].If, eth.If, "EthIfs.If is not the expected one")
				assert.Equal(t, test.EthIfs[i].Ipv4, eth.Ipv4, "EthIfs.Ipv4 is not the expected one")
				assert.Equal(t, test.EthIfs[i].Ipv6, eth.Ipv6, "EthIfs.Ipv6 is not the expected one")
			}
			assert.Equal(t, test.Platform, identity.Platform, "Platform is not the expected one")
			assert.Equal(t, test.Platform, identity.Platform, "Platform is not the expected one")

		})
	}

}
