package db

import("github.com/rs/xid")

func GenerateId() string{

	guid := xid.New()
	id := guid.String()

	return (id)
}
