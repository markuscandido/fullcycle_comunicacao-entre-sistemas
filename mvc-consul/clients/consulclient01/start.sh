#!/bin/sh

# Instalar Nginx e outros pacotes
apk -U add nginx bind-tools vim

# Iniciar o Nginx
nginx

# Rodar o Consul
consul agent -client 0.0.0.0 -config-dir=/etc/consul.d
