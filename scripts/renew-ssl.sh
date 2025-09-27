#!/bin/bash

# SSL Certificate Renewal Script for QRify
# This script renews SSL certificates and reloads nginx

# Change to project directory
cd /var/lib/jenkins/QRify

# Renew certificates using webroot method
sudo certbot renew --webroot -w /var/www/html --quiet

# Log the result
echo "$(date): SSL certificate renewal completed" >> /var/log/ssl-renewal.log
