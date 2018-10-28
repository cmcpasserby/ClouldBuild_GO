package cli

import (
	"flag"
	"fmt"
	"github.com/cmcpasserby/CloudBuild_GO/pkg/cloudbuild"
	"gopkg.in/AlecAivazis/survey.v1"
	"reflect"
	"regexp"
)

var reProjectId = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
var reApiKey = regexp.MustCompile(`[0-9a-f]{32}`)

type Command struct {
	Name     string
	HelpText string
	Flags    *flag.FlagSet
	Action   func(flags map[string]string) error
}

func PopulateArgs(flags map[string]string, data interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	tt := v.Type()
	fCount := v.NumField()

	qs := make([]*survey.Question, 0, fCount)

	for i := 0; i < fCount; i++ {
		fName := tt.Field(i).Tag.Get("survey")
		if fName == "" {
			fName = tt.Field(i).Name
		}

		if val, ok := flags[fName]; ok {
			v.Field(i).SetString(val)
		} else {
			qs = append(qs, &survey.Question{
				Name:      fName,
				Prompt:    &survey.Input{Message: fName},
				Validate:  survey.Required,
				Transform: survey.ToLower,
			})
		}
	}

	if err := survey.Ask(qs, data); err != nil {
		return err
	}

	return nil
}

var CommandOrder = [...]string{"getCred", "listCreds", "updateCred", "uploadCred", "deleteCred"}

var Commands = map[string]Command{

	"getCred": {
		"getCred",
		"Get IOS Credential Detials",
		func() *flag.FlagSet {
			flags := CreateFlagSet("getCred")
			flags.String("projectId", "", "Project Id")
			flags.String("credId", "", "Credential Id")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey    string `survey:"apiKey"`
				OrgId     string `survey:"orgId"`
				ProjectId string `survey:"projectId"`
				CredId    string `survey:"credId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			cred, err := credsService.GetIOS(results.ProjectId, results.CredId)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", cred)

			return nil
		},
	},

	"listCreds": {
		"listCreds",
		"List all IOS Credentials",
		func() *flag.FlagSet {
			flags := CreateFlagSet("listCreds")
			flags.String("projectId", "", "Project Id")
			return flags
		}(),
		func(flags map[string]string) error {
			// parse args and settings, and question if needed
			results := struct {
				ApiKey    string `survey:"apiKey"`
				OrgId     string `survey:"orgId"`
				ProjectId string `survey:"projectId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			creds, err := credsService.GetAllIOS(results.ProjectId)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", creds)

			return nil
		},
	},

	"updateCred": {
		"updateCred",
		"Update a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("updateCred")
			flags.String("projectId", "", "Project Id")
			flags.String("certId", "", "Certificate Id")
			flags.String("label", "", "Label")
			flags.String("certPath", "", "Certificate Path")
			flags.String("profilePath", "", "Provisioning Profile Path")
			flags.String("certPass", "", "Certificate password")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey      string `survey:"apiKey"`
				OrgId       string `survey:"orgId"`
				ProjectId   string `survey:"projectId"`
				CertId      string `survey:"certId"`
				Label       string `survey:"label"`
				CertPath    string `survey:"certPath"`
				ProfilePath string `survey:"profilePath"`
				CertPass    string `survey:"certPass"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			cred, err := credsService.UpdateIOS(results.ProjectId, results.CertId, results.Label, results.CertId, results.ProfilePath, results.CertPass)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", cred)

			return nil
		},
	},

	"uploadCred": {
		"uploadCred",
		"Upload a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("uploadCred")
			flags.String("projectId", "", "Project Id")
			flags.String("label", "", "Label")
			flags.String("certPath", "", "Certificate Path")
			flags.String("profilePath", "", "Provisioning Profile Path")
			flags.String("certPass", "", "Certificate password")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey      string `survey:"apiKey"`
				OrgId       string `survey:"orgId"`
				ProjectId   string `survey:"projectId"`
				Label       string `survey:"label"`
				CertPath    string `survey:"certPath"`
				ProfilePath string `survey:"profilePath"`
				CertPass    string `survey:"certPass"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			cred, err := credsService.UploadIOS(results.ProjectId, results.Label, results.CertPath, results.ProfilePath, results.CertPass)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", cred)

			return nil
		},
	},

	"deleteCred": {
		"DeleteCred",
		"Delete a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("deleteCred")
			flags.String("projectId", "", "Project Id")
			flags.String("credId", "", "Credential Id")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey    string `survey:"apiKey"`
				OrgId     string `survey:"orgId"`
				ProjectId string `survey:"projectId"`
				CertId    string `survey:"certId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			resp, err := credsService.DeleteIOS(results.ProjectId, results.CertId)
			if err != nil {
				return err
			}

			fmt.Println(resp.Status)

			return nil
		},
	},
}
