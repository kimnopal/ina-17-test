'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import type { Booking } from '@/types';
import { bookingAPI } from '@/lib/api';
import { useAuth } from '@/lib/auth-context';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { ArrowLeft, CreditCard, Loader2, Ticket } from 'lucide-react';

const statusVariants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
  PENDING: 'secondary',
  PAID: 'default',
  CONFIRMED: 'default',
  CANCELLED: 'destructive',
};

export default function BookingDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [booking, setBooking] = useState<Booking | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const bookingId = params.id as string;

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/login');
      return;
    }

    const fetchBooking = async () => {
      if (!isAuthenticated) return;

      try {
        const response = await bookingAPI.getBookingById(bookingId);
        setBooking(response.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load booking');
      } finally {
        setIsLoading(false);
      }
    };

    if (isAuthenticated && bookingId) {
      fetchBooking();
    }
  }, [isAuthenticated, authLoading, bookingId, router]);

  if (authLoading || isLoading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !booking) {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-4 text-muted-foreground">
        <Ticket className="h-12 w-12" />
        <p className="text-lg font-medium">Booking not found</p>
        <p className="text-sm">{error}</p>
        <Button asChild variant="outline">
          <Link href="/bookings">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Bookings
          </Link>
        </Button>
      </div>
    );
  }

  const formattedAmount = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(booking.total_amount);

  const canPay = booking.status === 'PENDING';

  return (
    <div className="space-y-6">
      <Button asChild variant="ghost" className="-ml-4">
        <Link href="/bookings">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Bookings
        </Link>
      </Button>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Ticket className="h-5 w-5" />
              Booking Details
            </CardTitle>
            <CardDescription>
              Booking ID: {booking.id}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Status</span>
              <Badge variant={statusVariants[booking.status] || 'outline'}>
                {booking.status}
              </Badge>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Quantity</span>
              <span className="font-medium">{booking.quantity} ticket(s)</span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Total Amount</span>
              <span className="text-xl font-bold">{formattedAmount}</span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Created At</span>
              <span className="font-medium">
                {new Date(booking.created_at).toLocaleDateString('id-ID', {
                  day: '2-digit',
                  month: 'long',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </span>
            </div>
            {booking.expired_at && (
              <>
                <Separator />
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Expires At</span>
                  <span className="font-medium text-destructive">
                    {new Date(booking.expired_at).toLocaleDateString('id-ID', {
                      day: '2-digit',
                      month: 'long',
                      year: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </span>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CreditCard className="h-5 w-5" />
              Payment
            </CardTitle>
            <CardDescription>
              {canPay
                ? 'Complete your payment to confirm the booking'
                : booking.status === 'CONFIRMED' || booking.status === 'PAID'
                ? 'Payment completed'
                : 'This booking is no longer available for payment'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {canPay ? (
              <div className="space-y-4">
                <div className="rounded-lg border bg-muted/50 p-4 text-center">
                  <p className="text-sm text-muted-foreground mb-1">Amount to pay</p>
                  <p className="text-2xl font-bold">{formattedAmount}</p>
                </div>
                <Button asChild className="w-full" size="lg">
                  <Link href={`/payments/${booking.id}`}>
                    <CreditCard className="mr-2 h-4 w-4" />
                    Proceed to Payment
                  </Link>
                </Button>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
                {(booking.status === 'CONFIRMED' || booking.status === 'PAID') ? (
                  <>
                    <div className="rounded-full bg-green-100 p-3 mb-3">
                      <CreditCard className="h-6 w-6 text-green-600" />
                    </div>
                    <p className="font-medium text-green-600">Payment Successful</p>
                    <p className="text-sm">Your booking has been confirmed</p>
                  </>
                ) : (
                  <>
                    <CreditCard className="h-12 w-12 mb-2" />
                    <p>Payment is not available</p>
                  </>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

