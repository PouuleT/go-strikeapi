Strike API client
=========

This is a wrapper around the _Strike API v2_ (https://getstrike.net/) written in go and based on their documentation [documentation](https://getstrike.net/api/)

## Search Torrents

```
	// Search, SearchWithCategory and SearchWithCategoryAndSubCategory will search torrents
	torrentList, err = strikeapi.SearchWithCategory("Slackware 14.1 x86_64 DVD ISO", strikeapi.Applications)
	for _, t := range torrentList {
		log.Printf("Got : %+v", t)
	}
```

## Get informations from torrent hash

```
	// Get torrent informations by its hash
	torrent, err := strikeapi.GetTorrentInfos("B425907E5755031BDA4A8D1B6DCCACA97DA14C04")
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	log.Printf("Informations about the Torrent : %+v", torrent)

	// Get torrents informations by a list of hashes
	torrentList, err := strikeapi.GetTorrentsInfos([]string{"B425907E5755031BDA4A8D1B6DCCACA97DA14C04", "156B69B8643BD11849A5D8F2122E13F"})
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	for _, t := range torrentList {
		log.Printf("Informations about the Torrent : %+v", t)
	}
```

## Get Description of a torrent

```
	// GetDescription will return the description of a given hash or of a given torrent
	desc, err := strikeapi.GetDescription("B425907E5755031BDA4A8D1B6DCCACA97DA14C04")
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	// And the equivalent method on a Torrent object :
	desc, err = torrent.GetDescription()
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	log.Printf("Description : %s", desc)
```

## Get Download link

```
	// GetDownloadLink will return the download link of a .torrent of a given hash or of a given torrent
	downloadLink, err := strikeapi.GetDownloadLink("B425907E5755031BDA4A8D1B6DCCACA97DA14C04")
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	// And the equivalent method on a Torrent object :
	downloadLink, err = torrent.GetDownloadLink()
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	log.Printf("Download link : %s", downloadLink)
```

## Get the number of indexed torrents

```
	// Get how many torrents are indexed
	torrentNb, err := strikeapi.CountTorrents()
	if err != nil {
		log.Fatal("Got error : ", err)
	}
	log.Printf("Number of Torrents indexed : %+v", torrentNb)
```
