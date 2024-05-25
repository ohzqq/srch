package srch

import "net/url"

func NewData(u *url.URL) {
	switch u.Scheme {
	case "file":
	case "http", "https":
	case "":
	}
}
