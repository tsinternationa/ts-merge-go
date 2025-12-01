package main

import (
	"strings"
)

// 地区前缀配置
type AreaPrefix struct {
	Name     string   // 地区名称
	Prefixes []string // 号码前缀列表
}

// 国家配置接口
type CountryConfig struct {
	Name        string       // 国家名称
	CountryCode []string     // 国家代码（如 "1", "86"）
	Areas       []AreaPrefix // 地区配置列表
}

// 判断是否是该国家的号码
func (c *CountryConfig) IsCountryNumber(phoneNumber string) bool {
	for _, code := range c.CountryCode {
		if strings.HasPrefix(phoneNumber, code) {
			return true
		}
	}
	return false
}

// 识别号码所属地区
func (c *CountryConfig) IdentifyArea(phoneNumber string) string {
	// 按前缀长度从长到短排序，优先匹配更长的前缀
	maxLen := 0
	for _, area := range c.Areas {
		for _, prefix := range area.Prefixes {
			if len(prefix) > maxLen {
				maxLen = len(prefix)
			}
		}
	}

	// 从最长前缀开始匹配
	for length := maxLen; length >= 1; length-- {
		for _, area := range c.Areas {
			for _, prefix := range area.Prefixes {
				if len(prefix) == length && strings.HasPrefix(phoneNumber, prefix) {
					return area.Name
				}
			}
		}
	}

	return "" // 未识别出地区
}

// 获取可用的国家列表
func getAvailableCountries() []string {
	return []string{
		"美国",
		"加拿大",
		"英国",
		"澳大利亚",
		"日本",
		"韩国",
		"德国",
		"法国",
		"印度",
	}
}

// 获取国家配置
func getCountryConfig(countryName string) *CountryConfig {
	switch countryName {
	case "美国":
		return getUSAConfig()
	case "加拿大":
		return getCanadaConfig()
	case "英国":
		return getUKConfig()
	case "澳大利亚":
		return getAustraliaConfig()
	case "日本":
		return getJapanConfig()
	case "韩国":
		return getKoreaConfig()
	case "德国":
		return getGermanyConfig()
	case "法国":
		return getFranceConfig()
	case "印度":
		return getIndiaConfig()
	default:
		return nil
	}
}

