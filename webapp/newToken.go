package main

import (
	"bitbucket.org/ansijax/rfidlab_telegramdi_backend/auth"

	"github.com/garyburd/redigo/redis"
)

// Generates a new len-character string in base 36 that's checked to be unique in the DB
func newUniqueToken(len int, base string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	token := auth.NewBase36(len)
	res, err := redis.String(conn.Do("GET", base+token))
	if err != nil && err != redis.ErrNil {
		return "", ErrDB
	}

	for i := 0; i < 1e3; i++ {
		token = auth.NewBase36(len)
		res, err = redis.String(conn.Do("GET", base+token))
		if err != nil && err != redis.ErrNil {
			return "", ErrDB
		}

		if res == "" {
			return token, nil
		}
	}

	return "", ErrGeneric
}

func newActivationToken(len int) (string, error) {
	return newUniqueToken(len, "webapp:temp:token:")
}

func newLoginToken(len int) (string, error) {
	return newUniqueToken(len, "webapp:auth:login:")
}

func newResetToken(len int) (string, error) {
	return newUniqueToken(len, "webapp:auth:reset:")
}

func newSessionToken(len int) (string, error) {
	return newUniqueToken(len, "webapp:auth:session:")
}
