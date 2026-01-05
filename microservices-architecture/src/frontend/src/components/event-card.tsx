import Link from 'next/link';
import type { Event } from '@/types';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { CalendarDays, ArrowRight } from 'lucide-react';

interface EventCardProps {
  event: Event;
}

export function EventCard({ event }: EventCardProps) {
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
    <Card className="group hover:shadow-lg transition-shadow">
      <CardHeader>
        <CardTitle className="line-clamp-1">{event.name}</CardTitle>
        <CardDescription className="flex items-center gap-2">
          <CalendarDays className="h-4 w-4" />
          {formattedDate} - {formattedTime}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground line-clamp-2">
          {event.description}
        </p>
        <Button asChild className="w-full group-hover:bg-primary/90">
          <Link href={`/events/${event.id}`}>
            View Details
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}

