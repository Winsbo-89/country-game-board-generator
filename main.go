package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
)

type Country struct {
	Name    string
	Code    string
	AreaKm2 float64
}

type Rule struct {
	Score int
	Dice  int
}

type PageData struct {
	Countries    []Country
	Difficulties []string
	Tiles        []int
}

type FormData struct {
	Country     Country
	Difficulty  string
	Territories int
	Reward      int
}

var (
	countries = []Country{
		{Name: "Afghanistan", Code: "AF", AreaKm2: 652230},
		{Name: "Albania", Code: "AL", AreaKm2: 28748},
		{Name: "Algeria", Code: "DZ", AreaKm2: 2381741},
		{Name: "Andorra", Code: "AD", AreaKm2: 468},
		{Name: "Angola", Code: "AO", AreaKm2: 1246700},
		{Name: "Antigua and Barbuda", Code: "AG", AreaKm2: 442},
		{Name: "Argentina", Code: "AR", AreaKm2: 2780400},
		{Name: "Armenia", Code: "AM", AreaKm2: 29743},
		{Name: "Australia", Code: "AU", AreaKm2: 7692024},
		{Name: "Austria", Code: "AT", AreaKm2: 83871},
		{Name: "Azerbaijan", Code: "AZ", AreaKm2: 86600},
		{Name: "Bahamas", Code: "BS", AreaKm2: 13943},
		{Name: "Bahrain", Code: "BH", AreaKm2: 765},
		{Name: "Bangladesh", Code: "BD", AreaKm2: 148460},
		{Name: "Barbados", Code: "BB", AreaKm2: 430},
		{Name: "Belarus", Code: "BY", AreaKm2: 207600},
		{Name: "Belgium", Code: "BE", AreaKm2: 30528},
		{Name: "Belize", Code: "BZ", AreaKm2: 22966},
		{Name: "Benin", Code: "BJ", AreaKm2: 112622},
		{Name: "Bhutan", Code: "BT", AreaKm2: 38394},
		{Name: "Bolivia", Code: "BO", AreaKm2: 1098581},
		{Name: "Bosnia and Herzegovina", Code: "BA", AreaKm2: 51209},
		{Name: "Botswana", Code: "BW", AreaKm2: 581730},
		{Name: "Brazil", Code: "BR", AreaKm2: 8515767},
		{Name: "Brunei", Code: "BN", AreaKm2: 5765},
		{Name: "Bulgaria", Code: "BG", AreaKm2: 110879},
		{Name: "Burkina Faso", Code: "BF", AreaKm2: 274200},
		{Name: "Burundi", Code: "BI", AreaKm2: 27834},
		{Name: "Cabo Verde", Code: "CV", AreaKm2: 4033},
		{Name: "Cambodia", Code: "KH", AreaKm2: 181035},
		{Name: "Cameroon", Code: "CM", AreaKm2: 475442},
		{Name: "Canada", Code: "CA", AreaKm2: 9984670},
		{Name: "Central African Republic", Code: "CF", AreaKm2: 622984},
		{Name: "Chad", Code: "TD", AreaKm2: 1284000},
		{Name: "Chile", Code: "CL", AreaKm2: 756102},
		{Name: "China", Code: "CN", AreaKm2: 9706961},
		{Name: "Colombia", Code: "CO", AreaKm2: 1141748},
		{Name: "Comoros", Code: "KM", AreaKm2: 1862},
		{Name: "Congo, Democratic Republic of the", Code: "CD", AreaKm2: 2344858},
		{Name: "Congo, Republic of the", Code: "CG", AreaKm2: 342000},
		{Name: "Costa Rica", Code: "CR", AreaKm2: 51100},
		{Name: "Côte d'Ivoire", Code: "CI", AreaKm2: 322463},
		{Name: "Croatia", Code: "HR", AreaKm2: 56594},
		{Name: "Cuba", Code: "CU", AreaKm2: 109884},
		{Name: "Cyprus", Code: "CY", AreaKm2: 9251},
		{Name: "Czechia", Code: "CZ", AreaKm2: 78865},
		{Name: "Denmark", Code: "DK", AreaKm2: 43094},
		{Name: "Djibouti", Code: "DJ", AreaKm2: 23200},
		{Name: "Dominica", Code: "DM", AreaKm2: 751},
		{Name: "Dominican Republic", Code: "DO", AreaKm2: 48671},
		{Name: "Ecuador", Code: "EC", AreaKm2: 276841},
		{Name: "Egypt", Code: "EG", AreaKm2: 1002450},
		{Name: "El Salvador", Code: "SV", AreaKm2: 21041},
		{Name: "Equatorial Guinea", Code: "GQ", AreaKm2: 28051},
		{Name: "Eritrea", Code: "ER", AreaKm2: 117600},
		{Name: "Estonia", Code: "EE", AreaKm2: 45227},
		{Name: "Eswatini", Code: "SZ", AreaKm2: 17364},
		{Name: "Ethiopia", Code: "ET", AreaKm2: 1104300},
		{Name: "Fiji", Code: "FJ", AreaKm2: 18272},
		{Name: "Finland", Code: "FI", AreaKm2: 338455},
		{Name: "France", Code: "FR", AreaKm2: 551695},
		{Name: "Gabon", Code: "GA", AreaKm2: 267668},
		{Name: "Gambia", Code: "GM", AreaKm2: 11295},
		{Name: "Georgia", Code: "GE", AreaKm2: 69700},
		{Name: "Germany", Code: "DE", AreaKm2: 357022},
		{Name: "Ghana", Code: "GH", AreaKm2: 238533},
		{Name: "Greece", Code: "GR", AreaKm2: 131957},
		{Name: "Grenada", Code: "GD", AreaKm2: 344},
		{Name: "Guatemala", Code: "GT", AreaKm2: 108889},
		{Name: "Guinea", Code: "GN", AreaKm2: 245857},
		{Name: "Guinea-Bissau", Code: "GW", AreaKm2: 36125},
		{Name: "Guyana", Code: "GY", AreaKm2: 214969},
		{Name: "Haiti", Code: "HT", AreaKm2: 27750},
		{Name: "Honduras", Code: "HN", AreaKm2: 112492},
		{Name: "Hungary", Code: "HU", AreaKm2: 93028},
		{Name: "Iceland", Code: "IS", AreaKm2: 103000},
		{Name: "India", Code: "IN", AreaKm2: 3287263},
		{Name: "Indonesia", Code: "ID", AreaKm2: 1904569},
		{Name: "Iran", Code: "IR", AreaKm2: 1648195},
		{Name: "Iraq", Code: "IQ", AreaKm2: 438317},
		{Name: "Ireland", Code: "IE", AreaKm2: 70273},
		{Name: "Israel", Code: "IL", AreaKm2: 22072},
		{Name: "Italy", Code: "IT", AreaKm2: 301340},
		{Name: "Jamaica", Code: "JM", AreaKm2: 10991},
		{Name: "Japan", Code: "JP", AreaKm2: 377975},
		{Name: "Jordan", Code: "JO", AreaKm2: 89342},
		{Name: "Kazakhstan", Code: "KZ", AreaKm2: 2724900},
		{Name: "Kenya", Code: "KE", AreaKm2: 580367},
		{Name: "Kiribati", Code: "KI", AreaKm2: 811},
		{Name: "Kuwait", Code: "KW", AreaKm2: 17818},
		{Name: "Kyrgyzstan", Code: "KG", AreaKm2: 199951},
		{Name: "Laos", Code: "LA", AreaKm2: 236800},
		{Name: "Latvia", Code: "LV", AreaKm2: 64559},
		{Name: "Lebanon", Code: "LB", AreaKm2: 10452},
		{Name: "Lesotho", Code: "LS", AreaKm2: 30355},
		{Name: "Liberia", Code: "LR", AreaKm2: 111369},
		{Name: "Libya", Code: "LY", AreaKm2: 1759540},
		{Name: "Liechtenstein", Code: "LI", AreaKm2: 160},
		{Name: "Lithuania", Code: "LT", AreaKm2: 65300},
		{Name: "Luxembourg", Code: "LU", AreaKm2: 2586},
		{Name: "Madagascar", Code: "MG", AreaKm2: 587041},
		{Name: "Malawi", Code: "MW", AreaKm2: 118484},
		{Name: "Malaysia", Code: "MY", AreaKm2: 330803},
		{Name: "Maldives", Code: "MV", AreaKm2: 300},
		{Name: "Mali", Code: "ML", AreaKm2: 1240192},
		{Name: "Malta", Code: "MT", AreaKm2: 316},
		{Name: "Marshall Islands", Code: "MH", AreaKm2: 181},
		{Name: "Mauritania", Code: "MR", AreaKm2: 1030700},
		{Name: "Mauritius", Code: "MU", AreaKm2: 2040},
		{Name: "Mexico", Code: "MX", AreaKm2: 1964375},
		{Name: "Micronesia", Code: "FM", AreaKm2: 702},
		{Name: "Moldova", Code: "MD", AreaKm2: 33846},
		{Name: "Monaco", Code: "MC", AreaKm2: 2.02},
		{Name: "Mongolia", Code: "MN", AreaKm2: 1564110},
		{Name: "Montenegro", Code: "ME", AreaKm2: 13812},
		{Name: "Morocco", Code: "MA", AreaKm2: 446550},
		{Name: "Mozambique", Code: "MZ", AreaKm2: 801590},
		{Name: "Myanmar", Code: "MM", AreaKm2: 676578},
		{Name: "Namibia", Code: "NA", AreaKm2: 825615},
		{Name: "Nauru", Code: "NR", AreaKm2: 21},
		{Name: "Nepal", Code: "NP", AreaKm2: 147516},
		{Name: "Netherlands", Code: "NL", AreaKm2: 41850},
		{Name: "New Zealand", Code: "NZ", AreaKm2: 268021},
		{Name: "Nicaragua", Code: "NI", AreaKm2: 130373},
		{Name: "Niger", Code: "NE", AreaKm2: 1267000},
		{Name: "Nigeria", Code: "NG", AreaKm2: 923768},
		{Name: "North Korea", Code: "KP", AreaKm2: 120538},
		{Name: "North Macedonia", Code: "MK", AreaKm2: 25713},
		{Name: "Norway", Code: "NO", AreaKm2: 385207},
		{Name: "Oman", Code: "OM", AreaKm2: 309500},
		{Name: "Pakistan", Code: "PK", AreaKm2: 881912},
		{Name: "Palau", Code: "PW", AreaKm2: 459},
		{Name: "Palestine", Code: "PS", AreaKm2: 6020},
		{Name: "Panama", Code: "PA", AreaKm2: 75420},
		{Name: "Papua New Guinea", Code: "PG", AreaKm2: 462840},
		{Name: "Paraguay", Code: "PY", AreaKm2: 406752},
		{Name: "Peru", Code: "PE", AreaKm2: 1285216},
		{Name: "Philippines", Code: "PH", AreaKm2: 300000},
		{Name: "Poland", Code: "PL", AreaKm2: 312696},
		{Name: "Portugal", Code: "PT", AreaKm2: 92212},
		{Name: "Qatar", Code: "QA", AreaKm2: 11586},
		{Name: "Romania", Code: "RO", AreaKm2: 238391},
		{Name: "Russia", Code: "RU", AreaKm2: 17098242},
		{Name: "Rwanda", Code: "RW", AreaKm2: 26338},
		{Name: "Saint Kitts and Nevis", Code: "KN", AreaKm2: 261},
		{Name: "Saint Lucia", Code: "LC", AreaKm2: 616},
		{Name: "Saint Vincent and the Grenadines", Code: "VC", AreaKm2: 389},
		{Name: "Samoa", Code: "WS", AreaKm2: 2842},
		{Name: "San Marino", Code: "SM", AreaKm2: 61},
		{Name: "Sao Tome and Principe", Code: "ST", AreaKm2: 964},
		{Name: "Saudi Arabia", Code: "SA", AreaKm2: 2149690},
		{Name: "Senegal", Code: "SN", AreaKm2: 196722},
		{Name: "Serbia", Code: "RS", AreaKm2: 77474},
		{Name: "Seychelles", Code: "SC", AreaKm2: 459},
		{Name: "Sierra Leone", Code: "SL", AreaKm2: 71740},
		{Name: "Singapore", Code: "SG", AreaKm2: 728},
		{Name: "Slovakia", Code: "SK", AreaKm2: 49035},
		{Name: "Slovenia", Code: "SI", AreaKm2: 20273},
		{Name: "Solomon Islands", Code: "SB", AreaKm2: 28896},
		{Name: "Somalia", Code: "SO", AreaKm2: 637657},
		{Name: "South Africa", Code: "ZA", AreaKm2: 1221037},
		{Name: "South Korea", Code: "KR", AreaKm2: 100210},
		{Name: "South Sudan", Code: "SS", AreaKm2: 644329},
		{Name: "Spain", Code: "ES", AreaKm2: 505990},
		{Name: "Sri Lanka", Code: "LK", AreaKm2: 65610},
		{Name: "Sudan", Code: "SD", AreaKm2: 1861484},
		{Name: "Suriname", Code: "SR", AreaKm2: 163820},
		{Name: "Sweden", Code: "SE", AreaKm2: 450295},
		{Name: "Switzerland", Code: "CH", AreaKm2: 41284},
		{Name: "Syria", Code: "SY", AreaKm2: 185180},
		{Name: "Tajikistan", Code: "TJ", AreaKm2: 143100},
		{Name: "Tanzania", Code: "TZ", AreaKm2: 945087},
		{Name: "Thailand", Code: "TH", AreaKm2: 513120},
		{Name: "Timor-Leste", Code: "TL", AreaKm2: 14919},
		{Name: "Togo", Code: "TG", AreaKm2: 56785},
		{Name: "Tonga", Code: "TO", AreaKm2: 747},
		{Name: "Trinidad and Tobago", Code: "TT", AreaKm2: 5131},
		{Name: "Tunisia", Code: "TN", AreaKm2: 163610},
		{Name: "Turkey", Code: "TR", AreaKm2: 783562},
		{Name: "Turkmenistan", Code: "TM", AreaKm2: 488100},
		{Name: "Tuvalu", Code: "TV", AreaKm2: 26},
		{Name: "Uganda", Code: "UG", AreaKm2: 241038},
		{Name: "Ukraine", Code: "UA", AreaKm2: 603500},
		{Name: "United Arab Emirates", Code: "AE", AreaKm2: 83600},
		{Name: "United Kingdom", Code: "GB", AreaKm2: 243610},
		{Name: "United States", Code: "US", AreaKm2: 9833517},
		{Name: "Uruguay", Code: "UY", AreaKm2: 176215},
		{Name: "Uzbekistan", Code: "UZ", AreaKm2: 447400},
		{Name: "Vanuatu", Code: "VU", AreaKm2: 12189},
		{Name: "Vatican City", Code: "VA", AreaKm2: 0.49},
		{Name: "Venezuela", Code: "VE", AreaKm2: 916445},
		{Name: "Vietnam", Code: "VN", AreaKm2: 331212},
		{Name: "Yemen", Code: "YE", AreaKm2: 527968},
		{Name: "Zambia", Code: "ZM", AreaKm2: 752618},
		{Name: "Zimbabwe", Code: "ZW", AreaKm2: 390757},
	}

	difficulties = []string{"Beginner", "Easy", "Medium", "Hard", "Insane"}

	difficultyRules = map[string][]Rule{
		"Beginner": {
			{Score: 3, Dice: 3}, {Score: 4, Dice: 3}, {Score: 5, Dice: 3}, {Score: 4, Dice: 4}, {Score: 5, Dice: 4}, {Score: 6, Dice: 4}, {Score: 7, Dice: 4}, {Score: 5, Dice: 5}, {Score: 6, Dice: 5}, {Score: 7, Dice: 5}, {Score: 8, Dice: 5}, {Score: 9, Dice: 5}, {Score: 6, Dice: 6}, {Score: 7, Dice: 6}, {Score: 8, Dice: 6}, {Score: 9, Dice: 6}, {Score: 10, Dice: 6}, {Score: 11, Dice: 6}, {Score: 7, Dice: 7}, {Score: 8, Dice: 7}, {Score: 9, Dice: 7}, {Score: 10, Dice: 7}, {Score: 11, Dice: 7}, {Score: 12, Dice: 7}, {Score: 13, Dice: 7}, {Score: 8, Dice: 8}, {Score: 9, Dice: 8}, {Score: 10, Dice: 8}, {Score: 11, Dice: 8}, {Score: 12, Dice: 8}, {Score: 13, Dice: 8}, {Score: 14, Dice: 8}, {Score: 15, Dice: 8}, {Score: 9, Dice: 9}, {Score: 10, Dice: 9}, {Score: 11, Dice: 9}, {Score: 12, Dice: 9}, {Score: 13, Dice: 9}, {Score: 14, Dice: 9}, {Score: 15, Dice: 9}, {Score: 16, Dice: 9}, {Score: 17, Dice: 9}, {Score: 10, Dice: 10}, {Score: 11, Dice: 10}, {Score: 12, Dice: 10},
		},
		"Easy": {
			{Score: 3, Dice: 3}, {Score: 4, Dice: 3}, {Score: 5, Dice: 3}, {Score: 6, Dice: 3}, {Score: 7, Dice: 3}, {Score: 8, Dice: 3}, {Score: 4, Dice: 4}, {Score: 5, Dice: 4}, {Score: 6, Dice: 4}, {Score: 7, Dice: 4}, {Score: 8, Dice: 4}, {Score: 9, Dice: 4}, {Score: 10, Dice: 4}, {Score: 11, Dice: 4}, {Score: 5, Dice: 5}, {Score: 6, Dice: 5}, {Score: 7, Dice: 5}, {Score: 8, Dice: 5}, {Score: 9, Dice: 5}, {Score: 10, Dice: 5}, {Score: 11, Dice: 5}, {Score: 12, Dice: 5}, {Score: 13, Dice: 5}, {Score: 14, Dice: 5}, {Score: 6, Dice: 6}, {Score: 7, Dice: 6}, {Score: 8, Dice: 6}, {Score: 9, Dice: 6}, {Score: 10, Dice: 6}, {Score: 11, Dice: 6}, {Score: 12, Dice: 6}, {Score: 13, Dice: 6}, {Score: 14, Dice: 6}, {Score: 15, Dice: 6}, {Score: 16, Dice: 6}, {Score: 17, Dice: 6}, {Score: 7, Dice: 7}, {Score: 8, Dice: 7}, {Score: 9, Dice: 7}, {Score: 10, Dice: 7}, {Score: 11, Dice: 7}, {Score: 12, Dice: 7}, {Score: 13, Dice: 7}, {Score: 14, Dice: 7}, {Score: 15, Dice: 7},
		},
		"Medium": {
			{Score: 6, Dice: 3}, {Score: 7, Dice: 3}, {Score: 8, Dice: 3}, {Score: 9, Dice: 3}, {Score: 10, Dice: 3}, {Score: 11, Dice: 3}, {Score: 8, Dice: 4}, {Score: 9, Dice: 4}, {Score: 10, Dice: 4}, {Score: 11, Dice: 4}, {Score: 12, Dice: 4}, {Score: 13, Dice: 4}, {Score: 14, Dice: 4}, {Score: 15, Dice: 4}, {Score: 10, Dice: 5}, {Score: 11, Dice: 5}, {Score: 12, Dice: 5}, {Score: 13, Dice: 5}, {Score: 14, Dice: 5}, {Score: 15, Dice: 5}, {Score: 16, Dice: 5}, {Score: 17, Dice: 5}, {Score: 18, Dice: 5}, {Score: 19, Dice: 5}, {Score: 12, Dice: 6}, {Score: 13, Dice: 6}, {Score: 14, Dice: 6}, {Score: 15, Dice: 6}, {Score: 16, Dice: 6}, {Score: 17, Dice: 6}, {Score: 18, Dice: 6}, {Score: 19, Dice: 6}, {Score: 20, Dice: 6}, {Score: 21, Dice: 6}, {Score: 22, Dice: 6}, {Score: 23, Dice: 6}, {Score: 14, Dice: 7}, {Score: 15, Dice: 7}, {Score: 16, Dice: 7}, {Score: 17, Dice: 7}, {Score: 18, Dice: 7}, {Score: 19, Dice: 7}, {Score: 20, Dice: 7}, {Score: 21, Dice: 7}, {Score: 22, Dice: 7},
		},
		"Hard": {
			{Score: 9, Dice: 3}, {Score: 10, Dice: 3}, {Score: 11, Dice: 3}, {Score: 12, Dice: 3}, {Score: 13, Dice: 3}, {Score: 14, Dice: 3}, {Score: 12, Dice: 4}, {Score: 13, Dice: 4}, {Score: 14, Dice: 4}, {Score: 15, Dice: 4}, {Score: 16, Dice: 4}, {Score: 17, Dice: 4}, {Score: 18, Dice: 4}, {Score: 19, Dice: 4}, {Score: 15, Dice: 5}, {Score: 16, Dice: 5}, {Score: 17, Dice: 5}, {Score: 18, Dice: 5}, {Score: 19, Dice: 5}, {Score: 20, Dice: 5}, {Score: 21, Dice: 5}, {Score: 22, Dice: 5}, {Score: 23, Dice: 5}, {Score: 24, Dice: 5}, {Score: 18, Dice: 6}, {Score: 19, Dice: 6}, {Score: 20, Dice: 6}, {Score: 21, Dice: 6}, {Score: 22, Dice: 6}, {Score: 23, Dice: 6}, {Score: 24, Dice: 6}, {Score: 25, Dice: 6}, {Score: 26, Dice: 6}, {Score: 27, Dice: 6}, {Score: 28, Dice: 6}, {Score: 29, Dice: 6}, {Score: 21, Dice: 7}, {Score: 22, Dice: 7}, {Score: 23, Dice: 7}, {Score: 24, Dice: 7}, {Score: 25, Dice: 7}, {Score: 26, Dice: 7}, {Score: 27, Dice: 7}, {Score: 28, Dice: 7}, {Score: 29, Dice: 7},
		},
		"Insane": {
			{Score: 12, Dice: 3}, {Score: 13, Dice: 3}, {Score: 14, Dice: 3}, {Score: 15, Dice: 3}, {Score: 16, Dice: 3}, {Score: 17, Dice: 3}, {Score: 18, Dice: 3}, {Score: 16, Dice: 4}, {Score: 17, Dice: 4}, {Score: 18, Dice: 4}, {Score: 19, Dice: 4}, {Score: 20, Dice: 4}, {Score: 21, Dice: 4}, {Score: 22, Dice: 4}, {Score: 23, Dice: 4}, {Score: 24, Dice: 4}, {Score: 20, Dice: 5}, {Score: 21, Dice: 5}, {Score: 22, Dice: 5}, {Score: 23, Dice: 5}, {Score: 24, Dice: 5}, {Score: 25, Dice: 5}, {Score: 26, Dice: 5}, {Score: 27, Dice: 5}, {Score: 28, Dice: 5}, {Score: 29, Dice: 5}, {Score: 30, Dice: 5}, {Score: 24, Dice: 6}, {Score: 25, Dice: 6}, {Score: 26, Dice: 6}, {Score: 27, Dice: 6}, {Score: 28, Dice: 6}, {Score: 29, Dice: 6}, {Score: 30, Dice: 6}, {Score: 31, Dice: 6}, {Score: 32, Dice: 6}, {Score: 33, Dice: 6}, {Score: 34, Dice: 6}, {Score: 35, Dice: 6}, {Score: 36, Dice: 6}, {Score: 28, Dice: 7}, {Score: 29, Dice: 7}, {Score: 30, Dice: 7}, {Score: 31, Dice: 7}, {Score: 32, Dice: 7},
		},
	}

	territoryPayouts = map[string][]int{
		"Beginner": {5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
		"Easy":     {5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15},
		"Medium":   {5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20},
		"Hard":     {10, 10, 10, 10, 10, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
		"Insane":   {15, 15, 15, 15, 15, 20, 20, 20, 20, 20, 25, 25, 25, 25, 25, 30, 30, 30, 30, 30, 35, 35, 35, 35, 35, 40, 40, 40, 40, 40, 45, 45, 45, 45, 45, 50, 50, 50, 50, 50, 55, 55, 55, 55, 55},
	}

	// Country rewards from the supplied template, for 5/10/15/.../45 territories.
	countryRewards = map[string][]int{
		"Beginner": {5, 10, 15, 20, 25, 30, 35, 40, 45},
		"Easy":     {20, 30, 40, 50, 100, 150, 200, 250, 300},
		"Medium":   {50, 100, 150, 200, 250, 300, 350, 400, 500},
		"Hard":     {100, 200, 250, 300, 400, 500, 600, 700, 800},
		"Insane":   {200, 300, 400, 500, 600, 700, 800, 900, 1000},
	}

	indexTemplate = template.Must(template.New("index").Parse(indexHTML))
)

func main() {
	sort.Slice(countries, func(i, j int) bool {
		return strings.ToLower(countries[i].Name) < strings.ToLower(countries[j].Name)
	})

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/generate", handleGenerate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Country board generator listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := PageData{
		Countries:    countries,
		Difficulties: difficulties,
		Tiles:        makeRange(6, 45),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTemplate.Execute(w, data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not read form.", http.StatusBadRequest)
		return
	}

	countryCode := strings.TrimSpace(r.FormValue("country"))
	difficulty := strings.TrimSpace(r.FormValue("difficulty"))
	territories, err := strconv.Atoi(r.FormValue("territories"))
	if err != nil || territories < 6 || territories > 45 {
		http.Error(w, "Territories must be between 6 and 45.", http.StatusBadRequest)
		return
	}

	country, ok := findCountry(countryCode)
	if !ok {
		http.Error(w, "Unknown country.", http.StatusBadRequest)
		return
	}

	if _, ok := difficultyRules[difficulty]; !ok {
		http.Error(w, "Unknown difficulty.", http.StatusBadRequest)
		return
	}

	pdfBytes, err := buildPDF(FormData{
		Country:     country,
		Difficulty:  difficulty,
		Territories: territories,
		Reward:      countryReward(difficulty, territories),
	})
	if err != nil {
		log.Printf("PDF generation error: %v", err)
		http.Error(w, "Could not generate the PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := safeFilename(country.Name) + "-game-board.pdf"

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	_, _ = w.Write(pdfBytes)
}

func findCountry(code string) (Country, bool) {
	for _, c := range countries {
		if strings.EqualFold(c.Code, code) {
			return c, true
		}
	}
	return Country{}, false
}

func makeRange(start, end int) []int {
	result := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func countryReward(difficulty string, territories int) int {
	// The supplied country-payout template only gives values at 5, 10, ... 45.
	// For the new 6-45 selector, this first version uses the next 5-territory
	// reward bracket. This is deliberately isolated here so it can be changed
	// easily after we test the game.
	bracket := int(math.Ceil(float64(territories)/5.0)) * 5
	if bracket < 5 {
		bracket = 5
	}
	if bracket > 45 {
		bracket = 45
	}

	values := countryRewards[difficulty]
	return values[(bracket/5)-1]
}

func safeFilename(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	result := strings.Trim(b.String(), "-")
	if result == "" {
		return "country-game-board"
	}
	return result
}

func formatArea(area float64) string {
	n := int64(math.Round(area))
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}

	var out []string
	for len(s) > 3 {
		out = append([]string{s[len(s)-3:]}, out...)
		s = s[:len(s)-3]
	}
	out = append([]string{s}, out...)
	return strings.Join(out, ",")
}

func buildPDF(data FormData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Legal", "")
	pdf.SetMargins(10, 8, 10)
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetTitle(data.Country.Name+" Game Board", true)
	pdf.SetAuthor("Country Game Board Generator", true)
	pdf.SetSubject("Printable country game board", true)

	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 10)

	pageW, pageH := pdf.GetPageSize()

	// Header.
	y := 9.0
	flagOK := addFlag(pdf, data.Country.Code, 10, y, 30, 20)

	titleX := 45.0
	if !flagOK {
		titleX = 10.0
	}

	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(20, 20, 20)
	pdf.Text(titleX, 19, data.Country.Name)

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.Text(titleX, 25, fmt.Sprintf("%s difficulty - %d territories", data.Difficulty, data.Territories))

	// Board.
	boardX := 10.0
	boardY := 32.0
	boardW := pageW - 20.0
	boardH := 224.0

	cols := chooseColumns(data.Territories)
	rows := int(math.Ceil(float64(data.Territories) / float64(cols)))
	cellW := boardW / float64(cols)
	cellH := boardH / float64(rows)

	rules := difficultyRules[data.Difficulty]
	payouts := territoryPayouts[data.Difficulty]

	for i := 0; i < data.Territories; i++ {
		row := i / cols
		col := i % cols

		x := boardX + float64(col)*cellW
		cellY := boardY + float64(row)*cellH

		pdf.SetDrawColor(30, 30, 30)
		pdf.SetLineWidth(0.5)
		pdf.Rect(x, cellY, cellW, cellH, "D")

		// Territory number: top-left.
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetTextColor(30, 30, 30)
		pdf.Text(x+2.0, cellY+5.0, strconv.Itoa(i+1))

		// Payout: top-right.
		payout := fmt.Sprintf("$%d", payouts[i])
		pdf.SetFont("Helvetica", "B", 10)
		pdf.Text(x+cellW-2.0-pdf.GetStringWidth(payout), cellY+5.0, payout)

		// Territory name in the center.
		name := fmt.Sprintf("Territory %d", i+1)
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetTextColor(60, 60, 60)
		centerText(pdf, name, x+cellW/2, cellY+cellH/2+2.0)

		// Score: bottom-left.
		rule := rules[i]
		scoreText := fmt.Sprintf("Roll: %d", rule.Score)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(30, 30, 30)
		pdf.Text(x+2.0, cellY+cellH-3.0, scoreText)

		// Dice: bottom-right.
		diceText := fmt.Sprintf("Dice: %d", rule.Dice)
		pdf.Text(x+cellW-2.0-pdf.GetStringWidth(diceText), cellY+cellH-3.0, diceText)
	}

	// Information row under board.
	infoY := boardY + boardH + 8.0

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(20, 20, 20)
	pdf.Text(10, infoY, "Country Information")

	pdf.SetFont("Helvetica", "", 10)
	pdf.Text(10, infoY+6, "Area: "+formatArea(data.Country.AreaKm2)+" km2")
	pdf.Text(10, infoY+12, "Difficulty: "+data.Difficulty)
	pdf.Text(10, infoY+18, "Conquest reward: $"+strconv.Itoa(data.Reward))

	// Legend at lower left.
	legendX := 10.0
	legendY := infoY + 28.0
	pdf.SetFont("Helvetica", "B", 11)
	pdf.Text(legendX, legendY, "Legend")

	pdf.SetFont("Helvetica", "", 8.5)
	legend := []string{
		"Top-left = Territory number",
		"Top-right = Payout for conquering that territory",
		"Bottom-left = Number you must roll",
		"Bottom-right = Number of dice to roll",
	}

	for i, line := range legend {
		pdf.Text(legendX, legendY+5.5+float64(i)*5.0, line)
	}

	// Small footer.
	pdf.SetFont("Helvetica", "", 7)
	pdf.SetTextColor(110, 110, 110)
	footer := "Generated by Country Game Board Generator"
	pdf.Text(pageW-pdf.GetStringWidth(footer)-10, pageH-7, footer)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func chooseColumns(territories int) int {
	// Five columns gives a clean, consistent board for most sizes.
	switch {
	case territories <= 8:
		return 3
	case territories <= 15:
		return 4
	default:
		return 5
	}
}

func centerText(pdf *gofpdf.Fpdf, text string, centerX, baselineY float64) {
	pdf.Text(centerX-pdf.GetStringWidth(text)/2, baselineY, text)
}

func addFlag(pdf *gofpdf.Fpdf, countryCode string, x, y, w, h float64) bool {
	client := &http.Client{Timeout: 4 * time.Second}
	url := "https://flagcdn.com/w160/" + strings.ToLower(countryCode) + ".png"

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	reader := bytes.NewReader(data)
	key := "flag-" + strings.ToLower(countryCode)

	opt := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}

	if err := pdf.RegisterImageOptionsReader(key, opt, reader); err != nil {
		return false
	}

	pdf.ImageOptions(key, x, y, w, h, false, opt, 0, "")
	return pdf.Ok()
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Country Game Board Generator</title>
<style>
	* { box-sizing: border-box; }
	body {
		margin: 0;
		font-family: Arial, Helvetica, sans-serif;
		background: #f3f4f6;
		color: #1f2937;
	}
	.container {
		max-width: 760px;
		margin: 45px auto;
		padding: 0 20px;
	}
	.card {
		background: white;
		border-radius: 14px;
		padding: 30px;
		box-shadow: 0 8px 30px rgba(0,0,0,.08);
	}
	h1 { margin: 0 0 8px; font-size: 30px; }
	.subtitle { margin: 0 0 28px; color: #6b7280; line-height: 1.5; }
	.field { margin-bottom: 20px; }
	label {
		display: block;
		font-weight: 700;
		margin-bottom: 7px;
	}
	select {
		width: 100%;
		padding: 12px 14px;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		background: white;
		font-size: 16px;
	}
	button {
		width: 100%;
		padding: 14px;
		border: 0;
		border-radius: 8px;
		background: #111827;
		color: white;
		font-size: 17px;
		font-weight: 700;
		cursor: pointer;
	}
	button:hover { background: #374151; }
	.note {
		margin-top: 22px;
		padding: 14px;
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		font-size: 14px;
		line-height: 1.5;
	}
	strong { color: #111827; }
</style>
</head>
<body>
<div class="container">
	<div class="card">
		<h1>Country Game Board Generator</h1>
		<p class="subtitle">
			Choose a country, difficulty, and number of territories. The generator
			will create a single-sheet legal-size PDF ready to print.
		</p>

		<form method="post" action="/generate">
			<div class="field">
				<label for="country">Country</label>
				<select id="country" name="country" required>
					{{range .Countries}}
					<option value="{{.Code}}">{{.Name}} — {{printf "%.0f" .AreaKm2}} km²</option>
					{{end}}
				</select>
			</div>

			<div class="field">
				<label for="difficulty">Difficulty</label>
				<select id="difficulty" name="difficulty" required>
					{{range .Difficulties}}
					<option value="{{.}}">{{.}}</option>
					{{end}}
				</select>
			</div>

			<div class="field">
				<label for="territories">Number of Territories</label>
				<select id="territories" name="territories" required>
					{{range .Tiles}}
					<option value="{{.}}">{{.}}</option>
					{{end}}
				</select>
			</div>

			<button type="submit">Generate Printable PDF</button>
		</form>

		<div class="note">
			<strong>How the generator works:</strong>
			country size is displayed as educational information, while difficulty
			and territory count are independent game settings. No combination is
			locked or tracked by this app.
		</div>
	</div>
</div>
</body>
</html>`
