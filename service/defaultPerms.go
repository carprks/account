package service

import (
	permissions "github.com/carprks/permissions/service"
)

func getAccountPerms(ident string) []permissions.Permission {
	return []permissions.Permission{
		{
			Name:       "account",
			Action:     "login",
			Identifier: ident,
		},
		{
			Name:       "account",
			Action:     "edit",
			Identifier: ident,
		},
		{
			Name:       "account",
			Action:     "view",
			Identifier: ident,
		},
		{
			Name:       "payments",
			Action:     "create",
			Identifier: ident,
		},
		{
			Name:       "payments",
			Action:     "view",
			Identifier: ident,
		},
		{
			Name:       "payments",
			Action:     "report",
			Identifier: ident,
		},
	}
}

func getBookingPerms(ident string) []permissions.Permission {
	return []permissions.Permission{
		{
			Name:       "carparks",
			Action:     "book",
			Identifier: "*",
		},
		{
			Name:       "carparks",
			Action:     "report",
			Identifier: "*",
		},
	}
}

func getDefaultPerms(ident string) []permissions.Permission {
	perms := []permissions.Permission{}

	// Account perms
	accountPerms := getAccountPerms(ident)
	for _, perm := range accountPerms {
		perms = append(perms, perm)
	}

	// Carpark Perms
	carparkPerms := getBookingPerms(ident)
	for _, perm := range carparkPerms {
		perms = append(perms, perm)
	}

	return perms
}
