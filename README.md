# ipxeblue : iPXE management

> Manage network boot over HTTPS with iPXE with admin WebUI and API

## features

- auto create Computer object on boot
- manage bootentry
  - ipxe script field
  - upload files for boot 
- manage iPXE account for basic auth
- https support
- iPXE menu that list bootentries

## config

Supported environment variables

- `PORT` used to bind http server
  - default: `8080`
- `BASE_URL` used in template to defined the public URL
  - default: `http://127.0.0.1:8080`
- `ENABLE_API_AUTH` can be switch to `false` to used SSO proxy in front of API
  - default: `true`
- `MINIO_ENDPOINT` S3 compatible endpoint
  - default: `127.0.0.1:9000`
- `MINIO_ACCESS_KEY`
- `MINIO_SECRET_KEY`
- `MINIO_SECURE`
- `MINIO_BUCKETNAME`
  - default: `ipxeblue`
  
## Exemple of ipxe config

```shell
ifstat ||
dhcp ||
route ||
set crosscert http://ca.ipxe.org/auto
chain https://USERNAME:PASSWORD@FQDN/?asset=${asset}&buildarch=${buildarch}&hostname=${hostname}&mac=${mac:hexhyp}&ip=${ip}&manufacturer=${manufacturer}&platform=${platform}&product=${product}&serial=${serial}&uuid=${uuid}&version=${version}
```