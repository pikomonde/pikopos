sudo apt-get update
sudo apt-get install mysql-server

## Run Project Locally

You might want to run the service in your local machine. This backend service runs on port `:1235`, meanhile another frontend service (pikopos-frontend) runs on port `:8080`.

1. Init MySQL

    `mysql -h localhost -u root -p < ./setup/deploy_00.00.001_init_schema.sql`

2. Setting Up NGINX

    Setting up NGINX so opening `localhost:1111` will redirect to `localhost:1235` for backend endpoints and `localhost:8080` for frontend endpoints

    1. Create a new file called pikopos in nginx's sites-available:

        `sudo touch /etc/nginx/sites-available/pikopos`

        `sudo nano /etc/nginx/sites-available/pikopos`

        ```
        upstream pikopos {
            server localhost:57672;
        }

        upstream frontend {
            server localhost:8080;
        }

        server {
            listen 1111;

            location / {
            proxy_pass http://frontend;
            }

            location ~ ^/(ping|auth|employee) {
                proxy_pass http://pikopos;
            }
        }

        ```
    2. Create a symbolic link in sites-enabled to `pikopos` file in sites-available:

        `sudo ln -s /etc/nginx/sites-available/pikopos /etc/nginx/sites-enabled/pikopos`

    3. Test and Restart nginx service

        `sudo nginx -t && sudo service nginx restart`

        The first command `sudo nginx -t` is used to test whether we have an error in the config or not. If there is no error, then we restart the nginx service by using `sudo service nginx restart`

        We might want check whether the nginx service successfully running or not by checking it's status

        `sudo service nginx status`

        pres `q` to quit

3. TODO: create tunneling for login / or just inject cookie?

4. Run Backend and Frontend services

    1. Run Backend Service

        `go run app.go`

    2. Run Frontend Service

        See [PikoPOS Frontend Repository](https://github.com/pikomonde/pikopos-frontend) for running the frontend service

5. Test whether backend is running

    `curl -G localhost:1111/ping`




