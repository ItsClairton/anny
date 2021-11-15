package emojis

var (
	PingPong = "🏓"

	Peer    = "<:KannaPeer:838567821205176340>"
	Yeah    = "<:yeah:838568353139916850>"
	Cry     = "<:paimonpanic:909597749970874408>"
	Loading = "<a:1180staff:836709984909525032>"
	Sleep   = "<:keqingsleep:909567537778421810>"
	Sip     = "<:nejiresip:909600499131646063>"

	Twitch  = "<:twitch:896600475833606154>"
	Youtube = "<:youtube:896600900909559868>"

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
