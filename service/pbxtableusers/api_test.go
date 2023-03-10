package pbxtableusers_test

import (
	"encoding/json"
	"testing"

	"github.com/ricoschulte/go-myapps/service/pbxtableusers"
)

func TestPbxTableUsers_ReplicateUpdate(t *testing.T) {

	tests := []struct {
		name string
		json string
	}{

		{
			name: "wakeup",
			json: `{"mt":"ReplicateUpdate","src":"src_1677997325518944961","api":"PbxTableUsers","columns":{"apps-my":"users","cn":"Charlotte Maihoff","config":"Config User","dn":"Maihoff, Charlotte","e164":"403","fax":false,"grps":[{"name":"Support"},{"name":"pats_haus1"}],"guid":"66eaaee678464cffb9a3c8ae259fedd2","h323":"charlottemaihoff","hide":false,"loc":"pbx-main","node":"master","pwd":"f8c26f062fedf2ff2980b8b8bdbb01a6df91add1f605e7f5","t-allows":[{"name":"@fritz.box","online":true,"presence":true,"otf":true,"note":true,"visible":true}],"wakeups":[{"h":2,"m":3,"s":4,"name":"44444","num":"44333","retry":2,"mult":true,"to":34,"fallback":"3333","bool":"my_new_boolean","bool-not":true}]}}`,
		},
		{
			name: "fork",
			json: `{"mt":"ReplicateUpdate","src":"src_1677997916352198068","api":"PbxTableUsers","columns":{"apps-my":"users","cn":"Charlotte Maihoff","config":"Config User","dn":"Maihoff, Charlotte","e164":"403","fax":false,"forks":[{"h323":"aaaaaaaaaaaa","bool":"my_new_boolean","bool-not":true,"mobility":"vvvvvvvvvvvvvvv","delay":444,"hw":"vvvvvvvvvvvvvvvv","app":"vvvvvvvvvvvvv","off":true,"cw":true,"min":5,"max":5}],"grps":[{"name":"Support"},{"name":"pats_haus1"}],"guid":"66eaaee678464cffb9a3c8ae259fedd2","h323":"charlottemaihoff","hide":false,"loc":"pbx-main","node":"master","pwd":"f8c26f062fedf2ff2980b8b8bdbb01a6df91add1f605e7f5","t-allows":[{"name":"@fritz.box","online":true,"presence":true,"otf":true,"note":true,"visible":true}],"wakeups":[]}}`,
		},
		{
			name: "cds",
			json: `{"mt":"ReplicateUpdate","src":"src_1677998086155540423","api":"PbxTableUsers","columns":{"apps-my":"users","cds":[{"type":"cfu","bool":"my_new_boolean","e164":"111233","h323":"sdsdsd","src":"<src type=\"do\"><ep e164=\"3333\" h323=\"dddd\" ext=\"true\" fwd=\"both\"/></src>","precedence":true},{"type":"cfb","bool":"Business hours","bool-not":true,"e164":"4444","h323":"ffff","src":"<src type=\"dont\"><ep e164=\"4444\" h323=\"tggggg\" ext=\"false\" fwd=\"fwd\"/></src>"},{"type":"cfnr","h323":"44444444","precedence":true}],"cn":"Charlotte Maihoff","config":"Config User","dn":"Maihoff, Charlotte","e164":"403","fax":false,"grps":[{"name":"Support"},{"name":"pats_haus1"}],"guid":"66eaaee678464cffb9a3c8ae259fedd2","h323":"charlottemaihoff","hide":false,"loc":"pbx-main","node":"master","pwd":"f8c26f062fedf2ff2980b8b8bdbb01a6df91add1f605e7f5","t-allows":[{"name":"@fritz.box","online":true,"presence":true,"otf":true,"note":true,"visible":true}],"wakeups":[]}}`,
		},
		{
			name: "device",
			json: `{"mt":"ReplicateUpdate","src":"src_1677998086155540423","api":"PbxTableUsers","columns":{"apps-my":"users","cn":"Charlotte Maihoff","config":"Config User","devices":[{"hw":"hwid","text":"name","app":"app","admin":true,"no-filter":true,"tls":true,"no-mob":true,"trusted":true,"sreg":true,"mr":true}],"dn":"Maihoff, Charlotte","e164":"403","fax":false,"grps":[{"name":"Support"},{"name":"pats_haus1"}],"guid":"66eaaee678464cffb9a3c8ae259fedd2","h323":"charlottemaihoff","hide":false,"loc":"pbx-main","node":"master","pwd":"f8c26f062fedf2ff2980b8b8bdbb01a6df91add1f605e7f5","t-allows":[{"name":"@fritz.box","online":true,"presence":true,"otf":true,"note":true,"visible":true}],"wakeups":[]}}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := pbxtableusers.ReplicateUpdate{}
			if err := json.Unmarshal([]byte(test.json), &msg); err != nil {
				t.Fatalf("error unmarshalling message: %v", err)
			}

		})
	}
}
