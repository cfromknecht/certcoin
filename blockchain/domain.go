package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"

	"errors"
	"regexp"
	"strings"
)

type Identity struct {
	Domain    crypto.SHA256Sum
	Subdomain crypto.SHA256Sum
}

func NewIdentity(domain, subdomain string) (Identity, error) {
	if len(domain) > 32 || len(subdomain) > 32 {
		return Identity{}, errors.New("Domain and subdomain cannot be longer than 32 bytes")
	}

	domain = strings.ToUpper(domain)
	subdomain = strings.ToUpper(subdomain)

	if !validDomain(domain) {
		return Identity{}, errors.New("Invalid domain")
	}

	if !validSubdomain(subdomain) {
		return Identity{}, errors.New("Invalid subdomain")
	}

	domainBytes := crypto.SHA256Sum{}
	subdomainBytes := crypto.SHA256Sum{}

	copy(domainBytes[:], []byte(domain))
	copy(subdomainBytes[:], []byte(subdomain))

	return Identity{
		Domain:    domainBytes,
		Subdomain: subdomainBytes,
	}, nil
}

func (i Identity) DomainStr() string {
	return string(i.Domain[:])
}

func (i Identity) SubdomainStr() string {
	return string(i.Subdomain[:])
}

func (i Identity) FullNameStr() string {
	return string(i.Domain[:]) + string(i.Subdomain[:])
}

func validDomain(d string) bool {
	r, err := regexp.Compile(`^[A-Z0-9.-]+\.[A-Z]{2,}$`)
	if err != nil {
		panic(err)
	}

	return r.MatchString(d)
}

func validSubdomain(s string) bool {
	r, err := regexp.Compile(`^(([A-Z0-9._%+-]+@)|([A-Z0-9.-]+))?$`)
	if err != nil {
		panic(err)
	}

	return r.MatchString(s)
}
