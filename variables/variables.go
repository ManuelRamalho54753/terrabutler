package variables

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"terrabutler/logger"
	"terrabutler/settings"
	"terrabutler/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type TemplateData struct {
	Env                 string
	Site                string
	Password            string
	Organization        string
	Region              string
	Profile             string
	FirebaseCredentials string
	MailPassword        string
}

func GeneratePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

func EncryptPassword(password, keyID, region, profile string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := kms.NewFromConfig(cfg)
	input := &kms.EncryptInput{
		KeyId:     aws.String(keyID),
		Plaintext: []byte(password),
	}
	output, err := client.Encrypt(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(output.CiphertextBlob)
	return encoded, nil
}

func GenerateEncryptedPassword(size int, keyID, region, profile string) (string, error) {
	pw := GeneratePassword(size)
	return EncryptPassword(pw, keyID, region, profile)
}

func GenerateVarFiles(env string) {
	root := os.Getenv("TERRABUTLER_ROOT")
	paths := utils.GetPaths(root)
	settingsPath := paths["settings"]

	cfg, err := settings.LoadSettings(settingsPath)
	if err != nil {
		logger.Log.Fatal("Error loading settings", zapError(err))
	}

	sites := cfg.Sites.Ordered
	if contains(sites, "inception") {
		sites = remove(sites, "inception")
	}

	firebase := cfg.Environments.Temporary.Secrets.FirebaseCredentials
	mail := cfg.Environments.Temporary.Secrets.MailPassword
	org := cfg.General.Organization
	region := cfg.Environments.Default.Region
	profile := cfg.Environments.Default.ProfileName
	keyID := cfg.General.SecretsKeyID

	encryptedPassword, err := GenerateEncryptedPassword(16, keyID, region, profile)
	if err != nil {
		logger.Log.Fatal("Error encrypting password", zapError(err))
	}

	tmplDir := filepath.Join(paths["root"], "configs", "templates")
	tmplEnv, err := template.ParseFiles(filepath.Join(tmplDir, "env.tpl"))
	if err != nil {
		logger.Log.Fatal("Error parsing env.tpl", zapError(err))
	}
	tmplSite, err := template.ParseFiles(filepath.Join(tmplDir, "site.tpl"))
	if err != nil {
		logger.Log.Fatal("Error parsing site.tpl", zapError(err))
	}

	// Generate env.tfvars
	envFile := filepath.Join(paths["variables"], fmt.Sprintf("%s-%s.tfvars", org, env))
	f1, err := os.Create(envFile)
	if err != nil {
		logger.Log.Fatal("Error creating env tfvars file", zapError(err))
	}
	defer f1.Close()
	tmplEnv.Execute(f1, TemplateData{
		Env:          env,
		Organization: org,
		Region:       region,
		Profile:      profile,
	})

	// Generate site tfvars
	for _, site := range sites {
		file := filepath.Join(paths["variables"], fmt.Sprintf("%s-%s-%s.tfvars", org, env, site))
		fh, err := os.Create(file)
		if err != nil {
			logger.Log.Error("Error creating site tfvars file", zapError(err))
			continue
		}
		defer fh.Close()

		tmplSite.Execute(fh, TemplateData{
			Env:                 env,
			Site:                site,
			Password:            encryptedPassword,
			Organization:        org,
			FirebaseCredentials: firebase,
			MailPassword:        mail,
		})
	}
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func remove(slice []string, s string) []string {
	var result []string
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

func zapError(err error) any {
	return map[string]any{"error": err.Error()}
}
