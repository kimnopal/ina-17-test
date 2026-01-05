'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import type { Event, Ticket } from '@/types';
import { bookingAPI } from '@/lib/api';
import { useAuth } from '@/lib/auth-context';
import { TicketCard } from '@/components/ticket-card';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { toast } from 'sonner';
import { ArrowLeft, CalendarDays, Loader2, TicketX } from 'lucide-react';

export default function EventDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const [event, setEvent] = useState<Event | null>(null);
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isBooking, setIsBooking] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const eventId = params.id as string;

  useEffect(() => {
    const fetchEventData = async () => {
      try {
        const [eventResponse, ticketsResponse] = await Promise.all([
          bookingAPI.getEventById(eventId),
          bookingAPI.getEventTickets(eventId),
        ]);
        setEvent(eventResponse.data);
        setTickets(ticketsResponse.data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load event');
      } finally {
        setIsLoading(false);
      }
    };

    if (eventId) {
      fetchEventData();
    }
  }, [eventId]);

  const handleBook = async (ticketId: string, quantity: number) => {
    if (!isAuthenticated) {
      toast.error('Please login to book tickets');
      router.push('/login');
      return;
    }

    setIsBooking(true);
    try {
      const response = await bookingAPI.createBooking({
        event_id: eventId,
        ticket_id: ticketId,
        quantity,
      });
      toast.success('Booking created successfully!');
      router.push(`/bookings/${response.data.id}`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to create booking');
    } finally {
      setIsBooking(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !event) {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-4 text-muted-foreground">
        <TicketX className="h-12 w-12" />
        <p className="text-lg font-medium">Event not found</p>
        <p className="text-sm">{error}</p>
        <Button asChild variant="outline">
          <Link href="/events">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Events
          </Link>
        </Button>
      </div>
    );
  }

  const eventDate = new Date(event.event_date);
  const formattedDate = eventDate.toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
  const formattedTime = eventDate.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div className="space-y-6">
      <Button asChild variant="ghost" className="-ml-4">
        <Link href="/events">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Events
        </Link>
      </Button>

      <div className="space-y-4">
        <h1 className="text-3xl font-bold tracking-tight">{event.name}</h1>
        <div className="flex items-center gap-2 text-muted-foreground">
          <CalendarDays className="h-5 w-5" />
          <span>{formattedDate} - {formattedTime}</span>
        </div>
        <p className="text-muted-foreground max-w-3xl">{event.description}</p>
      </div>

      <Separator />

      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Available Tickets</h2>
        {tickets.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
            <TicketX className="h-12 w-12 mb-2" />
            <p>No tickets available for this event</p>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {tickets.map((ticket) => (
              <TicketCard
                key={ticket.id}
                ticket={ticket}
                onBook={handleBook}
                isBooking={isBooking}
                isAuthenticated={isAuthenticated}
              />
            ))}
          </div>
        )}
      </div>

      {!isAuthenticated && tickets.length > 0 && (
        <div className="rounded-lg border bg-muted/50 p-4 text-center">
          <p className="text-sm text-muted-foreground mb-2">
            Please login to book tickets
          </p>
          <Button asChild size="sm">
            <Link href="/login">Login</Link>
          </Button>
        </div>
      )}
    </div>
  );
}

