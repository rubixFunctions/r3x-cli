package cmd

const (
	LicenseMIT      = "MIT"
	LicenseApache20 = "Apache-2.0"
)

var Licenses = make(map[string]License)

var KnownLicenses = []string{
	LicenseMIT,
	LicenseApache20,
}

type License struct {
	Name            string
	PossibleMatches []string
	Text            string
}

func init() {

}

// func getLicense() License {

// }

func findLicense(name string) bool {
	// save knownlicenses to map
	l := make(map[string]bool)
	for i := 0; i < len(KnownLicenses); i++ {
		l[KnownLicenses[i]] = true
	}
	// check is license exists
	if _, ok := l[name]; ok {
		return true
	}
	return false
}

// func matchLicense(name string) string {

// }
