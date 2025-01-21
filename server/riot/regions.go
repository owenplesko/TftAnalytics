package riot

var RegionToCluster = map[string]string{
	"NA1":  "americas",
	"BR1":  "americas",
	"LA1":  "americas",
	"LA2":  "americas",
	"KR":   "asia",
	"JP1":  "asia",
	"EUN1": "europe",
	"EUW1": "europe",
	"TR1":  "europe",
	"RU":   "europe",
	"PH2":  "sea",
	"SG2":  "sea",
	"TH2":  "sea",
	"TW2":  "sea",
	"VN2":  "sea",
	"OC1":  "sea",
}

var ClusterToRegions = map[string][]string{}
