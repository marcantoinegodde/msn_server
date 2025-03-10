package web

import (
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

var validCountryCodes = map[string]bool{
	"AF": true, "AL": true, "DZ": true, "AS": true, "AD": true, "AO": true, "AI": true, "AQ": true,
	"AG": true, "AR": true, "AM": true, "AW": true, "AU": true, "AT": true, "AZ": true, "BS": true,
	"BH": true, "BD": true, "BB": true, "BY": true, "BE": true, "BZ": true, "BJ": true, "BM": true,
	"BT": true, "BO": true, "BA": true, "BW": true, "BV": true, "BR": true, "IO": true, "BN": true,
	"BG": true, "BF": true, "BI": true, "KH": true, "CM": true, "CA": true, "CV": true, "KY": true,
	"CF": true, "TD": true, "CL": true, "CN": true, "CX": true, "CC": true, "CO": true, "KM": true,
	"CG": true, "CK": true, "CR": true, "CI": true, "HR": true, "CU": true, "CY": true, "CZ": true,
	"DK": true, "DJ": true, "DM": true, "DO": true, "TP": true, "EC": true, "EG": true, "SV": true,
	"GQ": true, "ER": true, "EE": true, "ET": true, "FK": true, "FO": true, "FJ": true, "FI": true,
	"FR": true, "FX": true, "GF": true, "PF": true, "TF": true, "GA": true, "GM": true, "GE": true,
	"DE": true, "GH": true, "GI": true, "GR": true, "GL": true, "GD": true, "GP": true, "GU": true,
	"GT": true, "GN": true, "GW": true, "GY": true, "HT": true, "HM": true, "HN": true, "HK": true,
	"HU": true, "IS": true, "IN": true, "ID": true, "IR": true, "IQ": true, "IE": true, "IL": true,
	"IT": true, "JM": true, "JP": true, "JO": true, "KZ": true, "KE": true, "KI": true, "KP": true,
	"KR": true, "KW": true, "KG": true, "LA": true, "LV": true, "LB": true, "LS": true, "LR": true,
	"LY": true, "LI": true, "LT": true, "LU": true, "MO": true, "MK": true, "MG": true, "MW": true,
	"MY": true, "MV": true, "ML": true, "MT": true, "MH": true, "MQ": true, "MR": true, "MU": true,
	"YT": true, "MX": true, "FM": true, "MD": true, "MC": true, "MN": true, "MS": true, "MA": true,
	"MZ": true, "MM": true, "NA": true, "NR": true, "NP": true, "NL": true, "AN": true, "NC": true,
	"NZ": true, "NI": true, "NE": true, "NG": true, "NU": true, "NF": true, "MP": true, "NO": true,
	"OM": true, "PK": true, "PW": true, "PA": true, "PG": true, "PY": true, "PE": true, "PH": true,
	"PN": true, "PL": true, "PT": true, "PR": true, "QA": true, "RE": true, "RO": true, "RU": true,
	"RW": true, "KN": true, "LC": true, "VC": true, "WS": true, "SM": true, "ST": true, "SA": true,
	"SN": true, "SC": true, "SL": true, "SG": true, "SK": true, "SI": true, "SB": true, "SO": true,
	"ZA": true, "GS": true, "ES": true, "LK": true, "SH": true, "PM": true, "SD": true, "SR": true,
	"SJ": true, "SZ": true, "SE": true, "CH": true, "SY": true, "TW": true, "TJ": true, "TZ": true,
	"TH": true, "TG": true, "TK": true, "TO": true, "TT": true, "TN": true, "TR": true, "TM": true,
	"TC": true, "TV": true, "UG": true, "UA": true, "AE": true, "GB": true, "US": true, "UM": true,
	"UY": true, "UZ": true, "VU": true, "VA": true, "VE": true, "VN": true, "VG": true, "VI": true,
	"WF": true, "EH": true, "YE": true, "YU": true, "ZR": true, "ZM": true, "ZW": true,
}

var validUSStates = map[string]bool{
	"AL": true, "AK": true, "AZ": true, "AR": true, "CA": true, "CO": true, "CT": true, "DC": true,
	"DE": true, "FL": true, "GA": true, "HI": true, "ID": true, "IL": true, "IN": true, "IA": true,
	"KS": true, "KY": true, "LA": true, "ME": true, "MD": true, "MA": true, "MI": true, "MN": true,
	"MS": true, "MO": true, "MT": true, "NE": true, "NV": true, "NH": true, "NJ": true, "NM": true,
	"NY": true, "NC": true, "ND": true, "OH": true, "OK": true, "OR": true, "PA": true, "PR": true,
	"RI": true, "SC": true, "SD": true, "TN": true, "TX": true, "UT": true, "VT": true, "VA": true,
	"WA": true, "WV": true, "WI": true, "WY": true,
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func validateName(fl validator.FieldLevel) bool {
	nameRegex := regexp.MustCompile(`^[a-zA-Z' -]+$`)
	return nameRegex.MatchString(fl.Field().String())
}

func validateCountry(fl validator.FieldLevel) bool {
	_, ok := validCountryCodes[fl.Field().String()]
	return ok
}

func validateUSState(fl validator.FieldLevel) bool {
	_, ok := validUSStates[fl.Field().String()]
	return ok
}