// 美国配置（排除加拿大区号）
func getUSAConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "美国",
		CountryCode: []string{"1"},
		Areas: []AreaPrefix{
			{Name: "纽约州", Prefixes: []string{"1212", "1315", "1347", "1516", "1518", "1585", "1607", "1631", "1646", "1716", "1718", "1845", "1914", "1917", "1929"}},
			{Name: "加利福尼亚州-洛杉矶", Prefixes: []string{"1213", "1310", "1323", "1424", "1562", "1626", "1747", "1818"}},
			{Name: "伊利诺伊州-芝加哥", Prefixes: []string{"1312", "1224", "1331", "1630", "1708", "1773", "1779", "1847", "1872"}},
			{Name: "德克萨斯州-休斯顿", Prefixes: []string{"1281", "1713", "1832"}},
			{Name: "亚利桑那州-凤凰城", Prefixes: []string{"1480", "1602", "1623", "1928"}},
			{Name: "宾夕法尼亚州-费城", Prefixes: []string{"1215", "1267", "1445", "1484", "1610", "1717", "1724", "1814", "1878"}},
			{Name: "德克萨斯州-圣安东尼奥", Prefixes: []string{"1210", "1726", "1830"}},
			{Name: "加利福尼亚州-圣地亚哥", Prefixes: []string{"1619", "1858", "1935"}},
			{Name: "德克萨斯州-达拉斯", Prefixes: []string{"1214", "1469", "1972"}},
			{Name: "加利福尼亚州-圣何塞", Prefixes: []string{"1408", "1669"}},
			{Name: "德克萨斯州-奥斯汀", Prefixes: []string{"1512", "1737"}},
			{Name: "佛罗里达州-杰克逊维尔", Prefixes: []string{"1904"}},
			{Name: "加利福尼亚州-旧金山", Prefixes: []string{"1415", "1628", "1650"}},
			{Name: "印第安纳州", Prefixes: []string{"1317", "1463", "1765", "1812"}},
			{Name: "俄亥俄州-哥伦布", Prefixes: []string{"1380", "1614", "1740"}},
			{Name: "德克萨斯州-沃斯堡", Prefixes: []string{"1682", "1817"}},
			{Name: "北卡罗来纳州-夏洛特", Prefixes: []string{"1704", "1980"}},
			{Name: "华盛顿州-西雅图", Prefixes: []string{"1206", "1253", "1360", "1425", "1564"}},
			{Name: "科罗拉多州-丹佛", Prefixes: []string{"1303", "1720", "1970"}},
			{Name: "华盛顿特区", Prefixes: []string{"1202"}},
			{Name: "马萨诸塞州-波士顿", Prefixes: []string{"1339", "1351", "1413", "1508", "1617", "1774", "1781", "1857", "1978"}},
			{Name: "密歇根州-底特律", Prefixes: []string{"1248", "1313", "1586", "1734", "1810", "1947"}},
			{Name: "田纳西州", Prefixes: []string{"1423", "1615", "1629", "1731", "1865", "1901", "1931"}},
			{Name: "俄克拉荷马州", Prefixes: []string{"1405", "1539", "1580", "1918"}},
			{Name: "俄勒冈州-波特兰", Prefixes: []string{"1503", "1971"}},
			{Name: "内华达州-拉斯维加斯", Prefixes: []string{"1702", "1725"}},
			{Name: "威斯康星州-密尔沃基", Prefixes: []string{"1262", "1414", "1534"}},
			{Name: "新墨西哥州-阿尔伯克基", Prefixes: []string{"1505", "1575"}},
			{Name: "亚利桑那州-图森", Prefixes: []string{"1520"}},
			{Name: "加利福尼亚州-弗雷斯诺", Prefixes: []string{"1559"}},
			{Name: "加利福尼亚州-萨克拉门托", Prefixes: []string{"1279", "1530", "1916"}},
			{Name: "密苏里州-堪萨斯城", Prefixes: []string{"1816"}},
			{Name: "亚利桑那州-梅萨", Prefixes: []string{"1480"}},
			{Name: "佐治亚州-亚特兰大", Prefixes: []string{"1404", "1470", "1678", "1770", "1943"}},
			{Name: "科罗拉多州-科罗拉多斯普林斯", Prefixes: []string{"1719"}},
			{Name: "北卡罗来纳州-罗利", Prefixes: []string{"1919", "1984"}},
			{Name: "佛罗里达州-迈阿密", Prefixes: []string{"1305", "1786", "1954"}},
			{Name: "加利福尼亚州-长滩", Prefixes: []string{"1562"}},
			{Name: "弗吉尼亚州", Prefixes: []string{"1757", "1703", "1571", "1804"}},
			{Name: "内布拉斯加州-奥马哈", Prefixes: []string{"1402", "1531"}},
			{Name: "加利福尼亚州-奥克兰", Prefixes: []string{"1510"}},
			{Name: "明尼苏达州", Prefixes: []string{"1320", "1507", "1612", "1651", "1763", "1952"}},
			{Name: "俄克拉荷马州-塔尔萨", Prefixes: []string{"1918"}},
			{Name: "德克萨斯州-阿灵顿", Prefixes: []string{"1817"}},
			{Name: "路易斯安那州-新奥尔良", Prefixes: []string{"1504"}},
			{Name: "俄亥俄州-克利夫兰", Prefixes: []string{"1216", "1440"}},
			{Name: "俄亥俄州-辛辛那提", Prefixes: []string{"1513"}},
			{Name: "堪萨斯州", Prefixes: []string{"1316", "1620", "1785", "1913"}},
			{Name: "佛罗里达州-坦帕", Prefixes: []string{"1813", "1727"}},
			{Name: "佛罗里达州-奥兰多", Prefixes: []string{"1407", "1321", "1689"}},
			{Name: "路易斯安那州-巴吞鲁日", Prefixes: []string{"1225"}},
			{Name: "密西西比州", Prefixes: []string{"1228", "1601", "1662", "1769"}},
			{Name: "阿拉巴马州", Prefixes: []string{"1205", "1251", "1256", "1334", "1938"}},
			{Name: "南卡罗来纳州", Prefixes: []string{"1803", "1843", "1854", "1864"}},
			{Name: "肯塔基州", Prefixes: []string{"1270", "1502", "1606", "1859"}},
			{Name: "爱荷华州", Prefixes: []string{"1319", "1515", "1563", "1641", "1712"}},
			{Name: "阿肯色州", Prefixes: []string{"1479", "1501", "1870"}},
			{Name: "犹他州", Prefixes: []string{"1385", "1801"}},
			{Name: "内华达州-雷诺", Prefixes: []string{"1775"}},
			{Name: "康涅狄格州", Prefixes: []string{"1203", "1475", "1860", "1959"}},
			{Name: "新泽西州", Prefixes: []string{"1201", "1551", "1609", "1732", "1848", "1856", "1862", "1908", "1973"}},
			{Name: "罗德岛州", Prefixes: []string{"1401"}},
			{Name: "新罕布什尔州", Prefixes: []string{"1603"}},
			{Name: "缅因州", Prefixes: []string{"1207"}},
			{Name: "佛蒙特州", Prefixes: []string{"1802"}},
			{Name: "特拉华州", Prefixes: []string{"1302"}},
			{Name: "西弗吉尼亚州", Prefixes: []string{"1304"}},
			{Name: "怀俄明州", Prefixes: []string{"1307"}},
			{Name: "蒙大拿州", Prefixes: []string{"1406"}},
			{Name: "南达科他州", Prefixes: []string{"1605"}},
			{Name: "北达科他州", Prefixes: []string{"1701"}},
			{Name: "阿拉斯加州", Prefixes: []string{"1907"}},
			{Name: "夏威夷州", Prefixes: []string{"1808"}},
			{Name: "爱达荷州", Prefixes: []string{"1208"}},
		},
	}
}

