package emojis

var (
	PingPong = "🏓"

	PepeArt       = "<:pepeart:857771343633580113>"
	KannaPeer     = "<:KannaPeer:838567821205176340>"
	MikuCry       = "<:mikuCry:830091923129237554>"
	AnimatedStaff = "<a:1180staff:836709984909525032>"
)

func GetNumberAsEmoji(num int) string {
	switch num {
	case 0:
		return "0⃣"
	case 1:
		return "1⃣"
	case 2:
		return "2⃣"
	case 3:
		return "3⃣"
	case 4:
		return "4⃣"
	case 5:
		return "5⃣"
	case 6:
		return "6⃣"
	case 7:
		return "7⃣"
	case 8:
		return "8⃣"
	case 9:
		return "9⃣"
	case 10:
		return "🔟"
	default:
		return ""
	}
}
