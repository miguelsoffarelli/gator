# Gator - Blog Aggregator

CLI tool that allows users to:

* Add RSS feeds from across the internet to be collected
* Store the collected posts in a PostgreSQL database
* Follow and unfollow RSS feeds that other users have added
* View summaries of the aggregated posts in the terminal, with a link to the full post
Built with educational purposes to practice client-side HTTP requests in GO and PostgreSQL database implementation.

## Dependencies

* Go
* PostgreSQL (running locally on port 5432)
    - Create database with name "gator"
    - Default credentials should be `postgres:postgres`

## Installation

Run
```
go install github.com/miguelsoffarelli/gator@latest
```

## Usage

### Setup
1. Create a file with name .gatorconfig.json at home directory.

2. Open created file .gatorconfig.json and insert the following content:
    ```
    {
        "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
        "current_user_name": ""
    }
    ```

3. Use "register" command (see commands below) to sign up and login.

4. Use "addfeed" command to add feeds.

5. Use "follow" command to follow desired feeds.

6. Run
    ```
    gator agg <time between requests>
    ```
    in a separate terminal. Agg is the background process that fetches the posts from the followed feeds with a given <time between requests> interval.

7. Use "browse" command to display the posts fetched in the previous step.
### Commands
* #### Register:

Register and log in with the given username.
```
gator register <username>
```

* #### Login:

Log in with the given username. User must be previously registered.
```
gator login <username>
```

* #### Add Feed:

Save the given feed in database.
```
gator addfeed <feed name> <url>
```

* #### Feeds:

Show all feeds currently in database.
```
gator feeds
```

* #### Follow:

Follow feed with the given URL.
```
gator follow <url>
```

* #### Following:

Display followed feeds for current user.
```
gator following
```

* #### Unfollow:

Unfollow feed with the given URL.
```
gator unfollow <url>
```

* #### Browse:

Display all the fetched posts for the followed feeds. 
Takes an optional <limit> argument that determines how many posts will be displayed. Defaults to 2 if not provided.
```
gator browse <limit>
```