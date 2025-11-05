package evaluate

var (
	CategorySystemInfo = CategorySpec{
		ID:              "information.system",
		Title:           "Information Gathering - System Info",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           100,
	}
	CategoryServices = CategorySpec{
		ID:              "information.services",
		Title:           "Information Gathering - Services",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           200,
	}
	CategoryCommands = CategorySpec{
		ID:              "information.commands",
		Title:           "Information Gathering - Commands and Capabilities",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           300,
	}
	CategoryMounts = CategorySpec{
		ID:              "information.mounts",
		Title:           "Information Gathering - Mounts",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           400,
	}
	CategoryNetNamespace = CategorySpec{
		ID:              "information.netns",
		Title:           "Information Gathering - Net Namespace",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           500,
	}
	CategorySysctl = CategorySpec{
		ID:              "information.sysctl",
		Title:           "Information Gathering - Sysctl Variables",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           600,
	}
	CategoryDNS = CategorySpec{
		ID:              "information.dns",
		Title:           "Information Gathering - DNS-Based Service Discovery",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           700,
	}
	CategoryK8sAPIServer = CategorySpec{
		ID:              "discovery.k8s_api",
		Title:           "Discovery - K8s API Server",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           800,
	}
	CategoryK8sServiceAccount = CategorySpec{
		ID:              "discovery.k8s_sa",
		Title:           "Discovery - K8s Service Account",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           900,
	}
	CategoryCloudMetadata = CategorySpec{
		ID:              "discovery.cloud_metadata",
		Title:           "Discovery - Cloud Provider Metadata API",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           1000,
	}
	CategoryKernel = CategorySpec{
		ID:              "exploit.kernel",
		Title:           "Exploit Pre - Kernel Exploits",
		DefaultProfiles: []string{ProfileBasic, ProfileExtended},
		Order:           1100,
	}
	CategorySensitiveFiles = CategorySpec{
		ID:              "information.sensitive_files",
		Title:           "Information Gathering - Sensitive Files",
		DefaultProfiles: []string{ProfileExtended, ProfileAdditional},
		Order:           1200,
	}
	CategoryASLR = CategorySpec{
		ID:              "information.aslr",
		Title:           "Information Gathering - ASLR",
		DefaultProfiles: []string{ProfileExtended, ProfileAdditional},
		Order:           1300,
	}
	CategoryCgroups = CategorySpec{
		ID:              "information.cgroups",
		Title:           "Information Gathering - Cgroups",
		DefaultProfiles: []string{ProfileExtended, ProfileAdditional},
		Order:           1400,
	}
)
