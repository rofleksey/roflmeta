package roflmeta

// EpisodeMetadata best attempt at extracting metadata from filename alone
// SHOULD follow these rules:
// * Season should either be a show title, a season name/number or empty if provided filename lacks information
// * Episode should be as short as possible, usually 0*\\d+ or non-numerical episode name
//
// MUST follow these rules:
// * Episode MUST be displayable to end user
// * Episode MUST NOT be blank for video files
// * Episode MUST BE BLANK for non-video files (as well as season)
type EpisodeMetadata struct {
	Season  string
	Episode string
}
