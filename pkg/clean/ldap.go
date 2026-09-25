package clean

import "strings"

func (c *Cleaner) redactLDAP(text string) string {
	if !c.enableLDAP && c.ldapMode == "" {
		// default on when DefaultConfig includes type LDAP; also apply if enableLDAP set
	}
	text = ldapDNRe.ReplaceAllStringFunc(text, func(v string) string {
		return c.Store.consistent("ldapdn", strings.ToLower(v), "ldapdn", modeOr(c.ldapMode), "cn=x-redacted-dn-x")
	})
	text = ldapFilterRe.ReplaceAllStringFunc(text, func(v string) string {
		m := ldapFilterRe.FindStringSubmatch(v)
		if len(m) != 4 {
			return v
		}
		tok := c.Store.consistent("ldapsearch", strings.ToLower(m[2]), "ldapsearch", modeOr(c.ldapMode), "x-redacted-ldap-x")
		return m[1] + tok + m[3]
	})
	return text
}

func modeOr(m string) string {
	if m == "" {
		return "Consistent"
	}
	return m
}
