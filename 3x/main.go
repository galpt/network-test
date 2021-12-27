package main

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"encoding/base64"

	goaway "github.com/TwiN/go-away"
	"github.com/bwmarrin/discordgo"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/lus/dgc"
	totalmem "github.com/pbnjay/memory"
	"github.com/spf13/afero"
	"github.com/tidwall/gjson"
	xurls "mvdan.cc/xurls/v2"
)

const (
	Gigabyte      = 1 << 30
	Megabyte      = 1 << 20
	Kilobyte      = 1 << 10
	giperfexeLoc  = "H:\\giperf\\giperf\\giperf.exe"
	presetsbatLoc = "H:\\giperf\\giperf\\Presets.bat"
)

var (
	b64katmon = "T0RVME1EY3hNVGt6T0RNek56QXhOREUyLllNZWx5QS5FSGFkZFFPYUVQZXpoWnRGQk5mRjRGaEU2WHc="
	b64katliy = "T1RBek5UVXdORE01TnpjeU1ERTJOamN4LllYdW02dy5hcGxIRFpUSUFyS0R6RWNhUDl0T0JHRHVrdG8="
	b64katinz = "T1RBME16QTNNemcyTWpNME16STNNRGN3LllYNW40Zy5yLTF1M24zc1JlYWlMTmR0UXVuWkx6aGdHT1E="

	tlsConf = &tls.Config{
		InsecureSkipVerify:          true,
		PreferServerCipherSuites:    false,
		SessionTicketsDisabled:      false,
		DynamicRecordSizingDisabled: false,
	}

	nhCached         = ""
	nhGetGallIDSplit []string
	nhGetGallID1     string
	nhGetGallID2     string
	nhImgLink        string
	nhImgLinkLocal   string
	nhImgName        string
	nhImgNames       []string
	nhImgLinks       []string
	nhCode           string
	nhTotalPage      int

	getMaxRender  = 1
	getImgs       []string
	getFileFormat = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	ckImgs        []string
	vmgMaxRender  = 1
	scaleOpts     = []string{"1", "2", "4", "8", "16", "32"}
	scaleOK       bool

	// katInz GET feature
	katInzGETCachedURL      = ""
	katInzGETCachedFileName = ""

	// katInz YTDL feature
	katInzVidID = ""

	universalLogs      []string
	universalLogsLimit = 50

	xurlsRelaxed = xurls.Relaxed()

	botName          = "Katheryne"
	katMondstadtSess *discordgo.Session
	katLiyueSess     *discordgo.Session
	katInazumaSess   *discordgo.Session

	//go:embed serverRules.txt
	serverRules string

	kokonattomilkuGuildID       = "893138943334297682"
	kokonattomilkuBackupGuildID = "904497628874682378"

	ucoverModsDB   string
	ucoverUsername string
	ucoveruserID   string
	ucoverInfo     []string
	ucoverNewData  []string
	ucoverNewAdded []string

	noBanStaff    bool
	staffDetected bool
	staffID       = []string{"631418827841863712", "149228888403214337", "320455208524316672", "682274986987356184", "856073889847574538", "726577226023436392"}

	botID              = []string{"854071193833701416", "903550439772016671", "904307386234327070"}
	maidchanID         = "903550439772016671"
	maidsanID          = "854071193833701416"
	katheryneInazumaID = "904307386234327070"

	blacklistedID = []string{"485113382547226645", "818007831641980928"}

	giperfChangelog string
	giperfExeSHA512 string
	osFS            = afero.NewOsFs()
	memFS           = afero.NewMemMapFs()
	httpFs          = afero.NewHttpFs(memFS)
	mem             runtime.MemStats
	duration        = time.Now()
	ReqLogs         string
	RespLogs        string
	ConnReqLogs     string
	totalMem        string
	HeapAlloc       string
	SysMem          string
	Frees           string
	NumGCMem        string
	timeElapsed     string
	latestLog       string
	winLogs         string
	tempDirLoc      string

	bannedWords         = []string{"cheat", "inject", "hack"}
	falsePosWords       = []string{"kokonatto", "milku", "hakku"}
	falseNegWords       = []string{"injector", "cheating", "hacking", "cheats", "injection", "hacks"}
	bannedWordsDetector = goaway.NewProfanityDetector().WithCustomDictionary(bannedWords, falsePosWords, falseNegWords)

	lastMsgTimestamp   string
	lastMsgUsername    string
	lastMsgUserID      string
	lastMsgpfp         string
	lastMsgAccType     string
	lastMsgID          string
	lastMsgContent     string
	lastMsgTranslation string

	maidsanErrorMsg         string
	maidsanLastMsgChannelID string
	maidsanLastMsgID        string
	maidsanLowercaseLastMsg string
	maidsanEditedLastMsg    string
	maidsanTranslatedMsg    string
	maidsanBanUserMsg       string
	maidsanWarnMsg          string

	maidchanLastMsgChannelID string
	maidchanLastMsgID        string

	katInzBlacklist               []string
	katInzBlacklistReadable       string
	katInzBlacklistLinkDetected   bool
	katInzCustomBlacklist         = []string{"discordf.gift"}
	katInzCustomBlacklistReadable string
	katInzAddCustomBlacklist      []string
	katInzNewAppended             string

	editedGETData string

	maidsanLogs         []string
	maidsanLogsLimit    = 100
	maidsanLogsTemplate string
	timestampLogs       []string
	useridLogs          []string
	profpicLogs         []string
	acctypeLogs         []string
	msgidLogs           []string
	msgLogs             []string
	translateLogs       []string

	maidsanBanList           []string
	maidsanEmojiInfo         []string
	maidsanWatchCurrentUser  string
	maidsanWatchPreviousUser string
	maidsanWelcomeMsg        string

	replyremoveNewLines string
	replyremoveSpaces   string
	replysplitEmojiInfo []string
	customEmojiReply    string
	customEmojiDetected bool

	customEmojiSlice []string
	customEmojiIdx   = 0

	welcomeradarChannelID = "894459541566136330"
	updatesChannelID      = "893140731848425492"
	genENChannelID        = "895285059106533446"
	genCNChannelID        = "903627523047440384"
	genRUChannelID        = "914707143423311903"
	offtopicChannelID     = "893140319200235561"

	h1Tr = &http.Transport{
		DisableKeepAlives:      false,
		DisableCompression:     false,
		ForceAttemptHTTP2:      false,
		TLSClientConfig:        tlsConf,
		TLSHandshakeTimeout:    20 * time.Second,
		ResponseHeaderTimeout:  20 * time.Second,
		IdleConnTimeout:        30 * time.Minute,
		ExpectContinueTimeout:  1 * time.Second,
		MaxIdleConns:           10000,    // Prevents resource exhaustion
		MaxIdleConnsPerHost:    100,      // Increases performance and prevents resource exhaustion
		MaxConnsPerHost:        0,        // 0 for no limit
		MaxResponseHeaderBytes: 32 << 10, // 32k
		WriteBufferSize:        1 << 20,  // 1m to minimize I/O writes and increase performance
		ReadBufferSize:         1 << 20,  // 1m to minimize I/O writes and increase performance
	}

	httpclient = &http.Client{
		Timeout:   30 * time.Minute,
		Transport: h1Tr,
	}
)

// =========================================
// HTTP server with customizable port (default is 9999)
func proxyServer() {

	duration := time.Now()

	// Use Gin as the HTTP router
	gin.SetMode(gin.ReleaseMode)
	ginroute := gin.New()
	ginroute.Use(gin.Recovery())
	ginroute.LoadHTMLGlob("templates/*.html")

	// print universalLogs slice
	ginroute.GET("/logs", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n • [UNIVERSAL LOGS] \n •===========================• \n \n%v \n\n", time.Now().UTC().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, universalLogs)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// Custom NotFound handler
	ginroute.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// Print homepage.
	ginroute.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "home.html", gin.H{
			"ServerTime":     fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850)),
			"TotalCPU":       fmt.Sprintf("%v", runtime.NumCPU()),
			"TotalMem":       fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte)),
			"MsgTimestamp":   lastMsgTimestamp,
			"MsgUsername":    lastMsgUsername,
			"MsgUserID":      lastMsgUserID,
			"Msgpfp":         lastMsgpfp,
			"MsgAccType":     lastMsgAccType,
			"MsgID":          lastMsgID,
			"MsgContent":     lastMsgContent,
			"MsgTranslation": lastMsgTranslation,
			"AllMsg":         fmt.Sprintf("%v", maidsanLogs),
		})
	})

	// Control Windows OS through proxy
	ginroute.StaticFS("/temp", http.Dir(os.TempDir()))
	ginroute.GET("/gettemp", func(c *gin.Context) {

		// Get the location of the TEMP dir
		tempDirLoc = fmt.Sprintf(" [DONE] Detected TEMP folder location \n >> %v", os.TempDir())
		c.String(http.StatusOK, tempDirLoc)
	})
	ginroute.GET("/deltemp", func(c *gin.Context) {

		// Delete the entire TEMP folder.
		// If it gets deleted properly, create a new TEMP folder.
		delTemp := osFS.RemoveAll(os.TempDir())
		if delTemp == nil {
			mkTemp := osFS.MkdirAll(os.TempDir(), 0777)
			if mkTemp != nil {
				winLogs = "\n • [ERROR] Failed to recreate TEMP folder. \n • Timestamp >> " + fmt.Sprintf("%v", time.Now()) + "\n • Reason >> " + fmt.Sprintf("%v", mkTemp)
				c.String(http.StatusOK, winLogs)
			}
			winLogs = "\n • [DONE] TEMP folder has been cleaned. \n • Timestamp >> " + fmt.Sprintf("%v", time.Now()) + "\n • Reason >> " + fmt.Sprintf("%v", mkTemp)
			c.String(http.StatusOK, winLogs)
		} else {
			winLogs = "\n • [ERROR] Failed to delete some files. \n • Timestamp >> " + fmt.Sprintf("%v", time.Now()) + "\n • Reason >> " + fmt.Sprintf("%v", delTemp)
			c.String(http.StatusOK, winLogs)
		}
	})

	// get Maid-san's available emoji info
	ginroute.GET("/emoji", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n • [AVAILABLE EMOJI LIST] \n • Total Available Emoji: %v \n •===========================• \n \n[Name —— Emoji ID —— Animated (true/false) —— Guild Name —— Guild ID]\n\n%v \n\n", time.Now().UTC().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, len(maidsanEmojiInfo), maidsanEmojiInfo)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// get Maid-san's URL blacklist
	ginroute.GET("/blacklist", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n • [BLACKLISTED LINKS] \n •===========================• \n\n [CUSTOM BLACKLIST] \n%v \n\n [AUTO BLACKLIST] \n%v \n\n", time.Now().UTC().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, katInzCustomBlacklistReadable, katInzBlacklistReadable)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// get Maid-san to get the Guild ban list
	ginroute.GET("/banlist", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n • [BAN LIST] \n •===========================• \n \n[Username : User ID : Reason]\n\n%v \n\n", time.Now().UTC().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, maidsanBanList)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// get Maid-san to get the Guild undercover mod list
	ginroute.GET("/undercover", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n • [UNDERCOVER LIST] \n • Total Undercover Mods: %v \n •===========================• \n \n[Username —— User ID —— Undercover ID —— Source]\n\n%v \n\n", time.Now().UTC().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, len(ucoverInfo), ucoverInfo)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// support for nhen
	nhRelax := xurls.Relaxed()
	ginroute.GET("/nh/:nhcode", func(c *gin.Context) {

		// get url param
		nhCode = c.Param("nhcode")

		// check whether the manga has been cached or not
		if nhCode == nhCached {

			// get manga data from cache
			c.HTML(http.StatusOK, "nh.html", nhImgNames)
		} else {

			// fetch data directly from nhen server
			memFS.RemoveAll("./nh/")
			memFS.MkdirAll("./nh/", 0777)

			// Get the gallery ID
			nhGalleryID := fmt.Sprintf("https://nhentai.net/g/%v", nhCode)
			getGalleryID, err := httpclient.Get(nhGalleryID)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			bodyGalleryID, err := ioutil.ReadAll(bufio.NewReader(getGalleryID.Body))
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			scanPg1 := strings.Split(string(bodyGalleryID), "Pages:")
			scanPg2 := strings.Split(scanPg1[1], `<span class="name">`)
			scanPg3 := strings.Split(scanPg2[1], `</span></a></span></div><div class="tag-container field-name">`)
			nhTotalPage, err = strconv.Atoi(scanPg3[0])
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			scanGalleryID := nhRelax.FindAllString(string(bodyGalleryID), -1)
			for nhlinkIdx := range scanGalleryID {
				if strings.Contains(scanGalleryID[nhlinkIdx], "https://t.nhentai.net/galleries/") {
					nhGetGallIDSplit = nil
					nhGetGallID1 = strings.ReplaceAll(scanGalleryID[nhlinkIdx], "https://t.nhentai.net/galleries/", "")
					nhGetGallIDSplit = strings.Split(nhGetGallID1, "/cover")
					nhGetGallID2 = nhGetGallIDSplit[0]
					break
				}
			}

			nhImgNames = nil
			nhImgLinks = nil
			for nhCurrPg := 1; nhCurrPg <= nhTotalPage; nhCurrPg++ {
				nhImgLink = fmt.Sprintf("https://i.nhentai.net/galleries/%v/%v.jpg", nhGetGallID2, nhCurrPg)

				// Get the image and write it to memory
				getImg, err := httpclient.Get(nhImgLink)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				}

				// ==================================
				// Create a new image based on bodyImg
				nhImgName = fmt.Sprintf("./nh/%v.jpg", nhCurrPg)
				createImgFile, err := memFS.Create(nhImgName)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				} else {
					// Write to the file
					writeImgFile, err := io.Copy(createImgFile, getImg.Body)
					if err != nil {
						fmt.Println(" [ERROR] ", err)
						getImg.Body.Close()

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						break
					} else {
						// Close the file
						getImg.Body.Close()
						if err := createImgFile.Close(); err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							break
						} else {
							winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", nhImgName, (writeImgFile / Kilobyte), (writeImgFile / Megabyte))
							fmt.Println(winLogs)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
							}
						}
					}
				}

				nhImgLinkLocal = fmt.Sprintf("https://x.castella.network/img/%v.jpg", nhCurrPg)
				nhImgNames = append(nhImgNames, nhImgLinkLocal)
				nhImgLinks = append(nhImgLinks, fmt.Sprintf("\n%v", nhImgLink))
			}

			getGalleryID.Body.Close()

			// update cache
			nhCached = nhCode

			c.HTML(http.StatusOK, "nh.html", nhImgNames)
		}

	})

	// get data from memory
	ginroute.StaticFS("/get", httpFs.Dir("./get/"))
	ginroute.StaticFS("/img", httpFs.Dir("./nh/"))
	ginroute.StaticFS("/ai", httpFs.Dir("./pics/"))
	ginroute.GET("/nhlinks", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n โ€ข===========================โ€ข \n โ€ข [SERVER STATUS] \n โ€ข Last Modified: %v \n โ€ข Total OS Memory: %v \n โ€ข Completed GC Cycles: %v \n โ€ข Time Elapsed: %v \n โ€ข===========================โ€ข \n โ€ข [NH IMAGE LINKS] \n โ€ข===========================โ€ข \n \n%v \n\n", time.Now().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, nhImgLinks)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// HTTP proxy server
	httpserver := &http.Server{
		Addr:              ":80",
		Handler:           ginroute,
		TLSConfig:         tlsConf,
		MaxHeaderBytes:    32 << 10,      // 32k
		WriteTimeout:      1 * time.Hour, // required for anime streaming
		ReadTimeout:       1 * time.Hour, // required for anime streaming
		ReadHeaderTimeout: 15 * time.Second,
		IdleTimeout:       1 * time.Hour,
	}
	httpserver.SetKeepAlivesEnabled(true)
	httpserver.ListenAndServe()
}

