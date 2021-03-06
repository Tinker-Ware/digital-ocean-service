package usecases

import (
	"github.com/Tinker-Ware/digital-ocean-service/domain"
	"github.com/digitalocean/godo"
	"github.com/jinzhu/copier"
)

// Instance represens a created instance in any provider
type Instance struct {
	Provider string `json:"provider,omitempty"`
	domain.Droplet
}

// CreateDroplet creates a doplet in Digital Ocean
func (interactor DOInteractor) CreateDroplet(droplet domain.DropletRequest, token string) (*Instance, error) {
	client := getClient(token)

	dropletRequest := &godo.DropletCreateRequest{
		Name:              droplet.Name,
		Region:            droplet.Region,
		Size:              droplet.Size,
		Backups:           droplet.Backups,
		IPv6:              droplet.IPv6,
		PrivateNetworking: droplet.PrivateNetworking,
		UserData:          droplet.UserData,
		Image: godo.DropletCreateImage{
			Slug: droplet.Image,
		},
	}

	sshkeys := []godo.DropletCreateSSHKey{}
	for _, key := range droplet.SSHKeys {
		k := godo.DropletCreateSSHKey{
			Fingerprint: key.Fingerprint,
		}
		sshkeys = append(sshkeys, k)
	}
	dropletRequest.SSHKeys = sshkeys

	drop, _, err := client.Droplets.Create(dropletRequest)

	if err != nil {
		return nil, err
	}

	inst := &Instance{
		Provider: "digital_ocean",
		Droplet: domain.Droplet{
			ID:                drop.ID,
			Name:              droplet.Name,
			Region:            droplet.Region,
			OperatingSystem:   drop.Image.Slug,
			PrivateNetworking: false,
			InstanceName:      drop.Size.Slug,
		},
	}

	networksV4 := []domain.NetworkV4{}
	for _, net := range drop.Networks.V4 {
		n := domain.NetworkV4{}
		copier.Copy(n, net)
		networksV4 = append(networksV4, n)
	}

	networksV6 := []domain.NetworkV6{}
	for _, net := range drop.Networks.V6 {
		n := domain.NetworkV6{}
		copier.Copy(n, net)
		networksV6 = append(networksV6, n)
	}
	networks := domain.Networks{
		V4: networksV4,
		V6: networksV6,
	}

	inst.Networks = networks
	inst.SSHKeys = droplet.SSHKeys
	return inst, nil

}

// ListDroplets lists all the droplets a user has in Digital Ocean
func (interactor DOInteractor) ListDroplets(token string) ([]domain.Droplet, error) {

	client := getClient(token)

	doDrops, _, err := client.Droplets.List(nil)
	if err != nil {
		return nil, err
	}
	droplets := []domain.Droplet{}

	for _, drops := range doDrops {
		drp := domain.Droplet{
			Name:         drops.Name,
			Region:       drops.Region.String(),
			InstanceName: drops.Size.String(),
		}

		networksV4 := []domain.NetworkV4{}
		for _, net := range drops.Networks.V4 {
			n := domain.NetworkV4{}
			copier.Copy(n, net)
			networksV4 = append(networksV4, n)
		}

		networksV6 := []domain.NetworkV6{}
		for _, net := range drops.Networks.V6 {
			n := domain.NetworkV6{}
			copier.Copy(n, net)
			networksV6 = append(networksV6, n)
		}

		droplets = append(droplets, drp)

	}

	return droplets, nil
}

// GetDroplet gets a single droplet
func (interactor DOInteractor) GetDroplet(id int, token string) (*Instance, error) {
	client := getClient(token)

	droplet, _, err := client.Droplets.Get(id)
	if err != nil {
		return nil, err
	}

	instance := Instance{
		Provider: "digital_ocean",
		Droplet: domain.Droplet{
			ID:                droplet.ID,
			Name:              droplet.Name,
			Region:            droplet.Region.Slug,
			OperatingSystem:   droplet.Image.Slug,
			PrivateNetworking: false,
			InstanceName:      droplet.Size.Slug,
		},
	}

	networksV4 := []domain.NetworkV4{}
	for _, net := range droplet.Networks.V4 {
		n := domain.NetworkV4{
			IPAddress: net.IPAddress,
			Netmask:   net.Netmask,
			Gateway:   net.Gateway,
			Type:      net.Type,
		}

		networksV4 = append(networksV4, n)
	}

	networksV6 := []domain.NetworkV6{}
	for _, net := range droplet.Networks.V6 {
		n := domain.NetworkV6{
			IPAddress: net.IPAddress,
			Netmask:   net.Netmask,
			Gateway:   net.Gateway,
			Type:      net.Type,
		}

		networksV6 = append(networksV6, n)

	}
	networks := domain.Networks{
		V4: networksV4,
		V6: networksV6,
	}

	instance.Networks = networks

	return &instance, nil
}

// DestroyDroplet destroys a droplet
func (interactor DOInteractor) DestroyDroplet(id int, token string) error {
	client := getClient(token)
	_, err := client.Droplets.Delete(id)
	return err
}
