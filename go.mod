module github.com/ItsClairton/Anny

go 1.17

// +heroku goVersion go1.17
require (
	github.com/ItsClairton/gonius v0.0.0-20220104172845-2f30ca4d472d
	github.com/Pauloo27/searchtube v0.0.0-20211210213129-1828077b9033
	github.com/diamondburned/arikawa/v3 v3.0.0-rc.3
	github.com/diamondburned/oggreader v0.0.0-20201118014549-87df9534b647
	github.com/gofiber/fiber/v2 v2.24.0
	github.com/joho/godotenv v1.4.0
	github.com/kkdai/youtube/v2 v2.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
)

require (
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.4.1-0.20201116162257-a2a8dda75c91 // indirect
	github.com/dop251/goja v0.0.0-20211211112501-fb27c91c26ed // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/tidwall/gjson v1.12.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20211214234402-4825e8c3871d // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)

replace github.com/kkdai/youtube/v2 => github.com/ItsClairton/youtube/v2 v2.7.7-0.20220213025140-32db65e01853