// =========================================
// The main function of Maid-san bot
func main() {

	// Automatically set GOMAXPROCS to the number of your CPU cores.
	// Increase performance by allowing Golang to use multiple processors.
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs) // Sets the GOMAXPROCS value
	totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
	fmt.Println()

	fmt.Println(numCPUs)
	fmt.Println(totalMem)
	fmt.Println(lastMsgTimestamp)
	fmt.Println(lastMsgUsername)
	fmt.Println(lastMsgUserID)
	fmt.Println(lastMsgpfp)
	fmt.Println(lastMsgAccType)
	fmt.Println(lastMsgID)
	fmt.Println(lastMsgContent)
	fmt.Println(lastMsgTranslation)
	fmt.Println(katInzBlacklistReadable)
	fmt.Println(katInzCustomBlacklistReadable)

	// Set max number of operating system threads based on the number of CPUs.
	debug.SetMaxThreads(numCPUs * 10000)

	// Run proxy server in a separated goroutine
	go proxyServer()

	// Create the logs folder
	osFS.RemoveAll("./logs/")
	createLogFolder := osFS.MkdirAll("./logs/", 0777)
	if createLogFolder != nil {
		fmt.Println(" [ERROR] ", createLogFolder)
	}
	fmt.Println(` [DONE] New "logs" folder has been created. \n >> `, createLogFolder)

	createDBFolder := osFS.MkdirAll("./db/", 0777)
	if createDBFolder != nil {
		fmt.Println(" [ERROR] ", createDBFolder)
	}
	fmt.Println(` [DONE] New "db" folder has been created. \n >> `, createDBFolder)

	// Get GIPerf changelog file and update the old data
	readChangelog, err := afero.ReadFile(osFS, "./changelog.txt")
	if err != nil {
		fmt.Println(" [ERROR] ", err)
	}
	giperfChangelog = fmt.Sprintf("%v", string(readChangelog))

	// Get GIPerf 1-undercover-mods.txt file and update the old data
	readUndercoverData, err := afero.ReadFile(osFS, "./1-undercover-mods.txt")
	if err != nil {
		fmt.Println(" [ERROR] ", err)
	}
	ucoverModsDB = fmt.Sprintf("%v", string(readUndercoverData))

	// handle decoding
	decodekatmon, err := base64.StdEncoding.DecodeString(b64katmon)
	if err != nil {
		fmt.Println(" [ERROR] ", err)
	}
	decodekatliy, err := base64.StdEncoding.DecodeString(b64katliy)
	if err != nil {
		fmt.Println(" [ERROR] ", err)
	}
	decodekatinz, err := base64.StdEncoding.DecodeString(b64katinz)
	if err != nil {
		fmt.Println(" [ERROR] ", err)
	}

	// Create a new Discord session using the provided login information.
	// Maid-san session.
	maidsanSession, err := discordgo.New("Bot " + string(decodekatmon))
	if err != nil {
		fmt.Println(" [MAID-SAN CREATE SESSION] ", err)
		return
	}
	// Maid-chan session.
	maidchanSession, err := discordgo.New("Bot " + string(decodekatliy))
	if err != nil {
		fmt.Println(" [MAID-CHAN CREATE SESSION] ", err)
		return
	}
	// Katheryne Inazuma session.
	katInzSession, err := discordgo.New("Bot " + string(decodekatinz))
	if err != nil {
		fmt.Println(" [Katheryne Inazuma CREATE SESSION] ", err)
		return
	}

	// Use custom HTTP client
	maidsanSession.Client = httpclient
	maidchanSession.Client = httpclient
	katInzSession.Client = httpclient

	// err = maidsanSession.Open()
	// if err != nil {
	// 	fmt.Println(" [Katheryne Mondstadt OPEN SESSION] ", err)
	// 	return
	// }
	// err = maidchanSession.Open()
	// if err != nil {
	// 	fmt.Println(" [Katheryne Liyue OPEN SESSION] ", err)
	// 	return
	// }
	err = katInzSession.Open()
	if err != nil {
		fmt.Println(" [Katheryne Inazuma OPEN SESSION] ", err)
		return
	}
	fmt.Println()
	// fmt.Println(" [Katheryne Mondstadt is Up and Running] ")
	// fmt.Println(" [Katheryne Liyue is Up and Running] ")
	fmt.Println(" [Katheryne Inazuma is Up and Running] ")

	katMondstadtSess = maidsanSession
	katLiyueSess = maidchanSession
	katInazumaSess = katInzSession

	maidsanSession.AddHandler(maidsanEmojiReact)
	maidsanSession.AddHandler(maidsanAutoCheck)
	maidchanSession.AddHandler(maidchanAutoCheck)
	katInzSession.AddHandler(katInzAutoCheck)

	// Add emoji reactions handler
	maidsanSession.AddHandler(emojiReactions)
	maidchanSession.AddHandler(emojiReactions)
	katInzSession.AddHandler(emojiReactions)

	// Add sayKatherynes handler
	maidsanSession.AddHandler(sayKatherynes)
	maidchanSession.AddHandler(sayKatherynes)
	katInzSession.AddHandler(sayKatherynes)

	// Wait for the user to cancel the process
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}()

	// Router for Katheryne Mondstadt
	aoiRouter := dgc.Create(&dgc.Router{
		Prefixes:         []string{""},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands:         []*dgc.Command{},
		PingHandler: func(ctx *dgc.Ctx) {
			ctx.RespondText("I'm here~ \nYou can use `!help` to see the available commands, Master.")
		},
	})

	aoiRouter.RegisterDefaultHelpCommand(maidsanSession, nil)

	// Get the bot's current server status
	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "status",
		Description: "I'll inform you about my current status, Master.",
		Usage:       "status\nstatus update\nstatus push",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     getServerStatus,
	})

	// Get latest COVID-19 data based on user's argument
	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "covid19",
		Description: "I'll give you a brief information about the latest COVID-19 data from a certain country.",
		Usage:       "covid19 <country>",
		Example:     "covid19 japan",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     getCovidData,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "check",
		Description: "I'll check the mentioned user for you, Master.",
		Usage:       "check @<user>",
		Example:     "check @Katheryne",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     getUserInfo,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "role",
		Description: "I'll add the specific role to the mentioned user, Master.",
		Usage:       "role @<user> <role>\n\nAvailable Roles\n• kokomember\n• updates\n• releases\n• announcements\n• betatester",
		Example:     "role @Katheryne kokomember",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     setRole,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "say",
		Description: "I'll repeat after you, Master.",
		Usage:       "say <mention katheryne> <mention channel> :: <any text>",
		Example:     "say <@!854071193833701416> <#897231580689485824> :: Hi! I'm Katheryne.",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     sayKatherynes,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "emoji",
		Description: "I'll reply with emoji, Master.",
		Usage:       "emoji <emoji name>",
		Example:     "emoji ganyustare",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     getEmojiMaidsan,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "rules",
		Description: "I'll tell the people about the rules, Master.",
		Usage:       "rules",
		Example:     "rules",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     getRules,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "blacklist",
		Description: "I'll update the blacklisted links, Master.",
		Usage:       "blacklist <link>\nblacklist <link-1>:<link-2>:<link-3>",
		Example:     "blacklist discorda.gift\nblacklist discorda.gift:discordb.gift:discordc.gift",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     addBlacklistLinks,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "ban",
		Description: "I'll ban/unban the given UserID, Master.",
		Usage:       "ban add@@<UserID>@@<reason>\n\nban remove@@<UserID>",
		Example:     "ban add@@631418827841863712@@infracted rules\n\nban remove@@631418827841863712",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     banUser,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "ucover",
		Description: "I'll add/remove the given UserID, Master.",
		Usage:       "ucover add@@<UserID>\nucover add@@<UserID1>:<UserID2>:...\n\nucover remove@@<UserID>",
		Example:     "ucover add@@112233445566778899\n\nucover remove@@112233445566778899",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     ucoverModsHandler,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "delmsg",
		Description: "I'll delete the message from the given Channel & Message ID, Master.",
		Usage:       "delmsg <Message Link>",
		Example:     "delmsg https://discord.com/channels/893138943334297682/895285059106533446/924216768111734815",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     ucoverModsDelMsg,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "msg1",
		Description: "I'll send the message to Genshin-EN channel, Master.",
		Usage:       "msg1 <Your Message>",
		Example:     "msg1 EN test",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     ucoverModsMsgGenEN,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "msg2",
		Description: "I'll send the message to Genshin-CN channel, Master.",
		Usage:       "msg2 <Your Message>",
		Example:     "msg2 CN test",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     ucoverModsMsgGenCN,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "msg3",
		Description: "I'll send the message to Genshin-RU channel, Master.",
		Usage:       "msg3 <Your Message>",
		Example:     "msg3 RU test",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     ucoverModsMsgGenRU,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "msg4",
		Description: "I'll send the message to Off-Topic channel, Master.",
		Usage:       "msg4 <Your Message>",
		Example:     "msg4 Off-Topic test",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     ucoverModsMsgOffTopic,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "warn",
		Description: "I'll convert your message as a warning message, Master.",
		Usage:       "warn <any text>",
		Example:     "warn This is a **warning message test**.",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     warnMsg,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "go",
		Description: "I'll simulate your Go code, Master.",
		Usage:       "go <your Go code>",
		Example:     "go <your Go code>",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katMonGoRun,
	})

	aoiRouter.RegisterCmd(&dgc.Command{
		Name:        "aipic",
		Description: "I'll improve the given image URL using my AI, Master.",
		Usage:       "aipic <any image URL>\n\naipic <any image URL> [<mode>::<scale>]\n\n<mode> -> anime/photo/auto\n\n<scale> -> 1/2/4/8/16/32",
		Example:     "aipic https://www.anime-planet.com/images/manga/covers/thumbs/42635.jpg?t=1573224888\n\naipic https://www.anime-planet.com/images/manga/covers/thumbs/42635.jpg?t=1573224888 [auto::4]",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katMonWaifu2x,
	})

	// Initialize the router for Katherye Mondstadt
	aoiRouter.Initialize(maidsanSession)

	// Router for Katheryne Inazuma
	katInzRouter := dgc.Create(&dgc.Router{
		Prefixes:         []string{""},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands:         []*dgc.Command{},
		PingHandler: func(ctx *dgc.Ctx) {
			ctx.RespondText("Yes~")
		},
	})

	katInzRouter.RegisterCmd(&dgc.Command{
		Name:        "get",
		Description: "I'll get the content from the given link, Master.",
		Usage:       "get <URL>",
		Example:     "get https://castella.network/hotlink-ok/qiqi1.png",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katInzGet,
	})

	katInzRouter.RegisterCmd(&dgc.Command{
		Name:        "nh",
		Description: "I'll get the content from the given nh ID, Master.",
		Usage:       "nh <nh ID>",
		Example:     "nh 114712",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katInzNH,
	})

	katInzRouter.RegisterCmd(&dgc.Command{
		Name:        "vmg",
		Description: "I'll get the content from the given vmg ID, Master.",
		Usage:       "vmg <page ID>",
		Example:     "vmg 18172",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katInzVMG,
	})

	katInzRouter.RegisterCmd(&dgc.Command{
		Name:        "ck",
		Description: "I'll get the content from the given ck101.com URL, Master.",
		Usage:       "ck <CK101 URL>",
		Example:     "ck https://ck101.com/thread-5369086-1-1.html?ref=index_starcontent",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katInzCK101,
	})

	katInzRouter.RegisterCmd(&dgc.Command{
		Name:        "yt",
		Description: "I'll get the content from the given YouTube URL, Master.",
		Usage:       "yt <YouTube URL>",
		Example:     "yt https://www.youtube.com/watch?v=kTJbE3sfvlI",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katInzYTDL,
	})

	katInzRouter.RegisterCmd(&dgc.Command{
		Name:        "lastsender",
		Description: "I'll get the last sender data from my memory, Master.",
		Usage:       "lastsender",
		Example:     "lastsender",
		IgnoreCase:  true,
		RateLimiter: dgc.NewRateLimiter(10*time.Second, 1*time.Second, ratelimitReply),
		Handler:     katInzShowLastSender,
	})

	// Initialize the router for Katherye Mondstadt
	katInzRouter.Initialize(katInzSession)

	var (
		statusInt   int
		statusSlice []string
	)
	statusInt = 0
	statusSlice = []string{"dnd", "idle", "online"}

	// Bot support with separated goroutines
	go func() {
		for range time.Tick(1 * time.Second) {
			setActivityText := discordgo.Activity{
				Name: maidsanWatchCurrentUser,
				Type: 3,
			}

			botStatusData := discordgo.UpdateStatusData{
				Activities: []*discordgo.Activity{&setActivityText},
				Status:     statusSlice[statusInt],
				AFK:        true,
			}
			maidsanSession.UpdateStatusComplex(botStatusData)
			maidchanSession.UpdateStatusComplex(botStatusData)
			katInzSession.UpdateStatusComplex(botStatusData)

			if statusInt != 2 {
				statusInt++
			} else {
				statusInt = 0
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// autocheck all emojis from the guilds the bot is in
	go func() {
		// clear slices
		maidsanEmojiInfo = nil
		customEmojiSlice = nil

		// get guild list
		getGuilds, err := maidsanSession.UserGuilds(100, "", "")
		if err != nil {
			fmt.Println(" [getGuilds] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		for guildIdx := range getGuilds {

			// Check the available emoji list
			getEmoji, err := maidsanSession.GuildEmojis(getGuilds[guildIdx].ID)
			if err != nil {
				fmt.Println(" [getEmoji] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			for idxEmoji := range getEmoji {
				maidsanEmojiInfo = append(maidsanEmojiInfo, fmt.Sprintf("\n%v —— %v —— %v —— %v —— %v", getEmoji[idxEmoji].Name, getEmoji[idxEmoji].ID, getEmoji[idxEmoji].Animated, getGuilds[guildIdx].Name, getGuilds[guildIdx].ID))

				customEmojiSlice = append(customEmojiSlice, fmt.Sprintf("%v:%v", getEmoji[idxEmoji].Name, getEmoji[idxEmoji].ID))
			}
		}
	}()

	// autocheck ban list from KokonattoMilku guild
	go func() {
		for range time.Tick(1 * time.Second) {

			// Check KokonattoMilku guild ban list
			getBanList, err := maidsanSession.GuildBans(kokonattomilkuGuildID)
			if err != nil {
				fmt.Println(" [getBanList] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			maidsanBanList = nil
			for idxBans := range getBanList {
				maidsanBanList = append(maidsanBanList, fmt.Sprintf("\n%v#%v : %v : %v", getBanList[idxBans].User.Username, getBanList[idxBans].User.Discriminator, getBanList[idxBans].User.ID, getBanList[idxBans].Reason))
			}

			time.Sleep(60 * time.Second)
		}
	}()

	// autocheck undercover mod list from KokonattoMilku guild
	go func() {
		for range time.Tick(1 * time.Second) {

			// convert string to slice
			convIDSlice := strings.Split(ucoverModsDB, ":")
			ucoverNewAdded = nil
			ucoverNewAdded = append(ucoverNewAdded, convIDSlice...)

			ucoverInfo = nil
			for umodIdx1, umodID1 := range convIDSlice {
				userData, err := katMondstadtSess.User(umodID1)
				if err != nil {
					fmt.Println(" [userData] ", err)
					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				// Reformat user data before printed out
				ucoverUsername = userData.Username + "#" + userData.Discriminator
				ucoveruserID = userData.ID

				// append mods info from database to slice
				ucoverInfo = append(ucoverInfo, fmt.Sprintf("\n%v —— %v —— U-%v —— From Database", ucoverUsername, ucoveruserID, umodIdx1))

			}

			time.Sleep(60 * time.Second)
		}
	}()

	// Katheryne Inazuma goroutines
	go func() {
		// Get the latest blocklist for dnscrypt
		getBlocklist, err := httpclient.Get("https://raw.githubusercontent.com/notracking/hosts-blocklists/master/dnscrypt-proxy/dnscrypt-proxy.blacklist.txt")
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		bodyBlocklist, err := ioutil.ReadAll(bufio.NewReader(getBlocklist.Body))
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		readCustomBlacklist, err := afero.ReadFile(osFS, "./customblacklist-galpt.txt")
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		katInzBlacklistReadable = fmt.Sprintf("\n%v\n", string(bodyBlocklist))
		katInzBlacklist = strings.Split(string(bodyBlocklist), "\n")

		if strings.Contains(string(readCustomBlacklist), ":") {
			katInzCustomBlacklist = strings.Split(string(readCustomBlacklist), ":")
			katInzBlacklist = append(katInzBlacklist, katInzCustomBlacklist...)
			katInzCustomBlacklistReadable = strings.ReplaceAll(string(readCustomBlacklist), ":", "\n")
		} else {
			katInzCustomBlacklistReadable = fmt.Sprintf("\n%v\n", string(readCustomBlacklist))
			katInzBlacklist = append(katInzBlacklist, katInzCustomBlacklist...)
		}
	}()

}

// react with the available server emojis
func emojiReactions(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.Author.ID == maidsanID {
		return
	} else if m.Author.ID == maidchanID {
		return
	}

	// React with ganyustare emoji
	// if the m.Content contains "geez" word
	if strings.Contains(strings.ToLower(m.Content), "geez") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "ganyustare:903098908966785024")
	} else if strings.Contains(strings.ToLower(m.Content), "<:ganyustare:903098908966785024>") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "ganyustare:903098908966785024")
	}

}

// Maid-san's emoji reactions handler
func maidsanEmojiReact(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.Author.ID == katheryneInazumaID {
		return
	} else if m.Author.ID == maidchanID {
		return
	}

	customEmojiDetected = false

	// Reply with custom emoji if the message contains the keyword
	for currIdx := range maidsanEmojiInfo {
		replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
		replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
		replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

		if strings.EqualFold(replysplitEmojiInfo[0], strings.ToLower(m.Content)) {
			customEmojiDetected = true
			if replysplitEmojiInfo[2] != "false" {
				customEmojiReply = fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1])
			} else {
				customEmojiReply = fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1])
			}
		}
	}

	if customEmojiDetected {
		s.ChannelMessageSend(m.ChannelID, customEmojiReply)
	} else {
		katMondstadtSess.MessageReactionAdd(m.ChannelID, m.ID, customEmojiSlice[customEmojiIdx])
		katLiyueSess.MessageReactionAdd(m.ChannelID, m.ID, customEmojiSlice[customEmojiIdx])
		katInazumaSess.MessageReactionAdd(m.ChannelID, m.ID, customEmojiSlice[customEmojiIdx])
		if customEmojiIdx == (len(customEmojiSlice) - 1) {
			customEmojiIdx = 0
		} else {
			customEmojiIdx++
		}
	}

}

