package services

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendarUrl(url string) (string, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return "", nil
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return "", nil
	}

	config.RedirectURL = url

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	return authURL, nil
}

func GetCalendarToken(code string, url string) (*oauth2.Token, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	config.RedirectURL = url

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetCalendarEvents(token *oauth2.Token) ([]*calendar.Event, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := config.Client(ctx, token)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	timeMin := now.Format(time.RFC3339)
	timeMax := now.AddDate(0, 0, 7).Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(timeMin).TimeMax(timeMax).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("No se pudieron obtener los eventos: %v", err)
	}

	return events.Items, nil
}

func CreateCalendarEvent(token *oauth2.Token, eventCreate *calendar.Event) (error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return err
	}

	client := config.Client(context.Background(), token)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	event := &calendar.Event{
		Summary:     "Reunión importante",
		Location:    "Oficina",
		Description: "Reunión para discutir el proyecto",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(3 * time.Hour).Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		},
	}

	calendarId := "primary"
	_, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return err
	}

	return nil
}

func UpdateCalendarEvent(token *oauth2.Token, eventUpdate *calendar.Event, eventId string) (error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return err
	}

	client := config.Client(context.Background(), token)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	// Obtenemos el evento existente
	event, err := srv.Events.Get("primary", eventId).Do()
	if err != nil {
		return err
	}

	// Editamos lo que queremos
	event.Summary = "Reunión modificada"
	event.Description = "Reunión actualizada"

	_, err = srv.Events.Update("primary", event.Id, event).Do()
	if err != nil {
		return err
	}

	return nil
}


func DeleteCalendarEvent(token *oauth2.Token, eventId string) error {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return err
	}

	client := config.Client(context.Background(), token)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	err = srv.Events.Delete("primary", eventId).Do()
	if err != nil {
		return err
	}

	return nil
}
