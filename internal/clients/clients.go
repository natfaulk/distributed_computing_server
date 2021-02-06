package clients

import "time"

// Client is a single client object
type Client struct {
	ID            string `json:"id"`
	Lastseen      int64  `json:"lastseen"`
	Hostname      string `json:"hostname"`
	ClientVersion string `json:"version"`
	LastIP        string `json:"ip"`
}

// ClientList holds all the clients
type ClientList struct {
	Clients []Client
}

// NewClient creates a new client
func NewClient(_id string, _hostname string, _version string, _ip string) Client {
	return Client{_id, time.Now().Unix(), _hostname, _version, _ip}
}

// Ping updates the last time the device was seen
func (cl *ClientList) Ping(_id string, _hostname string, _version string, _ip string) {
	for i, c := range cl.Clients {
		if c.ID == _id {
			cl.Clients[i].Lastseen = time.Now().Unix()

			if _hostname != "" {
				cl.Clients[i].Hostname = _hostname
			}

			if _version != "" {
				cl.Clients[i].ClientVersion = _version
			}

			if _ip != "" {
				cl.Clients[i].LastIP = _ip
			}

			return
		}
	}

	cl.AddClient(_id, _hostname, _version, _ip)
}

// AddClient adds another client
func (cl *ClientList) AddClient(_id string, _hostname string, _version string, _ip string) {
	c := NewClient(_id, _hostname, _version, _ip)
	cl.Clients = append(cl.Clients, c)
}