// 英国配置
func getUKConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "英国",
		CountryCode: []string{"44"},
		Areas: []AreaPrefix{
			{Name: "伦敦", Prefixes: []string{"4420"}},
			{Name: "曼彻斯特", Prefixes: []string{"44161"}},
			{Name: "伯明翰", Prefixes: []string{"44121"}},
			{Name: "利兹", Prefixes: []string{"44113"}},
			{Name: "格拉斯哥", Prefixes: []string{"44141"}},
			{Name: "爱丁堡", Prefixes: []string{"44131"}},
			{Name: "利物浦", Prefixes: []string{"44151"}},
			{Name: "布里斯托", Prefixes: []string{"44117"}},
			{Name: "谢菲尔德", Prefixes: []string{"44114"}},
			{Name: "纽卡斯尔", Prefixes: []string{"44191"}},
			{Name: "贝尔法斯特", Prefixes: []string{"4428"}},
			{Name: "卡迪夫", Prefixes: []string{"4429"}},
		},
	}
}

// 加拿大配置（与美国完全区分）
func getCanadaConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "加拿大",
		CountryCode: []string{"1"},
		Areas: []AreaPrefix{
			{Name: "安大略省-多伦多", Prefixes: []string{"1416", "1437", "1647", "1905"}},
			{Name: "魁北克省-蒙特利尔", Prefixes: []string{"1438", "1514", "1450"}},
			{Name: "不列颠哥伦比亚省-温哥华", Prefixes: []string{"1604", "1778", "1236"}},
			{Name: "阿尔伯塔省-卡尔加里", Prefixes: []string{"1403", "1587", "1825"}},
			{Name: "阿尔伯塔省-埃德蒙顿", Prefixes: []string{"1780", "1825"}},
			{Name: "安大略省-渥太华", Prefixes: []string{"1343", "1613", "1819"}},
			{Name: "魁北克省-魁北克城", Prefixes: []string{"1418", "1581", "1873"}},
			{Name: "马尼托巴省-温尼伯", Prefixes: []string{"1204", "1431"}},
			{Name: "新斯科舍省-哈利法克斯", Prefixes: []string{"1782", "1902"}},
			{Name: "不列颠哥伦比亚省-维多利亚", Prefixes: []string{"1250", "1672"}},
			{Name: "安大略省-汉密尔顿", Prefixes: []string{"1289", "1905"}},
			{Name: "安大略省-伦敦", Prefixes: []string{"1226", "1519"}},
			{Name: "安大略省-温莎", Prefixes: []string{"1519", "1226"}},
			{Name: "萨斯喀彻温省-里贾纳", Prefixes: []string{"1306", "1639"}},
			{Name: "萨斯喀彻温省-萨斯卡通", Prefixes: []string{"1306", "1639"}},
			{Name: "新不伦瑞克省", Prefixes: []string{"1506"}},
			{Name: "纽芬兰和拉布拉多省", Prefixes: []string{"1709"}},
			{Name: "爱德华王子岛省", Prefixes: []string{"1902"}},
			{Name: "育空地区", Prefixes: []string{"1867"}},
			{Name: "西北地区", Prefixes: []string{"1867"}},
			{Name: "努纳武特地区", Prefixes: []string{"1867"}},
		},
	}
}

