package clients

import "time"

// Client is a single client object
type Client struct {
	ID       string `json:"id"`
	Lastseen int64  `json:"lastseen"`
}

// ClientList holds all the clients
type ClientList struct {
	Clients []Client
}

// NewClient creates a new client
func NewClient(_id string) Client {
	return Client{_id, time.Now().Unix()}
}

// Ping updates the last time the device was seen
func (cl *ClientList) Ping(_id string) {
	for i, c := range cl.Clients {
		if c.ID == _id {
			cl.Clients[i].Lastseen = time.Now().Unix()
			return
		}
	}

	cl.AddClient(_id)
}

// AddClient adds another client
func (cl *ClientList) AddClient(_id string) {
	c := NewClient(_id)
	cl.Clients = append(cl.Clients, c)
}