// Maid-san's handle to auto-check for banned words & auto-add roles
func maidsanAutoCheck(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.Author.ID == maidchanID {
		return
	} else if m.Author.ID == katheryneInazumaID {
		return
	} else if m.Author.ID == staffID[0] {
		maidsanWatchCurrentUser = maidsanWatchPreviousUser
	} else {
		maidsanWatchCurrentUser = m.Author.Username + "#" + m.Author.Discriminator
		maidsanWatchPreviousUser = m.Author.Username + "#" + m.Author.Discriminator
	}

	// Get channel last message IDs
	senderUserID := m.Author.ID
	senderUsername := m.Author.Username + "#" + m.Author.Discriminator
	maidsanLastMsgChannelID = m.ChannelID

	maidsanLastMsgID = m.ID
	maidsanLowercaseLastMsg = strings.ToLower(m.Content)

	// Check if it's a new member or not
	if maidsanLastMsgChannelID == welcomeradarChannelID {
		maidsanWelcomeMsg = fmt.Sprintf("A new member! Welcome <@!%v> 👋", senderUserID)
		s.ChannelMessageSend(welcomeradarChannelID, maidsanWelcomeMsg)
	}

	// Add default roles for all members
	// KokoMember
	s.GuildMemberRoleAdd(kokonattomilkuGuildID, senderUserID, "894892275363115008")

	if bannedWordsDetector.IsProfane(maidsanLowercaseLastMsg) {

		for idx, banWords := range bannedWords {
			if strings.Contains(maidsanLowercaseLastMsg, banWords) {
				maidsanLowercaseLastMsg = strings.ReplaceAll(maidsanLowercaseLastMsg, bannedWords[idx], " [EDITED] ")
			}
		}

		scanLinks := xurlsRelaxed.FindAllString(maidsanLowercaseLastMsg, -1)

		katInzBlacklistLinkDetected = false
		for atIdx := range katInzBlacklist {
			for linkIdx := range scanLinks {
				if strings.EqualFold(scanLinks[linkIdx], strings.ToLower(katInzBlacklist[atIdx])) {
					maidsanLowercaseLastMsg = strings.ReplaceAll(maidsanLowercaseLastMsg, katInzBlacklist[atIdx], " [EDITED] ")
					katInzBlacklistLinkDetected = true
				}
			}
		}
		maidsanEditedLastMsg = maidsanLowercaseLastMsg

		if katInzBlacklistLinkDetected {
			// Create the embed templates
			senderField := discordgo.MessageEmbedField{
				Name:   "Sender",
				Value:  fmt.Sprintf("<@%v>", senderUserID),
				Inline: true,
			}
			senderUserIDField := discordgo.MessageEmbedField{
				Name:   "User ID",
				Value:  fmt.Sprintf("%v", senderUserID),
				Inline: true,
			}
			reasonField := discordgo.MessageEmbedField{
				Name:   "Reason",
				Value:  "Blacklisted Links/Banned Words",
				Inline: true,
			}
			editedMsgField := discordgo.MessageEmbedField{
				Name:   "Edited Message",
				Value:  fmt.Sprintf("%v", maidsanEditedLastMsg),
				Inline: false,
			}
			messageFields := []*discordgo.MessageEmbedField{&senderField, &senderUserIDField, &reasonField, &editedMsgField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  fmt.Sprintf("Edited by %v ❤️", botName),
				Color:  0x4287f5,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			s.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
			s.ChannelMessageSendEmbed(maidsanLastMsgChannelID, &aoiEmbeds)
		} else {
			// Create the embed templates
			senderField := discordgo.MessageEmbedField{
				Name:   "Sender",
				Value:  fmt.Sprintf("<@%v>", senderUserID),
				Inline: true,
			}
			senderUserIDField := discordgo.MessageEmbedField{
				Name:   "User ID",
				Value:  fmt.Sprintf("%v", senderUserID),
				Inline: true,
			}
			reasonField := discordgo.MessageEmbedField{
				Name:   "Reason",
				Value:  "Banned Words",
				Inline: true,
			}
			editedMsgField := discordgo.MessageEmbedField{
				Name:   "Edited Message",
				Value:  fmt.Sprintf("%v", maidsanEditedLastMsg),
				Inline: false,
			}
			messageFields := []*discordgo.MessageEmbedField{&senderField, &senderUserIDField, &reasonField, &editedMsgField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  fmt.Sprintf("Edited by %v ❤️", botName),
				Color:  0x4287f5,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			s.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
			s.ChannelMessageSendEmbed(maidsanLastMsgChannelID, &aoiEmbeds)
		}

		// Reformat user data before printed out
		userAvatar := m.Author.Avatar
		userisBot := fmt.Sprintf("%v", m.Author.Bot)
		userAccType := ""
		userAvaEmbedImgURL := ""

		// Check whether the user's avatar type is GIF or not
		if strings.Contains(userAvatar, "a_") {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".gif?size=4096"
		} else {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".jpg?size=4096"
		}

		// Check the user's account type
		if userisBot == "true" {
			userAccType = "Bot Account"
		} else {
			userAccType = "Standard User Account"
		}

		// copy logs to Maid-san's memory
		maidsanTranslatedMsg = fmt.Sprintf("https://translate.google.com/?sl=auto&tl=en&text=%v&op=translate", url.QueryEscape(maidsanEditedLastMsg))

		maidsanLogsTemplate = fmt.Sprintf("\n •===========================• \n • Timestamp: %v \n •===========================• \n \n Username: %v \n User ID: %v \n Profile Picture: %v \n Account Type: %v \n Message ID: %v \n Message:\n%v \n Translation:\n%v \n\n", time.Now().UTC().Format(time.RFC850), senderUsername, senderUserID, userAvaEmbedImgURL, userAccType, m.ID, maidsanEditedLastMsg, maidsanTranslatedMsg)

		lastMsgTimestamp = fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850))
		lastMsgUsername = fmt.Sprintf("%v", senderUsername)
		lastMsgUserID = fmt.Sprintf("%v", senderUserID)
		lastMsgpfp = fmt.Sprintf("%v", userAvaEmbedImgURL)
		lastMsgAccType = fmt.Sprintf("%v", userAccType)
		lastMsgID = fmt.Sprintf("%v", m.ID)
		lastMsgContent = fmt.Sprintf("%v", maidsanEditedLastMsg)
		lastMsgTranslation = fmt.Sprintf("%v", maidsanTranslatedMsg)

		if len(maidsanLogs) != maidsanLogsLimit {
			maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
			timestampLogs = append(timestampLogs, lastMsgTimestamp)
			useridLogs = append(useridLogs, lastMsgUserID)
			profpicLogs = append(profpicLogs, lastMsgpfp)
			acctypeLogs = append(acctypeLogs, lastMsgAccType)
			msgidLogs = append(msgidLogs, lastMsgID)
			msgLogs = append(msgLogs, lastMsgContent)
			translateLogs = append(translateLogs, lastMsgTranslation)
		} else {
			maidsanLogs = nil
			timestampLogs = nil
			useridLogs = nil
			profpicLogs = nil
			acctypeLogs = nil
			msgidLogs = nil
			msgLogs = nil
			translateLogs = nil
			maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
			timestampLogs = append(timestampLogs, lastMsgTimestamp)
			useridLogs = append(useridLogs, lastMsgUserID)
			profpicLogs = append(profpicLogs, lastMsgpfp)
			acctypeLogs = append(acctypeLogs, lastMsgAccType)
			msgidLogs = append(msgidLogs, lastMsgID)
			msgLogs = append(msgLogs, lastMsgContent)
			translateLogs = append(translateLogs, lastMsgTranslation)
		}
	} else {

		scanLinks := xurlsRelaxed.FindAllString(maidsanLowercaseLastMsg, -1)

		katInzBlacklistLinkDetected = false
		for atIdx := range katInzBlacklist {
			for linkIdx := range scanLinks {
				if strings.EqualFold(scanLinks[linkIdx], strings.ToLower(katInzBlacklist[atIdx])) {
					maidsanLowercaseLastMsg = strings.ReplaceAll(maidsanLowercaseLastMsg, katInzBlacklist[atIdx], " [EDITED] ")
					katInzBlacklistLinkDetected = true
				}
			}
		}
		maidsanEditedLastMsg = maidsanLowercaseLastMsg

		if katInzBlacklistLinkDetected {
			// Create the embed templates
			senderField := discordgo.MessageEmbedField{
				Name:   "Sender",
				Value:  fmt.Sprintf("<@%v>", senderUserID),
				Inline: true,
			}
			senderUserIDField := discordgo.MessageEmbedField{
				Name:   "User ID",
				Value:  fmt.Sprintf("%v", senderUserID),
				Inline: true,
			}
			reasonField := discordgo.MessageEmbedField{
				Name:   "Reason",
				Value:  "Blacklisted Links/Banned Words",
				Inline: true,
			}
			editedMsgField := discordgo.MessageEmbedField{
				Name:   "Edited Message",
				Value:  fmt.Sprintf("%v", maidsanEditedLastMsg),
				Inline: false,
			}
			messageFields := []*discordgo.MessageEmbedField{&senderField, &senderUserIDField, &reasonField, &editedMsgField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  fmt.Sprintf("Edited by %v ❤️", botName),
				Color:  0x4287f5,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			s.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
			s.ChannelMessageSendEmbed(maidsanLastMsgChannelID, &aoiEmbeds)

			// Reformat user data before printed out
			userAvatar := m.Author.Avatar
			userisBot := fmt.Sprintf("%v", m.Author.Bot)
			userAccType := ""
			userAvaEmbedImgURL := ""

			// Check whether the user's avatar type is GIF or not
			if strings.Contains(userAvatar, "a_") {
				userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".gif?size=4096"
			} else {
				userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".jpg?size=4096"
			}

			// Check the user's account type
			if userisBot == "true" {
				userAccType = "Bot Account"
			} else {
				userAccType = "Standard User Account"
			}

			// copy logs to Maid-san's memory
			maidsanTranslatedMsg = fmt.Sprintf("https://translate.google.com/?sl=auto&tl=en&text=%v&op=translate", url.QueryEscape(maidsanEditedLastMsg))

			maidsanLogsTemplate = fmt.Sprintf("\n •===========================• \n • Timestamp: %v \n •===========================• \n \n Username: %v \n User ID: %v \n Profile Picture: %v \n Account Type: %v \n Message ID: %v \n Message:\n%v \n Translation:\n%v \n\n", time.Now().UTC().Format(time.RFC850), senderUsername, senderUserID, userAvaEmbedImgURL, userAccType, m.ID, maidsanEditedLastMsg, maidsanTranslatedMsg)

			lastMsgTimestamp = fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850))
			lastMsgUsername = fmt.Sprintf("%v", senderUsername)
			lastMsgUserID = fmt.Sprintf("%v", senderUserID)
			lastMsgpfp = fmt.Sprintf("%v", userAvaEmbedImgURL)
			lastMsgAccType = fmt.Sprintf("%v", userAccType)
			lastMsgID = fmt.Sprintf("%v", m.ID)
			lastMsgContent = fmt.Sprintf("%v", maidsanEditedLastMsg)
			lastMsgTranslation = fmt.Sprintf("%v", maidsanTranslatedMsg)

			if len(maidsanLogs) != maidsanLogsLimit {
				maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
				timestampLogs = append(timestampLogs, lastMsgTimestamp)
				useridLogs = append(useridLogs, lastMsgUserID)
				profpicLogs = append(profpicLogs, lastMsgpfp)
				acctypeLogs = append(acctypeLogs, lastMsgAccType)
				msgidLogs = append(msgidLogs, lastMsgID)
				msgLogs = append(msgLogs, lastMsgContent)
				translateLogs = append(translateLogs, lastMsgTranslation)
			} else {
				maidsanLogs = nil
				timestampLogs = nil
				useridLogs = nil
				profpicLogs = nil
				acctypeLogs = nil
				msgidLogs = nil
				msgLogs = nil
				translateLogs = nil
				maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
				timestampLogs = append(timestampLogs, lastMsgTimestamp)
				useridLogs = append(useridLogs, lastMsgUserID)
				profpicLogs = append(profpicLogs, lastMsgpfp)
				acctypeLogs = append(acctypeLogs, lastMsgAccType)
				msgidLogs = append(msgidLogs, lastMsgID)
				msgLogs = append(msgLogs, lastMsgContent)
				translateLogs = append(translateLogs, lastMsgTranslation)
			}
		} else {

			// Reformat user data before printed out
			userAvatar := m.Author.Avatar
			userisBot := fmt.Sprintf("%v", m.Author.Bot)
			userAccType := ""
			userAvaEmbedImgURL := ""

			// Check whether the user's avatar type is GIF or not
			if strings.Contains(userAvatar, "a_") {
				userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".gif?size=4096"
			} else {
				userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".jpg?size=4096"
			}

			// Check the user's account type
			if userisBot == "true" {
				userAccType = "Bot Account"
			} else {
				userAccType = "Standard User Account"
			}

			// copy logs to Maid-san's memory
			maidsanTranslatedMsg = fmt.Sprintf("https://translate.google.com/?sl=auto&tl=en&text=%v&op=translate", url.QueryEscape(m.Content))

			maidsanLogsTemplate = fmt.Sprintf("\n •===========================• \n • Timestamp: %v \n •===========================• \n \n Username: %v \n User ID: %v \n Profile Picture: %v \n Account Type: %v \n Message ID: %v \n Message:\n%v \n Translation:\n%v \n\n", time.Now().UTC().Format(time.RFC850), senderUsername, senderUserID, userAvaEmbedImgURL, userAccType, m.ID, maidsanEditedLastMsg, maidsanTranslatedMsg)

			lastMsgTimestamp = fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850))
			lastMsgUsername = fmt.Sprintf("%v", senderUsername)
			lastMsgUserID = fmt.Sprintf("%v", senderUserID)
			lastMsgpfp = fmt.Sprintf("%v", userAvaEmbedImgURL)
			lastMsgAccType = fmt.Sprintf("%v", userAccType)
			lastMsgID = fmt.Sprintf("%v", m.ID)
			lastMsgContent = fmt.Sprintf("%v", m.Content)
			lastMsgTranslation = fmt.Sprintf("%v", maidsanTranslatedMsg)

			if len(maidsanLogs) != maidsanLogsLimit {
				maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
				timestampLogs = append(timestampLogs, lastMsgTimestamp)
				useridLogs = append(useridLogs, lastMsgUserID)
				profpicLogs = append(profpicLogs, lastMsgpfp)
				acctypeLogs = append(acctypeLogs, lastMsgAccType)
				msgidLogs = append(msgidLogs, lastMsgID)
				msgLogs = append(msgLogs, lastMsgContent)
				translateLogs = append(translateLogs, lastMsgTranslation)
			} else {
				maidsanLogs = nil
				timestampLogs = nil
				useridLogs = nil
				profpicLogs = nil
				acctypeLogs = nil
				msgidLogs = nil
				msgLogs = nil
				translateLogs = nil
				maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
				timestampLogs = append(timestampLogs, lastMsgTimestamp)
				useridLogs = append(useridLogs, lastMsgUserID)
				profpicLogs = append(profpicLogs, lastMsgpfp)
				acctypeLogs = append(acctypeLogs, lastMsgAccType)
				msgidLogs = append(msgidLogs, lastMsgID)
				msgLogs = append(msgLogs, lastMsgContent)
				translateLogs = append(translateLogs, lastMsgTranslation)
			}
		}

	}
}

// Maid-chan's handle to auto-check for banned words & auto-add roles
func maidchanAutoCheck(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.Author.ID == maidsanID {
		maidchanLastMsgChannelID = m.ChannelID
		maidchanLastMsgID = m.ID
		if strings.Contains(m.Content, "A new member!") {
			s.ChannelMessageSend(welcomeradarChannelID, "Welcome~ \nPlease make sure to read the <#894462808736010250>. \nCheck these channels too so you don't miss anything. \n<#893139038138167316> <#893139006395678760> <#893140731848425492> <#893140762903072808>")
		} else {
			return
		}
	} else if m.Author.ID == katheryneInazumaID {
		return
	}

}

// Inazuma Katheryne's handle to auto-check for banned words & auto-add roles
func katInzAutoCheck(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.Author.ID == maidsanID {
		return
	} else if m.Author.ID == maidchanID {
		return
	}

	// Get channel last message IDs
	senderUserID := m.Author.ID

	// check if userID is one of the staff members
	for _, checkStaff := range staffID {
		if senderUserID == checkStaff {
			// Add roles for staff members
			// Kacho
			s.GuildMemberRoleAdd(kokonattomilkuGuildID, senderUserID, "893141284787736656")
			s.GuildMemberRoleAdd(kokonattomilkuBackupGuildID, senderUserID, "904497628887285802")
		}
	}

	// check if the userID is one of the maid bots
	for _, checkBots := range botID {
		if senderUserID == checkBots {
			// Add roles for maid bots
			// Blessed by Castella
			s.GuildMemberRoleAdd(kokonattomilkuGuildID, senderUserID, "899557703502946335")
			s.GuildMemberRoleAdd(kokonattomilkuBackupGuildID, senderUserID, "904497628887285803")
			// KokoMember
			s.GuildMemberRoleAdd(kokonattomilkuGuildID, senderUserID, "894892275363115008")
			s.GuildMemberRoleAdd(kokonattomilkuBackupGuildID, senderUserID, "904497628874682386")
		}
	}

	// check if the userID is in blacklistedID
	for _, blacklistedUser := range blacklistedID {
		if senderUserID == blacklistedUser {
			// Delete or kick that user from the server immediately with the reason
			s.GuildMemberDeleteWithReason(kokonattomilkuGuildID, senderUserID, "You've been blacklisted.")
			s.GuildMemberDeleteWithReason(kokonattomilkuBackupGuildID, senderUserID, "You've been blacklisted.")
		}
	}

}

// =========================================
// Available bot commands

// RateLimiter respond message
func ratelimitReply(ctx *dgc.Ctx) {
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "❤️")
	ctx.RespondText("Please let me rest for **10 seconds**, Master!")
}

// Get a brief information about the mentioned user
func getUserInfo(ctx *dgc.Ctx) {

	// Open the database with 1 GB index cache size
	db, err := badger.Open(badger.DefaultOptions("./db").WithIndexCacheSize(1 << 30))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {
		argument := arguments.Get(0)
		userMention := argument.AsUserMentionID()
		userData, err := ctx.Session.User(userMention)
		if err != nil {
			fmt.Println(" [userData] ", err)
			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
			return
		}

		// Reformat user data before printed out
		userUsername := userData.Username + "#" + userData.Discriminator
		userID := userData.ID
		userAvatar := userData.Avatar
		userisBot := fmt.Sprintf("%v", userData.Bot)
		userAccType := ""
		userAvatarURLFullSize := ""
		userAvaEmbedImgURL := ""

		// Check whether the user has image sharing data or not.
		dbAccess := db.Update(func(txn *badger.Txn) error {
			dbData, err := txn.Get([]byte(userID))
			if err != nil {

				// Check whether the user's avatar type is GIF or not
				if strings.Contains(userAvatar, "a_") {
					userAvatarURLFullSize = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".gif?size=4096"
					userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".gif?size=256"
				} else {
					userAvatarURLFullSize = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".jpg?size=4096"
					userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".jpg?size=256"
				}

				// Check the user's account type
				if userisBot == "true" {
					userAccType = "Bot Account"
				} else {
					userAccType = "Standard User Account"
				}

				// Create the embed templates
				usernameField := discordgo.MessageEmbedField{
					Name:   "Username",
					Value:  userUsername,
					Inline: true,
				}
				userIDField := discordgo.MessageEmbedField{
					Name:   "User ID",
					Value:  userID,
					Inline: true,
				}
				userAvatarField := discordgo.MessageEmbedField{
					Name:   "Profile Picture URL",
					Value:  userAvatarURLFullSize,
					Inline: false,
				}
				userAccTypeField := discordgo.MessageEmbedField{
					Name:   "Account Type",
					Value:  userAccType,
					Inline: true,
				}
				messageFields := []*discordgo.MessageEmbedField{&usernameField, &userIDField, &userAvatarField, &userAccTypeField}

				aoiEmbedFooter := discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
				}

				aoiEmbedThumbnail := discordgo.MessageEmbedThumbnail{
					URL: userAvaEmbedImgURL,
				}

				aoiEmbeds := discordgo.MessageEmbed{
					Title:     "About User",
					Color:     0x00D2FF,
					Thumbnail: &aoiEmbedThumbnail,
					Footer:    &aoiEmbedFooter,
					Fields:    messageFields,
				}

				ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)

			} else {

				dbValue := dbData.Value(func(userVal []byte) error {

					// Check whether the user's avatar type is GIF or not
					if strings.Contains(userAvatar, "a_") {
						userAvatarURLFullSize = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".gif?size=4096"
						userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".gif?size=256"
					} else {
						userAvatarURLFullSize = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".jpg?size=4096"
						userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".jpg?size=256"
					}

					// Check the user's account type
					if userisBot == "true" {
						userAccType = "Bot Account"
					} else {
						userAccType = "Standard User Account"
					}

					// Create the embed templates
					usernameField := discordgo.MessageEmbedField{
						Name:   "Username",
						Value:  userUsername,
						Inline: true,
					}
					userIDField := discordgo.MessageEmbedField{
						Name:   "User ID",
						Value:  userID,
						Inline: true,
					}
					userAvatarField := discordgo.MessageEmbedField{
						Name:   "Profile Picture URL",
						Value:  userAvatarURLFullSize,
						Inline: false,
					}
					userAccTypeField := discordgo.MessageEmbedField{
						Name:   "Account Type",
						Value:  userAccType,
						Inline: true,
					}
					messageFields := []*discordgo.MessageEmbedField{&usernameField, &userIDField, &userAvatarField, &userAccTypeField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbedThumbnail := discordgo.MessageEmbedThumbnail{
						URL: userAvaEmbedImgURL,
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:     "About User",
						Color:     0x00D2FF,
						Thumbnail: &aoiEmbedThumbnail,
						Footer:    &aoiEmbedFooter,
						Fields:    messageFields,
					}

					ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)

					return nil
				})
				if dbValue != nil {
					log.Fatal(dbValue)
				}

			}

			return nil
		})
		if dbAccess != nil {
			log.Fatal(err)
		}

	} else {
		ctx.RespondText("Please use the **!check** command properly, Master! \nType **!help check** if you need more information about it.")
	}

}

// Get current COVID-19 data for Indonesia country
func getCovidData(ctx *dgc.Ctx) {

	arguments := ctx.Arguments
	countryArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// countryArgs shouldn't be empty
	if len(countryArgs) != 0 {

		if countryArgs == "indonesia" {
			// Get covid-19 json data Indonesia
			covIndo, err := httpclient.Get("https://data.covid19.go.id/public/api/update.json")
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			bodyCovIndo, err := ioutil.ReadAll(covIndo.Body)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			// Indonesia - Reformat JSON before printed out
			indoCreatedVal := gjson.Get(string(bodyCovIndo), `update.penambahan.created`)
			indoPosVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_positif`)
			indoMeninggalVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_meninggal`)
			indoSembuhVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_sembuh`)
			indoDirawatVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_dirawat`)
			indoTotalPosVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_positif`)
			indoTotalMeninggalVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_meninggal`)
			indoTotalSembuhVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_sembuh`)
			indoTotalDirawatVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_dirawat`)

			// Create the embed templates
			createdField := discordgo.MessageEmbedField{
				Name:   "Date Created",
				Value:  indoCreatedVal.String(),
				Inline: true,
			}
			countryField := discordgo.MessageEmbedField{
				Name:   "Country",
				Value:  strings.ToUpper(countryArgs),
				Inline: true,
			}
			totalConfirmedField := discordgo.MessageEmbedField{
				Name:   "Total Confirmed",
				Value:  fmt.Sprintf("%v", indoTotalPosVal.Int()),
				Inline: true,
			}
			totalDeathsField := discordgo.MessageEmbedField{
				Name:   "Total Deaths",
				Value:  fmt.Sprintf("%v", indoTotalMeninggalVal.Int()),
				Inline: true,
			}
			totalRecoveredField := discordgo.MessageEmbedField{
				Name:   "Total Recovered",
				Value:  fmt.Sprintf("%v", indoTotalSembuhVal.Int()),
				Inline: true,
			}
			totalTreatedField := discordgo.MessageEmbedField{
				Name:   "Total Treated",
				Value:  fmt.Sprintf("%v", indoTotalDirawatVal.Int()),
				Inline: true,
			}
			additionalConfirmedField := discordgo.MessageEmbedField{
				Name:   "Additional Confirmed",
				Value:  fmt.Sprintf("%v", indoPosVal.Int()),
				Inline: true,
			}
			additionalDeathsField := discordgo.MessageEmbedField{
				Name:   "Additional Deaths",
				Value:  fmt.Sprintf("%v", indoMeninggalVal.Int()),
				Inline: true,
			}
			additionalRecoveredField := discordgo.MessageEmbedField{
				Name:   "Additional Recovered",
				Value:  fmt.Sprintf("%v", indoSembuhVal.Int()),
				Inline: true,
			}
			additionalTreatedField := discordgo.MessageEmbedField{
				Name:   "Additional Treated",
				Value:  fmt.Sprintf("%v", indoDirawatVal.Int()),
				Inline: true,
			}
			messageFields := []*discordgo.MessageEmbedField{&createdField, &countryField, &totalConfirmedField, &totalDeathsField, &totalRecoveredField, &totalTreatedField, &additionalConfirmedField, &additionalDeathsField, &additionalRecoveredField, &additionalTreatedField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  "Latest COVID-19 Data",
				Color:  0xE06666,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
			covIndo.Body.Close()
		} else {
			// Get covid-19 json data from a certain country
			// based on the user's argument
			urlCountry := "https://covid19.mathdro.id/api/countries/" + countryArgs
			covData, err := httpclient.Get(urlCountry)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			bodyCovData, err := ioutil.ReadAll(covData.Body)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			// Reformat JSON before printed out
			countryCreatedVal := gjson.Get(string(bodyCovData), `lastUpdate`)
			countryTotalPosVal := gjson.Get(string(bodyCovData), `confirmed.value`)
			countryTotalSembuhVal := gjson.Get(string(bodyCovData), `recovered.value`)
			countryTotalMeninggalVal := gjson.Get(string(bodyCovData), `deaths.value`)

			// Create the embed templates
			createdField := discordgo.MessageEmbedField{
				Name:   "Date Created",
				Value:  countryCreatedVal.String(),
				Inline: true,
			}
			countryField := discordgo.MessageEmbedField{
				Name:   "Country",
				Value:  strings.ToUpper(countryArgs),
				Inline: true,
			}
			totalConfirmedField := discordgo.MessageEmbedField{
				Name:   "Total Confirmed",
				Value:  fmt.Sprintf("%v", countryTotalPosVal.Int()),
				Inline: true,
			}
			totalDeathsField := discordgo.MessageEmbedField{
				Name:   "Total Deaths",
				Value:  fmt.Sprintf("%v", countryTotalMeninggalVal.Int()),
				Inline: true,
			}
			totalRecoveredField := discordgo.MessageEmbedField{
				Name:   "Total Recovered",
				Value:  fmt.Sprintf("%v", countryTotalSembuhVal.Int()),
				Inline: true,
			}
			messageFields := []*discordgo.MessageEmbedField{&createdField, &countryField, &totalConfirmedField, &totalDeathsField, &totalRecoveredField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  "Latest COVID-19 Data",
				Color:  0xE06666,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
			covData.Body.Close()
		}

	} else {
		ctx.RespondText("Please use the **!covid19** command properly, Master! \nType **!help covid19** if you need more information about it.")
	}

}

