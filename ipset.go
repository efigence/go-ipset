package ipset

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
)

var bin = ""

func init() {
	bin, _ = exec.LookPath("ipset")
}

func runError(cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"error running %s: %s[%s])",
			strings.Join(cmd.Args, " "),
			string(out),
			err,
		)
	}
	return nil
}

// https://www.onlinetool.io/xmltogo/
type Ipsets struct {
	XMLName xml.Name `xml:"ipsets"`
	Text    string   `xml:",chardata"`
	Ipset   []struct {
		Text     string `xml:",chardata"`
		Name     string `xml:"name,attr"`
		Type     string `xml:"type"`
		Revision string `xml:"revision"`
		Header   struct {
			Text       string `xml:",chardata"`
			Family     string `xml:"family"`
			Hashsize   string `xml:"hashsize"`
			Maxelem    string `xml:"maxelem"`
			Memsize    string `xml:"memsize"`
			References string `xml:"references"`
			Numentries string `xml:"numentries"`
		} `xml:"header"`
		Members struct {
			Text   string `xml:",chardata"`
			Member []struct {
				Text string `xml:",chardata"`
				Elem string `xml:"elem"`
			} `xml:"member"`
		} `xml:"members"`
	} `xml:"ipset"`
}

type IPSet interface {
	Add(string) error
	Delete(string) error
	Create() error
	Destroy() error
	Swap(string) error
	Exist() bool
	Name() string
}
