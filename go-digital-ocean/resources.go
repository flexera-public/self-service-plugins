package main

import (
	"fmt"

	"github.com/rightscale/godo"
)

// Links
type Link struct {
	href string
}
type Links map[string]*Link

// Droplet resource, see https://github.com/digitalocean/godo/blob/master/droplets.go
// Add links to the droplet kernels, backups etc.
type Droplet struct {
	ID          int            `json:"id,float64,omitempty"`
	Name        string         `json:"name,omitempty"`
	Memory      int            `json:"memory,omitempty"`
	Vcpus       int            `json:"vcpus,omitempty"`
	Disk        int            `json:"disk,omitempty"`
	Region      *godo.Region   `json:"region,omitempty"`
	Image       *godo.Image    `json:"image,omitempty"`
	Size        *godo.Size     `json:"size,omitempty"`
	BackupIDs   []int          `json:"backup_ids,omitempty"`
	SnapshotIDs []int          `json:"snapshot_ids,omitempty"`
	Locked      bool           `json:"locked,bool,omitempty"`
	Status      string         `json:"status,omitempty"`
	Networks    *godo.Networks `json:"networks,omitempty"`
	ActionIDs   []int          `json:"action_ids,omitempty"`
	Created     string         `json:"created_at,omitempty"`
	Links       Links          `json:"links"`
}

// Droplet factory
func DropletFromApi(do *godo.Droplet) *Droplet {
	href := dropletHref(do.ID)
	links := Links{
		"kernels":   &Link{href: fmt.Sprintf("%s/kernels", href)},
		"snapshots": &Link{href: fmt.Sprintf("%s/snapshots", href)},
		"backups":   &Link{href: fmt.Sprintf("%s/backups", href)},
		"actions":   &Link{href: fmt.Sprintf("%s/actions", href)},
		"neighbors": &Link{href: fmt.Sprintf("%s/neighbors", href)},
	}
	d := Droplet{
		ID:          do.ID,
		Name:        do.Name,
		Memory:      do.Memory,
		Vcpus:       do.Vcpus,
		Disk:        do.Disk,
		Region:      do.Region,
		Image:       do.Image,
		Size:        do.Size,
		BackupIDs:   do.BackupIDs,
		SnapshotIDs: do.SnapshotIDs,
		Locked:      do.Locked,
		Status:      do.Status,
		Networks:    do.Networks,
		ActionIDs:   do.ActionIDs,
		Created:     do.Created,
		Links:       links,
	}
	return &d
}

// Image resource, see https://github.com/digitalocean/godo/blob/master/images.go
// Add links to the image actions
type Image struct {
	ID           int      `json:"id,float64,omitempty"`
	Name         string   `json:"name,omitempty"`
	Distribution string   `json:"distribution,omitempty"`
	Slug         string   `json:"slug,omitempty"`
	Public       bool     `json:"public,omitempty"`
	Regions      []string `json:"regions,omitempty"`
	Links        Links    `json:"links"`
}

// Image factory
func ImageFromApi(do *godo.Image) *Image {
	href := imageHref(do.ID)
	links := Links{
		"actions": &Link{href: fmt.Sprintf("%s/actions", href)},
	}
	d := Image{
		ID:           do.ID,
		Name:         do.Name,
		Distribution: do.Distribution,
		Slug:         do.Slug,
		Public:       do.Public,
		Links:        links,
	}
	return &d
}