// Get realtime server status
func getServerStatus(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	statusArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only Creator-sama who has the permission
	if strings.Contains(userID, staffID[0]) {

		// statusArgs shouldn't be empty
		if len(statusArgs) != 0 {
			if statusArgs == "update" {
				// Get GIPerf changelog file and update the old data
				//getChangelog, err := httpclient.Get("https://x.galpt.xyz/shared/changelog.txt")
				//if err != nil {
				//	fmt.Println(" [ERROR] ", err)
				//}

				//readChangelog, err := ioutil.ReadAll(bufio.NewReader(getChangelog.Body))
				//if err != nil {
				//	fmt.Println(" [ERROR] ", err)
				//}
				//getChangelog.Body.Close()

				// Get GIPerf files SHA256
				readChangelog, err := afero.ReadFile(osFS, "./changelog.txt")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				readgiperfExe, err := afero.ReadFile(osFS, "./giperf.exe")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				giperfHash := sha512.Sum512(readgiperfExe)
				giperfChangelog = fmt.Sprintf("%v", string(readChangelog))
				giperfExeSHA512 = fmt.Sprintf("%v", hex.EncodeToString(giperfHash[:]))

				ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "The reports have been updated, Master!")
			} else if statusArgs == "push" {

				// init the loc
				loc, err := time.LoadLocation("Asia/Seoul")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				// Create the embed templates
				changelogSHA256Field := discordgo.MessageEmbedField{
					Name:   "SHA-512",
					Value:  fmt.Sprintf("```giperf.exe: %v```", giperfExeSHA512),
					Inline: false,
				}
				timeLastUpdateField := discordgo.MessageEmbedField{
					Name:   "Last Updated",
					Value:  fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850)),
					Inline: false,
				}
				timeTestedField := discordgo.MessageEmbedField{
					Name:   "Tested On",
					Value:  fmt.Sprintf("%v", time.Now().In(loc).Format(time.RFC850)),
					Inline: false,
				}
				changelogContentField := discordgo.MessageEmbedField{
					Name:   "What's New",
					Value:  fmt.Sprintf("%v", giperfChangelog),
					Inline: false,
				}
				messageFields := []*discordgo.MessageEmbedField{&changelogSHA256Field, &timeLastUpdateField, &timeTestedField, &changelogContentField}

				aoiEmbedFooter := discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
				}

				aoiEmbeds := discordgo.MessageEmbed{
					Title:  fmt.Sprintf("%v's Reports", botName),
					Color:  0xF6B26B,
					Footer: &aoiEmbedFooter,
					Fields: messageFields,
				}

				ctx.Session.ChannelMessageSendEmbed(updatesChannelID, &aoiEmbeds)
			}
		} else {

			runtime.ReadMemStats(&mem)
			timeSince := time.Since(duration)

			// Create the embed templates
			cpuCoresField := discordgo.MessageEmbedField{
				Name:   "Available CPU Cores",
				Value:  fmt.Sprintf("%v", runtime.NumCPU()),
				Inline: true,
			}
			osMemoryField := discordgo.MessageEmbedField{
				Name:   "Available OS Memory",
				Value:  fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte)),
				Inline: true,
			}
			timeElapsedField := discordgo.MessageEmbedField{
				Name:   "Time Elapsed",
				Value:  fmt.Sprintf("%v", timeSince),
				Inline: true,
			}
			messageFields := []*discordgo.MessageEmbedField{&cpuCoresField, &osMemoryField, &timeElapsedField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  fmt.Sprintf("%v's Reports", botName),
				Color:  0xF6B26B,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
		}

	} else {

		runtime.ReadMemStats(&mem)
		timeSince := time.Since(duration)

		// Create the embed templates
		cpuCoresField := discordgo.MessageEmbedField{
			Name:   "Available CPU Cores",
			Value:  fmt.Sprintf("%v", runtime.NumCPU()),
			Inline: true,
		}
		osMemoryField := discordgo.MessageEmbedField{
			Name:   "Available OS Memory",
			Value:  fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte)),
			Inline: true,
		}
		timeElapsedField := discordgo.MessageEmbedField{
			Name:   "Time Elapsed",
			Value:  fmt.Sprintf("%v", timeSince),
			Inline: true,
		}
		messageFields := []*discordgo.MessageEmbedField{&cpuCoresField, &osMemoryField, &timeElapsedField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  fmt.Sprintf("%v's Reports", botName),
			Color:  0xF6B26B,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
	}
}

// Set a specific role to the mentioned user
func setRole(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only Creator-sama who has the permission to use this command
	if strings.Contains(userID, staffID[0]) {

		// rawArgs shouldn't be empty
		if len(rawArgs) != 0 {
			argument := arguments.Get(0)
			roleType := arguments.Get(1)
			userMention := argument.AsUserMentionID()
			userData, err := ctx.Session.User(userMention)
			if err != nil {
				fmt.Println(" [userData] ", err)
				return
			}

			// Set role based on the given roleType
			if strings.Contains(roleType.Raw(), "updates") {
				ctx.Session.GuildMemberRoleAdd("893138943334297682", userData.ID, "894895194489622548")
			} else if strings.Contains(roleType.Raw(), "kokomember") {
				ctx.Session.GuildMemberRoleAdd("893138943334297682", userData.ID, "894892275363115008")
			} else if strings.Contains(roleType.Raw(), "releases") {
				ctx.Session.GuildMemberRoleAdd("893138943334297682", userData.ID, "894895372068065291")
			} else if strings.Contains(roleType.Raw(), "announcements") {
				ctx.Session.GuildMemberRoleAdd("893138943334297682", userData.ID, "894895170934435860")
			}

			// Send a confirmation message
			ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "I've added the role, Master!")

		} else {
			// Send a confirmation message
			ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "Please use the **!role** command properly, Master! \nType **!help role** if you need more information about it.")
		}

	} else {
		maidsanErrorMsg = fmt.Sprintf("`role %v`\nI'm sorry, Master. But only Creator-sama who's allowed to use this command.", rawArgs)
		ctx.RespondText(maidsanErrorMsg)
	}

}

// Get Maid-san to reply with the available emoji
func getEmojiMaidsan(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {
		customEmojiDetected = false

		// Reply with custom emoji if the message contains the keyword
		for currIdx := range maidsanEmojiInfo {
			replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
			replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
			replysplitEmojiInfo = strings.Split(replyremoveSpaces, ":")

			if strings.ToLower(replysplitEmojiInfo[0]) == rawArgs {
				customEmojiDetected = true
				if replysplitEmojiInfo[2] != "false" {
					customEmojiReply = fmt.Sprintf("<@!%v>:  <a:%v:%v>", userID, replysplitEmojiInfo[0], replysplitEmojiInfo[1])
				} else {
					customEmojiReply = fmt.Sprintf("<@!%v>:  <:%v:%v>", userID, replysplitEmojiInfo[0], replysplitEmojiInfo[1])
				}
			}
		}

		if customEmojiDetected {
			ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
			ctx.Session.ChannelMessageSend(maidsanLastMsgChannelID, customEmojiReply)
		}
	}
}

// Get Maid-san to the rules for you
func getRules(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only Creator-sama who has the permission
	if strings.Contains(userID, staffID[0]) {

		// Create the embed templates
		rulesField := discordgo.MessageEmbedField{
			Name:   "DON'Ts",
			Value:  fmt.Sprintf("%v", serverRules),
			Inline: true,
		}
		messageFields := []*discordgo.MessageEmbedField{&rulesField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  "The Rules",
			Color:  0x4287f5,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
		ctx.Session.ChannelMessageSendEmbed(maidsanLastMsgChannelID, &aoiEmbeds)
		katMondstadtSess.MessageReactionAdd(maidchanLastMsgChannelID, maidchanLastMsgID, "👍")
		katLiyueSess.MessageReactionAdd(maidchanLastMsgChannelID, maidchanLastMsgID, "👍")
		katInazumaSess.MessageReactionAdd(maidchanLastMsgChannelID, maidchanLastMsgID, "👍")

	} else {
		ctx.RespondText("I'm sorry, Master. But only Creator-sama who's allowed to use this command.")
	}
}

// Get Maid-san to update the blacklisted links
func addBlacklistLinks(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only staff members who can access this command
	staffDetected = false
	for _, isStaff := range staffID {
		if userID == isStaff {
			staffDetected = true
			// rawArgs shouldn't be empty
			if len(rawArgs) != 0 {
				// check if rawArgs contains multiple links or not
				if strings.Contains(rawArgs, ":") {
					katInzAddCustomBlacklist = strings.Split(strings.ToLower(rawArgs), ":")
					katInzCustomBlacklist = append(katInzCustomBlacklist, katInzAddCustomBlacklist...)
					katInzBlacklist = append(katInzBlacklist, katInzCustomBlacklist...)

					katInzNewAppended = strings.Join(katInzCustomBlacklist, ":")

					// ==================================
					// update the customblacklist-galpt.txt file
					createNewBlacklist, err := osFS.Create("customblacklist-galpt.txt")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					} else {
						// Write to the file
						writeNewBlacklist, err := createNewBlacklist.WriteString(katInzNewAppended)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						} else {
							// Close the file
							if err := createNewBlacklist.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							} else {
								fmt.Println()
								winLogs = fmt.Sprintf(" [DONE] customblacklist-galpt.txt has been created. \n >> Size: %v KB (%v MB)", (writeNewBlacklist / Kilobyte), (writeNewBlacklist / Megabyte))
								fmt.Println(winLogs)
							}
						}
					}

					ctx.RespondText("Blacklisted links have been updated, Master!\nCheck `https://x.galpt.xyz/blacklist` for the details.")
				} else {
					katInzCustomBlacklist = append(katInzCustomBlacklist, rawArgs)
					katInzBlacklist = append(katInzBlacklist, katInzCustomBlacklist...)

					katInzNewAppended = strings.Join(katInzCustomBlacklist, ":")

					// ==================================
					// update the customblacklist-galpt.txt file
					createNewBlacklist, err := osFS.Create("customblacklist-galpt.txt")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					} else {
						// Write to the file
						writeNewBlacklist, err := createNewBlacklist.WriteString(katInzNewAppended)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						} else {
							// Close the file
							if err := createNewBlacklist.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							} else {
								fmt.Println()
								winLogs = fmt.Sprintf(" [DONE] customblacklist-galpt.txt has been created. \n >> Size: %v KB (%v MB)", (writeNewBlacklist / Kilobyte), (writeNewBlacklist / Megabyte))
								fmt.Println(winLogs)
							}
						}
					}

					ctx.RespondText("Blacklisted links have been updated, Master!\nCheck `https://x.galpt.xyz/blacklist` for the details.")
				}

			}
		}
	}
}

// Get Maid-san to ban the given User ID
func banUser(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	guildID := ctx.Event.GuildID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only staff members who can access this command
	staffDetected = false
	noBanStaff = false
	for _, isStaff := range staffID {
		if userID == isStaff {
			staffDetected = true
			// rawArgs shouldn't be empty
			if len(rawArgs) != 0 {
				// check if rawArgs contains multiple links or not
				if strings.Contains(rawArgs, "add@@") {
					getBanData := strings.Split(strings.ToLower(rawArgs), "@@")

					if len(getBanData) == 3 {
						for _, protectStaff := range staffID {
							if getBanData[1] != protectStaff {
								noBanStaff = true
							}
						}

						if noBanStaff {
							// ban the given User ID
							// GuildBanCreateWithReason(guildID, userID, reason string, days int)
							ctx.Session.GuildBanCreateWithReason(guildID, getBanData[1], getBanData[2], 7)

							maidsanBanUserMsg = fmt.Sprintf("I've banned <@!%v>\nwith the following reason \n```%v```", getBanData[1], getBanData[2])

							ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
							ctx.RespondText(maidsanBanUserMsg)
						} else {
							ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
							ctx.RespondText("I'm sorry, Master. But Creator-sama didn't allow me to ban staff members.")
						}
					}

				} else if strings.Contains(rawArgs, "remove@@") {
					getDelBanData := strings.Split(strings.ToLower(rawArgs), "@@")

					if len(getDelBanData) == 2 {
						// unban the given User ID
						// GuildBanDelete(guildID, userID string)
						ctx.Session.GuildBanDelete(guildID, getDelBanData[1])

						maidsanBanUserMsg = fmt.Sprintf("I've unbanned <@!%v> and removed him/her from the ban list, Master.", getDelBanData[1])

						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						ctx.RespondText(maidsanBanUserMsg)
					}
				}

			}
		}
	}
}

