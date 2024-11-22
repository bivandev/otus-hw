package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainSuffix := "." + strings.ToLower(domain)

	decoder := json.NewDecoder(r)
	for {
		var user User
		if err := decoder.Decode(&user); err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		emailDomain := getEmailDomain(user.Email)
		if emailDomain != "" && strings.HasSuffix(emailDomain, domainSuffix) {
			result[emailDomain]++
		}
	}

	return result, nil
}

func getEmailDomain(email string) string {
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 || atIndex == len(email)-1 {
		return ""
	}

	return strings.ToLower(email[atIndex+1:])
}
