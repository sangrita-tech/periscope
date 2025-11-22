package preprocess

func (c *Chain) AddStripComments() *Chain {
	return c.Add(newPreprocessStripComments())
}

func (c *Chain) AddMaskURL() *Chain {
	return c.Add(newPreprocessMaskURL())
}

func (c *Chain) AddCollapseEmptyLines() *Chain {
	return c.Add(newPreprocessCollapseEmptyLines())
}