// use Katherynes in chats
func sayKatherynes(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	katcode := arguments.Get(0).AsUserMentionID()
	channelID := arguments.Get(1).AsChannelMentionID()
	msgSplit := strings.Split(arguments.Raw(), "::")
	msgContent := msgSplit[1]
	editedRawArgs := ""
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only Creator-sama who has the permission
	if strings.Contains(userID, staffID[0]) {

		if katcode == maidsanID {

			// Reply with custom emoji if the message contains the keyword
			for currIdx := range maidsanEmojiInfo {
				replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
				replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
				replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

				if strings.Contains(strings.ToLower(msgContent), replysplitEmojiInfo[0]) {
					customEmojiDetected = true
					if replysplitEmojiInfo[2] != "false" {
						editedRawArgs = strings.ReplaceAll(strings.ToLower(msgContent), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

						katMondstadtSess.ChannelMessageSend(channelID, editedRawArgs)
						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						break // break here
					} else {
						editedRawArgs = strings.ReplaceAll(strings.ToLower(msgContent), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

						katMondstadtSess.ChannelMessageSend(channelID, editedRawArgs)
						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						break // break here
					}
				} else {
					katMondstadtSess.ChannelMessageSend(channelID, msgContent)
					ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
					break // break here
				}
			}

		} else if katcode == maidchanID {

			// Reply with custom emoji if the message contains the keyword
			for currIdx := range maidsanEmojiInfo {
				replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
				replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
				replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

				if strings.Contains(strings.ToLower(msgContent), replysplitEmojiInfo[0]) {
					customEmojiDetected = true
					if replysplitEmojiInfo[2] != "false" {
						editedRawArgs = strings.ReplaceAll(strings.ToLower(msgContent), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

						katLiyueSess.ChannelMessageSend(channelID, editedRawArgs)
						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						break // break here
					} else {
						editedRawArgs = strings.ReplaceAll(strings.ToLower(msgContent), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

						katLiyueSess.ChannelMessageSend(channelID, editedRawArgs)
						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						break // break here
					}
				} else {
					katLiyueSess.ChannelMessageSend(channelID, msgContent)
					ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
					break // break here
				}
			}

		} else if katcode == katheryneInazumaID {

			// Reply with custom emoji if the message contains the keyword
			for currIdx := range maidsanEmojiInfo {
				replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
				replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
				replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

				if strings.Contains(strings.ToLower(msgContent), replysplitEmojiInfo[0]) {
					customEmojiDetected = true
					if replysplitEmojiInfo[2] != "false" {
						editedRawArgs = strings.ReplaceAll(strings.ToLower(msgContent), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

						katInazumaSess.ChannelMessageSend(channelID, editedRawArgs)
						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						break // break here
					} else {
						editedRawArgs = strings.ReplaceAll(strings.ToLower(msgContent), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

						katInazumaSess.ChannelMessageSend(channelID, editedRawArgs)
						ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
						break // break here
					}
				} else {
					katInazumaSess.ChannelMessageSend(channelID, msgContent)
					ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
					break // break here
				}
			}

		}

	}
}

// Get Maid-san to wrap your message in a warning template
func warnMsg(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only staff members who can access this command
	staffDetected = false
	noBanStaff = false
	for _, isStaff := range staffID {
		if userID == isStaff {
			staffDetected = true
			// rawArgs shouldn't be empty
			if len(rawArgs) != 0 {
				maidsanWarnMsg = fmt.Sprintf("%v", rawArgs)
			}
		}
	}

	if staffDetected {
		// Get the sender information
		senderAvatar := ctx.Event.Message.Author.Avatar
		userAvaEmbedImgURL := ""

		// Check whether the user's avatar type is GIF or not
		if strings.Contains(senderAvatar, "a_") {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + senderAvatar + ".gif?size=4096"
		} else {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + senderAvatar + ".jpg?size=4096"
		}

		// Create the embed templates
		senderUsernameField := discordgo.MessageEmbedField{
			Name:   "From",
			Value:  fmt.Sprintf("<@!%v>", userID),
			Inline: false,
		}
		warningMsgField := discordgo.MessageEmbedField{
			Name:   "Message",
			Value:  maidsanWarnMsg,
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&senderUsernameField, &warningMsgField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		logEmbedThumbnail := discordgo.MessageEmbedThumbnail{
			URL: userAvaEmbedImgURL,
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:     "Warning ⚠️",
			Color:     0xfffa69,
			Thumbnail: &logEmbedThumbnail,
			Footer:    &aoiEmbedFooter,
			Fields:    messageFields,
		}

		ctx.Session.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
		ctx.Session.ChannelMessageSendEmbed(maidsanLastMsgChannelID, &aoiEmbeds)
	}
}

// Get Maid-san to add/remove the given User ID
// to the undercover mod list.
func ucoverModsHandler(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// Only galpt has access to this command
	if userID == staffID[0] {

		// rawArgs shouldn't be empty
		if len(rawArgs) != 0 {
			// check if rawArgs contains multiple links or not
			if strings.Contains(rawArgs, "add@@") {
				getUserID := strings.Split(strings.ToLower(rawArgs), "@@")

				if strings.Contains(getUserID[1], ":") {
					splitIDs := strings.Split(getUserID[1], ":")

					for modIndex, modID := range splitIDs {

						// Create the embed templates
						notifField := discordgo.MessageEmbedField{
							Name:   "Congratulations!",
							Value:  "You are now an Undercover Moderator.",
							Inline: false,
						}
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", modID),
							Inline: false,
						}
						modIDField := discordgo.MessageEmbedField{
							Name:   "Undercover ID",
							Value:  fmt.Sprintf("U-%v", modIndex),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&notifField, &usernameField, &modIDField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "INFO",
							Color:  0x42f5ce,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						// Send notification to each added ID.
						// We create the private channel with the user who sent the message.
						channel, err := ctx.Session.UserChannelCreate(modID)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							return
						}
						// Then we send the message through the channel we created.
						_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						}

					}

					// combine old data with new data
					ucoverNewData = nil
					oldData := strings.Split(ucoverModsDB, ":")
					ucoverNewData = append(ucoverNewData, oldData...)
					ucoverNewData = append(ucoverNewData, splitIDs...)

					// ==================================
					// Create the 1-undercover-mods.txt file
					createNewList, err := osFS.Create("1-undercover-mods.txt")
					if err != nil {
						log.Fatal(err)
					} else {
						// Write to the file
						writeListFile, err := createNewList.WriteString(strings.Join(ucoverNewData, ":"))
						if err != nil {
							log.Fatal(err)
						} else {
							// Close the file
							if err := createNewList.Close(); err != nil {
								log.Fatal(err)
							} else {
								fmt.Println()
								winLogs = fmt.Sprintf(" [DONE] 1-undercover-mods.txt has been created. \n >> Size: %v KB (%v MB)", (writeListFile / Kilobyte), (writeListFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}
								ucoverModsDB = strings.Join(ucoverNewAdded, ":")

								splitModIDs := strings.Split(ucoverModsDB, ":")
								ucoverNewAdded = nil
								ucoverNewAdded = append(ucoverNewAdded, splitModIDs...)
							}
						}
					}

					ctx.RespondText("I've updated the Undercover Mods list, Master.")
				} else {

					// Create the embed templates
					notifField := discordgo.MessageEmbedField{
						Name:   "Congratulations!",
						Value:  "You are now an Undercover Moderator.",
						Inline: false,
					}
					usernameField := discordgo.MessageEmbedField{
						Name:   "Username",
						Value:  fmt.Sprintf("<@!%v>", getUserID[1]),
						Inline: false,
					}
					modIDField := discordgo.MessageEmbedField{
						Name:   "Undercover ID",
						Value:  fmt.Sprintf("U-%v", len(ucoverNewAdded)-1),
						Inline: false,
					}
					messageFields := []*discordgo.MessageEmbedField{&notifField, &usernameField, &modIDField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  "INFO",
						Color:  0x42f5ce,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					// Send notification to each added ID.
					// We create the private channel with the user who sent the message.
					channel, err := ctx.Session.UserChannelCreate(getUserID[1])
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
						return
					}
					// Then we send the message through the channel we created.
					_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					}

					// combine old data with new data
					ucoverNewData = nil
					if strings.Contains(ucoverModsDB, ":") {
						oldData := strings.Split(ucoverModsDB, ":")
						ucoverNewData = append(ucoverNewData, oldData...)
						ucoverNewData = append(ucoverNewData, getUserID[1])
					} else {
						ucoverNewData = append(ucoverNewData, ucoverNewAdded...)
						ucoverNewData = append(ucoverNewData, getUserID[1])
					}

					// ==================================
					// Create the 1-undercover-mods.txt file
					createNewList, err := osFS.Create("1-undercover-mods.txt")
					if err != nil {
						log.Fatal(err)
					} else {
						// Write to the file
						writeListFile, err := createNewList.WriteString(strings.Join(ucoverNewData, ":"))
						if err != nil {
							log.Fatal(err)
						} else {
							// Close the file
							if err := createNewList.Close(); err != nil {
								log.Fatal(err)
							} else {
								fmt.Println()
								winLogs = fmt.Sprintf(" [DONE] 1-undercover-mods.txt has been created. \n >> Size: %v KB (%v MB)", (writeListFile / Kilobyte), (writeListFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}
								ucoverModsDB = strings.Join(ucoverNewAdded, ":")

								splitModIDs := strings.Split(ucoverModsDB, ":")
								ucoverNewAdded = nil
								ucoverNewAdded = append(ucoverNewAdded, splitModIDs...)
							}
						}
					}

					maidsanBanUserMsg = fmt.Sprintf("I've added <@!%v> to the Undercover Mods list, Master.", getUserID[1])
					ctx.RespondText(maidsanBanUserMsg)
				}

			} else if strings.Contains(rawArgs, "remove@@") {
				getDelUcoverMod := strings.Split(strings.ToLower(rawArgs), "@@")

				for getModIndex, getModID := range ucoverNewAdded {
					if getDelUcoverMod[1] == getModID {
						if getModIndex == 0 {
							dataConvStr := strings.Join(ucoverNewAdded, ":")
							strReplace := strings.ReplaceAll(dataConvStr, fmt.Sprintf("%v:", getModID), "")
							dataConvSlice := strings.Split(strReplace, ":")
							ucoverNewAdded = nil
							ucoverNewAdded = append(ucoverNewAdded, dataConvSlice...)

							// Create the embed templates
							notifField := discordgo.MessageEmbedField{
								Name:   "Notice",
								Value:  "You've been removed from Undercover Moderator list.",
								Inline: false,
							}
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", getDelUcoverMod[1]),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", len(ucoverNewAdded)-1),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&notifField, &usernameField, &modIDField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "INFO",
								Color:  0x42f5ce,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to the given User ID.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(getDelUcoverMod[1])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							// combine old data with new data
							ucoverNewData = nil
							if strings.Contains(ucoverModsDB, ":") {
								oldData := strings.Split(ucoverModsDB, ":")
								ucoverNewData = append(ucoverNewData, oldData...)
								ucoverNewData = append(ucoverNewData, getDelUcoverMod[1])
							} else {
								ucoverNewData = append(ucoverNewData, ucoverNewAdded...)
								ucoverNewData = append(ucoverNewData, getDelUcoverMod[1])
							}

							// ==================================
							// Create the 1-undercover-mods.txt file
							createNewList, err := osFS.Create("1-undercover-mods.txt")
							if err != nil {
								log.Fatal(err)
							} else {
								// Write to the file
								writeListFile, err := createNewList.WriteString(strings.Join(ucoverNewData, ":"))
								if err != nil {
									log.Fatal(err)
								} else {
									// Close the file
									if err := createNewList.Close(); err != nil {
										log.Fatal(err)
									} else {
										fmt.Println()
										winLogs = fmt.Sprintf(" [DONE] 1-undercover-mods.txt has been created. \n >> Size: %v KB (%v MB)", (writeListFile / Kilobyte), (writeListFile / Megabyte))
										fmt.Println(winLogs)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
										}
										ucoverModsDB = strings.Join(ucoverNewAdded, ":")

										splitModIDs := strings.Split(ucoverModsDB, ":")
										ucoverNewAdded = nil
										ucoverNewAdded = append(ucoverNewAdded, splitModIDs...)
									}
								}
							}

							maidsanBanUserMsg = fmt.Sprintf("I've removed <@!%v> from the Undercover Mods list, Master.", getDelUcoverMod[1])
							ctx.RespondText(maidsanBanUserMsg)

							break // break here
						} else if getModIndex == len(ucoverNewAdded)-1 {
							dataConvStr := strings.Join(ucoverNewAdded, ":")
							strReplace := strings.ReplaceAll(dataConvStr, fmt.Sprintf(":%v", getModID), "")
							dataConvSlice := strings.Split(strReplace, ":")
							ucoverNewAdded = nil
							ucoverNewAdded = append(ucoverNewAdded, dataConvSlice...)

							// Create the embed templates
							notifField := discordgo.MessageEmbedField{
								Name:   "Notice",
								Value:  "You've been removed from Undercover Moderator list.",
								Inline: false,
							}
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", getDelUcoverMod[1]),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", len(ucoverNewAdded)-1),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&notifField, &usernameField, &modIDField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "INFO",
								Color:  0x42f5ce,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to the given User ID.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(getDelUcoverMod[1])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							// combine old data with new data
							ucoverNewData = nil
							if strings.Contains(ucoverModsDB, ":") {
								oldData := strings.Split(ucoverModsDB, ":")
								ucoverNewData = append(ucoverNewData, oldData...)
								ucoverNewData = append(ucoverNewData, getDelUcoverMod[1])
							} else {
								ucoverNewData = append(ucoverNewData, ucoverNewAdded...)
								ucoverNewData = append(ucoverNewData, getDelUcoverMod[1])
							}

							// ==================================
							// Create the 1-undercover-mods.txt file
							createNewList, err := osFS.Create("1-undercover-mods.txt")
							if err != nil {
								log.Fatal(err)
							} else {
								// Write to the file
								writeListFile, err := createNewList.WriteString(strings.Join(ucoverNewData, ":"))
								if err != nil {
									log.Fatal(err)
								} else {
									// Close the file
									if err := createNewList.Close(); err != nil {
										log.Fatal(err)
									} else {
										fmt.Println()
										winLogs = fmt.Sprintf(" [DONE] 1-undercover-mods.txt has been created. \n >> Size: %v KB (%v MB)", (writeListFile / Kilobyte), (writeListFile / Megabyte))
										fmt.Println(winLogs)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
										}
										ucoverModsDB = strings.Join(ucoverNewAdded, ":")

										splitModIDs := strings.Split(ucoverModsDB, ":")
										ucoverNewAdded = nil
										ucoverNewAdded = append(ucoverNewAdded, splitModIDs...)
									}
								}
							}

							maidsanBanUserMsg = fmt.Sprintf("I've removed <@!%v> from the Undercover Mods list, Master.", getDelUcoverMod[1])
							ctx.RespondText(maidsanBanUserMsg)

							break // break here
						} else {
							dataConvStr := strings.Join(ucoverNewAdded, ":")
							strReplace := strings.ReplaceAll(dataConvStr, fmt.Sprintf(":%v:", getModID), ":")
							dataConvSlice := strings.Split(strReplace, ":")
							ucoverNewAdded = nil
							ucoverNewAdded = append(ucoverNewAdded, dataConvSlice...)

							// Create the embed templates
							notifField := discordgo.MessageEmbedField{
								Name:   "Notice",
								Value:  "You've been removed from Undercover Moderator list.",
								Inline: false,
							}
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", getDelUcoverMod[1]),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", len(ucoverNewAdded)-1),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&notifField, &usernameField, &modIDField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "INFO",
								Color:  0x42f5ce,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to the given User ID.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(getDelUcoverMod[1])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							// combine old data with new data
							ucoverNewData = nil
							if strings.Contains(ucoverModsDB, ":") {
								oldData := strings.Split(ucoverModsDB, ":")
								ucoverNewData = append(ucoverNewData, oldData...)
								ucoverNewData = append(ucoverNewData, dataConvSlice...)
							} else {
								ucoverNewData = append(ucoverNewData, ucoverNewAdded...)
								ucoverNewData = append(ucoverNewData, dataConvSlice...)
							}

							// ==================================
							// Create the 1-undercover-mods.txt file
							createNewList, err := osFS.Create("1-undercover-mods.txt")
							if err != nil {
								log.Fatal(err)
							} else {
								// Write to the file
								writeListFile, err := createNewList.WriteString(strings.Join(ucoverNewData, ":"))
								if err != nil {
									log.Fatal(err)
								} else {
									// Close the file
									if err := createNewList.Close(); err != nil {
										log.Fatal(err)
									} else {
										fmt.Println()
										winLogs = fmt.Sprintf(" [DONE] 1-undercover-mods.txt has been created. \n >> Size: %v KB (%v MB)", (writeListFile / Kilobyte), (writeListFile / Megabyte))
										fmt.Println(winLogs)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
										}
										ucoverModsDB = strings.Join(ucoverNewAdded, ":")

										splitModIDs := strings.Split(ucoverModsDB, ":")
										ucoverNewAdded = nil
										ucoverNewAdded = append(ucoverNewAdded, splitModIDs...)
									}
								}
							}

							maidsanBanUserMsg = fmt.Sprintf("I've removed <@!%v> from the Undercover Mods list, Master.", getDelUcoverMod[1])
							ctx.RespondText(maidsanBanUserMsg)

							break // break here
						}
					}
				}

			}

		}

	}

}

// Undercover Mods are allowed to delete inappropriate messages.
func ucoverModsDelMsg(ctx *dgc.Ctx) {

	delmsgRelax := xurls.Relaxed()
	channelID := ctx.Event.ChannelID
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")
	maidsanWatchCurrentUser = "@everyone" // to keep undermods hidden

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// Check userID in ucoverNewAdded slice
		for chkIdx, chkID := range ucoverNewAdded {
			if userID == chkID {

				scanLinks := delmsgRelax.FindAllString(rawArgs, -1)
				splitChanMsgIDs := strings.Split(scanLinks[0], fmt.Sprintf("channels/%v/", channelID))
				getDelMsgData := strings.Split(splitChanMsgIDs[1], "/")
				ctx.Session.ChannelMessageDelete(getDelMsgData[0], getDelMsgData[1])
				maidsanBanUserMsg = fmt.Sprintf("I've deleted MessageID(%v) from <#%v>, Master.", getDelMsgData[1], getDelMsgData[0])
				ctx.RespondText(maidsanBanUserMsg)

				// Create the embed templates
				usernameField := discordgo.MessageEmbedField{
					Name:   "Username",
					Value:  fmt.Sprintf("<@!%v>", userID),
					Inline: false,
				}
				modIDField := discordgo.MessageEmbedField{
					Name:   "Undercover ID",
					Value:  fmt.Sprintf("U-%v", chkIdx),
					Inline: false,
				}
				delmsgIDField := discordgo.MessageEmbedField{
					Name:   "Deleted Message ID",
					Value:  fmt.Sprintf("%v", splitChanMsgIDs[1]),
					Inline: false,
				}
				delmsgChanField := discordgo.MessageEmbedField{
					Name:   "Deleted Message Channel",
					Value:  fmt.Sprintf("<#%v>", splitChanMsgIDs[0]),
					Inline: false,
				}
				messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

				aoiEmbedFooter := discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
				}

				aoiEmbeds := discordgo.MessageEmbed{
					Title:  "Usage Information",
					Color:  0x32a852,
					Footer: &aoiEmbedFooter,
					Fields: messageFields,
				}

				// Send notification to galpt.
				// We create the private channel with the user who sent the message.
				channel, err := ctx.Session.UserChannelCreate(staffID[0])
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}
				// Then we send the message through the channel we created.
				_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
				}

				break // break here
			}
		}
	}

}

// Undercover Mods (Phase 1 & Phase 2) are allowed
// to send messages using Katheryne Mondstadt.
func ucoverModsMsgGenEN(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	editedRawArgs := ""
	copyeditedRawArgs := ""
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")
	maidsanWatchCurrentUser = "@everyone" // to keep undermods hidden

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// Check userID in ucoverNewAdded slice
		for chkIdx, chkID := range ucoverNewAdded {
			if userID == chkID {
				// Reply with custom emoji if the message contains the keyword
				for currIdx := range maidsanEmojiInfo {
					replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
					replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
					replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

					if strings.Contains(strings.ToLower(rawArgs), replysplitEmojiInfo[0]) {
						customEmojiDetected = true
						if replysplitEmojiInfo[2] != "false" {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(genENChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", genENChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						} else {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(genENChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", genENChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						}
					} else {
						// copy content
						copyeditedRawArgs = rawArgs

						ctx.Session.ChannelMessageSend(genENChannelID, rawArgs)

						// Create the embed templates
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", userID),
							Inline: false,
						}
						modIDField := discordgo.MessageEmbedField{
							Name:   "Undercover ID",
							Value:  fmt.Sprintf("U-%v", chkIdx),
							Inline: false,
						}
						delmsgIDField := discordgo.MessageEmbedField{
							Name:   "Message Channel",
							Value:  fmt.Sprintf("<#%v>", genENChannelID),
							Inline: false,
						}
						delmsgChanField := discordgo.MessageEmbedField{
							Name:   "Message Content",
							Value:  fmt.Sprintf("%v", copyeditedRawArgs),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "Usage Information",
							Color:  0x3bbfbd,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						// Send notification to galpt.
						// We create the private channel with the user who sent the message.
						channel, err := ctx.Session.UserChannelCreate(staffID[0])
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							return
						}
						// Then we send the message through the channel we created.
						_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						}

						break // break here
					}
				}
			}
		}
	}

}

// Undercover Mods (Phase 1 & Phase 2) are allowed
// to send messages using Katheryne Mondstadt.
func ucoverModsMsgGenCN(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	editedRawArgs := ""
	copyeditedRawArgs := ""
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")
	maidsanWatchCurrentUser = "@everyone" // to keep undermods hidden

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// Check userID in ucoverNewAdded slice
		for chkIdx, chkID := range ucoverNewAdded {
			if userID == chkID {
				// Reply with custom emoji if the message contains the keyword
				for currIdx := range maidsanEmojiInfo {
					replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
					replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
					replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

					if strings.Contains(strings.ToLower(rawArgs), replysplitEmojiInfo[0]) {
						customEmojiDetected = true
						if replysplitEmojiInfo[2] != "false" {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(genCNChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", genCNChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						} else {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(genCNChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", genCNChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						}
					} else {
						// copy content
						copyeditedRawArgs = rawArgs

						ctx.Session.ChannelMessageSend(genCNChannelID, rawArgs)

						// Create the embed templates
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", userID),
							Inline: false,
						}
						modIDField := discordgo.MessageEmbedField{
							Name:   "Undercover ID",
							Value:  fmt.Sprintf("U-%v", chkIdx),
							Inline: false,
						}
						delmsgIDField := discordgo.MessageEmbedField{
							Name:   "Message Channel",
							Value:  fmt.Sprintf("<#%v>", genCNChannelID),
							Inline: false,
						}
						delmsgChanField := discordgo.MessageEmbedField{
							Name:   "Message Content",
							Value:  fmt.Sprintf("%v", copyeditedRawArgs),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "Usage Information",
							Color:  0x3bbfbd,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						// Send notification to galpt.
						// We create the private channel with the user who sent the message.
						channel, err := ctx.Session.UserChannelCreate(staffID[0])
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							return
						}
						// Then we send the message through the channel we created.
						_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						}

						break // break here
					}
				}
			}
		}
	}

}

// Undercover Mods (Phase 1 & Phase 2) are allowed
// to send messages using Katheryne Mondstadt.
func ucoverModsMsgGenRU(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	editedRawArgs := ""
	copyeditedRawArgs := ""
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")
	maidsanWatchCurrentUser = "@everyone" // to keep undermods hidden

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// Check userID in ucoverNewAdded slice
		for chkIdx, chkID := range ucoverNewAdded {
			if userID == chkID {
				// Reply with custom emoji if the message contains the keyword
				for currIdx := range maidsanEmojiInfo {
					replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
					replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
					replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

					if strings.Contains(strings.ToLower(rawArgs), replysplitEmojiInfo[0]) {
						customEmojiDetected = true
						if replysplitEmojiInfo[2] != "false" {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(genRUChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", genRUChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						} else {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(genRUChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", genRUChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						}
					} else {
						// copy content
						copyeditedRawArgs = rawArgs

						ctx.Session.ChannelMessageSend(genRUChannelID, rawArgs)

						// Create the embed templates
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", userID),
							Inline: false,
						}
						modIDField := discordgo.MessageEmbedField{
							Name:   "Undercover ID",
							Value:  fmt.Sprintf("U-%v", chkIdx),
							Inline: false,
						}
						delmsgIDField := discordgo.MessageEmbedField{
							Name:   "Message Channel",
							Value:  fmt.Sprintf("<#%v>", genRUChannelID),
							Inline: false,
						}
						delmsgChanField := discordgo.MessageEmbedField{
							Name:   "Message Content",
							Value:  fmt.Sprintf("%v", copyeditedRawArgs),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "Usage Information",
							Color:  0x3bbfbd,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						// Send notification to galpt.
						// We create the private channel with the user who sent the message.
						channel, err := ctx.Session.UserChannelCreate(staffID[0])
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							return
						}
						// Then we send the message through the channel we created.
						_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						}

						break // break here
					}
				}
			}
		}
	}

}

// Undercover Mods (Phase 1 & Phase 2) are allowed
// to send messages using Katheryne Mondstadt.
func ucoverModsMsgOffTopic(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	editedRawArgs := ""
	copyeditedRawArgs := ""
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")
	maidsanWatchCurrentUser = "@everyone" // to keep undermods hidden

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// Check userID in ucoverNewAdded slice
		for chkIdx, chkID := range ucoverNewAdded {
			if userID == chkID {
				// Reply with custom emoji if the message contains the keyword
				for currIdx := range maidsanEmojiInfo {
					replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
					replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
					replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

					if strings.Contains(strings.ToLower(rawArgs), replysplitEmojiInfo[0]) {
						customEmojiDetected = true
						if replysplitEmojiInfo[2] != "false" {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(offtopicChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", offtopicChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						} else {
							editedRawArgs = strings.ReplaceAll(strings.ToLower(rawArgs), replysplitEmojiInfo[0], fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1]))

							// copy content
							copyeditedRawArgs = editedRawArgs

							ctx.Session.ChannelMessageSend(offtopicChannelID, editedRawArgs)

							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							modIDField := discordgo.MessageEmbedField{
								Name:   "Undercover ID",
								Value:  fmt.Sprintf("U-%v", chkIdx),
								Inline: false,
							}
							delmsgIDField := discordgo.MessageEmbedField{
								Name:   "Message Channel",
								Value:  fmt.Sprintf("<#%v>", offtopicChannelID),
								Inline: false,
							}
							delmsgChanField := discordgo.MessageEmbedField{
								Name:   "Message Content",
								Value:  fmt.Sprintf("%v", copyeditedRawArgs),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Usage Information",
								Color:  0x3bbfbd,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							// Send notification to galpt.
							// We create the private channel with the user who sent the message.
							channel, err := ctx.Session.UserChannelCreate(staffID[0])
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
								return
							}
							// Then we send the message through the channel we created.
							_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}
							}

							break // break here
						}
					} else {
						// copy content
						copyeditedRawArgs = rawArgs

						ctx.Session.ChannelMessageSend(offtopicChannelID, rawArgs)

						// Create the embed templates
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", userID),
							Inline: false,
						}
						modIDField := discordgo.MessageEmbedField{
							Name:   "Undercover ID",
							Value:  fmt.Sprintf("U-%v", chkIdx),
							Inline: false,
						}
						delmsgIDField := discordgo.MessageEmbedField{
							Name:   "Message Channel",
							Value:  fmt.Sprintf("<#%v>", offtopicChannelID),
							Inline: false,
						}
						delmsgChanField := discordgo.MessageEmbedField{
							Name:   "Message Content",
							Value:  fmt.Sprintf("%v", copyeditedRawArgs),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "Usage Information",
							Color:  0x3bbfbd,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						// Send notification to galpt.
						// We create the private channel with the user who sent the message.
						channel, err := ctx.Session.UserChannelCreate(staffID[0])
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							return
						}
						// Then we send the message through the channel we created.
						_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						}

						break // break here
					}
				}
			}
		}
	}

}

