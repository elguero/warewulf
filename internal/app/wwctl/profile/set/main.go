package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/warewulfd"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	var err error

	nodeDB, err := node.New()
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not open node configuration: %s\n", err)
		os.Exit(1)
	}

	profiles, err := nodeDB.FindAllProfiles()
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "%s\n", err)
		os.Exit(1)
	}

	if SetAll {
		fmt.Printf("\n*** WARNING: This command will modify all profiles! ***\n\n")
	} else if len(args) > 0 {
		profiles = node.FilterByName(profiles, args)
	} else {
		wwlog.Printf(wwlog.INFO, "No profile specified, selecting the 'default' profile\n")
		profiles = node.FilterByName(profiles, []string{"default"})
	}

	if len(profiles) == 0 {
		fmt.Printf("No profiles found\n")
		os.Exit(1)
	}

	for _, p := range profiles {
		wwlog.Printf(wwlog.VERBOSE, "Modifying profile: %s\n", p.Id.Get())

		if SetComment != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting comment to: %s\n", p.Id.Get(), SetComment)
			p.Comment.Set(SetComment)
		}

		if SetClusterName != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting cluster name to: %s\n", p.Id.Get(), SetClusterName)
			p.ClusterName.Set(SetClusterName)
		}

		if SetContainer != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting Container name to: %s\n", p.Id.Get(), SetContainer)
			p.ContainerName.Set(SetContainer)
		}

		if SetInit != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting init command to: %s\n", p.Id.Get(), SetInit)
			p.Init.Set(SetInit)
		}

		if SetRoot != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting root to: %s\n", p.Id.Get(), SetRoot)
			p.Root.Set(SetRoot)
		}

		if SetAssetKey != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting asset key to: %s\n", p.Id.Get(), SetAssetKey)
			p.AssetKey.Set(SetAssetKey)
		}

		if SetKernelOverride != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting Kernel override version to: %s\n", p.Id.Get(), SetKernelOverride)
			p.Kernel.Override.Set(SetKernelOverride)
		}

		if SetKernelArgs != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting Kernel args to: %s\n", p.Id.Get(), SetKernelArgs)
			p.Kernel.Args.Set(SetKernelArgs)
		}

		if SetIpxe != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting iPXE template to: %s\n", p.Id.Get(), SetIpxe)
			p.Ipxe.Set(SetIpxe)
		}

		if len(SetRuntimeOverlay) != 0 {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting runtime overlay to: %s\n", p.Id.Get(), SetRuntimeOverlay)
			p.RuntimeOverlay.SetSlice(SetRuntimeOverlay)
		}

		if len(SetSystemOverlay) != 0 {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting system overlay to: %s\n", p.Id.Get(), SetSystemOverlay)
			p.SystemOverlay.SetSlice(SetSystemOverlay)
		}

		if SetIpmiNetmask != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting IPMI netmask to: %s\n", p.Id.Get(), SetIpmiNetmask)
			p.Ipmi.Netmask.Set(SetIpmiNetmask)
		}

		if SetIpmiPort != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting IPMI port to: %s\n", p.Id.Get(), SetIpmiPort)
			p.Ipmi.Port.Set(SetIpmiPort)
		}

		if SetIpmiGateway != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting IPMI gateway to: %s\n", p.Id.Get(), SetIpmiGateway)
			p.Ipmi.Gateway.Set(SetIpmiGateway)
		}

		if SetIpmiUsername != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting IPMI username to: %s\n", p.Id.Get(), SetIpmiUsername)
			p.Ipmi.UserName.Set(SetIpmiUsername)
		}

		if SetIpmiPassword != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting IPMI password to: %s\n", p.Id.Get(), SetIpmiPassword)
			p.Ipmi.Password.Set(SetIpmiPassword)
		}

		if SetIpmiInterface != "" {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting IPMI interface to: %s\n", p.Id.Get(), SetIpmiInterface)
			p.Ipmi.Interface.Set(SetIpmiInterface)
		}

		if SetIpmiWrite == "yes" || SetNetOnBoot == "y" || SetNetOnBoot == "1" || SetNetOnBoot == "true" {
			wwlog.Printf(wwlog.VERBOSE, "Node: %s, Setting Ipmiwrite to %s\n", p.Id.Get(), SetIpmiWrite)
			p.Ipmi.Write.SetB(true)
		} else {
			wwlog.Printf(wwlog.VERBOSE, "Node: %s, Setting Ipmiwrite to %s\n", p.Id.Get(), SetIpmiWrite)
			p.Ipmi.Write.SetB(false)
		}

		if SetDiscoverable {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting all nodes to discoverable\n", p.Id.Get())
			p.Discoverable.SetB(true)
		}

		if SetUndiscoverable {
			wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting all nodes to undiscoverable\n", p.Id.Get())
			p.Discoverable.SetB(false)
		}

		if SetNetName != "" {
			if _, ok := p.NetDevs[SetNetName]; !ok {
				var nd node.NetDevEntry
				nd.Tags = make(map[string]*node.Entry)
				p.NetDevs[SetNetName] = &nd
			}
		}

		if SetNetDev != "" {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			wwlog.Printf(wwlog.VERBOSE, "Node: %s:%s, Setting net Device to: %s\n", p.Id.Get(), SetNetName, SetNetDev)
			p.NetDevs[SetNetName].Device.Set(SetNetDev)
		}

		if SetNetmask != "" {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			wwlog.Printf(wwlog.VERBOSE, "Profile '%s': Setting netmask to: %s\n", p.Id.Get(), SetNetName)
			p.NetDevs[SetNetName].Netmask.Set(SetNetmask)
		}

		if SetGateway != "" {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			wwlog.Printf(wwlog.VERBOSE, "Profile '%s': Setting gateway to: %s\n", p.Id.Get(), SetNetName)
			p.NetDevs[SetNetName].Gateway.Set(SetGateway)
		}

		if SetType != "" {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			wwlog.Printf(wwlog.VERBOSE, "Profile '%s': Setting HW address to: %s:%s\n", p.Id.Get(), SetNetName, SetType)
			p.NetDevs[SetNetName].Type.Set(SetType)
		}

		if SetNetOnBoot != "" {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			if SetNetOnBoot == "yes" || SetNetOnBoot == "y" || SetNetOnBoot == "1" || SetNetOnBoot == "true" {
				wwlog.Printf(wwlog.VERBOSE, "Profile: %s:%s, Setting ONBOOT\n", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].OnBoot.SetB(true)
			} else {
				wwlog.Printf(wwlog.VERBOSE, "Profile: %s:%s, Unsetting ONBOOT\n", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].OnBoot.SetB(false)
			}
		}

		if SetNetPrimary != "" {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			if SetNetPrimary == "yes" || SetNetPrimary == "y" || SetNetPrimary == "1" || SetNetPrimary == "true" {

				// Set all other networks to non-default
				for _, n := range p.NetDevs {
					n.Primary.SetB(false)
				}

				wwlog.Printf(wwlog.VERBOSE, "Profile: %s:%s, Setting PRIMARY\n", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].Primary.SetB(true)
			} else {
				wwlog.Printf(wwlog.VERBOSE, "Profile: %s:%s, Unsetting PRIMARY\n", p.Id.Get(), SetNetName)
				p.NetDevs[SetNetName].Primary.SetB(false)
			}
		}

		if SetNetDevDel {
			if SetNetName == "" {
				wwlog.Printf(wwlog.ERROR, "You must include the '--netname' option\n")
				os.Exit(1)
			}

			if _, ok := p.NetDevs[SetNetName]; !ok {
				wwlog.Printf(wwlog.ERROR, "Profile '%s': network name doesn't exist: %s\n", p.Id.Get(), SetNetName)
				os.Exit(1)
			}

			wwlog.Printf(wwlog.VERBOSE, "Profile %s: Deleting network: %s\n", p.Id.Get(), SetNetName)
			delete(p.NetDevs, SetNetName)
		}

		if len(SetTags) > 0 {
			for _, t := range SetTags {
				keyval := strings.SplitN(t, "=", 2)
				key := keyval[0]
				val := keyval[1]

				if _, ok := p.Tags[key]; !ok {
					var nd node.Entry
					p.Tags[key] = &nd
				}

				wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Setting Tag '%s'='%s'\n", p.Id.Get(), key, val)
				p.Tags[key].Set(val)
			}
		}
		if len(SetDelTags) > 0 {
			for _, t := range SetDelTags {
				keyval := strings.SplitN(t, "=", 1)
				key := keyval[0]

				if _, ok := p.Tags[key]; !ok {
					wwlog.Printf(wwlog.WARN, "Key does not exist: %s\n", key)
					os.Exit(1)
				}

				wwlog.Printf(wwlog.VERBOSE, "Profile: %s, Deleting tag: %s\n", p.Id.Get(), key)
				delete(p.Tags, key)
			}
		}
		if len(SetNetTags) > 0 {
			for _, t := range SetNetTags {
				keyval := strings.SplitN(t, "=", 2)
				key := keyval[0]
				val := keyval[1]
				if _, ok := p.NetDevs[SetNetName].Tags[key]; !ok {
					var nd node.Entry
					p.NetDevs[SetNetName].Tags[key] = &nd
				}

				wwlog.Printf(wwlog.VERBOSE, "Profile: %s:%s, Setting NETTAG '%s'='%s'\n", p.Id.Get(), SetNetName, key, val)
				p.NetDevs[SetNetName].Tags[key].Set(val)
			}

		}
		if len(SetNetDelTags) > 0 {
			for _, t := range SetNetDelTags {
				keyval := strings.SplitN(t, "=", 1)
				key := keyval[0]
				if _, ok := p.NetDevs[SetNetName].Tags[key]; !ok {
					wwlog.Printf(wwlog.WARN, "Profile: %s,%s, Key %s does not exist\n", p.Id.Get(), SetNetName, key)
					os.Exit(1)
				}

				wwlog.Printf(wwlog.VERBOSE, "Profile: %s,%s Deleting NETTAG: %s\n", p.Id.Get(), SetNetName, key)
				delete(p.NetDevs[SetNetName].Tags, key)
			}
		}

		err := nodeDB.ProfileUpdate(p)
		if err != nil {
			wwlog.Printf(wwlog.ERROR, "%s\n", err)
			os.Exit(1)
		}
	}

	if len(profiles) > 0 {
		if SetYes {
			err := nodeDB.Persist()
			if err != nil {
				return errors.Wrap(err, "failed to persist nodedb")
			}

			err = warewulfd.DaemonReload()
			if err != nil {
				return errors.Wrap(err, "failed to reload warewulf daemon")
			}
		} else {
			q := fmt.Sprintf("Are you sure you want to modify %d profile(s)", len(profiles))

			prompt := promptui.Prompt{
				Label:     q,
				IsConfirm: true,
			}

			result, _ := prompt.Run()

			if result == "y" || result == "yes" {
				err := nodeDB.Persist()
				if err != nil {
					return errors.Wrap(err, "failed to persist nodedb")
				}

				err = warewulfd.DaemonReload()
				if err != nil {
					return errors.Wrap(err, "failed to reload daemon")
				}
			}
		}
	} else {
		fmt.Printf("No profiles found\n")
	}

	return nil
}
