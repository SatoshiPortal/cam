package dockerCompose

import (
  "fmt"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/utils"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "regexp"
  "strings"
)

// based on https://github.com/digibib/docker-compose-dot/blob/master/docker-compose-dot.go

type DockerComposeTemplate struct {
  Version  *string
  Networks *map[string]Network `yaml:"networks,omitempty"`
  Volumes  *map[string]Volume `yaml:"volumes,omitempty"`
  Services *map[string]Service `yaml:"services,omitempty"`
  Replacements *map[string]string `yaml:"-"`
}

type Network struct {
  Driver           *string `yaml:"driver,omitempty"`
  External         *string `yaml:"external,omitempty"`
  DriverOpts       *map[string]string `yaml:"driver_opts,omitempty"`
}

type Volume struct {
  Driver           *string `yaml:"driver,omitempty"`
  External         *string `yaml:"external,omitempty"`
  DriverOpts       *map[string]string `yaml:"driver_opts,omitempty"`
}

type Service struct {
  Image             *string `yaml:"image,omitempty"`
  Networks          *[]string `yaml:"networks,omitempty"`
  Volumes           *[]string `yaml:"volumes,omitempty"`
  Labels            *[]string `yaml:"labels,omitempty"`
  Ports             *[]string `yaml:"ports,omitempty"`
  Command           *interface{} `yaml:"command,omitempty"`
  ContainerName     *string `yaml:"container_name,omitempty"`
  DependsOn         *[]string `yaml:"depends_on,omitempty"`
  Environment       *interface{} `yaml:"environment,omitempty"`
  Restart           *string `yaml:"restart,omitempty"`
}

func LoadDockerComposeTemplate( path string ) (*DockerComposeTemplate, error) {
  dockerComposeTemplateBytes, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  var dockerComposeTemplate DockerComposeTemplate

  err = yaml.Unmarshal(dockerComposeTemplateBytes, &dockerComposeTemplate)
  if err != nil {
    return nil,err
  }
  return &dockerComposeTemplate,nil
}

func ( dockerComposeTemplate *DockerComposeTemplate ) CheckVolumes() error {
  if dockerComposeTemplate.Services == nil {
    return nil
  }
  for _, service := range *dockerComposeTemplate.Services {
    if service.Volumes == nil {
      continue
    }
    for _, volume := range *service.Volumes {
      arr := strings.Split( volume, ":" )
      hostDirectory := strings.Trim( arr[0], " \n" )
      output.Noticef( "Checking: %s\n", hostDirectory )
      if utils.SliceIndex( len(globals.DockerVolumeWhitelist), func(i int) bool {
        pattern := globals.DockerVolumeWhitelist[i]
        match, err := regexp.MatchString(pattern, hostDirectory)
        return match && err == nil
      } ) == -1 {
        return errors.VOLUME_NOT_IN_WHITELIST
      }

      if strings.HasPrefix(hostDirectory, "$UNSAFE__" ) {
        output.Warningf( "Volume %s is marked as unsafe. Please make sure, this app is not malicious.\n", hostDirectory)
      }

      if utils.SliceIndex( len(globals.DockerVolumeElementBlacklist), func(i int) bool {
        return strings.Contains( hostDirectory, globals.DockerVolumeElementBlacklist[i] )
      } ) != -1 {
        return errors.VOLUME_HAS_ILLEGAL_ELEMENTS
      }
    }
  }
  return nil
}

func ( dockerComposeTemplate *DockerComposeTemplate ) CleanUnsafeEntries() {
  if dockerComposeTemplate.Services == nil {
    return
  }
  for _, service := range *dockerComposeTemplate.Services {
    if service.ContainerName != nil {
      service.ContainerName = nil
    }
  }
}

func ( dockerComposeTemplate *DockerComposeTemplate ) CheckServiceNames() error {
  if dockerComposeTemplate.Services == nil {
    return nil
  }
  for name, _ := range *dockerComposeTemplate.Services {
    r := regexp.MustCompile(fmt.Sprintf( globals.DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE, "(APP_ID|APP_UPSTREAM_HOST)" ))
    if r != nil && !r.MatchString( name ) {
      return errors.SERVICE_NAME_NOT_UNIQUE
    }
  }
  return nil
}

func ( dockerComposeTemplate *DockerComposeTemplate ) Replace( dockerComposeTemplateBytes []byte ) ([]byte, error) {

  if dockerComposeTemplate.Replacements == nil {
    return dockerComposeTemplateBytes, nil
  }

  for toReplace, replaceWith := range *dockerComposeTemplate.Replacements {
    r := regexp.MustCompile(fmt.Sprintf( globals.DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE, toReplace ))
    if r != nil {
      dockerComposeTemplateBytes = r.ReplaceAll( dockerComposeTemplateBytes, []byte(replaceWith) )
    }
  }
  return dockerComposeTemplateBytes, nil
}

func ( dockerComposeTemplate *DockerComposeTemplate ) SaveAsDockerCompose( path string ) error {
  dockerComposeTemplateBytes, err := yaml.Marshal(&dockerComposeTemplate)
  if err != nil {
    return err
  }
  dockerComposeTemplateBytes, err = dockerComposeTemplate.Replace( dockerComposeTemplateBytes )
  err = ioutil.WriteFile(path, dockerComposeTemplateBytes, 0655 )
  if err != nil {
    return err
  }
  return nil
}