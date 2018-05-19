package main

import (
	"github.com/urfave/cli"
	"firebase.google.com/go/auth"
	"context"
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"
	"strings"
	"strconv"
	"reflect"
)

var au *auth.Client

var MARKS map[bool]string

func setupUsers(ctx *cli.Context) (err error) {
	au, err = fb.Auth(context.Background())
	if err != nil {
		return err
	}

	MARKS = make(map[bool]string)
	MARKS[true] = "✓"
	MARKS[false] = "✗"

	return nil
}

func usersList(ctx *cli.Context) (err error) {
	users := au.Users(context.Background(), "")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"UID", MARKS[true] + "|" + MARKS[false] + " Email", "Name", "Active", "Claims"})

	for user, err := users.Next(); user != nil && err == nil; user, err = users.Next() {
		var claims string
		if user.CustomClaims != nil && len(user.CustomClaims) != 0 {
			for key, value := range user.CustomClaims {
				claims = fmt.Sprintf("%s%s=%v(%v)&", claims, key, value, reflect.TypeOf(value))
			}
			claims = claims[:len(claims)-1]
		}
		table.Append([]string{
			user.UID,
			MARKS[user.EmailVerified] + " " + user.Email,
			user.DisplayName,
			MARKS[!user.Disabled],
			claims,
		})
	}

	table.Render()

	return nil
}

func userUpdate(ctx *cli.Context) (err error) {
	uid := ctx.String("uid")
	if uid == "" || len(uid) < 28 {
		return fmt.Errorf("invalid UID supplied")
	}

	user, err := au.GetUser(context.Background(), uid)
	if err != nil {
		return err
	}

	updates := &auth.UserToUpdate{}

	if name := ctx.String("name"); name != "" {
		updates.DisplayName(name)
	}

	if email := ctx.String("email"); email != "" {
		updates.Email(email)
	}

	if disabled := ctx.Bool("disabled"); disabled != user.Disabled {
		updates.Disabled(disabled)
	}

	if verified := ctx.Bool("verified"); verified != user.EmailVerified {
		updates.EmailVerified(verified)
	}

	if claims := ctx.StringSlice("claim"); len(claims) > 0 {
		claimsMap := make(map[string]interface{})

		for _, claim := range claims {
			keyValue := strings.Split(claim, "=")
			if len(keyValue) != 2 {
				return fmt.Errorf("claim %s is not formatted correctly (key=value)\n", claim)
			}

			var value interface{}
			switch keyValue[1] {
			case "true":
				fallthrough
			case "t":
				value = true
			case "false":
				fallthrough
			case "f":
				value = false
			default:
				if i, err := strconv.Atoi(keyValue[1]); err == nil {
					value = i
				} else {
					value = keyValue[1]
				}
			}
			claimsMap[keyValue[0]] = value
		}

		updates.CustomClaims(claimsMap)
	}

	_, err = au.UpdateUser(context.Background(), uid, updates)

	return usersList(ctx)
}
