package database

type Container struct {
	DbObject `bson:",inline"`
	Cash     int
	Capacity int
	Weight   int
}

func (c *Container) SetCash(cash int) {
	c.writeLock(func() {
		c.Cash = cash
	})
}

func (c *Container) GetCash() int {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.Cash
}

func (c *Container) AddCash(amount int) {
	c.SetCash(c.GetCash() + amount)
}

func (c *Container) RemoveCash(amount int) {
	c.SetCash(c.GetCash() - amount)
}

func (c *Container) GetCapacity() int {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.Capacity
}

func (c *Container) SetCapacity(limit int) {
	c.writeLock(func() {
		c.Capacity = limit
	})
}
