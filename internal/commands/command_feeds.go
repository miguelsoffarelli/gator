package commands

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/miguelsoffarelli/gator/internal/database"
	"github.com/miguelsoffarelli/gator/internal/rss"
)

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		fmt.Println("Parameters ignored: command feeds takes no parameters")
	}

	ctx := context.Background()
	feeds, err := s.Db.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found in database. Use command addfeed <feed name> <url> to add a new feed.")
	}

	for _, feed := range feeds {
		name, URL := feed.Name, feed.Url
		user, err := s.Db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting name of the user that created the feed: %v", err)
		}
		userName := user.Name

		fmt.Printf("Name: %v\n", name)
		fmt.Printf("URL: %v\n", URL)
		fmt.Printf("Created By: %v\n", userName)
		fmt.Println("==============================================")
	}

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("must provide feed name and url")
	}

	feedName, feedURL := cmd.Args[0], cmd.Args[1]

	ctx := context.Background()

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, createFeedParams)
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	if _, err = s.Db.CreateFeedFollow(ctx, createFeedFollowParams); err != nil {
		return fmt.Errorf("error following created feed: %v", err)
	}

	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Updated At: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("User ID: %v\n", feed.UserID)

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("must provide time between requests argument")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time between requests: %v", err)
	}

	fmt.Printf("Collecting feeds every %v", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("error scraping feed: %v", err)
		}
	}
}

func scrapeFeeds(s *State) error {
	ctx := context.Background()
	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %v", err)
	}

	params := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	if err := s.Db.MarkFeedFetched(ctx, params); err != nil {
		return fmt.Errorf("error marking feed as fetched: %v", err)
	}

	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	for _, item := range feed.Channel.Item {
		publishedAt, err := parseTime(item.PubDate)
		if err != nil {
			publishedAt = nil
		}

		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  *publishedAt,
				Valid: true,
			},
			FeedID: nextFeed.ID,
		}

		if _, err := s.Db.CreatePost(ctx, params); err != nil && !isUniqueConstraintError(err) {
			fmt.Printf("error adding post %s to database: %v", item.Title, err)
		}
	}

	return nil
}

func parseTime(timeStr string) (*time.Time, error) {
	// To avoid trying unnecessary formats that are less common on RSS's,
	// first iterate over the most common formats.
	commonFormats := []string{
		time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
	}

	for _, format := range commonFormats {
		// if error is nil it means the parsing was successful
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t, nil
		}
	}

	// ONLY if that fails, try every format provided in go's time package documentation
	allFormats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, format := range allFormats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("error: couldn't parse time of publication")
}

// Check for SQL State 23505 for duplicate url's
func isUniqueConstraintError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := int32(2) // default value
	if len(cmd.Args) >= 1 {
		if intLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = int32(intLimit)
		}
	}

	ctx := context.Background()
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.Db.GetPostsForUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error getting posts for current user: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %v\n", post.Title.String)
		fmt.Printf("URL: %v\n", post.Url)
		fmt.Printf("Published At: %v\n", post.PublishedAt.Time.Format(time.DateTime))
		fmt.Printf("Description: %v\n", post.Description.String)
		fmt.Println("================================================================")
		fmt.Println("")
	}

	return nil
}
