package dockerCompose

// based on https://github.com/digibib/docker-compose-dot/blob/master/docker-compose-dot.go

type Config struct {
  Version  string
  Networks map[string]Network
  Volumes  map[string]Volume
  Services map[string]Service
}

type Network struct {
  Driver, External string
  DriverOpts       map[string]string "driver_opts"
}

type Volume struct {
  Driver, External string
  DriverOpts       map[string]string "driver_opts"
}

type Service struct {
  Networks, Volumes []string
  VolumesFrom       []string "volumes_from"
}
