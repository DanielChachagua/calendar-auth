package services

import (
	"calendar_auth/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func GetCalendarUrl(url string) (string, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return "", nil
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
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

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
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

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
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

func CreateCalendarEvent(token *oauth2.Token, eventCreate *models.CalendarCreate) (*calendar.Event, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background(), token)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	event := &calendar.Event{
		Summary:     eventCreate.Summary,
		Location:    *eventCreate.Location,
		Description: *eventCreate.Description,
		Start: &calendar.EventDateTime{
			Date:     time.Now().Format("2006-01-02"),
			DateTime: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		},
	}

	if eventCreate.Time != nil {
		date := eventCreate.Date.ToTime()
		timePart := eventCreate.Time.ToTime()

		loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
		combined := time.Date(
			date.Year(), date.Month(), date.Day(),
			timePart.Hour(), timePart.Minute(), timePart.Second(), timePart.Nanosecond(),
			loc,
		)
		event.Start = &calendar.EventDateTime{
			DateTime: combined.Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
		event.End = &calendar.EventDateTime{
			DateTime: combined.Add(1 * time.Hour).Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
	} else {
		event.Start = &calendar.EventDateTime{
			Date:     eventCreate.Date.ToTime().Format("2006-01-02"),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
		event.End = &calendar.EventDateTime{
			Date:     eventCreate.Date.ToTime().Add(24 * time.Hour).Format("2006-01-02"),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
	}

	calendarId := "primary"
	eventCreated, err := srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		return nil, err
	}

	return eventCreated, nil
}

func UpdateCalendarEvent(token *oauth2.Token, eventUpdate *models.CalendarUpdate) (*calendar.Event, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background(), token)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Obtenemos el evento existente
	event, err := srv.Events.Get("primary", eventUpdate.ID).Do()
	if err != nil {
		return nil, err
	}

	// Editamos lo que queremos
	event.Summary = eventUpdate.Summary

	if eventUpdate.Location != nil {
		event.Location = *eventUpdate.Location
	}

	if eventUpdate.Description != nil {
		event.Description = *eventUpdate.Description
	}

	if eventUpdate.Time != nil {
		date := eventUpdate.Date.ToTime()
		timePart := eventUpdate.Time.ToTime()

		loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
		combined := time.Date(
			date.Year(), date.Month(), date.Day(),
			timePart.Hour(), timePart.Minute(), timePart.Second(), timePart.Nanosecond(),
			loc,
		)
		event.Start = &calendar.EventDateTime{
			DateTime: combined.Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
		event.End = &calendar.EventDateTime{
			DateTime: combined.Add(1 * time.Hour).Format(time.RFC3339),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
	} else {
		event.Start = &calendar.EventDateTime{
			Date:     eventUpdate.Date.ToTime().Format("2006-01-02"),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
		event.End = &calendar.EventDateTime{
			Date:     eventUpdate.Date.ToTime().Add(24 * time.Hour).Format("2006-01-02"),
			TimeZone: "America/Argentina/Buenos_Aires",
		}
	}

	updated, err := srv.Events.Update("primary", event.Id, event).Do()
	if err != nil {
		if gErr, ok := err.(*googleapi.Error); ok {
			switch gErr.Code {
			case 404:
				fmt.Printf("⚠️ El evento no existe o ya fue eliminado: %s\n", gErr.Message)
				return nil, fmt.Errorf("el evento no existe o ya fue eliminado: %s", gErr.Message)

			case 410:
				fmt.Printf("⚠️ El recurso fue eliminado permanentemente: %s\n", gErr.Message)
				return nil, fmt.Errorf("recurso eliminado permanentemente: %s", gErr.Message)

			default:
				fmt.Printf("❌ Error de Google API (código %d): %s\n", gErr.Code, gErr.Message)
				return nil, gErr
			}
		}
		return nil, err
	}

	return updated, nil
}

func DeleteCalendarEvents(token *oauth2.Token, eventIds []string) error {
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

	for _, eventId := range eventIds {
		err := srv.Events.Delete("primary", eventId).Do()
		if err != nil {
			if gErr, ok := err.(*googleapi.Error); ok {
				switch gErr.Code {
				case 404:
					fmt.Printf("⚠️ El evento no existe o ya fue eliminado: %s\n", gErr.Message)
					return fmt.Errorf("el evento no existe o ya fue eliminado: %s", gErr.Message)

				case 410:
					fmt.Printf("⚠️ El recurso fue eliminado permanentemente: %s\n", gErr.Message)
					return fmt.Errorf("recurso eliminado permanentemente: %s", gErr.Message)

				default:
					fmt.Printf("❌ Error de Google API (código %d): %s\n", gErr.Code, gErr.Message)
					return gErr
				}
			}

			return fmt.Errorf("error eliminando evento %s: %w", eventId, err)
		}
	}

	return nil
}
