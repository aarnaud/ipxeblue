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
  
## DHCP or ipxe config for connection to ipxeblue

For embed ipxe script or chainload
```shell
ifstat ||
dhcp ||
route ||
set crosscert http://ca.ipxe.org/auto
chain https://USERNAME:PASSWORD@FQDN/?asset=${asset}&buildarch=${buildarch}&hostname=${hostname}&mac=${mac:hexhyp}&ip=${ip}&manufacturer=${manufacturer}&platform=${platform}&product=${product}&serial=${serial}&uuid=${uuid}&version=${version}
```

For isc-dhcp-server

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

## screenshots 

![Computer List](docs/images/computer-list.png?raw=true "Computer List")
![Computer Edit](docs/images/computer-edit.png?raw=true "Computer Edit")
![Account List](docs/images/account-list.png?raw=true "Account List")
![Bootentry List](docs/images/bootentry-list.png?raw=true "Bootentry List")
![Bootentry Edit](docs/images/bootentry-edit.png?raw=true "Bootentry List")