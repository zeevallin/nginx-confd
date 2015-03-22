# nginx-confd-proxy

## Extend with custom modules and static files
If the default modules does not meet your needs, it's very easy to to create your own.
Create a new repository and inherit from our image in your `Dockerfile`

```dockerfile
# Dockerfile
FROM chatspry/nginx-confd-proxy:latest
```

Say you want to create a maintenance page. Let's create an HTML file in the `static` directory.
All assets you add to your `static` folder will appear in `/etc/nginx/static/custom` when you build your docker image.

```html
<!-- static/maintenance.html -->
<!DOCTYPE html>
<head>
  <title>Under maintenance</title>
</head>
<body>
  Service is under maintenance
</body>
```

Now you can create a directory called `include` where you can place your custom NGINX config.
All modules you add to your `include` folder will appear in `/etc/nginx/include/custom` when you build your docker image.

```nginx
# include/maintenance.conf
location = / {
  root /etc/nginx/static/custom/maintenance.html;
}
```

Now when you're creating your site record in etcd, you can reference your maintenance module in the custom namespace.

```json
{
  "listen": "80",
  "server_name": "example.com",
  "includes": [
    "custom/maintenance"
  ]
}
```