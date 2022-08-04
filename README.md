<p align="center">
    <img src="logo.png" alt="roflmeta logo">
</p>
<h1 align="center">roflmeta</h1>
<p align="center">
      Episode Metadata extractor from filenames in Go.<br><br>
</p>
<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#installation">Installation</a> •
  <a href="#license">License</a>
</p>

--------

## Overview

This library allows to extract season and episode names from video file names.
This might be particularly useful if dealing with video files in torrents, as they almost never come presorted
in order of broadcast and there is no universal regex to parse episode and season info from them.

## How To Use

All functions return single or multiple instances of this struct:

```go
type EpisodeMetadata struct {
Season  string
Episode string
}
```

Season should either be a show title, a season name/number or empty if filename completely lacks information.

Episode should be as short as possible, usually `0*\\d+` or non-numerical episode name. It is never blank and can be
generally displayed in frontend as is.

There are two functions: "single" and "multiple". The former is straightforward:

```go
import "github.com/rofleksey/roflmeta"

metadata := roflmeta.ParseSingleEpisodeMetadata("[Judas] Hunter x Hunter (2011) - S01E012.mkv")
// EpisodeMetadata{season=01, episode=012}

metadata := roflmeta.ParseSingleEpisodeMetadata("[Samir755] Hellsing Ultimate 02.mkv")
// EpisodeMetadata{season=Hellsing Ultimate, episode=02}
```

The "multiple" function tries to figure out information by restoring template used to generate torrent filenames.
It fallbacks to single parser on failure:

```go
import "github.com/rofleksey/roflmeta"

filenames := []string{"S01E001.mkv", "S01E002.mkv", ..., "S01E148.mkv"}

metadataSlice, noFallback := roflmeta.ParseMultipleEpisodeMetadata(filenames)
// []EpisodeMetadata where season = 01, episode = 001...148
```

All functions ignore non-video files and return empty struct for them.

## Installation

```
go get github.com/rofleksey/roflmeta
```

## License

Apache 2.0, see [LICENSE](LICENSE). It allows you to use this code in proprietary projects.