// ======
// Handlers for Katheryne Inazuma
// ======
// KatInz will get the data from the given URL
func katInzGet(ctx *dgc.Ctx) {

	getRelax := xurls.Relaxed()
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		getImgs = nil
		getMaxRender = 1

		// support for getting all images on a webpage
		if strings.Contains(rawArgs, "img") {

			// get the link
			scanLinks := getRelax.FindAllString(rawArgs, -1)

			// Get the webpage data
			getPage, err := httpclient.Get(scanLinks[0])
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			bodyPage, err := ioutil.ReadAll(bufio.NewReader(getPage.Body))
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			scanLinks = nil
			scanLinks = getRelax.FindAllString(string(bodyPage), -1)

			for getCurrIdx := range scanLinks {

				for formatIdx := range getFileFormat {

					if strings.Contains(scanLinks[getCurrIdx], getFileFormat[formatIdx]) {

						// add only image links
						getImgs = append(getImgs, scanLinks[getCurrIdx])

						// Get the image and write it to memory
						getImg, err := httpclient.Get(fmt.Sprintf("%v", scanLinks[getCurrIdx]))
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							break
						}

						// convert http response to io.Reader
						bodyIMG, err := ioutil.ReadAll(bufio.NewReader(getImg.Body))
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							break
						}
						reader := bytes.NewReader(bodyIMG)

						// Send image thru DM.
						// We create the private channel with the user who sent the message.
						channel, err := ctx.Session.UserChannelCreate(userID)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							break
						}
						// Then we send the message through the channel we created.
						_, err = ctx.Session.ChannelFileSend(channel.ID, fmt.Sprintf("%v%v", getCurrIdx, getFileFormat[formatIdx]), reader)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							break
						}

						getImg.Body.Close()
						getMaxRender++

						// limit max render to 50 images
						if getMaxRender == 50 {
							break
						}
					}
				}

			}

			getPage.Body.Close()
			scanLinks = nil

			// send manga info to user via DM.
			// Create the embed templates.
			oriURLField := discordgo.MessageEmbedField{
				Name:   "Original URL",
				Value:  fmt.Sprintf("%v", rawArgs),
				Inline: false,
			}
			showURLField := discordgo.MessageEmbedField{
				Name:   "Total Images",
				Value:  fmt.Sprintf("%v", len(getImgs)),
				Inline: false,
			}
			messageFields := []*discordgo.MessageEmbedField{&oriURLField, &showURLField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  "GET Information",
				Color:  0x03fcad,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			// Send image thru DM.
			// We create the private channel with the user who sent the message.
			channel, err := ctx.Session.UserChannelCreate(userID)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
				return
			}
			// Then we send the message through the channel we created.
			_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
			}

			// add a quick reply for Go to Page 1
			ctx.Session.ChannelMessageSendReply(channel.ID, "**Go to Image 1**", ctx.Event.Reference())
		} else {

			// check whether user input contains the file name
			if strings.Contains(rawArgs, "@@") {

				// split file name from the URL
				splitInput := strings.Split(rawArgs, "@@")

				// check whether the data has been cached or not
				if splitInput[1] == katInzGETCachedURL {

					// get data from cache.
					// Create the embed templates
					filenameField := discordgo.MessageEmbedField{
						Name:   "File Name",
						Value:  fmt.Sprintf("%v", katInzGETCachedFileName),
						Inline: false,
					}
					oriURLField := discordgo.MessageEmbedField{
						Name:   "Original URL",
						Value:  fmt.Sprintf("%v", katInzGETCachedURL),
						Inline: false,
					}
					showURLField := discordgo.MessageEmbedField{
						Name:   "Data Location in Memory",
						Value:  fmt.Sprintf("https://x.castella.network/get/%v", katInzGETCachedFileName),
						Inline: false,
					}
					messageFields := []*discordgo.MessageEmbedField{&filenameField, &oriURLField, &showURLField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  "Data from Katheryne's Memory",
						Color:  0x34c0eb,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
				} else if splitInput[1] != katInzGETCachedURL {

					// fetch data directly from the given URL
					memFS.RemoveAll("./get/")
					memFS.MkdirAll("./get/", 0777)

					getDataFromURL, err := httpclient.Get(splitInput[1])
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					}

					// Detect suspicious links inside the response body
					bodyBytes, err := io.ReadAll(getDataFromURL.Body)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						// Close the response body
						getDataFromURL.Body.Close()

						return
					}

					bodyString := string(bodyBytes)
					editedGETData = bodyString
					for linkIdx := range katInzBlacklist {
						if strings.Contains(bodyString, katInzBlacklist[linkIdx]) {
							editedGETData = strings.ReplaceAll(bodyString, katInzBlacklist[linkIdx], "")
						}
					}

					// Create a new file based on the body
					createNewFile, err := memFS.Create(fmt.Sprintf("./get/%v", splitInput[0]))
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						// Close the response body
						getDataFromURL.Body.Close()

						return
					} else {
						// Write to the file
						writeNewFile, err := createNewFile.Write([]byte(editedGETData))
						if err != nil {
							fmt.Println(" [ERROR] ", err)
							getDataFromURL.Body.Close()

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							// Close the response body
							getDataFromURL.Body.Close()

							return
						} else {

							// Close the response body
							getDataFromURL.Body.Close()

							if err := createNewFile.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							} else {
								winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", splitInput[0], (writeNewFile / Kilobyte), (writeNewFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}

								katInzGETCachedFileName = splitInput[0]
								katInzGETCachedURL = splitInput[1]

								// get data from cache.
								// Create the embed templates
								filenameField := discordgo.MessageEmbedField{
									Name:   "File Name",
									Value:  fmt.Sprintf("%v", katInzGETCachedFileName),
									Inline: false,
								}
								oriURLField := discordgo.MessageEmbedField{
									Name:   "Original URL",
									Value:  fmt.Sprintf("%v", katInzGETCachedURL),
									Inline: false,
								}
								showURLField := discordgo.MessageEmbedField{
									Name:   "Data Location in Memory",
									Value:  fmt.Sprintf("https://x.castella.network/get/%v", katInzGETCachedFileName),
									Inline: false,
								}
								messageFields := []*discordgo.MessageEmbedField{&filenameField, &oriURLField, &showURLField}

								aoiEmbedFooter := discordgo.MessageEmbedFooter{
									Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
								}

								aoiEmbeds := discordgo.MessageEmbed{
									Title:  "Data Has Been Cached",
									Color:  0x34c0eb,
									Footer: &aoiEmbedFooter,
									Fields: messageFields,
								}

								ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
							}
						}
					}

				}
			} else {
				if rawArgs == katInzGETCachedFileName {

					// get data from cache.
					// Create the embed templates
					filenameField := discordgo.MessageEmbedField{
						Name:   "File Name",
						Value:  fmt.Sprintf("%v", katInzGETCachedFileName),
						Inline: false,
					}
					oriURLField := discordgo.MessageEmbedField{
						Name:   "Original URL",
						Value:  fmt.Sprintf("%v", katInzGETCachedURL),
						Inline: false,
					}
					showURLField := discordgo.MessageEmbedField{
						Name:   "Data Location in Memory",
						Value:  fmt.Sprintf("https://x.castella.network/get/%v", katInzGETCachedFileName),
						Inline: false,
					}
					messageFields := []*discordgo.MessageEmbedField{&filenameField, &oriURLField, &showURLField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  "Data from Katheryne's Memory",
						Color:  0x34c0eb,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
				}
			}
		}

	}

}

// KatInz's NH feature
func katInzNH(ctx *dgc.Ctx) {

	nhRelax := xurls.Relaxed()
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// get url param
		nhCode = rawArgs

		// fetch data directly from nhen server
		memFS.RemoveAll("./nh/")
		memFS.MkdirAll("./nh/", 0777)

		// Get the gallery ID
		nhGalleryID := fmt.Sprintf("https://nhentai.net/g/%v", nhCode)
		getGalleryID, err := httpclient.Get(nhGalleryID)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		bodyGalleryID, err := ioutil.ReadAll(bufio.NewReader(getGalleryID.Body))
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		scanPg1 := strings.Split(string(bodyGalleryID), "Pages:")
		scanPg2 := strings.Split(scanPg1[1], `<span class="name">`)
		scanPg3 := strings.Split(scanPg2[1], `</span></a></span></div><div class="tag-container field-name">`)
		nhTotalPage, err = strconv.Atoi(scanPg3[0])
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		scanGalleryID := nhRelax.FindAllString(string(bodyGalleryID), -1)
		for nhlinkIdx := range scanGalleryID {
			if strings.Contains(scanGalleryID[nhlinkIdx], "https://t.nhentai.net/galleries/") {
				nhGetGallIDSplit = nil
				nhGetGallID1 = strings.ReplaceAll(scanGalleryID[nhlinkIdx], "https://t.nhentai.net/galleries/", "")
				nhGetGallIDSplit = strings.Split(nhGetGallID1, "/cover")
				nhGetGallID2 = nhGetGallIDSplit[0]
				break
			}
		}

		nhImgNames = nil
		nhImgLinks = nil
		for nhCurrPg := 1; nhCurrPg <= nhTotalPage; nhCurrPg++ {
			nhImgLink = fmt.Sprintf("https://i.nhentai.net/galleries/%v/%v.jpg", nhGetGallID2, nhCurrPg)

			// Get the image and write it to memory
			getImg, err := httpclient.Get(nhImgLink)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				break
			}

			// convert http response to io.Reader
			bodyIMG, err := ioutil.ReadAll(bufio.NewReader(getImg.Body))
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				break
			}
			reader := bytes.NewReader(bodyIMG)

			// Send image thru DM.
			// We create the private channel with the user who sent the message.
			channel, err := ctx.Session.UserChannelCreate(userID)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
				break
			}
			// Then we send the message through the channel we created.
			_, err = ctx.Session.ChannelFileSend(channel.ID, fmt.Sprintf("%v.jpg", nhCurrPg), reader)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}
				break
			}

			// ==================================
			// Create a new image based on bodyImg
			nhImgName = fmt.Sprintf("./nh/%v.jpg", nhCurrPg)
			createImgFile, err := memFS.Create(nhImgName)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				break
			} else {
				// Write to the file
				writeImgFile, err := io.Copy(createImgFile, getImg.Body)
				if err != nil {
					fmt.Println(" [ERROR] ", err)
					getImg.Body.Close()

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				} else {
					// Close the file
					getImg.Body.Close()
					if err := createImgFile.Close(); err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						break
					} else {
						winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", nhImgName, (writeImgFile / Kilobyte), (writeImgFile / Megabyte))
						fmt.Println(winLogs)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
						}
					}
				}
			}

			nhImgLinkLocal = fmt.Sprintf("https://x.castella.network/img/%v.jpg", nhCurrPg)
			nhImgNames = append(nhImgNames, nhImgLinkLocal)
			nhImgLinks = append(nhImgLinks, fmt.Sprintf("\n%v", nhImgLink))
		}

		getGalleryID.Body.Close()

		// send manga info to user via DM.
		// Create the embed templates.
		oriURLField := discordgo.MessageEmbedField{
			Name:   "Original URL",
			Value:  fmt.Sprintf("https://nhentai.net/g/%v", nhCode),
			Inline: false,
		}
		showURLField := discordgo.MessageEmbedField{
			Name:   "Total Pages",
			Value:  fmt.Sprintf("%v", len(nhImgNames)),
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&oriURLField, &showURLField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  fmt.Sprintf("About ID-%v", nhCode),
			Color:  0x03fcad,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		// Send image thru DM.
		// We create the private channel with the user who sent the message.
		channel, err := ctx.Session.UserChannelCreate(userID)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
			return
		}
		// Then we send the message through the channel we created.
		_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		// add a quick reply for Go to Page 1
		ctx.Session.ChannelMessageSendReply(channel.ID, "**Go to Page 1**", ctx.Event.Reference())

	}

}

// KatInz's VMG feature
func getPics(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func katInzVMG(ctx *dgc.Ctx) {

	vmgRelax := xurls.Relaxed()
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	vmgMaxRender = 1
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// Get the gallery ID
		vmgGalleryID := fmt.Sprintf("https://www.vmgirls.com/%v.html", rawArgs)
		getGalleryID, err := httpclient.Get(vmgGalleryID)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		bodyGalleryID, err := ioutil.ReadAll(bufio.NewReader(getGalleryID.Body))
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		scanGalleryID := vmgRelax.FindAllString(string(bodyGalleryID), -1)
		onlyPics := getPics(scanGalleryID)

		for vmgCurrImg := range onlyPics {

			// only handle image links
			if strings.Contains(onlyPics[vmgCurrImg], "t.cdn.ink/image/") {

				// Get the image and write it to memory
				getImg, err := httpclient.Get(fmt.Sprintf("https://%v", onlyPics[vmgCurrImg]))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				}

				// convert http response to io.Reader
				bodyIMG, err := ioutil.ReadAll(bufio.NewReader(getImg.Body))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				}
				reader := bytes.NewReader(bodyIMG)

				// Send image thru DM.
				// We create the private channel with the user who sent the message.
				channel, err := ctx.Session.UserChannelCreate(userID)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					break
				}
				// Then we send the message through the channel we created.
				_, err = ctx.Session.ChannelFileSend(channel.ID, fmt.Sprintf("%v.jpg", vmgCurrImg), reader)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					break
				}

				getImg.Body.Close()
				vmgMaxRender++

				if vmgMaxRender == 20 {
					break
				}
			}

		}

		getGalleryID.Body.Close()

		// send manga info to user via DM.
		// Create the embed templates.
		oriURLField := discordgo.MessageEmbedField{
			Name:   "Original URL",
			Value:  fmt.Sprintf("https://www.vmgirls.com/%v.html", rawArgs),
			Inline: false,
		}
		showURLField := discordgo.MessageEmbedField{
			Name:   "Total Images",
			Value:  fmt.Sprintf("%v", len(onlyPics)),
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&oriURLField, &showURLField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  fmt.Sprintf("About ID-%v", rawArgs),
			Color:  0x03fcad,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		// Send image thru DM.
		// We create the private channel with the user who sent the message.
		channel, err := ctx.Session.UserChannelCreate(userID)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
			return
		}
		// Then we send the message through the channel we created.
		_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		// add a quick reply for Go to Page 1
		ctx.Session.ChannelMessageSendReply(channel.ID, "**Go to Image 1**", ctx.Event.Reference())

	}

}

// KatMon run Go code
func katMonGoRun(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	msgAttachment := ctx.Event.Attachments
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// make a new empty folder
		osFS.RemoveAll("./gocode/")
		osFS.MkdirAll("./gocode/", 0777)

		for fileIdx := range msgAttachment {

			// Get the image and write it to memory
			getFile, err := httpclient.Get(msgAttachment[fileIdx].URL)
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				break
			}

			// ==================================
			// Create a new dummy.go file
			createGoFile, err := osFS.Create("./gocode/dummy.go")
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				return
			} else {
				// Write to the file
				writeGoFile, err := io.Copy(createGoFile, getFile.Body)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					getFile.Body.Close()

					return
				} else {

					getFile.Body.Close()

					if err := createGoFile.Close(); err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					} else {
						winLogs = fmt.Sprintf(" [DONE] `dummy.go` file has been created. \n >> Size: %v KB (%v MB)", (writeGoFile / Kilobyte), (writeGoFile / Megabyte))
						fmt.Println(winLogs)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
						}

						// run the code
						codeExec := time.Now()
						gofmt, err := exec.Command("go", "fmt", "./gocode/dummy.go").Output()
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						}
						fmt.Println(string(gofmt))

						gorun, err := exec.Command("go", "run", "./gocode/dummy.go").Output()
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						}
						execTime := time.Since(codeExec)

						// report after code execution has ended
						// Create the embed templates
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", userID),
							Inline: false,
						}
						timeElapsedField := discordgo.MessageEmbedField{
							Name:   "Execution Time",
							Value:  fmt.Sprintf("%v", execTime),
							Inline: false,
						}
						outputField := discordgo.MessageEmbedField{
							Name:   "Output",
							Value:  fmt.Sprintf("%v", string(gorun)),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &outputField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "Go Katheryne",
							Color:  0x9155fa,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
					}
				}
			}

		}

	}

}

