package alfred

type Items struct {
	Items    []Item `json:"items"`
	DebugMsg string `json:"debug_msg"`
}

type Item struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
}
