#!/bin/bash
sudo systemctl restart auth.service
sudo systemctl restart user.service
sudo systemctl restart news.service
sudo systemctl restart gitserver.service
sudo systemctl restart server.service