package main

import (
    "gopkg.in/yaml.v2"
    "fmt"
    "io/ioutil"
)

func (g *GitHubStruct)create_yaml(file string, req *CreateReq) error {
    // Write the github.yaml source file.
    g.github_source.Urls = make(map[string]string)
    g.github_source.Urls["github-base-url"] = g.Client.BaseURL.String()

    if orga := req.GithubOrganization; orga == "" {
        g.github_source.Organization = req.ForjjOrganization
    } else {
        g.github_source.Organization = req.GithubOrganization
    }

    // Ensure Infra is already in the list of repo managed.
    if g.github_source.Repos == nil {
        g.github_source.Repos = make(map[string]RepositoryStruct)
    }

    upstream := "git@" + g.Client.BaseURL.Host + ":" + g.github_source.Organization + "/" + req.ForjjInfra + ".git"
    infra, found := g.github_source.Repos[req.ForjjInfra]
    if ! found {
        infra = RepositoryStruct{
            Description: fmt.Sprintf("Infrastructure repository for Organization '%s' maintained by Forjj", g.github_source.Organization),
            UserGroups: make([]UserGroupStruct, 0),
        }
        infra.Name = req.ForjjInfra
    }
    infra.Upstream = upstream
    g.github_source.Repos[req.ForjjInfra] = infra


    d, err := yaml.Marshal(&g.github_source)
    if  err != nil {
        return fmt.Errorf("Unable to encode github data in yaml. %s", err)
    }

    if err := ioutil.WriteFile(file, d, 0644) ; err != nil {
        return fmt.Errorf("Unable to save '%s'. %s", file, err)
    }
    return nil
}

func (g *GitHubStruct)load_yaml(file string) error {
    d, err := ioutil.ReadFile(file)
    if err != nil {
        return fmt.Errorf("Unable to load '%s'. %s", file, err)
    }

    err = yaml.Unmarshal(d, &g.github_source)
    if  err != nil {
        return fmt.Errorf("Unable to decode github data in yaml. %s", err)
    }
    return nil
}