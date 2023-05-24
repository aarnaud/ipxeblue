# ipxeblue : iPXE management

> Manage network boot over HTTPS with iPXE with admin WebUI and API

[Go to screenshots](#screenshots)

## features

- auto create Computer object on first network boot
- manage bootentry
  - ipxe script field
  - upload files for boot 
- manage iPXE account for basic auth
- https support
- iPXE menu that list all bootentries

## config

Supported environment variables

- `PORT` used to bind http server
  - default: `8080`
- `BASE_URL` used in template to defined the public URL
  - default: `http://127.0.0.1:8080`
- `ENABLE_API_AUTH` can be switch to `false` to used SSO proxy in front of API
  - default: `true`
- `DATABASE_URL` postgres URL, example `postgres://user:passworkd@localhost:5432/ipxeblue?sslmode=disable`
- `MINIO_ENDPOINT` S3 compatible endpoint
  - default: `127.0.0.1:9000`
- `MINIO_ACCESS_KEY`
- `MINIO_SECRET_KEY`
- `MINIO_SECURE`
- `MINIO_BUCKETNAME`
  - default: `ipxeblue`
- `GRUB_SUPPORT_ENABLED`
  - default: `False`
- `TFTP_ENABLED`
  - default: `False`
- `DEFAULT_BOOTENTRY_NAME`
  - default: ``
  
## DHCP or ipxe config for connection to ipxeblue

For embed ipxe script or chainload
```shell
ifstat ||
dhcp ||
route ||
set crosscert http://ca.ipxe.org/auto
chain https://USERNAME:PASSWORD@FQDN/?asset=${asset}&buildarch=${buildarch}&hostname=${hostname}&mac=${mac:hexhyp}&ip=${ip}&manufacturer=${manufacturer}&platform=${platform}&product=${product}&serial=${serial}&uuid=${uuid}&version=${version}
```

### For isc-dhcp-server

you need to set `iPXE-specific options` see https://ipxe.org/howto/dhcpd

```text
  if option arch = 00:07 {
     filename "snponly.efi";
  } else {
     filename "undionly.kpxe";
  }
  if exists user-class and option user-class = "iPXE" and exists ipxe.https {
      option ipxe.crosscert "http://ca.ipxe.org/auto";
      option ipxe.username "demo";
      option ipxe.password "demo";
      filename "https://ipxeblue.yourdomain/";
  }
  # a TFTP server to load iPXE if not already load by default
  next-server 10.123.123.123;
```

### For kea-dhcp-server
```text
"Dhcp4": {
    ... 
    "option-def": [
                { "space": "dhcp4", "name": "ipxe-encap-opts", "code": 175, "type": "empty", "array": false, "record-types": "", "encapsulate": "ipxe" },
                { "space": "ipxe", "name": "crosscert", "code": 93, "type": "string" },
                { "space": "ipxe", "name": "username", "code": 190, "type": "string" },
                { "space": "ipxe", "name": "password", "code": 191, "type": "string" }
    ],
    "client-classes": [
        {
            "name": "XClient_iPXE",
            "test": "substring(option[77].hex,0,4) == 'iPXE'",
            "boot-file-name": "ipxeblue.ipxe",
            "option-data": [
                { "space": "dhcp4", "name": "ipxe-encap-opts", "code": 175 },
                { "space": "ipxe", "name": "crosscert", "data": "http://ca.ipxe.org/auto" },
                { "space": "ipxe", "name": "username", "data": "demo" },
                { "space": "ipxe", "name": "password", "data": "demo" }
            ]
        },
        {
            "name": "UEFI-64",
            "test": "substring(option[60].hex,0,20) == 'PXEClient:Arch:00007'",
             "boot-file-name": "snponly.efi"
        },
        {
            "name": "Legacy",
            "test": "substring(option[60].hex,0,20) == 'PXEClient:Arch:00000'",
            "boot-file-name": "undionly.kpxe"
        }
    ],
    "subnet4": [
        {
        ...
        "next-server": "10.123.123.123",
        ...
        }
    ]
    ...
}
```

### Grub over PXE:
> Secure Boot supported with signed binaries

/srv/tftp/bootx64.efi  (sha256sum 8c885fa9886ab668da267142c7226b8ce475e682b99e4f4afc1093c5f77ce275)
/srv/tftp/grubx64.efi  (sha256sum d0d6d85f44a0ffe07d6a856ad5a1871850c31af17b7779086b0b9384785d5449)
/srv/tftp/grub/grub.cfg  
```text
insmod http
source (http,192.168.32.7)/grub/
```

## screenshots 

![Computer List](docs/images/computer-list.png?raw=true "Computer List")
![Computer Edit](docs/images/computer-edit.png?raw=true "Computer Edit")
![Account List](docs/images/account-list.png?raw=true "Account List")
![Bootentry List](docs/images/bootentry-list.png?raw=true "Bootentry List")
![Bootentry Edit](docs/images/bootentry-edit.png?raw=true "Bootentry List")