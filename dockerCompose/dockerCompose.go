/*
 * MIT License
 *
 * Copyright (c) 2020 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package dockerCompose

import (
  "fmt"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/utils"
  "gopkg.in/yaml.v3"
  "io/ioutil"
  "regexp"
  "strings"
)

// based on https://github.com/digibib/docker-compose-dot/blob/master/docker-compose-dot.go

type DockerComposeTemplate struct {
  Version       *string
  Networks      *map[string]Network `yaml:"networks,omitempty"`
  Volumes       *map[string]Volume  `yaml:"volumes,omitempty"`
  Services      *map[string]Service `yaml:"services,omitempty"`
  Replacements  *map[string]string  `yaml:"-"`
  IsInSwarmMode bool                `yaml:"-"`
}

type Network struct {
  Driver           *string `yaml:"driver,omitempty"`
  External         *bool `yaml:"external,omitempty"`
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
  Deploy            *Deploy `yaml:"deploy,omitempty"`
  Ports             *[]string `yaml:"ports,omitempty"`
  Command           *interface{} `yaml:"command,omitempty"`
  ContainerName     *string `yaml:"container_name,omitempty"`
  DependsOn         *[]string `yaml:"depends_on,omitempty"`
  Environment       *interface{} `yaml:"environment,omitempty"`
  Restart           *string `yaml:"restart,omitempty"`
}

type RestartPolicy struct {
  Condition   *string `yaml:"condition,omitempty"`
  Delay       *string `yaml:"delay,omitempty"`
  MaxAttempts *int    `yaml:"max_attempts,omitempty"`
  Window      *string `yaml:"window,omitempty"`
}

// We only care about the labels and restart_policy right now
type Deploy struct {
  Labels        *[]string      `yaml:"labels,omitempty"`
  RestartPolicy *RestartPolicy `yaml:"restart_policy,omitempty"`
}

func LoadDockerComposeTemplate( path string, isInSwarmMode bool ) (*DockerComposeTemplate, error) {
  dockerComposeTemplateBytes, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  var dockerComposeTemplate DockerComposeTemplate

  err = yaml.Unmarshal(dockerComposeTemplateBytes, &dockerComposeTemplate)
  if err != nil {
    return nil,err
  }
  dockerComposeTemplate.IsInSwarmMode = isInSwarmMode

  for serviceKey, service := range *dockerComposeTemplate.Services {

    if isInSwarmMode {
      service.Restart = nil
      if service.Deploy == nil {
        service.Deploy = &Deploy{}
      }
      if service.Deploy.RestartPolicy == nil {
        service.Deploy.RestartPolicy = &RestartPolicy{}
      }
      con := "any"
      service.Deploy.RestartPolicy.Condition = &con
    } else {
      restart := "always"
      service.Restart = &restart
      service.Deploy = nil
    }

    (*dockerComposeTemplate.Services)[serviceKey]=service
  }


  dockerComposeTemplate.StripLabels( isInSwarmMode )
  return &dockerComposeTemplate,nil
}

func ( dockerComposeTemplate *DockerComposeTemplate ) StripLabels( isInSwarmMode bool ) {
  // remove all docker labels, except some allowed ones from template
  // might be security risk, since they are directly read by
  // docker and we want control over what gets passed to
  // docker

  // only labels in service.labels are used and rewritten for either
  // compose or swarm mode.
  // service.deploy.labels in the template are ignored

  allowedLabels := func( labels *[]string ) *[]string {
    var newLabels []string
    for _, label := range *labels {
      for _, allowedLabel := range globals.DOCKER_COMPOSE_ALLOWED_MAIN_SERVICE_LABELS {
        labelIsAllowed, _ := regexp.MatchString( allowedLabel, strings.Trim( label, " " ) )
        if labelIsAllowed {
          newLabels = append(newLabels, label)
        }
      }
    }
    return &newLabels
  }

  for serviceKey, service := range *dockerComposeTemplate.Services {
    foundMainService, err := regexp.MatchString(
      "^"+fmt.Sprintf(globals.DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE, "APP_UPSTREAM_HOST" )+"$" ,
      strings.Trim(serviceKey, " " ) )

    if foundMainService && err == nil {
      // allow certain labels for main service
      if isInSwarmMode {
        if service.Deploy == nil {
          service.Deploy = &Deploy{}
        }
        service.Deploy.Labels = allowedLabels( service.Labels )
        service.Labels = nil
      } else {
        if service.Labels == nil {
          service.Labels = &[]string{}
        }
        service.Labels = allowedLabels( service.Labels )
        if service.Deploy != nil {
          service.Deploy.Labels = nil
        }
      }

    } else {
      service.Labels = nil
      if service.Deploy != nil {
        service.Deploy.Labels = nil
      }
    }
    (*dockerComposeTemplate.Services)[serviceKey] = service
  }
}

func ( dockerComposeTemplate *DockerComposeTemplate ) CheckVolumes( trustZone string ) error {
  if dockerComposeTemplate.Services == nil {
    return nil
  }
  for _, service := range *dockerComposeTemplate.Services {
    if service.Volumes == nil {
      continue
    }
    output.Noticef( "Checking volumes for unallowed access\n" )
    for _, volume := range *service.Volumes {
      arr := strings.Split( volume, ":" )
      hostDirectory := strings.Trim( arr[0], " \n" )
      output.Noticef( "...%s\n", hostDirectory )
      if utils.SliceIndex( len(globals.DockerVolumeWhitelist), func(i int) bool {
        pattern := globals.DockerVolumeWhitelist[i]
        match, err := regexp.MatchString(pattern, hostDirectory)
        return match && err == nil
      } ) == -1 {
        return errors.VOLUME_NOT_IN_WHITELIST
      }

      needsCoreZone, err := regexp.MatchString(globals.TRUST_ZONE_CORE_PATTERN, hostDirectory)

      if err != nil {
        return err
      }

      needsTrustedZone, err := regexp.MatchString(globals.TRUST_ZONE_TRUSTED_PATTERN, hostDirectory)

      if err != nil {
        return err
      }

      needsServiceZone, err := regexp.MatchString(globals.TRUST_ZONE_SERVICE_PATTERN, hostDirectory)

      if err != nil {
        return err
      }

      if needsCoreZone &&
        trustZone != globals.TRUST_ZONE_CORE {
        return errors.APP_HAS_WRONG_TRUST_ZONE
      }

      if needsServiceZone &&
        trustZone != globals.TRUST_ZONE_SERVICE &&
        trustZone != globals.TRUST_ZONE_CORE {
        return errors.APP_HAS_WRONG_TRUST_ZONE
      }

      if needsTrustedZone &&
          trustZone != globals.TRUST_ZONE_TRUSTED &&
          trustZone != globals.TRUST_ZONE_SERVICE &&
          trustZone != globals.TRUST_ZONE_CORE {
        return errors.APP_HAS_WRONG_TRUST_ZONE
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

func ( dockerComposeTemplate *DockerComposeTemplate ) CheckNetworks( trustZone string, clientID string ) error {
  // TODO: implement
  if dockerComposeTemplate.Networks == nil {
    return nil
  }

  networkIsOk := func( networkName string, trustZone string, clientID string ) bool {

    // either use one of the predefined networks or
    // local networks prefixed with clientID

    if networkName == globals.CORE_NETWORK &&
        trustZone != globals.TRUST_ZONE_CORE {
      return false
    }
    if networkName == globals.SERVICE_NETWORK &&
      trustZone != globals.TRUST_ZONE_CORE &&
      trustZone != globals.TRUST_ZONE_SERVICE {
      return false
    }

    isLocalNetwork, _ := regexp.MatchString( "^"+clientID+"_", networkName )

    if networkName != globals.APPS_NETWORK &&
      !isLocalNetwork {
      return false
    }

    return true
  }

  // check network definitions
  output.Noticef( "Checking network definitions for unallowed access\n" )
  for networkName, _ := range *dockerComposeTemplate.Networks {
    output.Noticef( "...%s\n", networkName )
    if !networkIsOk(networkName, trustZone, clientID) {
      return errors.APP_HAS_WRONG_TRUST_ZONE
    }
  }

  if dockerComposeTemplate.Services == nil {
    return nil
  }

  // check service networks
  for _, service := range *dockerComposeTemplate.Services {
    if service.Networks == nil {
      continue
    }
    output.Noticef( "Checking services for unallowed access\n" )
    for _, networkName := range *service.Networks {
      output.Noticef( "...%s\n", networkName )
      if !networkIsOk(networkName, trustZone, clientID) {
        return errors.APP_HAS_WRONG_TRUST_ZONE
      }
    }
  }
  return nil
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

func ( dockerComposeTemplate *DockerComposeTemplate ) Replace( dockerComposeTemplateBytes []byte ) ([]byte, error) {

  if dockerComposeTemplate.Replacements == nil {
    return dockerComposeTemplateBytes, nil
  }

  for toReplace, replaceWith := range *dockerComposeTemplate.Replacements {
    r := regexp.MustCompile(fmt.Sprintf( globals.DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE, toReplace ))
    if r != nil {
      dockerComposeTemplateBytes = r.ReplaceAllLiteral( dockerComposeTemplateBytes, []byte(replaceWith) )
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
