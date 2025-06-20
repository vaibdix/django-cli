package main

type Config struct {
	DefaultDjangoVersion string   `json:"default_django_version"`
	DefaultFeatures      []string `json:"default_features"`
	CreateTemplates      bool     `json:"create_templates"`
	CreateAppTemplates   bool     `json:"create_app_templates"`
	RunServer            bool     `json:"run_server"`
	InitializeGit        bool     `json:"initialize_git"`
	PreferUV             bool     `json:"prefer_uv"`
}
