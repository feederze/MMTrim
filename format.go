package main

type Entry struct {
	Key         string `json:"key"`
	Original    string `json:"original"`
	Translation string `json:"translation"`
}

type tempEntry struct {
	// pageNum     int
	SentenceNum int
	str         string
}
