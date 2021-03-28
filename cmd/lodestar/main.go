package main

import (
	"errors"
	"fmt"
	app "github.com/lodestar-cli/lodestar/internal/app"
	"github.com/urfave/cli/v2"
	lodestarDir "github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
	"log"
	"os"
)

func main() {
	var tag string
	var username string
	var token string
	var srcEnv string
	var name string
	var appConfigPath string
	var destEnv string
	var environment string

	app := &cli.App{
		Name: "lodestar",
		Version: "0.1.0",
		Usage: "Help guide your applications through their environments",
		Commands: []*cli.Command{
			{
				Name:        "app",
				Usage:       "Manage application images",
				Subcommands: []*cli.Command{
					{
						Name:  "push",
						Usage: "Update an App Environment's configuration file with new tag",
						UsageText: "In order to push a tag to an environment, either a name for an App configured in ~/.lodestar\n\t"+
							       " needs to be provided with --name, or a path to an app needs to be provided with --config-path.\n\t"+
							       " Lodestar will then be able to find the App and pass the tag to the correct environment.",
						Flags: []cli.Flag {
							&cli.StringFlag{
								Name: "username",
								Hidden: true,
								Aliases: []string{"u"},
								Usage: "`username` for the version control account that can access the repository",
								Required: true,
								Destination: &username,
								EnvVars: []string{"GIT_USER"},
							},
							&cli.StringFlag{
								Name: "token",
								Hidden: true,
								Aliases: []string{"t"},
								Usage: "`token` for the version control account that can access the repository",
								Required: true,
								Destination: &token,
								EnvVars: []string{"GIT_TOKEN"},
							},
							&cli.StringFlag{
								Name: "name",
								Usage: "the `name` of an app",
								Destination: &name,
							},
							&cli.StringFlag{
								Name: "config-path",
								Usage: "the `path` to the app configuration file",
								Destination: &appConfigPath,
							},
							&cli.StringFlag{
								Name: "environment",
								Aliases: []string{"env"},
								Usage: "the `environment` the tag will be pushed to",
								Required: true,
								Destination: &environment,
							},
							&cli.StringFlag{
								Name: "tag",
								Usage: "the `tag` for the new image",
								Required: true,
								Destination: &tag,
							},
						},
						Action: func(c *cli.Context) error {
							var config *app.LodestarAppConfig

							if name == "" && appConfigPath == "" {
								return errors.New("Must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
							} else if appConfigPath != ""{
								config, err := app.GetAppConfig(appConfigPath)
								if err != nil {
									return err
								}
								if len(config.EnvGraph) < 1 {
									return errors.New("No environments are provided for "+config.AppInfo.Name)
								}

								for _, env := range config.EnvGraph{
									if env.Name == environment {
										err := app.Push(username,token,config.AppInfo.RepoUrl,env.SrcPath,tag)
										if err != nil {
											return err
										}
										break
									}
								}
								return nil
							} else{
								path, err := lodestarDir.GetConfigPath("app", name)
								if err != nil {
									return err
								}
								fmt.Printf("Retrieving config for %s...\n", name)
								config, err = app.GetAppConfig(path)
								if err != nil {
									return err
								}
								if len(config.EnvGraph) < 1 {
									return errors.New("No environments are provided for "+name)
								}

								for _, env := range config.EnvGraph{
									if env.Name == environment {
										err := app.Push(username,token,config.AppInfo.RepoUrl,env.SrcPath,tag)
										if err != nil {
											return err
										}
										break
									}
								}
								return nil
							}
						},
					},
					{
						Name:  "promote",
						Usage: "promote an image tag to the next environment",
						UsageText: "Retrieves an application tag specified in a source environment configuration file, and promotes it to a destination configuration file",
						Flags: []cli.Flag {
							&cli.StringFlag{
								Name: "username",
								Usage: "`username` for the version control account that can access the repository",
								Hidden: true,
								Required: true,
								Destination: &username,
								EnvVars: []string{"GIT_USER"},
							},
							&cli.StringFlag{
								Name: "token",
								Usage: "`token` for the version control account that can access the repository",
								Hidden: true,
								Required: true,
								Destination: &token,
								EnvVars: []string{"GIT_TOKEN"},
							},
							&cli.StringFlag{
								Name: "name",
								Usage: "the `name` of an app",
								Destination: &name,
							},
							&cli.StringFlag{
								Name: "config-path",
								Usage: "the `path` to the app configuration file",
								Destination: &appConfigPath,
							},
							&cli.StringFlag{
								Name: "src-env",
								Usage: "the `name` of the source environment",
								Required: true,
								Destination: &srcEnv,
							},
							&cli.StringFlag{
								Name: "dest-env",
								Usage: "the `name` of the destination",
								Required: true,
								Destination: &destEnv,
							},
						},
						Action: func(c *cli.Context) error {
							var config *app.LodestarAppConfig
							var srcPath string
							var destPath string
							if name == "" && appConfigPath == "" {
								return errors.New("Must provide an App name or a path to a configuration file. For more information, run: lodestar app push --help")
							} else if appConfigPath != ""{
								config, err := app.GetAppConfig(appConfigPath)
								if err != nil {
									return err
								}
								if len(config.EnvGraph) < 1 {
									return errors.New("No environments are provided for "+config.AppInfo.Name)
								}

								for _, env := range config.EnvGraph {
									if env.Name == srcEnv {
										srcPath=env.SrcPath
									}else if env.Name == destEnv {
										destPath=env.SrcPath
									}
									if srcPath != "" && destPath != "" {
										break
									}
								}
								err = app.Promote(username,token,config.AppInfo.RepoUrl,srcPath,destPath)
								return err
							} else {
								path, err := lodestarDir.GetConfigPath("app", name)
								if err != nil {
									return err
								}
								fmt.Printf("Retrieving config for %s...\n", name)
								config, err = app.GetAppConfig(path)
								if err != nil {
									return err
								}
								if len(config.EnvGraph) < 1 {
									return errors.New("No environments are provided for " + name)
								}

								for _, env := range config.EnvGraph {
									if env.Name == srcEnv {
										srcPath=env.SrcPath
									}else if env.Name == destEnv {
										destPath=env.SrcPath
									}
									if srcPath != "" && destPath != "" {
										break
									}
								}
								err = app.Promote(username,token,config.AppInfo.RepoUrl,srcPath,destPath)
								return err
							}
						},
					},
					{
						Name:  "list",
						Usage: "list all Apps in current context",
						Action: func(c *cli.Context) error {
							err := app.List()
							return err
						},
					},
					{
						Name:  "show",
						Usage: "prints the configuration file for the specified App",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "name",
								Usage:       "the `name` of the app",
								Required:    true,
								Destination: &name,
							},
						},
						Action: func(c *cli.Context) error {

							if name == "" {
								return errors.New("Must provide an App name. \n For more information, run: lodestar app push --help")
							}
							err := app.Show(name)
							return err
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}