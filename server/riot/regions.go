package riot

var RegionToCluster = map[string]string{
	"NA1":  "americas",
	"BR1":  "americas",
	"LA1":  "americas",
	"LA2":  "americas",
	"KR":   "asia",
	"JP1":  "asia",
	"ME1":  "asia",
	"SG2":  "sea",
	"TW2":  "sea",
	"VN2":  "sea",
	"OC1":  "sea",
	"EUN1": "europe",
	"EUW1": "europe",
	"TR1":  "europe",
	"RU":   "europe",
}

var ClusterToRegions map[string][]string

func init() {
	ClusterToRegions = make(map[string][]string)

	for region, cluster := range RegionToCluster {
		arr := ClusterToRegions[cluster]
		arr = append(arr, region)
		ClusterToRegions[cluster] = arr
	}
}
