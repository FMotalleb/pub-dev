/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultTokenSize = 32
	defaultTokenKind = "bcrypt"
)

// tokenCmd represents the token command.
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "generates a random token",
	RunE: func(cmd *cobra.Command, _ []string) error {
		var length uint
		var kind, token, matcher string
		var err error
		if length, err = cmd.Flags().GetUint("length"); err != nil {
			return err
		}
		if kind, err = cmd.Flags().GetString("kind"); err != nil {
			return err
		}
		if token, err = randomString(length); err != nil {
			return err
		}
		if matcher, err = generateHash(kind, token); err != nil {
			return err
		}
		matcher = strings.ReplaceAll(matcher, "$", "$$")
		fmt.Fprintf(os.Stdout, "client token:`%s`\n", token)
		fmt.Fprintf(os.Stdout, "server hash:`%s:%s`\n", kind, matcher)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	tokenCmd.Flags().StringP("kind", "k", defaultTokenKind, "server side hash format. (sha256,bcrypt,plain)")
	tokenCmd.Flags().UintP("length", "l", defaultTokenSize, "length of the generated token")
}

func generateHash(hashType, input string) (string, error) {
	switch hashType {
	case "sha256":
		// Removed "md5" case due to weak cryptographic primitive (G401)
		sum := sha256.Sum256([]byte(input))
		return hex.EncodeToString(sum[:]), nil
	case "bcrypt":
		hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(hash), nil
	default:
		return "", errors.New("unsupported hash type")
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._+/=-"

func randomString(length uint) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
