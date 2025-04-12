package ldap

// RendersLdapServers handles rendering of LDAP servers in GitLab
type RendersLdapServers struct {
	ldapConfig LdapConfig
}

// LdapConfig defines the interface for LDAP configuration
type LdapConfig interface {
	SignInEnabled() bool
	AvailableServers() []string
}

// NewRendersLdapServers creates a new instance of RendersLdapServers
func NewRendersLdapServers(ldapConfig LdapConfig) *RendersLdapServers {
	return &RendersLdapServers{
		ldapConfig: ldapConfig,
	}
}

// LdapServers returns a list of available LDAP servers if LDAP sign-in is enabled,
// or an empty array otherwise
func (r *RendersLdapServers) LdapServers() []string {
	if r.ldapConfig.SignInEnabled() {
		return r.ldapConfig.AvailableServers()
	}

	return []string{}
}
