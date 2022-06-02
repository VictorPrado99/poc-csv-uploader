package csvuploaderlogic

import "regexp"

var countriesRegex map[string]regexp.Regexp

// init the dictionary of regex when the package is initialize. That way, the regex won't slow down the system, because will be compiled just once.
// This dicstionary will be used to determine the country based on phone number
// To create new countries, is necessary just add a new reference and the regex
func init() {
	countriesRegex = make(map[string]regexp.Regexp)

	countriesRegex["Cameroon"] = *regexp.MustCompile(`(?m)\(237\)\ ?[2368]\d{7,8}$`)
	countriesRegex["Ethiopia"] = *regexp.MustCompile(`(?m)\(251\)\ ?[1-59]\d{8}$`)
	countriesRegex["Morocco"] = *regexp.MustCompile(`(?m)\(212\)\ ?[5-9]\d{8}$`)
	countriesRegex["Mozambique"] = *regexp.MustCompile(`(?m)\(258\)\ ?[28]\d{7,8}$`)
	countriesRegex["Uganda"] = *regexp.MustCompile(`(?m)\(256\)\ ?\d{9}$`)

}
