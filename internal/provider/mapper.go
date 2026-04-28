package provider

// MapAzure translates sovereign tiers to Azure VM specifications
func MapAzure(tier, region string) map[string]any {
	skus := map[string]string{
		"micro": "Standard_B1s",
		"small": "Standard_B2s", 
		"pro":   "Standard_D2s_v5",
	}
	
	regions := map[string]string{
		"us-east-1": "eastus",
		"west-europe": "westeurope",
		"central-europe": "germanywestcentral",
	}
	
	azureRegion := regions[region]
	if azureRegion == "" {
		azureRegion = "eastus" // default
	}
	
	return map[string]any{
		"apiVersion": "compute.azure.upbound.io/v1beta1",
		"kind":       "LinuxVirtualMachine",
		"metadata": map[string]any{
			"labels": map[string]any{
				"engine.io/provider": "azure",
				"engine.io/tier": tier,
			},
		},
		"spec": map[string]any{
			"forProvider": map[string]any{
				"location":            azureRegion,
				"size":               skus[tier],
				"adminUsername":       "engineadmin",
				"networkInterfaceIds":  []string{},
				"osDisk": map[string]any{
					"caching":              "ReadWrite",
					"storageAccountType":   "Premium_LRS",
					"diskSizeGB":           30,
				},
			},
		},
	}
}

// MapAWS translates sovereign tiers to AWS EC2 specifications
func MapAWS(tier, region string) map[string]any {
	skus := map[string]string{
		"micro": "t3.micro",
		"small": "t3.small",
		"pro":   "c6i.large",
	}
	
	amiMap := map[string]string{
		"us-east-1":      "ami-0c55b159cbfafe1f0", // Ubuntu 22.04 LTS
		"us-west-2":      "ami-0b5eea76982371e9b",
		"eu-west-1":      "ami-08ca3eda88f381d0f",
	}
	
	ami := amiMap[region]
	if ami == "" {
		ami = "ami-0c55b159cbfafe1f0" // default Ubuntu
	}
	
	return map[string]any{
		"apiVersion": "ec2.aws.upbound.io/v1beta1",
		"kind":       "Instance",
		"metadata": map[string]any{
			"labels": map[string]any{
				"engine.io/provider": "aws",
				"engine.io/tier": tier,
			},
		},
		"spec": map[string]any{
			"forProvider": map[string]any{
				"region":            region,
				"instanceType":      skus[tier],
				"ami":               ami,
				"rootBlockSize":     30,
				"associatePublicIpAddress": true,
				"subnetId":          "", // Will be filled by composition
			},
		},
	}
}

// MapGCP translates sovereign tiers to GCP GCE specifications
func MapGCP(tier, region string) map[string]any {
	skus := map[string]string{
		"micro": "e2-micro",
		"small": "e2-small",
		"pro":   "n2-standard-2",
	}
	
	zoneMap := map[string]string{
		"us-east-1": "us-east1-b",
		"us-west-2": "us-west2-a",
		"europe-west1": "europe-west1-b",
	}
	
	zone := zoneMap[region]
	if zone == "" {
		zone = "us-east1-b" // default
	}
	
	return map[string]any{
		"apiVersion": "compute.gcp.upbound.io/v1beta1",
		"kind":       "Instance",
		"metadata": map[string]any{
			"labels": map[string]any{
				"engine.io/provider": "gcp",
				"engine.io/tier": tier,
			},
		},
		"spec": map[string]any{
			"forProvider": map[string]any{
				"zone":        zone,
				"machineType": skus[tier],
				"bootDisk": map[string]any{
					"initializeParams": map[string]any{
						"sizeGb": 30,
						"type":   "pd-balanced",
						"image":  "projects/ubuntu-os-cloud/global/images/ubuntu-2204-jammy-v20240101",
					},
				},
				"networkInterface": []map[string]any{
					{
						"network": "default",
						"accessConfig": []map[string]any{
							{"type": "ONE_TO_ONE_NAT"},
						},
					},
				},
			},
		},
	}
}

// MapHetzner translates sovereign tiers to Hetzner Cloud specifications
func MapHetzner(tier, region string) map[string]any {
	skus := map[string]string{
		"micro": "cx11",
		"small": "cpx11",
		"pro":   "cpx21",
	}
	
	locationMap := map[string]string{
		"us-east-1": "ash",
		"europe":    "nbg1",
		"central":   "hel1",
	}
	
	location := locationMap[region]
	if location == "" {
		location = "nbg1" // default to Nuremberg
	}
	
	return map[string]any{
		"apiVersion": "server.hetzner.upbound.io/v1alpha1",
		"kind":       "Server",
		"metadata": map[string]any{
			"labels": map[string]any{
				"engine.io/provider": "hetzner",
				"engine.io/tier": tier,
			},
		},
		"spec": map[string]any{
			"forProvider": map[string]any{
				"serverType": skus[tier],
				"location":   location,
				"image":      "ubuntu-22.04",
				"sshKeys":    []string{},
			},
		},
	}
}

// MapOVH translates sovereign tiers to OVHcloud specifications
func MapOVH(tier, region string) map[string]any {
	skus := map[string]string{
		"micro": "s1-2",
		"small": "s1-4",
		"pro":   "s1-8",
	}
	
	regionMap := map[string]string{
		"us-east-1": "US-EAST-VA-1",
		"europe":    "EU-West-1",
		"canada":    "CA-East-1",
	}
	
	ovhRegion := regionMap[region]
	if ovhRegion == "" {
		ovhRegion = "EU-West-1" // default
	}
	
	return map[string]any{
		"apiVersion": "instance.ovh.upbound.io/v1beta1",
		"kind":       "Instance",
		"metadata": map[string]any{
			"labels": map[string]any{
				"engine.io/provider": "ovh",
				"engine.io/tier": tier,
			},
		},
		"spec": map[string]any{
			"forProvider": map[string]any{
				"region":     ovhRegion,
				"flavor":     skus[tier],
				"image":      "Ubuntu 22.04",
				"monthlyBilling": false,
			},
		},
	}
}

// MapDigitalOcean translates sovereign tiers to DigitalOcean specifications
func MapDigitalOcean(tier, region string) map[string]any {
	skus := map[string]string{
		"micro": "s-1vcpu-1gb",
		"small": "s-1vcpu-2gb",
		"pro":   "s-2vcpu-4gb",
	}
	
	regionMap := map[string]string{
		"us-east-1": "nyc1",
		"us-west":   "sfo2",
		"europe":    "ams3",
		"asia":      "sgp1",
	}
	
	doRegion := regionMap[region]
	if doRegion == "" {
		doRegion = "nyc1" // default
	}
	
	return map[string]any{
		"apiVersion": "digitalocean.upbound.io/v1beta1",
		"kind":       "Droplet",
		"metadata": map[string]any{
			"labels": map[string]any{
				"engine.io/provider": "digitalocean",
				"engine.io/tier": tier,
			},
		},
		"spec": map[string]any{
			"forProvider": map[string]any{
				"region":    doRegion,
				"size":      skus[tier],
				"image":     "ubuntu-22-04-x64",
				"monitoring": true,
			},
		},
	}
}
