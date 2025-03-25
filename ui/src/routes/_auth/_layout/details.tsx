import { useId, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";

import userDetails from "@/icons/user_details.png";
import { queryKeys } from "@/repositories/queryKeys";
import { getMe } from "@/repositories/user/repositories";

export const Route = createFileRoute("/_auth/_layout/details")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData({
      queryKey: queryKeys.me,
      queryFn: getMe,
    }),
  component: RouteComponent,
});

const countryOptions = {
  AF: "Afghanistan",
  AL: "Albania",
  DZ: "Algeria",
  AS: "American Samoa",
  AD: "Andorra",
  AO: "Angola",
  AI: "Anguilla",
  AQ: "Antarctica",
  AG: "Antigua and Barbuda",
  AR: "Argentina",
  AM: "Armenia",
  AW: "Aruba",
  AU: "Australia",
  AT: "Austria",
  AZ: "Azerbaijan",
  BS: "Bahamas, The",
  BH: "Bahrain",
  BD: "Bangladesh",
  BB: "Barbados",
  BY: "Belarus",
  BE: "Belgium",
  BZ: "Belize",
  BJ: "Benin",
  BM: "Bermuda",
  BT: "Bhutan",
  BO: "Bolivia",
  BA: "Bosnia and Herzegovina",
  BW: "Botswana",
  BV: "Bouvet Island",
  BR: "Brazil",
  IO: "British Indian Ocean Territory",
  BN: "Brunei",
  BG: "Bulgaria",
  BF: "Burkina Faso",
  BI: "Burundi",
  KH: "Cambodia",
  CM: "Cameroon",
  CA: "Canada",
  CV: "Cape Verde",
  KY: "Cayman Islands",
  CF: "Central African Republic",
  TD: "Chad",
  CL: "Chile",
  CN: "China",
  CX: "Christmas Island",
  CC: "Cocos (Keeling) Islands",
  CO: "Colombia",
  KM: "Comoros",
  CG: "Congo",
  CK: "Cook Islands",
  CR: "Costa Rica",
  CI: "Cote d'Ivoire",
  HR: "Croatia",
  CU: "Cuba",
  CY: "Cyprus",
  CZ: "Czech Republic",
  DK: "Denmark",
  DJ: "Djibouti",
  DM: "Dominica",
  DO: "Dominican Republic",
  TP: "East Timor",
  EC: "Ecuador",
  EG: "Egypt",
  SV: "El Salvador",
  GQ: "Equatorial Guinea",
  ER: "Eritrea",
  EE: "Estonia",
  ET: "Ethiopia",
  FK: "Falkland Islands (Islas Malvinas)",
  FO: "Faroe Islands",
  FJ: "Fiji Islands",
  FI: "Finland",
  FR: "France",
  FX: "France, Metropolitan",
  GF: "French Guiana",
  PF: "French Polynesia",
  TF: "French Southern and Antarctic Lands",
  GA: "Gabon",
  GM: "Gambia, The",
  GE: "Georgia",
  DE: "Germany",
  GH: "Ghana",
  GI: "Gibraltar",
  GR: "Greece",
  GL: "Greenland",
  GD: "Grenada",
  GP: "Guadeloupe",
  GU: "Guam",
  GT: "Guatemala",
  GN: "Guinea",
  GW: "Guinea-Bissau",
  GY: "Guyana",
  HT: "Haiti",
  HM: "Heard Island and McDonald Islands",
  HN: "Honduras",
  HK: "Hong Kong S.A.R.",
  HU: "Hungary",
  IS: "Iceland",
  IN: "India",
  ID: "Indonesia",
  IR: "Iran",
  IQ: "Iraq",
  IE: "Ireland",
  IL: "Israel",
  IT: "Italy",
  JM: "Jamaica",
  JP: "Japan",
  JO: "Jordan",
  KZ: "Kazakhstan",
  KE: "Kenya",
  KI: "Kiribati",
  KP: "North Korea",
  KR: "Korea",
  KW: "Kuwait",
  KG: "Kyrgyzstan",
  LA: "Laos",
  LV: "Latvia",
  LB: "Lebanon",
  LS: "Lesotho",
  LR: "Liberia",
  LY: "Libya",
  LI: "Liechtenstein",
  LT: "Lithuania",
  LU: "Luxembourg",
  MO: "Macau",
  MK: "Macedonia, Former Yugoslav Republic Of",
  MG: "Madagascar",
  MW: "Malawi",
  MY: "Malaysia",
  MV: "Maldives",
  ML: "Mali",
  MT: "Malta",
  MH: "Marshall Islands",
  MQ: "Martinique",
  MR: "Mauritania",
  MU: "Mauritius",
  YT: "Mayotte",
  MX: "Mexico",
  FM: "Micronesia, Federated States Of",
  MD: "Moldova",
  MC: "Monaco",
  MN: "Mongolia",
  MS: "Montserrat",
  MA: "Morocco",
  MZ: "Mozambique",
  MM: "Myanmar",
  NA: "Namibia",
  NR: "Nauru",
  NP: "Nepal",
  NL: "Netherlands, The",
  AN: "Netherlands Antilles",
  NC: "New Caledonia",
  NZ: "New Zealand",
  NI: "Nicaragua",
  NE: "Niger",
  NG: "Nigeria",
  NU: "Niue",
  NF: "Norfolk Island",
  MP: "Northern Mariana Islands",
  NO: "Norway",
  OM: "Oman",
  PK: "Pakistan",
  PW: "Palau",
  PA: "Panama",
  PG: "Papua New Guinea",
  PY: "Paraguay",
  PE: "Peru",
  PH: "Philippines",
  PN: "Pitcairn Islands",
  PL: "Poland",
  PT: "Portugal",
  PR: "Puerto Rico",
  QA: "Qatar",
  RE: "Reunion",
  RO: "Romania",
  RU: "Russia",
  RW: "Rwanda",
  KN: "St. Kitts And Nevis",
  LC: "St. Lucia",
  VC: "St. Vincent And The Grenadines",
  WS: "Samoa",
  SM: "San Marino",
  ST: "Sao Tome And Principe",
  SA: "Saudi Arabia",
  SN: "Senegal",
  SC: "Seychelles",
  SL: "Sierra Leone",
  SG: "Singapore",
  SK: "Slovakia",
  SI: "Slovenia",
  SB: "Solomon Islands",
  SO: "Somalia",
  ZA: "South Africa",
  GS: "South Georgia and the South Sandwich Islands",
  ES: "Spain",
  LK: "Sri Lanka",
  SH: "Saint Helena",
  PM: "Saint Pierre and Miquelon",
  SD: "Sudan",
  SR: "Suriname",
  SJ: "Svalbard",
  SZ: "Swaziland",
  SE: "Sweden",
  CH: "Switzerland",
  SY: "Syria",
  TW: "Taiwan",
  TJ: "Tajikistan",
  TZ: "Tanzania",
  TH: "Thailand",
  TG: "Togo",
  TK: "Tokelau",
  TO: "Tonga",
  TT: "Trinidad And Tobago",
  TN: "Tunisia",
  TR: "Turkey",
  TM: "Turkmenistan",
  TC: "Turks and Caicos Islands",
  TV: "Tuvalu",
  UG: "Uganda",
  UA: "Ukraine",
  AE: "United Arab Emirates",
  GB: "United Kingdom",
  US: "United States",
  UM: "United States Minor Outlying Islands",
  UY: "Uruguay",
  UZ: "Uzbekistan",
  VU: "Vanuatu",
  VA: "Vatican City",
  VE: "Venezuela",
  VN: "Vietnam",
  VG: "Virgin Islands (British)",
  VI: "Virgin Islands (US)",
  WF: "Wallis and Futuna",
  EH: "Western Sahara",
  YE: "Yemen",
  YU: "Yugoslavia",
  ZR: "Congo (DRC)",
  ZM: "Zambia",
  ZW: "Zimbabwe",
};