// 澳大利亚配置
func getAustraliaConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "澳大利亚",
		CountryCode: []string{"61"},
		Areas: []AreaPrefix{
			{Name: "悉尼", Prefixes: []string{"612"}},
			{Name: "墨尔本", Prefixes: []string{"613"}},
			{Name: "布里斯班", Prefixes: []string{"617"}},
			{Name: "珀斯", Prefixes: []string{"618"}},
			{Name: "阿德莱德", Prefixes: []string{"618"}},
			{Name: "堪培拉", Prefixes: []string{"612"}},
		},
	}
}

// 日本配置
func getJapanConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "日本",
		CountryCode: []string{"81"},
		Areas: []AreaPrefix{
			{Name: "东京", Prefixes: []string{"813"}},
			{Name: "大阪", Prefixes: []string{"816"}},
			{Name: "名古屋", Prefixes: []string{"8152"}},
			{Name: "札幌", Prefixes: []string{"8111"}},
			{Name: "福冈", Prefixes: []string{"8192"}},
			{Name: "京都", Prefixes: []string{"8175"}},
			{Name: "横滨", Prefixes: []string{"8145"}},
		},
	}
}

// 韩国配置
func getKoreaConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "韩国",
		CountryCode: []string{"82"},
		Areas: []AreaPrefix{
			{Name: "首尔", Prefixes: []string{"822"}},
			{Name: "釜山", Prefixes: []string{"8251"}},
			{Name: "仁川", Prefixes: []string{"8232"}},
			{Name: "大邱", Prefixes: []string{"8253"}},
			{Name: "大田", Prefixes: []string{"8242"}},
			{Name: "光州", Prefixes: []string{"8262"}},
			{Name: "蔚山", Prefixes: []string{"8252"}},
		},
	}
}

// 德国配置
func getGermanyConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "德国",
		CountryCode: []string{"49"},
		Areas: []AreaPrefix{
			{Name: "柏林", Prefixes: []string{"4930"}},
			{Name: "慕尼黑", Prefixes: []string{"4989"}},
			{Name: "汉堡", Prefixes: []string{"4940"}},
			{Name: "法兰克福", Prefixes: []string{"4969"}},
			{Name: "科隆", Prefixes: []string{"49221"}},
			{Name: "斯图加特", Prefixes: []string{"49711"}},
			{Name: "杜塞尔多夫", Prefixes: []string{"49211"}},
			{Name: "多特蒙德", Prefixes: []string{"49231"}},
		},
	}
}

// 法国配置
func getFranceConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "法国",
		CountryCode: []string{"33"},
		Areas: []AreaPrefix{
			{Name: "巴黎", Prefixes: []string{"331"}},
			{Name: "马赛", Prefixes: []string{"334"}},
			{Name: "里昂", Prefixes: []string{"334"}},
			{Name: "图卢兹", Prefixes: []string{"335"}},
			{Name: "尼斯", Prefixes: []string{"334"}},
			{Name: "南特", Prefixes: []string{"332"}},
			{Name: "斯特拉斯堡", Prefixes: []string{"333"}},
		},
	}
}

// 印度配置
func getIndiaConfig() *CountryConfig {
	return &CountryConfig{
		Name:        "印度",
		CountryCode: []string{"91"},
		Areas: []AreaPrefix{
			{Name: "德里", Prefixes: []string{"9111"}},
			{Name: "孟买", Prefixes: []string{"9122"}},
			{Name: "班加罗尔", Prefixes: []string{"9180"}},
			{Name: "海得拉巴", Prefixes: []string{"9140"}},
			{Name: "艾哈迈达巴德", Prefixes: []string{"9179"}},
			{Name: "金奈", Prefixes: []string{"9144"}},
			{Name: "加尔各答", Prefixes: []string{"9133"}},
			{Name: "浦那", Prefixes: []string{"9120"}},
		},
	}
}
