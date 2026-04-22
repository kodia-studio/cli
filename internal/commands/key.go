package commands

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var keyGenerateCmd = &cobra.Command{
	Use:   "key:generate",
	Short: "Generate secure APP_KEY and JWT_SECRET for your application",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("🔐 Generating Secure Application Keys...")

		// Generate 32-byte key for APP_KEY (64 hex chars)
		appKey := generateRandomKey(32)
		// Generate 32-byte key for JWT_SECRET
		jwtSecret := generateRandomKey(32)

		// Set them using the logic from env.go
		if err := updateEnvFile("backend/.env", "APP_KEY", appKey); err != nil {
			color.Red("Failed to set APP_KEY: %v", err)
			return
		}

		if err := updateEnvFile("backend/.env", "JWT_SECRET", jwtSecret); err != nil {
			color.Red("Failed to set JWT_SECRET: %v", err)
			return
		}

		color.Green("✅ APP_KEY and JWT_SECRET generated and set in .env")
		color.HiBlack("APP_KEY: %s", appKey)
		color.HiBlack("JWT_SECRET: %s", jwtSecret)
	},
}

func generateRandomKey(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func init() {
	rootCmd.AddCommand(keyGenerateCmd)
}
