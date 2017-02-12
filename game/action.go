package game

type Action struct {
	Requirements string
	Name string
	TextInstigatorF string
	TextInstigatorM *string
	TextReceiverF string
	TextReceiverM *string
	Action string
	Reactions []string
}
