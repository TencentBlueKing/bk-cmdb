package publisher

// New create a new publisher
func New() Publisher {

	pub := &publisher{senders: MapSender{}}

	return pub
}
