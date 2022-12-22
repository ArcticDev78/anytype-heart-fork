package uri

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
)

var (
	// RFC 5322 mail regex
	noPrefixEmailRegexp = regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$`)
	// RFC 3966 tel regex
	noPrefixTelRegexp      = regexp.MustCompile(`^((?:\+[\d().-]*\d[\d().-]*|[0-9A-F*#().-]*[0-9A-F*#][0-9A-F*#().-]*(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*;phone-context=(?:\+[\d().-]*\d[\d().-]*|(?:[a-z0-9]\.|[a-z0-9][a-z0-9-]*[a-z0-9]\.)*(?:[a-z]|[a-z][a-z0-9-]*[a-z0-9])))(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*(?:,(?:\+[\d().-]*\d[\d().-]*|[0-9A-F*#().-]*[0-9A-F*#][0-9A-F*#().-]*(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*;phone-context=\+[\d().-]*\d[\d().-]*)(?:;[a-z\d-]+(?:=(?:[a-z\d\[\]\/:&+$_!~*'().-]|%[\dA-F]{2})+)?)*)*)$`)
	noPrefixHttpRegex      = regexp.MustCompile(`^[\pL\d.-]+(?:\.[\pL\\d.-]+)+[\pL\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.\/\d]+$`)
	haveUriSchemeRegex     = regexp.MustCompile(`^([a-zA-Z][A-Za-z0-9+.-]*):[\S]+`)
	winFilepathPrefixRegex = regexp.MustCompile(`^[a-zA-Z]:[\\\/]`)
)

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}

	return noPrefixEmailRegexp.MatchString(email)
}

func ValidatePhone(phone string) bool {
	if len(phone) == 0 {
		return false
	}

	return noPrefixTelRegexp.MatchString(phone)
}

// ProcessURI tries to verify the web URI and return the normalized URI
func ProcessURI(url string) (urlOut string, err error) {
	if len(url) == 0 {
		return url, fmt.Errorf("url is empty")

	} else if noPrefixEmailRegexp.MatchString(url) {
		return "mailto:" + url, nil

	} else if noPrefixTelRegexp.MatchString(url) {
		return "tel:" + url, nil

	} else if winFilepathPrefixRegex.MatchString(url) {
		return "", fmt.Errorf("filepath not supported")

	} else if strings.HasPrefix(url, string(os.PathSeparator)) || strings.HasPrefix(url, ".") {
		return "", fmt.Errorf("filepath not supported")

	} else if noPrefixHttpRegex.MatchString(url) {
		return "http://" + url, nil

	} else if haveUriSchemeRegex.MatchString(url) {
		return url, nil
	}

	return url, fmt.Errorf("not a uri")
}

func ProcessAllURI(blocks []*model.Block) []*model.Block {
	for bI := range blocks {
		if blocks[bI].GetText() != nil && blocks[bI].GetText().Marks != nil && len(blocks[bI].GetText().Marks.Marks) > 0 {
			marks := blocks[bI].GetText().Marks.Marks

			for mI := range marks {
				if marks[mI].Type == model.BlockContentTextMark_Link {
					marks[mI].Param, _ = ProcessURI(marks[mI].Param)
				}
			}

			blocks[bI].GetText().Marks.Marks = marks
		}
	}

	return blocks
}
