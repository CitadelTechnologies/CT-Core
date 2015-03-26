package node

type Client struct{
    Id int
}

type Node struct{
    Id int
	Clients map[int]Client
}

func (n *Node) addClient(c *Client) bool{

    n.Clients[c.Id] = *c

    return true

}

func (n *Node) removeClient(id int) bool{

    delete(n.Clients, id)

    return true

}

func (n *Node) getClients() map[int]Client{

    return n.Clients

}

func (n *Node) getClient(id int) Client{

    return n.Clients[id]

}