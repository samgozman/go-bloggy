# go-bloggy

[![codecov](https://codecov.io/gh/samgozman/go-bloggy/graph/badge.svg?token=Gk2j7fAlI9)](https://codecov.io/gh/samgozman/go-bloggy)

Go Bloggy: A simple &amp; lightweight backend for developers' personal blogs.

Bloggy started as a small side-project for my own dev blog [gozman.space/](https://gozman.space/),
but you can use it too.
You can find API documentation in [openapi](api/openapi.yaml) config file.

All it does is serving posts from DB and providing a simple API to manage them;
subscribing readers to that blog and sending email notifications via Mailjet.