const usStateOptions = {
  AL: "AL-Alabama",
  AK: "AK-Alaska",
  AZ: "AZ-Arizona",
  AR: "AR-Arkansas",
  CA: "CA-California",
  CO: "CO-Colorado",
  CT: "CT-Connecticut",
  DC: "DC-Washington D.C.",
  DE: "DE-Delaware",
  FL: "FL-Florida",
  GA: "GA-Georgia",
  HI: "HI-Hawaii",
  ID: "ID-Idaho",
  IL: "IL-Illinois",
  IN: "IN-Indiana",
  IA: "IA-Iowa",
  KS: "KS-Kansas",
  KY: "KY-Kentucky",
  LA: "LA-Louisiana",
  ME: "ME-Maine",
  MD: "MD-Maryland",
  MA: "MA-Massachusetts",
  MI: "MI-Michigan",
  MN: "MN-Minnesota",
  MS: "MS-Mississippi",
  MO: "MO-Missouri",
  MT: "MT-Montana",
  NE: "NE-Nebraska",
  NV: "NV-Nevada",
  NH: "NH-New Hampshire",
  NJ: "NJ-New Jersey",
  NM: "NM-New Mexico",
  NY: "NY-New York",
  NC: "NC-North Carolina",
  ND: "ND-North Dakota",
  OH: "OH-Ohio",
  OK: "OK-Oklahoma",
  OR: "OR-Oregon",
  PA: "PA-Pennsylvania",
  PR: "PR-Puerto Rico",
  RI: "RI-Rhode Island",
  SC: "SC-South Carolina",
  SD: "SD-South Dakota",
  TN: "TN-Tennessee",
  TX: "TX-Texas",
  UT: "UT-Utah",
  VT: "VT-Vermont",
  VA: "VA-Virginia",
  WA: "WA-Washington",
  WV: "WV-West Virginia",
  WI: "WI-Wisconsin",
  WY: "WY-Wyoming",
};