// KatMon run Real-ESRGAN
func katMonWaifu2x(ctx *dgc.Ctx) {

	w2xRelax := xurls.Relaxed()
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	msgAttachment := ctx.Event.Attachments
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// default value
	scaleOK = false

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		// make a new empty folder
		osFS.RemoveAll("./img/")
		osFS.MkdirAll("./img/", 0777)
		memFS.RemoveAll("./pics/")
		memFS.MkdirAll("./pics/", 0777)

		if rawArgs == "file" {

			if strings.Contains(rawArgs, "[anime::") {

				// get input parameters
				getScale := strings.Split(rawArgs, "[anime::")
				splitScale := strings.Split(getScale[1], "]")
				scale := splitScale[0]

				// check user's desired scale with scaleOpts
				for optIdx := range scaleOpts {
					if scale == scaleOpts[optIdx] {
						scaleOK = true
						break
					}
				}

				if scaleOK {
					for fileIdx := range msgAttachment {

						// Get the image and write it to memory
						getFile, err := httpclient.Get(msgAttachment[fileIdx].URL)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							break
						}

						// ==================================
						// Create a new dummy.png file
						createIMGFile, err := osFS.Create("./img/dummy.png")
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						} else {
							// Write to the file
							writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								getFile.Body.Close()

								return
							} else {

								getFile.Body.Close()

								if err := createIMGFile.Close(); err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								} else {
									winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
									fmt.Println(winLogs)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
									}

									// use waifu2x
									codeExec := time.Now()

									w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-upconv_7_anime_style_art_rgb", "-n", "3", "-s", scale, "-j", "4:4:4").Output()
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									fmt.Println(string(w2x))

									// inform the new file size
									info, err := osFS.Stat("./img/output.png")
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									size := info.Size()
									name := info.Name()

									readOutput, err := afero.ReadFile(osFS, "./img/output.png")
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									reader := bytes.NewReader(readOutput)

									// copy it to memFS
									// ==================================
									// Create the output.png file
									createMemCopy, err := memFS.Create("./pics/output.png")
									if err != nil {
										log.Fatal(err)
									} else {
										// Write to the file
										writeMemCopy, err := createMemCopy.Write(readOutput)
										if err != nil {
											log.Fatal(err)
										} else {
											// Close the file
											if err := createMemCopy.Close(); err != nil {
												log.Fatal(err)
											} else {
												fmt.Println()
												winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
												fmt.Println(winLogs)
											}
										}
									}

									// send output image
									execTime := time.Since(codeExec)

									// report after code execution has ended
									// Create the embed templates
									usernameField := discordgo.MessageEmbedField{
										Name:   "Username",
										Value:  fmt.Sprintf("<@!%v>", userID),
										Inline: false,
									}
									timeElapsedField := discordgo.MessageEmbedField{
										Name:   "Execution Time",
										Value:  fmt.Sprintf("%v", execTime),
										Inline: false,
									}
									newsizeField := discordgo.MessageEmbedField{
										Name:   "New Size",
										Value:  fmt.Sprintf("x%v upscale [%v KB (%v MB)]", scale, (size / Kilobyte), (size / Megabyte)),
										Inline: false,
									}
									newlinkField := discordgo.MessageEmbedField{
										Name:   "Location in Memory",
										Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
										Inline: false,
									}
									messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

									aoiEmbedFooter := discordgo.MessageEmbedFooter{
										Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
									}

									aoiEmbeds := discordgo.MessageEmbed{
										Title:  "Katheryne's AI for Images",
										Color:  0x7dfa7a,
										Footer: &aoiEmbedFooter,
										Fields: messageFields,
									}

									ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
									ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

								}
							}
						}

					}
				}

			} else if strings.Contains(rawArgs, "[photo::") {

				// get input parameters
				getScale := strings.Split(rawArgs, "[photo::")
				splitScale := strings.Split(getScale[1], "]")
				scale := splitScale[0]

				// check user's desired scale with scaleOpts
				for optIdx := range scaleOpts {
					if scale == scaleOpts[optIdx] {
						scaleOK = true
						break
					}
				}

				if scaleOK {
					for fileIdx := range msgAttachment {

						// Get the image and write it to memory
						getFile, err := httpclient.Get(msgAttachment[fileIdx].URL)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							break
						}

						// ==================================
						// Create a new dummy.png file
						createIMGFile, err := osFS.Create("./img/dummy.png")
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						} else {
							// Write to the file
							writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								getFile.Body.Close()

								return
							} else {

								getFile.Body.Close()

								if err := createIMGFile.Close(); err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								} else {
									winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
									fmt.Println(winLogs)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
									}

									// use waifu2x
									codeExec := time.Now()

									w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-upconv_7_photo", "-n", "3", "-s", scale, "-j", "4:4:4").Output()
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									fmt.Println(string(w2x))

									// inform the new file size
									info, err := osFS.Stat("./img/output.png")
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									size := info.Size()
									name := info.Name()

									readOutput, err := afero.ReadFile(osFS, "./img/output.png")
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									reader := bytes.NewReader(readOutput)

									// copy it to memFS
									// ==================================
									// Create the output.png file
									createMemCopy, err := memFS.Create("./pics/output.png")
									if err != nil {
										log.Fatal(err)
									} else {
										// Write to the file
										writeMemCopy, err := createMemCopy.Write(readOutput)
										if err != nil {
											log.Fatal(err)
										} else {
											// Close the file
											if err := createMemCopy.Close(); err != nil {
												log.Fatal(err)
											} else {
												fmt.Println()
												winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
												fmt.Println(winLogs)
											}
										}
									}

									// send output image
									execTime := time.Since(codeExec)

									// report after code execution has ended
									// Create the embed templates
									usernameField := discordgo.MessageEmbedField{
										Name:   "Username",
										Value:  fmt.Sprintf("<@!%v>", userID),
										Inline: false,
									}
									timeElapsedField := discordgo.MessageEmbedField{
										Name:   "Execution Time",
										Value:  fmt.Sprintf("%v", execTime),
										Inline: false,
									}
									newsizeField := discordgo.MessageEmbedField{
										Name:   "New Size",
										Value:  fmt.Sprintf("x%v upscale [%v KB (%v MB)]", scale, (size / Kilobyte), (size / Megabyte)),
										Inline: false,
									}
									newlinkField := discordgo.MessageEmbedField{
										Name:   "Location in Memory",
										Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
										Inline: false,
									}
									messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

									aoiEmbedFooter := discordgo.MessageEmbedFooter{
										Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
									}

									aoiEmbeds := discordgo.MessageEmbed{
										Title:  "Katheryne's AI for Images",
										Color:  0x7dfa7a,
										Footer: &aoiEmbedFooter,
										Fields: messageFields,
									}

									ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
									ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

								}
							}
						}

					}
				}

			} else if strings.Contains(rawArgs, "[auto::") {

				// get input parameters
				getScale := strings.Split(rawArgs, "[auto::")
				splitScale := strings.Split(getScale[1], "]")
				scale := splitScale[0]

				// check user's desired scale with scaleOpts
				for optIdx := range scaleOpts {
					if scale == scaleOpts[optIdx] {
						scaleOK = true
						break
					}
				}

				if scaleOK {
					for fileIdx := range msgAttachment {

						// Get the image and write it to memory
						getFile, err := httpclient.Get(msgAttachment[fileIdx].URL)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							break
						}

						// ==================================
						// Create a new dummy.png file
						createIMGFile, err := osFS.Create("./img/dummy.png")
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						} else {
							// Write to the file
							writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								getFile.Body.Close()

								return
							} else {

								getFile.Body.Close()

								if err := createIMGFile.Close(); err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								} else {
									winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
									fmt.Println(winLogs)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
									}

									// use waifu2x
									codeExec := time.Now()

									w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-cunet", "-n", "3", "-s", scale, "-j", "4:4:4").Output()
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									fmt.Println(string(w2x))

									// inform the new file size
									info, err := osFS.Stat("./img/output.png")
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									size := info.Size()
									name := info.Name()

									readOutput, err := afero.ReadFile(osFS, "./img/output.png")
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) == universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										return
									}
									reader := bytes.NewReader(readOutput)

									// copy it to memFS
									// ==================================
									// Create the output.png file
									createMemCopy, err := memFS.Create("./pics/output.png")
									if err != nil {
										log.Fatal(err)
									} else {
										// Write to the file
										writeMemCopy, err := createMemCopy.Write(readOutput)
										if err != nil {
											log.Fatal(err)
										} else {
											// Close the file
											if err := createMemCopy.Close(); err != nil {
												log.Fatal(err)
											} else {
												fmt.Println()
												winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
												fmt.Println(winLogs)
											}
										}
									}

									// send output image
									execTime := time.Since(codeExec)

									// report after code execution has ended
									// Create the embed templates
									usernameField := discordgo.MessageEmbedField{
										Name:   "Username",
										Value:  fmt.Sprintf("<@!%v>", userID),
										Inline: false,
									}
									timeElapsedField := discordgo.MessageEmbedField{
										Name:   "Execution Time",
										Value:  fmt.Sprintf("%v", execTime),
										Inline: false,
									}
									newsizeField := discordgo.MessageEmbedField{
										Name:   "New Size",
										Value:  fmt.Sprintf("x%v upscale [%v KB (%v MB)]", scale, (size / Kilobyte), (size / Megabyte)),
										Inline: false,
									}
									newlinkField := discordgo.MessageEmbedField{
										Name:   "Location in Memory",
										Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
										Inline: false,
									}
									messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

									aoiEmbedFooter := discordgo.MessageEmbedFooter{
										Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
									}

									aoiEmbeds := discordgo.MessageEmbed{
										Title:  "Katheryne's AI for Images",
										Color:  0x7dfa7a,
										Footer: &aoiEmbedFooter,
										Fields: messageFields,
									}

									ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
									ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

								}
							}
						}

					}
				}

			} else {

				for fileIdx := range msgAttachment {

					// Get the image and write it to memory
					getFile, err := httpclient.Get(msgAttachment[fileIdx].URL)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						break
					}

					// ==================================
					// Create a new dummy.png file
					createIMGFile, err := osFS.Create("./img/dummy.png")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					} else {
						// Write to the file
						writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							getFile.Body.Close()

							return
						} else {

							getFile.Body.Close()

							if err := createIMGFile.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							} else {
								winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}

								// use waifu2x
								codeExec := time.Now()

								w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-cunet", "-n", "3", "-s", "2", "-j", "4:4:4").Output()
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								fmt.Println(string(w2x))

								// inform the new file size
								info, err := osFS.Stat("./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								size := info.Size()
								name := info.Name()

								readOutput, err := afero.ReadFile(osFS, "./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								reader := bytes.NewReader(readOutput)

								// copy it to memFS
								// ==================================
								// Create the output.png file
								createMemCopy, err := memFS.Create("./pics/output.png")
								if err != nil {
									log.Fatal(err)
								} else {
									// Write to the file
									writeMemCopy, err := createMemCopy.Write(readOutput)
									if err != nil {
										log.Fatal(err)
									} else {
										// Close the file
										if err := createMemCopy.Close(); err != nil {
											log.Fatal(err)
										} else {
											fmt.Println()
											winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
											fmt.Println(winLogs)
										}
									}
								}

								// send output image
								execTime := time.Since(codeExec)

								// report after code execution has ended
								// Create the embed templates
								usernameField := discordgo.MessageEmbedField{
									Name:   "Username",
									Value:  fmt.Sprintf("<@!%v>", userID),
									Inline: false,
								}
								timeElapsedField := discordgo.MessageEmbedField{
									Name:   "Execution Time",
									Value:  fmt.Sprintf("%v", execTime),
									Inline: false,
								}
								newsizeField := discordgo.MessageEmbedField{
									Name:   "New Size",
									Value:  fmt.Sprintf("x2 upscale [%v KB (%v MB)]", (size / Kilobyte), (size / Megabyte)),
									Inline: false,
								}
								newlinkField := discordgo.MessageEmbedField{
									Name:   "Location in Memory",
									Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
									Inline: false,
								}
								messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

								aoiEmbedFooter := discordgo.MessageEmbedFooter{
									Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
								}

								aoiEmbeds := discordgo.MessageEmbed{
									Title:  "Katheryne's AI for Images",
									Color:  0x7dfa7a,
									Footer: &aoiEmbedFooter,
									Fields: messageFields,
								}

								ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
								ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

							}
						}
					}

				}
			}

		} else if strings.Contains(rawArgs, "http://") || strings.Contains(rawArgs, "https://") {

			// scan the link
			scanLinks := w2xRelax.FindAllString(rawArgs, -1)

			if strings.Contains(rawArgs, "[anime::") {

				// get input parameters
				getScale := strings.Split(rawArgs, "[anime::")
				splitScale := strings.Split(getScale[1], "]")
				scale := splitScale[0]

				// check user's desired scale with scaleOpts
				for optIdx := range scaleOpts {
					if scale == scaleOpts[optIdx] {
						scaleOK = true
						break
					}
				}

				if scaleOK {

					// Get the image and write it to memory
					getFile, err := httpclient.Get(scanLinks[0])
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					}

					// ==================================
					// Create a new dummy.png file
					createIMGFile, err := osFS.Create("./img/dummy.png")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					} else {
						// Write to the file
						writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							getFile.Body.Close()

							return
						} else {

							getFile.Body.Close()

							if err := createIMGFile.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							} else {
								winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}

								// use waifu2x
								codeExec := time.Now()

								w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-upconv_7_anime_style_art_rgb", "-n", "3", "-s", scale, "-j", "4:4:4").Output()
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								fmt.Println(string(w2x))

								// inform the new file size
								info, err := osFS.Stat("./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								size := info.Size()
								name := info.Name()

								readOutput, err := afero.ReadFile(osFS, "./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								reader := bytes.NewReader(readOutput)

								// copy it to memFS
								// ==================================
								// Create the output.png file
								createMemCopy, err := memFS.Create("./pics/output.png")
								if err != nil {
									log.Fatal(err)
								} else {
									// Write to the file
									writeMemCopy, err := createMemCopy.Write(readOutput)
									if err != nil {
										log.Fatal(err)
									} else {
										// Close the file
										if err := createMemCopy.Close(); err != nil {
											log.Fatal(err)
										} else {
											fmt.Println()
											winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
											fmt.Println(winLogs)
										}
									}
								}

								// send output image
								execTime := time.Since(codeExec)

								// report after code execution has ended
								// Create the embed templates
								usernameField := discordgo.MessageEmbedField{
									Name:   "Username",
									Value:  fmt.Sprintf("<@!%v>", userID),
									Inline: false,
								}
								timeElapsedField := discordgo.MessageEmbedField{
									Name:   "Execution Time",
									Value:  fmt.Sprintf("%v", execTime),
									Inline: false,
								}
								newsizeField := discordgo.MessageEmbedField{
									Name:   "New Size",
									Value:  fmt.Sprintf("x%v upscale [%v KB (%v MB)]", scale, (size / Kilobyte), (size / Megabyte)),
									Inline: false,
								}
								newlinkField := discordgo.MessageEmbedField{
									Name:   "Location in Memory",
									Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
									Inline: false,
								}
								messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

								aoiEmbedFooter := discordgo.MessageEmbedFooter{
									Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
								}

								aoiEmbeds := discordgo.MessageEmbed{
									Title:  "Katheryne's AI for Images",
									Color:  0x7dfa7a,
									Footer: &aoiEmbedFooter,
									Fields: messageFields,
								}

								ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
								ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

							}
						}
					}

				}
			} else if strings.Contains(rawArgs, "[photo::") {

				// get input parameters
				getScale := strings.Split(rawArgs, "[photo::")
				splitScale := strings.Split(getScale[1], "]")
				scale := splitScale[0]

				// check user's desired scale with scaleOpts
				for optIdx := range scaleOpts {
					if scale == scaleOpts[optIdx] {
						scaleOK = true
						break
					}
				}

				if scaleOK {

					// Get the image and write it to memory
					getFile, err := httpclient.Get(scanLinks[0])
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					}

					// ==================================
					// Create a new dummy.png file
					createIMGFile, err := osFS.Create("./img/dummy.png")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					} else {
						// Write to the file
						writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							getFile.Body.Close()

							return
						} else {

							getFile.Body.Close()

							if err := createIMGFile.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							} else {
								winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}

								// use waifu2x
								codeExec := time.Now()

								w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-upconv_7_photo", "-n", "3", "-s", scale, "-j", "4:4:4").Output()
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								fmt.Println(string(w2x))

								// inform the new file size
								info, err := osFS.Stat("./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								size := info.Size()
								name := info.Name()

								readOutput, err := afero.ReadFile(osFS, "./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								reader := bytes.NewReader(readOutput)

								// copy it to memFS
								// ==================================
								// Create the output.png file
								createMemCopy, err := memFS.Create("./pics/output.png")
								if err != nil {
									log.Fatal(err)
								} else {
									// Write to the file
									writeMemCopy, err := createMemCopy.Write(readOutput)
									if err != nil {
										log.Fatal(err)
									} else {
										// Close the file
										if err := createMemCopy.Close(); err != nil {
											log.Fatal(err)
										} else {
											fmt.Println()
											winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
											fmt.Println(winLogs)
										}
									}
								}

								// send output image
								execTime := time.Since(codeExec)

								// report after code execution has ended
								// Create the embed templates
								usernameField := discordgo.MessageEmbedField{
									Name:   "Username",
									Value:  fmt.Sprintf("<@!%v>", userID),
									Inline: false,
								}
								timeElapsedField := discordgo.MessageEmbedField{
									Name:   "Execution Time",
									Value:  fmt.Sprintf("%v", execTime),
									Inline: false,
								}
								newsizeField := discordgo.MessageEmbedField{
									Name:   "New Size",
									Value:  fmt.Sprintf("x%v upscale [%v KB (%v MB)]", scale, (size / Kilobyte), (size / Megabyte)),
									Inline: false,
								}
								newlinkField := discordgo.MessageEmbedField{
									Name:   "Location in Memory",
									Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
									Inline: false,
								}
								messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

								aoiEmbedFooter := discordgo.MessageEmbedFooter{
									Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
								}

								aoiEmbeds := discordgo.MessageEmbed{
									Title:  "Katheryne's AI for Images",
									Color:  0x7dfa7a,
									Footer: &aoiEmbedFooter,
									Fields: messageFields,
								}

								ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
								ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

							}
						}
					}

				}

			} else if strings.Contains(rawArgs, "[auto::") {

				// get input parameters
				getScale := strings.Split(rawArgs, "[auto::")
				splitScale := strings.Split(getScale[1], "]")
				scale := splitScale[0]

				// check user's desired scale with scaleOpts
				for optIdx := range scaleOpts {
					if scale == scaleOpts[optIdx] {
						scaleOK = true
						break
					}
				}

				if scaleOK {

					// Get the image and write it to memory
					getFile, err := httpclient.Get(scanLinks[0])
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					}

					// ==================================
					// Create a new dummy.png file
					createIMGFile, err := osFS.Create("./img/dummy.png")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					} else {
						// Write to the file
						writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							getFile.Body.Close()

							return
						} else {

							getFile.Body.Close()

							if err := createIMGFile.Close(); err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							} else {
								winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
								fmt.Println(winLogs)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
								}

								// use waifu2x
								codeExec := time.Now()

								w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-cunet", "-n", "3", "-s", scale, "-j", "4:4:4").Output()
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								fmt.Println(string(w2x))

								// inform the new file size
								info, err := osFS.Stat("./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								size := info.Size()
								name := info.Name()

								readOutput, err := afero.ReadFile(osFS, "./img/output.png")
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) == universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									return
								}
								reader := bytes.NewReader(readOutput)

								// copy it to memFS
								// ==================================
								// Create the output.png file
								createMemCopy, err := memFS.Create("./pics/output.png")
								if err != nil {
									log.Fatal(err)
								} else {
									// Write to the file
									writeMemCopy, err := createMemCopy.Write(readOutput)
									if err != nil {
										log.Fatal(err)
									} else {
										// Close the file
										if err := createMemCopy.Close(); err != nil {
											log.Fatal(err)
										} else {
											fmt.Println()
											winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
											fmt.Println(winLogs)
										}
									}
								}

								// send output image
								execTime := time.Since(codeExec)

								// report after code execution has ended
								// Create the embed templates
								usernameField := discordgo.MessageEmbedField{
									Name:   "Username",
									Value:  fmt.Sprintf("<@!%v>", userID),
									Inline: false,
								}
								timeElapsedField := discordgo.MessageEmbedField{
									Name:   "Execution Time",
									Value:  fmt.Sprintf("%v", execTime),
									Inline: false,
								}
								newsizeField := discordgo.MessageEmbedField{
									Name:   "New Size",
									Value:  fmt.Sprintf("x%v upscale [%v KB (%v MB)]", scale, (size / Kilobyte), (size / Megabyte)),
									Inline: false,
								}
								newlinkField := discordgo.MessageEmbedField{
									Name:   "Location in Memory",
									Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
									Inline: false,
								}
								messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

								aoiEmbedFooter := discordgo.MessageEmbedFooter{
									Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
								}

								aoiEmbeds := discordgo.MessageEmbed{
									Title:  "Katheryne's AI for Images",
									Color:  0x7dfa7a,
									Footer: &aoiEmbedFooter,
									Fields: messageFields,
								}

								ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
								ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

							}
						}
					}

				}

			} else {

				// Get the image and write it to memory
				getFile, err := httpclient.Get(rawArgs)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}

				// ==================================
				// Create a new dummy.png file
				createIMGFile, err := osFS.Create("./img/dummy.png")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				} else {
					// Write to the file
					writeIMGFile, err := io.Copy(createIMGFile, getFile.Body)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) == universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						getFile.Body.Close()

						return
					} else {

						getFile.Body.Close()

						if err := createIMGFile.Close(); err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						} else {
							winLogs = fmt.Sprintf(" [DONE] `dummy.png` file has been created. \n >> Size: %v KB (%v MB)", (writeIMGFile / Kilobyte), (writeIMGFile / Megabyte))
							fmt.Println(winLogs)

							if len(universalLogs) == universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
							}

							// use waifu2x
							codeExec := time.Now()

							w2x, err := exec.Command("./1w2x", "-i", "./img/dummy.png", "-o", "./img/output.png", "-m", "models-cunet", "-n", "3", "-s", "2", "-j", "4:4:4").Output()
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							}
							fmt.Println(string(w2x))

							// inform the new file size
							info, err := osFS.Stat("./img/output.png")
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							}
							size := info.Size()
							name := info.Name()

							readOutput, err := afero.ReadFile(osFS, "./img/output.png")
							if err != nil {
								fmt.Println(" [ERROR] ", err)

								if len(universalLogs) == universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							}
							reader := bytes.NewReader(readOutput)

							// copy it to memFS
							// ==================================
							// Create the output.png file
							createMemCopy, err := memFS.Create("./pics/output.png")
							if err != nil {
								log.Fatal(err)
							} else {
								// Write to the file
								writeMemCopy, err := createMemCopy.Write(readOutput)
								if err != nil {
									log.Fatal(err)
								} else {
									// Close the file
									if err := createMemCopy.Close(); err != nil {
										log.Fatal(err)
									} else {
										fmt.Println()
										winLogs = fmt.Sprintf(" [DONE] output.png has been created. \n >> Size: %v KB (%v MB)", (writeMemCopy / Kilobyte), (writeMemCopy / Megabyte))
										fmt.Println(winLogs)
									}
								}
							}

							// send output image
							execTime := time.Since(codeExec)

							// report after code execution has ended
							// Create the embed templates
							usernameField := discordgo.MessageEmbedField{
								Name:   "Username",
								Value:  fmt.Sprintf("<@!%v>", userID),
								Inline: false,
							}
							timeElapsedField := discordgo.MessageEmbedField{
								Name:   "Execution Time",
								Value:  fmt.Sprintf("%v", execTime),
								Inline: false,
							}
							newsizeField := discordgo.MessageEmbedField{
								Name:   "New Size",
								Value:  fmt.Sprintf("x2 upscale [%v KB (%v MB)]", (size / Kilobyte), (size / Megabyte)),
								Inline: false,
							}
							newlinkField := discordgo.MessageEmbedField{
								Name:   "Location in Memory",
								Value:  fmt.Sprintf("https://x.castella.network/ai/%v", name),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &newlinkField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "Katheryne's AI for Images",
								Color:  0x7dfa7a,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
							}

							ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
							ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "output.png", reader)

						}
					}
				}

			}

		}

	}

}

