package utils

const NeoworksUsdAccount = "B72D0E55-F53A-4DB0-897E-2CE4A73CB94B"

const CoinbaseUsdAccount = "C4D0E14E-1B2B-4023-AFA6-8891AD1960C9"

func GetFeeAccounts(currency string) (string, string) {
	switch currency {
	case "USD":
		return NeoworksUsdAccount, CoinbaseUsdAccount
	default:
		return "", ""
	}
}