function RouteComponent() {
  const firstNameId = useId();
  const lastNameId = useId();
  const countryId = useId();
  const stateId = useId();
  const cityId = useId();

  const meQuery = useSuspenseQuery({
    queryKey: queryKeys.me,
    queryFn: getMe,
  });

  const [country, setCountry] = useState<string>(meQuery.data.country);

  return (
    <div className="flex flex-row gap-2.5 w-full">
      <div>
        <img src={userDetails} />
      </div>
      <form className="w-full">
        <fieldset>
          <div className="field-row-stacked">
            <label htmlFor={firstNameId}>First Name</label>
            <input
              type="text"
              id={firstNameId}
              defaultValue={meQuery.data.first_name}
            />
          </div>
          <div className="field-row-stacked">
            <label htmlFor={lastNameId}>Last Name</label>
            <input
              type="text"
              id={lastNameId}
              defaultValue={meQuery.data.last_name}
            />
          </div>
          <div className="field-row-stacked">
            <label htmlFor={countryId}>Country</label>
            <select
              id={countryId}
              defaultValue={meQuery.data.country}
              onChange={(e) => setCountry(e.target.value)}
            >
              {Object.entries(countryOptions).map(([code, name]) => (
                <option key={code} value={code}>
                  {name}
                </option>
              ))}
            </select>
          </div>
          {country === "US" && (
            <>
              <div className="field-row-stacked">
                <label htmlFor={stateId}>State</label>
                <select id={stateId} defaultValue={meQuery.data.state}>
                  {Object.entries(usStateOptions).map(([code, name]) => (
                    <option key={code} value={code}>
                      {name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="field-row-stacked">
                <label htmlFor={cityId}>City</label>
                <input
                  type="text"
                  id={cityId}
                  defaultValue={meQuery.data.city}
                />
              </div>
            </>
          )}
        </fieldset>
        <div className="flex justify-end mt-2.5">
          <button type="submit" className="cursor-pointer">
            OK
          </button>
        </div>
      </form>
    </div>
  );
}
