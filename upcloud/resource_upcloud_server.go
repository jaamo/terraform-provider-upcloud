package upcloud

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUpCloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceUpCloudServerCreate,
		Read:   resourceUpCloudServerRead,
		Update: resourceUpCloudServerUpdate,
		Delete: resourceUpCloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"mem": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"template": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"private_networking": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ipv4": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_address_private": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_devices": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"tier": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"login": {
				Type:     schema.TypeSet,
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"keys": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"create_password": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"password_delivery": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "none",
						},
					},
				},
			},
		},
	}
}

func resourceUpCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	r, err := buildServerOpts(d, meta)
	if err != nil {
		return err
	}
	server, err := client.CreateServer(r)
	if err != nil {
		return err
	}
	d.SetId(server.UUID)
	log.Printf("[INFO] Server %s with UUID %s created", server.Title, server.UUID)

	server, err = client.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         server.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      time.Minute * 5,
	})
	if err != nil {
		return err
	}
	return resourceUpCloudServerRead(d, meta)
}

func resourceUpCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	r := &request.GetServerDetailsRequest{
		UUID: d.Id(),
	}
	server, err := client.GetServerDetails(r)
	if err != nil {
		return err
	}
	d.Set("hostname", server.Hostname)
	d.Set("title", server.Title)
	d.Set("zone", server.Zone)
	d.Set("cpu", server.CoreNumber)
	storageDevices := d.Get("storage_devices").([]interface{})
	log.Print(storageDevices)
	log.Print(server.StorageDevices)
	for i, storageDevice := range storageDevices {
		storageDevice := storageDevice.(map[string]interface{})
		storageDevice["id"] = server.StorageDevices[i].UUID
		storageDevice["address"] = server.StorageDevices[i].Address
	}
	d.Set("storage_devices", storageDevices)

	return nil
}

func resourceUpCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	if d.HasChange("storage_devices") {
		storageDevicesI, oldStorageDevicesI := d.GetChange("storage_devices")
		d.Set("storage_devices", storageDevicesI)
		storageDevices := storageDevicesI.([]interface{})
		oldStorageDevices := oldStorageDevicesI.([]interface{})
		log.Printf("New devices: %v\nOld devices: %v\n", storageDevices, oldStorageDevices)
		for i, storageDevice := range storageDevices {
			storageDevice := storageDevice.(map[string]interface{})
			oldStorageDeviceN := sort.Search(len(oldStorageDevices), func(i int) bool {
				return oldStorageDevices[i].(map[string]interface{})["id"] == storageDevice["id"]
			})
			var oldStorageDevice map[string]interface{}
			if oldStorageDeviceN < len(oldStorageDevices) {
				oldStorageDevice = oldStorageDevices[oldStorageDeviceN].(map[string]interface{})
			}
			log.Printf("New device: %v\n", storageDevice)
			log.Printf("Old device: %v\n", oldStorageDevice)
			if oldStorageDevice == nil {
				var newStorageDeviceID string
				switch storageDevice["action"] {
				case upcloud.CreateServerStorageDeviceActionCreate:
					storage, err := buildStorage(storageDevice, i, meta)
					if err != nil {
						return err
					}
					newStorage, err := client.CreateStorage(&request.CreateStorageRequest{
						Size:  storage.Size,
						Tier:  storage.Tier,
						Title: storage.Title,
						Zone:  d.Get("zone").(string),
					})
					if err != nil {
						return err
					}
					newStorageDeviceID = newStorage.UUID
					break
				case upcloud.CreateServerStorageDeviceActionClone:
					// storage, err := buildStorage(storageDevice, i, meta)
					// if err != nil {
					// 	return err
					// }
					// newStorage, err := client.CloneStorage(&request.CloneStorageRequest{
					// 	UUID:  storageDevice["storage"].(string),
					// 	Tier:  storage.Tier,
					// 	Title: storage.Title,
					// 	Zone:  d.Get("zone").(string),
					// })
					// if err != nil {
					// 	return err
					// }
					newStorageDeviceID = storageDevice["storage"].(string)
					break
				case upcloud.CreateServerStorageDeviceActionAttach:
					newStorageDeviceID = storageDevice["storage"].(string)
					break
				}

				attachStorageRequest := request.AttachStorageRequest{
					ServerUUID:  d.Id(),
					StorageUUID: newStorageDeviceID,
				}

				if storageType := storageDevice["type"].(string); storageType != "" {
					attachStorageRequest.Type = storageType
				}

				client.AttachStorage(&attachStorageRequest)
			} else {
				if !reflect.DeepEqual(oldStorageDevice, storageDevice) {
					client.ModifyStorage(&request.ModifyStorageRequest{
						Size:  storageDevice["size"].(int),
						Title: storageDevice["title"].(string),
					})
				}

				oldStorageDevices = append(oldStorageDevices[:oldStorageDeviceN], oldStorageDevices[oldStorageDeviceN+1:]...)
			}
		}
		log.Printf("Old devices: %v\n", oldStorageDevices)
		for _, oldStorageDevice := range oldStorageDevices {
			oldStorageDevice := oldStorageDevice.(map[string]interface{})
			client.DetachStorage(&request.DetachStorageRequest{
				ServerUUID: d.Id(),
				Address:    oldStorageDevice["address"].(string),
			})
			if oldStorageDevice["action"] != upcloud.CreateServerStorageDeviceActionAttach {
				client.DeleteStorage(&request.DeleteStorageRequest{
					UUID: oldStorageDevice["id"].(string),
				})
			}
		}
	}
	if d.HasChange("mem") || d.HasChange("cpu") {
		_, newCPU := d.GetChange("cpu")
		_, newMem := d.GetChange("mem")
		if err := verifyServerStopped(d, meta); err != nil {
			return err
		}
		r := &request.ModifyServerRequest{
			UUID:         d.Id(),
			CoreNumber:   strconv.Itoa(newCPU.(int)),
			MemoryAmount: strconv.Itoa(newMem.(int)),
		}
		_, err := client.ModifyServer(r)
		if err != nil {
			return err
		}
		if err := verifyServerStarted(d, meta); err != nil {
			return err
		}

	}
	return resourceUpCloudServerRead(d, meta)
}

func resourceUpCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	// Verify server is stopped before deletion
	if err := verifyServerStopped(d, meta); err != nil {
		return err
	}
	// Delete server
	deleteServerRequest := &request.DeleteServerRequest{
		UUID: d.Id(),
	}
	log.Printf("[INFO] Deleting server (server UUID: %s)", d.Id())
	err := client.DeleteServer(deleteServerRequest)
	if err != nil {
		return err
	}

	storageDevices := d.Get("storage_devices").([]interface{})
	for _, storageDevice := range storageDevices {
		// Delete server root disk
		storageDevice := storageDevice.(map[string]interface{})
		id := storageDevice["id"].(string)
		action := storageDevice["action"].(string)
		if action != upcloud.CreateServerStorageDeviceActionAttach {
			deleteStorageRequest := &request.DeleteStorageRequest{
				UUID: id,
			}
			log.Printf("[INFO] Deleting server storage (storage UUID: %s)", id)
			err = client.DeleteStorage(deleteStorageRequest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func buildServerOpts(d *schema.ResourceData, meta interface{}) (*request.CreateServerRequest, error) {
	r := &request.CreateServerRequest{
		Zone:     d.Get("zone").(string),
		Hostname: d.Get("hostname").(string),
		Title:    fmt.Sprintf("%s (managed by terraform)", d.Get("hostname").(string)),
	}

	if attr, ok := d.GetOk("cpu"); ok {
		r.CoreNumber = attr.(int)
	}
	if attr, ok := d.GetOk("mem"); ok {
		r.MemoryAmount = attr.(int)
	}
	if attr, ok := d.GetOk("user_data"); ok {
		r.UserData = attr.(string)
	}
	if login, ok := d.GetOk("login"); ok {
		loginOpts, deliveryMethod, err := buildLoginOpts(login, meta)
		if err != nil {
			return nil, err
		}
		r.LoginUser = loginOpts
		r.PasswordDelivery = deliveryMethod
	}

	storageDevices := d.Get("storage_devices").([]interface{})
	storageOpts, err := buildStorageOpts(storageDevices, meta)
	if err != nil {
		return nil, err
	}
	r.StorageDevices = storageOpts

	networkOpts, err := buildNetworkOpts(d, meta)
	if err != nil {
		return nil, err
	}
	r.IPAddresses = networkOpts

	return r, nil
}

func buildStorage(storageDevice map[string]interface{}, i int, meta interface{}) (*upcloud.CreateServerStorageDevice, error) {
	osDisk := upcloud.CreateServerStorageDevice{
		Action: storageDevice["action"].(string),
	}

	if source := storageDevice["storage"].(string); source != "" {
		_, err := uuid.ParseUUID(source)
		// Assume template name is given and map name to UUID
		if err != nil {
			client := meta.(*service.Service)
			r := &request.GetStoragesRequest{
				Type: upcloud.StorageTypeTemplate,
			}
			l, err := client.GetStorages(r)
			if err != nil {
				return nil, err
			}
			for _, s := range l.Storages {
				if s.Title == source {
					source = s.UUID
					break
				}
			}
		}

		osDisk.Storage = source
	}

	// Set size or use the one defined by target template
	if size := storageDevice["size"]; size != -1 {
		osDisk.Size = size.(int)
	}

	// Autogenerate disk title
	osDisk.Title = fmt.Sprintf("terraform-os-disk-%d", i)

	// Set disk tier or use the one defined by target template
	if tier := storageDevice["tier"]; tier != "" {
		osDisk.Tier = tier.(string)
	}

	if storageType := storageDevice["type"].(string); storageType != "" {
		osDisk.Type = storageType
	}

	return &osDisk, nil
}

func buildStorageOpts(storageDevices []interface{}, meta interface{}) ([]upcloud.CreateServerStorageDevice, error) {
	storageCfg := make([]upcloud.CreateServerStorageDevice, 0)
	for i, storageDevice := range storageDevices {
		storageDevice, err := buildStorage(storageDevice.(map[string]interface{}), i, meta)

		if err != nil {
			return nil, err
		}

		storageCfg = append(storageCfg, *storageDevice)
	}

	return storageCfg, nil
}

func buildNetworkOpts(d *schema.ResourceData, meta interface{}) ([]request.CreateServerIPAddress, error) {
	ifaceCfg := make([]request.CreateServerIPAddress, 0)
	if attr, ok := d.GetOk("ipv4"); ok {
		publicIPv4 := attr.(bool)
		if publicIPv4 {
			publicIPv4 := request.CreateServerIPAddress{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv4,
			}
			ifaceCfg = append(ifaceCfg, publicIPv4)
		}
	}
	if attr, ok := d.GetOk("private_networking"); ok {
		setPrivateIP := attr.(bool)
		if setPrivateIP {
			privateIPv4 := request.CreateServerIPAddress{
				Access: upcloud.IPAddressAccessPrivate,
				Family: upcloud.IPAddressFamilyIPv4,
			}
			ifaceCfg = append(ifaceCfg, privateIPv4)
		}
	}
	if attr, ok := d.GetOk("ipv6"); ok {
		publicIPv6 := attr.(bool)
		if publicIPv6 {
			publicIPv6 := request.CreateServerIPAddress{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv6,
			}
			ifaceCfg = append(ifaceCfg, publicIPv6)
		}
	}
	return ifaceCfg, nil
}

func buildLoginOpts(v interface{}, meta interface{}) (*request.LoginUser, string, error) {
	// Construct LoginUser struct from the schema
	r := &request.LoginUser{}
	e := v.(*schema.Set).List()[0]
	m := e.(map[string]interface{})

	// Set username as is
	r.Username = m["user"].(string)

	// Set 'create_password' to "yes" or "no" depending on the bool value.
	// Would be nice if the API would just get a standard bool str.
	createPassword := "no"
	b := m["create_password"].(bool)
	if b {
		createPassword = "yes"
	}
	r.CreatePassword = createPassword

	// Handle SSH keys one by one
	keys := make([]string, 0)
	for _, k := range m["keys"].([]interface{}) {
		key := k.(string)
		keys = append(keys, key)
	}
	r.SSHKeys = keys

	// Define password delivery method none/email/sms
	deliveryMethod := m["password_delivery"].(string)

	return r, deliveryMethod, nil
}

func verifyServerStopped(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	// Get current server state
	r := &request.GetServerDetailsRequest{
		UUID: d.Id(),
	}
	server, err := client.GetServerDetails(r)
	if err != nil {
		return err
	}
	if server.State != upcloud.ServerStateStopped {
		// Soft stop with 2 minute timeout, after which hard stop occurs
		stopRequest := &request.StopServerRequest{
			UUID:     d.Id(),
			StopType: "soft",
			Timeout:  time.Minute * 2,
		}
		log.Printf("[INFO] Stopping server (server UUID: %s)", d.Id())
		_, err := client.StopServer(stopRequest)
		if err != nil {
			return err
		}
		_, err = client.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         d.Id(),
			DesiredState: upcloud.ServerStateStopped,
			Timeout:      time.Minute * 5,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyServerStarted(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*service.Service)
	// Get current server state
	r := &request.GetServerDetailsRequest{
		UUID: d.Id(),
	}
	server, err := client.GetServerDetails(r)
	if err != nil {
		return err
	}
	if server.State != upcloud.ServerStateStarted {
		// Soft stop with 2 minute timeout, after which hard stop occurs
		startRequest := &request.StartServerRequest{
			UUID:    d.Id(),
			Timeout: time.Minute * 2,
		}
		log.Printf("[INFO] Stopping server (server UUID: %s)", d.Id())
		_, err := client.StartServer(startRequest)
		if err != nil {
			return err
		}
		_, err = client.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         d.Id(),
			DesiredState: upcloud.ServerStateStarted,
			Timeout:      time.Minute * 5,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
