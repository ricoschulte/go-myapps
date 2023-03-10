package service_test

import (
	"strings"
	"testing"

	"github.com/ricoschulte/go-myapps/service"
	"gotest.tools/assert"
)

func TestGetDigestForAppLoginFromJson(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		password  string
		challenge string
		expected  string
	}{
		{
			"fritz.box go instance pbxobj go-search",
			`{"mt":"AppLogin","sip":"go-search","guid":"3954ae7854c96301c9dd009033400109","dn":"go-search","digest":"130f168773ff760701b74e629eca2544c6e62d5354e92e6f821001a982eb8951","domain":"fritz.box","app":"searchapi","info":{"appobj":"go-search","appdn":"go-search","appurl":"http://192.168.178.29:5000/fritz.box/go/instance/searchapi","pbx":"pbx-main","cn":"go-search","unlicensed":true,"apps":[]}}`,
			"go",
			"16d7dbbcccb63612",
			"130f168773ff760701b74e629eca2544c6e62d5354e92e6f821001a982eb8951",
		},
		{
			"fritz.box go instance pbxobj go",
			`{"mt":"AppLogin","sip":"go","guid":"8e7ca6fabbc5630170e9009033400109","dn":"go","digest":"8f802873dcdee7cef6144017a6112289298ad69e77474a4c37cfdb68a83640c0","domain":"fritz.box","app":"admin","info":{"appobj":"go","appdn":"go","appurl":"http://192.168.178.29:5000/fritz.box/go/instance/admin","pbx":"pbx-main","cn":"go","unlicensed":true,"apps":[]}}`,
			"go",
			"aeaaebe781e80289",
			"8f802873dcdee7cef6144017a6112289298ad69e77474a4c37cfdb68a83640c0",
		},
		{
			"fritz.box go instance user rico admin",
			`{"mt":"AppLogin","app":"admin","domain":"fritz.box","sip":"rico","guid":"f48b06484a8a61015853009033400109","dn":"Schulte, Rico","info":{"appobj":"go","appdn":"go","appurl":"http://192.168.178.29:5000/fritz.box/go/instance/admin","pbx":"pbx-main","cn":"rico","unlicensed":true,"groups":[],"apps":[]},"digest":"965b4dd4d4f127c99df796b1ed137b40454801f94f50aaf93b2ce98034f10d9c","pbxObj":"go"}`,
			"go",
			"277f41a7a5b7b452",
			"965b4dd4d4f127c99df796b1ed137b40454801f94f50aaf93b2ce98034f10d9c",
		},
		{
			"fritz.box go instance user rico searchapi",
			`{"mt":"AppLogin","app":"searchapi","domain":"fritz.box","sip":"rico","guid":"f48b06484a8a61015853009033400109","dn":"Schulte, Rico","info":{"appobj":"go-search","appdn":"go-search","appurl":"http://192.168.178.29:5000/fritz.box/go/instance/searchapi","pbx":"pbx-main","cn":"rico","testmode":true,"groups":[],"apps":[]},"digest":"378b82aa7508eb53866e3f4de054cdc0f7b5229e9742593e9359db3845681eb6","pbxObj":"go-search"}`,
			"go",
			"7c4fc2e477f69233",
			"378b82aa7508eb53866e3f4de054cdc0f7b5229e9742593e9359db3845681eb6",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mu := &service.MyAppsUtils{}
			result, err := mu.GetDigestForAppLoginFromJson(test.message, test.password, test.challenge)
			assert.NilError(t, err)
			assert.Equal(t, test.expected, result)

		})
	}
}

func TestGetRandomHexString(t *testing.T) {
	mu := &service.MyAppsUtils{}
	length := 100
	result := mu.GetRandomHexString(length)
	result2 := mu.GetRandomHexString(length)
	println("test result: ", result)
	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}
	for i := range result {
		if !strings.Contains("abcdef0123456789", string(result[i])) {
			t.Errorf("Expected hexadecimal characters, got %s", string(result[i]))
		}
	}
	if result == result2 {
		t.Errorf("two equal strings received, when two should be not equal aka random")
	}
}
