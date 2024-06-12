# Gitlab menu extender demo

This middleware project helps you to extend Gitlab information and automation pages.

## Before using and building

Please make a copy of repo _https://github.com/aborche/go-gitlab_ to directory _/projects/GolandProjects/go-gitlab_

Some pull requests for original _https://github.com/xanzy/go-gitlab_ repo is not approved. Until approve, please use my local repo.

go.mod file contains sublink to local repository.

## Building

Just install go 1.21 and run 
```shell
go mod tidy
go build .
```

## Environment
```bash
* HELPER_GITLABURLFROMHOST - Use Gitlab URL from Host headers. Values: 0/1. Default: 1 
* HELPER_URLPREFIX - Gitlab helper prefix. Ex: "/-/mainmenu". Default: "/-/helper"
* HELPER_MENUFILE - Custom menu file. Ex: "templates/sidebarmenu.json". Default: "templates/sidebarmenushort.json"
* HELPER_DEFAULTSECTION - Default gitlab start page for source template. Ex: "/dashboard/groups". Default: "/dashboard/projects"
* HELPER_CUSTOMGITLABURL - Custom gitlab url, if you cannot forward a Host header to service. Default: undefined
* PORT - Custom service port for publishing. Default: 8085
```

## Changing gitlab config

Create a file _helper.conf_ in directory _/etc/nginx/gitlab/_, contains
```
location /-/helper { # This line must be the same HELPER_URLPREFIX environment variable
	proxy_set_header Host $host;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header X-Forwarded-Proto $scheme;
	proxy_set_header Content-Length "";
	proxy_set_header X-Original-URI $request_uri;
	proxy_set_header X-Original-ARGS $args;
	proxy_set_header X-Remote-Addr $remote_addr;
	proxy_set_header X-Original-Host $host;
	proxy_pass http://localhost:8085; # Set upstream value to your published service
}
location = /api/v4/graphql {
    rewrite /api/v4/graphql /api/graphql last;

    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto https;
    proxy_set_header X-Forwarded-Ssl on;

    proxy_read_timeout                  900;
    proxy_cache off;
    proxy_buffering off;
    proxy_request_buffering off;
    proxy_http_version 1.1;

    proxy_pass         http://gitlab-workhorse;
}
```
Change /etc/gitlab/gitlab.rb file. Search line with **nginx['custom_gitlab_server_config']**, uncomment and replace it to:
```
nginx['custom_gitlab_server_config'] = "include /etc/nginx/gitlab/helper.conf;"
```
Make sure Gitlab has access to the file

Run command
```bash
sudo gitlab-ctl reconfigure
sudo gitlab-ctl restart nginx
```

Run your service with your custom environment and check out new gitlab page like
```
https://your.gitlab.server/-/helper/
```
where **'/-/helper'** is your prefix from nginx config and HELPER_URLPREFIX environment variable.

## Customize original gitlab menu

Clone repository https://github.com/aborche/gitlab-menu-changer and change your menu parameters in _postinst_ file.

Build own deb package with command
```
make build
```
Install it
```
sudo dpkg -i gitlab-menu-changer*.deb
```
If your gitlab instance already installed, run
```shell
$ sudo sh /var/lib/dpkg/triggers/gitlab-menu-changer.postinst triggered /opt/gitlab/embedded/service/gitlab-rails/lib/sidebars/your_work/panel.rb
$ sudo gitlab-ctl restart puma
```
wait some time and open your gitlab instance without page prefix.

If all is ok, you will see a new menu item targeted to your service.

Also read https://habr.com/ru/articles/814623/
