package main

type Message struct {
	Payload string `json:"payload"`
	Name    string `json:"name"`
	Time    string `json:"time"`
}

func (m *Message) SetTime(time string) {
	m.Time = time
}
