'use client';

import { useEffect, useState } from 'react';
import type { Event } from '@/types';
import { bookingAPI } from '@/lib/api';
import { EventCard } from '@/components/event-card';
import { Loader2, CalendarX } from 'lucide-react';

export default function EventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const response = await bookingAPI.getAllEvents();
        setEvents(response.data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load events');
      } finally {
        setIsLoading(false);
      }
    };

    fetchEvents();
  }, []);

  if (isLoading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-2 text-muted-foreground">
        <CalendarX className="h-12 w-12" />
        <p className="text-lg font-medium">Failed to load events</p>
        <p className="text-sm">{error}</p>
      </div>
    );
  }

  if (events.length === 0) {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-2 text-muted-foreground">
        <CalendarX className="h-12 w-12" />
        <p className="text-lg font-medium">No events available</p>
        <p className="text-sm">Check back later for upcoming events</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Upcoming Events</h1>
        <p className="text-muted-foreground">
          Browse and book tickets for your favorite events
        </p>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {events.map((event) => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
    </div>
  );
}

