package clean

func init() {
	// documented: ObfuscateText must call redactLDAP; see process wrap below.
}

func (c *Cleaner) wrapText(text string) string {
	return c.redactLDAP(text)
}