func katInzCK101(ctx *dgc.Ctx) {

	ckRelax := xurls.Relaxed()
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		ckImgs = nil

		// Get the webpage data
		getGalleryID, err := httpclient.Get(rawArgs)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		bodyGalleryID, err := ioutil.ReadAll(bufio.NewReader(getGalleryID.Body))
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		scanLinks := ckRelax.FindAllString(string(bodyGalleryID), -1)

		for ckCurrIdx := range scanLinks {

			// check if it's the user's pics
			if strings.Contains(scanLinks[ckCurrIdx], ".jpg") || strings.Contains(scanLinks[ckCurrIdx], ".png") {

				// add new data to slice
				ckImgs = append(ckImgs, scanLinks[ckCurrIdx])

				// Get the image and write it to memory
				getImg, err := httpclient.Get(fmt.Sprintf("%v", scanLinks[ckCurrIdx]))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				}

				// convert http response to io.Reader
				bodyIMG, err := ioutil.ReadAll(bufio.NewReader(getImg.Body))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					break
				}
				reader := bytes.NewReader(bodyIMG)

				// Send image thru DM.
				// We create the private channel with the user who sent the message.
				channel, err := ctx.Session.UserChannelCreate(userID)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					break
				}
				// Then we send the message through the channel we created.
				_, err = ctx.Session.ChannelFileSend(channel.ID, fmt.Sprintf("%v.jpg", ckCurrIdx), reader)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					break
				}

				getImg.Body.Close()
			}

		}

		getGalleryID.Body.Close()

		// send manga info to user via DM.
		// Create the embed templates.
		oriURLField := discordgo.MessageEmbedField{
			Name:   "Original URL",
			Value:  fmt.Sprintf("%v", rawArgs),
			Inline: false,
		}
		showURLField := discordgo.MessageEmbedField{
			Name:   "Total Images",
			Value:  fmt.Sprintf("%v", len(ckImgs)),
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&oriURLField, &showURLField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  "CK101 Information",
			Color:  0x03fcad,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		// Send image thru DM.
		// We create the private channel with the user who sent the message.
		channel, err := ctx.Session.UserChannelCreate(userID)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
			return
		}
		// Then we send the message through the channel we created.
		_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
		}

		// add a quick reply for Go to Page 1
		ctx.Session.ChannelMessageSendReply(channel.ID, "**Go to Image 1**", ctx.Event.Reference())

	}

}

func katInzYTDL(ctx *dgc.Ctx) {

	ytRelax := xurls.Relaxed()
	userID := ctx.Event.Author.ID
	arguments := ctx.Arguments
	rawArgs := arguments.Raw()
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	// rawArgs shouldn't be empty
	if len(rawArgs) != 0 {

		osFS.RemoveAll("./ytdl")
		osFS.MkdirAll("./ytdl", 0777)
		katInzVidID = ""

		// check if it's for DM or public
		if strings.Contains(strings.ToLower(rawArgs), "dm") {

			// check if input is a mention or a raw userID
			mentionedUser := arguments.Get(1)
			if strings.Contains(mentionedUser.Raw(), "<@!") || strings.Contains(mentionedUser.Raw(), "<@") {

				// get destination user data
				userMention := mentionedUser.AsUserMentionID()
				userData, err := ctx.Session.User(userMention)
				if err != nil {
					fmt.Println(" [userData] ", err)
					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				// delete user's message and send confirmation as a reply
				scanLinks := ytRelax.FindAllString(rawArgs, -1)

				// get the video ID
				if strings.Contains(scanLinks[0], "www.youtube.com") {
					// sample URL >> https://www.youtube.com/watch?v=J5x0tLiItVY
					splitVidID := strings.Split(scanLinks[0], "youtube.com/watch?v=")
					katInzVidID = splitVidID[1]
				} else if strings.Contains(scanLinks[0], "youtu.be") {
					// sample URL >> https://youtu.be/J5x0tLiItVY
					splitVidID := strings.Split(scanLinks[0], "youtu.be/")
					katInzVidID = splitVidID[1]
				}
				ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, ctx.Event.ID)
				ctx.RespondText(fmt.Sprintf("Processing `%v`. Please wait.", katInzVidID))

				// run the code
				codeExec := time.Now()
				katYT, err := exec.Command("yt-dlp", "--ignore-config", "--no-playlist", "--max-filesize", "20m", "-P", "$PWD/ytdl", "-o", "%(id)s.%(ext)s", "-x", "--audio-format", "mp3", scanLinks[0]).Output()
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				fmt.Println(string(katYT))
				execTime := time.Since(codeExec)

				outIdx, err := afero.ReadDir(osFS, "./ytdl")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}

				readOutput1, err := afero.ReadFile(osFS, fmt.Sprintf("./ytdl/%v", outIdx[0].Name()))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				reader1 := bytes.NewReader(readOutput1)

				readOutput2, err := afero.ReadFile(osFS, fmt.Sprintf("./ytdl/%v", outIdx[0].Name()))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				reader2 := bytes.NewReader(readOutput2)

				// report after code execution has ended
				// Create the embed templates
				fromField := discordgo.MessageEmbedField{
					Name:   "From",
					Value:  fmt.Sprintf("<@!%v>", userID),
					Inline: false,
				}
				toField := discordgo.MessageEmbedField{
					Name:   "To",
					Value:  fmt.Sprintf("<@!%v>", userData.ID),
					Inline: false,
				}
				timeElapsedField := discordgo.MessageEmbedField{
					Name:   "Execution Time",
					Value:  fmt.Sprintf("%v", execTime),
					Inline: false,
				}
				newsizeField := discordgo.MessageEmbedField{
					Name:   "New Size",
					Value:  fmt.Sprintf("%v KB (%v MB)", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
					Inline: false,
				}
				fileIDField := discordgo.MessageEmbedField{
					Name:   "File ID",
					Value:  fmt.Sprintf("`%v`", katInzVidID),
					Inline: false,
				}
				messageFields := []*discordgo.MessageEmbedField{&fromField, &toField, &timeElapsedField, &newsizeField, &fileIDField}

				aoiEmbedFooter := discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
				}

				aoiEmbeds := discordgo.MessageEmbed{
					Title:  "YT Katheryne",
					Color:  0x52ff91,
					Footer: &aoiEmbedFooter,
					Fields: messageFields,
				}

				// send a copy to Sender
				channel, err := ctx.Session.UserChannelCreate(userID)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}
				_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
				}
				ctx.Session.ChannelFileSend(channel.ID, outIdx[0].Name(), reader1)

				// send a copy to Receiver
				channel, err = ctx.Session.UserChannelCreate(userData.ID)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}
				_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
				}
				ctx.Session.ChannelFileSend(channel.ID, outIdx[0].Name(), reader2)
			} else {

				// get destination user data
				userData, err := ctx.Session.User(mentionedUser.Raw())
				if err != nil {
					fmt.Println(" [userData] ", err)
					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				// delete user's message and send confirmation as a reply
				scanLinks := ytRelax.FindAllString(rawArgs, -1)

				// get the video ID
				if strings.Contains(scanLinks[0], "www.youtube.com") {
					// sample URL >> https://www.youtube.com/watch?v=J5x0tLiItVY
					splitVidID := strings.Split(scanLinks[0], "youtube.com/watch?v=")
					katInzVidID = splitVidID[1]
				} else if strings.Contains(scanLinks[0], "youtu.be") {
					// sample URL >> https://youtu.be/J5x0tLiItVY
					splitVidID := strings.Split(scanLinks[0], "youtu.be/")
					katInzVidID = splitVidID[1]
				}
				ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, ctx.Event.ID)
				ctx.RespondText(fmt.Sprintf("Processing `%v`. Please wait.", katInzVidID))

				// run the code
				codeExec := time.Now()
				katYT, err := exec.Command("yt-dlp", "--ignore-config", "--no-playlist", "--max-filesize", "20m", "-P", "$PWD/ytdl", "-o", "%(id)s.%(ext)s", "-x", "--audio-format", "mp3", scanLinks[0]).Output()
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				fmt.Println(string(katYT))
				execTime := time.Since(codeExec)

				outIdx, err := afero.ReadDir(osFS, "./ytdl")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}

				readOutput1, err := afero.ReadFile(osFS, fmt.Sprintf("./ytdl/%v", outIdx[0].Name()))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				reader1 := bytes.NewReader(readOutput1)

				readOutput2, err := afero.ReadFile(osFS, fmt.Sprintf("./ytdl/%v", outIdx[0].Name()))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				reader2 := bytes.NewReader(readOutput2)

				// report after code execution has ended
				// Create the embed templates
				fromField := discordgo.MessageEmbedField{
					Name:   "From",
					Value:  fmt.Sprintf("<@!%v>", userID),
					Inline: false,
				}
				toField := discordgo.MessageEmbedField{
					Name:   "To",
					Value:  fmt.Sprintf("<@!%v>", userData.ID),
					Inline: false,
				}
				timeElapsedField := discordgo.MessageEmbedField{
					Name:   "Execution Time",
					Value:  fmt.Sprintf("%v", execTime),
					Inline: false,
				}
				newsizeField := discordgo.MessageEmbedField{
					Name:   "New Size",
					Value:  fmt.Sprintf("%v KB (%v MB)", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
					Inline: false,
				}
				fileIDField := discordgo.MessageEmbedField{
					Name:   "File ID",
					Value:  fmt.Sprintf("`%v`", katInzVidID),
					Inline: false,
				}
				messageFields := []*discordgo.MessageEmbedField{&fromField, &toField, &timeElapsedField, &newsizeField, &fileIDField}

				aoiEmbedFooter := discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
				}

				aoiEmbeds := discordgo.MessageEmbed{
					Title:  "YT Katheryne",
					Color:  0x52ff91,
					Footer: &aoiEmbedFooter,
					Fields: messageFields,
				}

				// send a copy to Sender
				channel, err := ctx.Session.UserChannelCreate(userID)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}
				_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
				}
				ctx.Session.ChannelFileSend(channel.ID, outIdx[0].Name(), reader1)

				// send a copy to Receiver
				channel, err = ctx.Session.UserChannelCreate(userData.ID)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}
				_, err = ctx.Session.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) == universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
				}
				ctx.Session.ChannelFileSend(channel.ID, outIdx[0].Name(), reader2)
			}

		} else {

			// delete user's message and send confirmation as a reply
			scanLinks := ytRelax.FindAllString(rawArgs, -1)

			// get the video ID
			if strings.Contains(scanLinks[0], "www.youtube.com") {
				// sample URL >> https://www.youtube.com/watch?v=J5x0tLiItVY
				splitVidID := strings.Split(scanLinks[0], "youtube.com/watch?v=")
				katInzVidID = splitVidID[1]
			} else if strings.Contains(scanLinks[0], "youtu.be") {
				// sample URL >> https://youtu.be/J5x0tLiItVY
				splitVidID := strings.Split(scanLinks[0], "youtu.be/")
				katInzVidID = splitVidID[1]
			}
			ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, ctx.Event.ID)
			ctx.RespondText(fmt.Sprintf("Processing `%v`. Please wait.", katInzVidID))

			// run the code
			codeExec := time.Now()
			katYT, err := exec.Command("yt-dlp", "--ignore-config", "--no-playlist", "--max-filesize", "20m", "-P", "$PWD/ytdl", "-o", "%(id)s.%(ext)s", "-x", "--audio-format", "mp3", scanLinks[0]).Output()
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				return
			}
			fmt.Println(string(katYT))
			execTime := time.Since(codeExec)

			outIdx, err := afero.ReadDir(osFS, "./ytdl")
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				return
			}

			readOutput, err := afero.ReadFile(osFS, fmt.Sprintf("./ytdl/%v", outIdx[0].Name()))
			if err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				return
			}
			reader := bytes.NewReader(readOutput)

			// report after code execution has ended
			// Create the embed templates
			usernameField := discordgo.MessageEmbedField{
				Name:   "Username",
				Value:  fmt.Sprintf("<@!%v>", userID),
				Inline: false,
			}
			timeElapsedField := discordgo.MessageEmbedField{
				Name:   "Execution Time",
				Value:  fmt.Sprintf("%v", execTime),
				Inline: false,
			}
			newsizeField := discordgo.MessageEmbedField{
				Name:   "New Size",
				Value:  fmt.Sprintf("%v KB (%v MB)", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
				Inline: false,
			}
			fileIDField := discordgo.MessageEmbedField{
				Name:   "File ID",
				Value:  fmt.Sprintf("`%v`", katInzVidID),
				Inline: false,
			}
			messageFields := []*discordgo.MessageEmbedField{&usernameField, &timeElapsedField, &newsizeField, &fileIDField}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  "YT Katheryne",
				Color:  0x52ff91,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
			ctx.Session.ChannelFileSend(ctx.Event.ChannelID, outIdx[0].Name(), reader)
		}

	}

}

func katInzShowLastSender(ctx *dgc.Ctx) {

	userID := ctx.Event.Author.ID
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, ctx.Event.ID, "✅")

	osFS.RemoveAll("./logs")
	osFS.MkdirAll("./logs", 0777)

	// ==================================
	// Create a new logs.txt
	createLogsFile, err := osFS.Create("./logs/logs.txt")
	if err != nil {
		fmt.Println(" [ERROR] ", err)

		if len(universalLogs) == universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	} else {
		// Write to the file
		writeLogsFile, err := createLogsFile.WriteString(fmt.Sprintf("%v", maidsanLogs))
		if err != nil {
			fmt.Println(" [ERROR] ", err)

			if len(universalLogs) == universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}

			return
		} else {
			// Close the file
			if err := createLogsFile.Close(); err != nil {
				fmt.Println(" [ERROR] ", err)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				return
			} else {
				winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", nhImgName, (writeLogsFile / Kilobyte), (writeLogsFile / Megabyte))
				fmt.Println(winLogs)

				if len(universalLogs) == universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
				}
			}
		}
	}

	outIdx, err := afero.ReadDir(osFS, "./logs")
	if err != nil {
		fmt.Println(" [ERROR] ", err)

		if len(universalLogs) == universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}

	readOutput, err := afero.ReadFile(osFS, fmt.Sprintf("./logs/%v", outIdx[0].Name()))
	if err != nil {
		fmt.Println(" [ERROR] ", err)

		if len(universalLogs) == universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}
	reader := bytes.NewReader(readOutput)

	// add some checks to prevent panics.
	// panic: runtime error: index out of range [-2]
	if len(maidsanLogs) >= 2 {

		// report after code execution has ended
		// Create the embed templates
		usernameField := discordgo.MessageEmbedField{
			Name:   "Data Issuer",
			Value:  fmt.Sprintf("<@!%v>", userID),
			Inline: false,
		}
		lastsenderField := discordgo.MessageEmbedField{
			Name:   "Last Sender",
			Value:  fmt.Sprintf("<@!%v>", useridLogs[(len(useridLogs)-2)]),
			Inline: false,
		}
		timestampField := discordgo.MessageEmbedField{
			Name:   "Timestamp",
			Value:  fmt.Sprintf("%v", timestampLogs[(len(timestampLogs)-2)]),
			Inline: false,
		}
		pfpField := discordgo.MessageEmbedField{
			Name:   "Profile Picture",
			Value:  fmt.Sprintf("```\n%v\n```", profpicLogs[(len(profpicLogs)-2)]),
			Inline: false,
		}
		acctypeField := discordgo.MessageEmbedField{
			Name:   "Account Type",
			Value:  fmt.Sprintf("%v", acctypeLogs[(len(acctypeLogs)-2)]),
			Inline: false,
		}
		msgidField := discordgo.MessageEmbedField{
			Name:   "Message ID",
			Value:  fmt.Sprintf("%v", msgidLogs[(len(msgidLogs)-2)]),
			Inline: false,
		}
		msgcontentField := discordgo.MessageEmbedField{
			Name:   "Message",
			Value:  fmt.Sprintf("```\n%v\n```", msgLogs[(len(msgLogs)-2)]),
			Inline: false,
		}
		translateField := discordgo.MessageEmbedField{
			Name:   "Translation",
			Value:  fmt.Sprintf("```\n%v\n```", translateLogs[(len(translateLogs)-2)]),
			Inline: false,
		}
		logsindexField := discordgo.MessageEmbedField{
			Name:   "Logs Limit",
			Value:  fmt.Sprintf("`%v / %v`", len(maidsanLogs), maidsanLogsLimit),
			Inline: false,
		}
		logssizeField := discordgo.MessageEmbedField{
			Name:   "Logs Size",
			Value:  fmt.Sprintf("%v KB (%v MB)", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&usernameField, &lastsenderField, &timestampField, &pfpField, &acctypeField, &msgidField, &msgcontentField, &translateField, &logsindexField, &logssizeField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  "All Seeing Eyes of Katheryne",
			Color:  0x4287f5,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
		ctx.Session.ChannelFileSend(ctx.Event.ChannelID, outIdx[0].Name(), reader)
	} else if len(maidsanLogs) >= 1 {

		// report after code execution has ended
		// Create the embed templates
		usernameField := discordgo.MessageEmbedField{
			Name:   "Data Issuer",
			Value:  fmt.Sprintf("<@!%v>", userID),
			Inline: false,
		}
		lastsenderField := discordgo.MessageEmbedField{
			Name:   "Last Sender",
			Value:  fmt.Sprintf("<@!%v>", useridLogs[(len(useridLogs)-1)]),
			Inline: false,
		}
		timestampField := discordgo.MessageEmbedField{
			Name:   "Timestamp",
			Value:  fmt.Sprintf("%v", timestampLogs[(len(timestampLogs)-1)]),
			Inline: false,
		}
		pfpField := discordgo.MessageEmbedField{
			Name:   "Profile Picture",
			Value:  fmt.Sprintf("```\n%v\n```", profpicLogs[(len(profpicLogs)-1)]),
			Inline: false,
		}
		acctypeField := discordgo.MessageEmbedField{
			Name:   "Account Type",
			Value:  fmt.Sprintf("%v", acctypeLogs[(len(acctypeLogs)-1)]),
			Inline: false,
		}
		msgidField := discordgo.MessageEmbedField{
			Name:   "Message ID",
			Value:  fmt.Sprintf("%v", msgidLogs[(len(msgidLogs)-1)]),
			Inline: false,
		}
		msgcontentField := discordgo.MessageEmbedField{
			Name:   "Message",
			Value:  fmt.Sprintf("```\n%v\n```", msgLogs[(len(msgLogs)-1)]),
			Inline: false,
		}
		translateField := discordgo.MessageEmbedField{
			Name:   "Translation",
			Value:  fmt.Sprintf("```\n%v\n```", translateLogs[(len(translateLogs)-1)]),
			Inline: false,
		}
		logsindexField := discordgo.MessageEmbedField{
			Name:   "Logs Limit",
			Value:  fmt.Sprintf("`%v / %v`", len(maidsanLogs), maidsanLogsLimit),
			Inline: false,
		}
		logssizeField := discordgo.MessageEmbedField{
			Name:   "Logs Size",
			Value:  fmt.Sprintf("%v KB (%v MB)", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&usernameField, &lastsenderField, &timestampField, &pfpField, &acctypeField, &msgidField, &msgcontentField, &translateField, &logsindexField, &logssizeField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  "All Seeing Eyes of Katheryne",
			Color:  0x4287f5,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, &aoiEmbeds)
		ctx.Session.ChannelFileSend(ctx.Event.ChannelID, outIdx[0].Name(), reader)
	} else {
		ctx.RespondText(fmt.Sprintf("I couldn't get any data from my memory.\n```\nLogs Data: %v / %v\n```", len(maidsanLogs), maidsanLogsLimit))
	}

}
