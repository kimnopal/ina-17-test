'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import type { Booking, Payment } from '@/types';
import { bookingAPI, paymentAPI } from '@/lib/api';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { toast } from 'sonner';
import {
  ArrowLeft,
  CheckCircle,
  CreditCard,
  Loader2,
  Smartphone,
  QrCode,
  Building,
} from 'lucide-react';

const paymentMethods = [
  { id: 'VA', name: 'Virtual Account', icon: Building, description: 'Pay via bank transfer' },
  { id: 'EWALLET', name: 'E-Wallet', icon: Smartphone, description: 'GoPay, OVO, Dana, etc.' },
  { id: 'QRIS', name: 'QRIS', icon: QrCode, description: 'Scan QR code to pay' },
];

export default function PaymentPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [booking, setBooking] = useState<Booking | null>(null);
  const [payment, setPayment] = useState<Payment | null>(null);
  const [selectedMethod, setSelectedMethod] = useState<string>('');
  const [isLoading, setIsLoading] = useState(true);
  const [isProcessing, setIsProcessing] = useState(false);
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

        // Check if booking is not pending, redirect back
        if (response.data.status !== 'PENDING') {
          router.push(`/bookings/${bookingId}`);
        }
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

  const handleCreatePayment = async () => {
    if (!selectedMethod) {
      toast.error('Please select a payment method');
      return;
    }

    if (!booking) return;

    setIsProcessing(true);
    try {
      const response = await paymentAPI.createPayment({
        booking_id: bookingId,
        amount: booking.total_amount,
        payment_method: selectedMethod,
      });
      setPayment(response.data);
      toast.success('Payment created! Please complete the payment.');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to create payment');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleSimulatePayment = async () => {
    if (!payment) return;

    setIsProcessing(true);
    try {
      await paymentAPI.simulatePaymentSuccess(payment.id);
      toast.success('Payment successful!');
      router.push(`/bookings/${bookingId}`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Payment failed');
    } finally {
      setIsProcessing(false);
    }
  };

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
        <CreditCard className="h-12 w-12" />
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

  return (
    <div className="space-y-6">
      <Button asChild variant="ghost" className="-ml-4">
        <Link href={`/bookings/${bookingId}`}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Booking
        </Link>
      </Button>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Payment Summary */}
        <Card>
          <CardHeader>
            <CardTitle>Payment Summary</CardTitle>
            <CardDescription>
              Booking ID: {booking.id.slice(0, 8)}...
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Tickets</span>
              <span className="font-medium">{booking.quantity}x</span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-lg font-medium">Total</span>
              <span className="text-2xl font-bold">{formattedAmount}</span>
            </div>
          </CardContent>
        </Card>

        {/* Payment Method / Status */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CreditCard className="h-5 w-5" />
              {payment ? 'Complete Payment' : 'Select Payment Method'}
            </CardTitle>
            <CardDescription>
              {payment
                ? 'Click the button below to simulate payment'
                : 'Choose your preferred payment method'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {payment ? (
              <div className="space-y-4">
                <div className="rounded-lg border bg-muted/50 p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-muted-foreground">Payment ID</span>
                    <span className="font-mono text-sm">{payment.id.slice(0, 8)}...</span>
                  </div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-muted-foreground">Method</span>
                    <span className="font-medium">{payment.payment_method}</span>
                  </div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-muted-foreground">Status</span>
                    <Badge variant="secondary">{payment.status}</Badge>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Amount</span>
                    <span className="font-bold">{formattedAmount}</span>
                  </div>
                </div>

                <div className="rounded-lg border border-dashed p-4 text-center">
                  <p className="text-sm text-muted-foreground mb-2">
                    In a real scenario, you would be redirected to the payment gateway.
                    Click below to simulate a successful payment.
                  </p>
                </div>

                <Button
                  className="w-full"
                  size="lg"
                  onClick={handleSimulatePayment}
                  disabled={isProcessing}
                >
                  {isProcessing ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <CheckCircle className="mr-2 h-4 w-4" />
                  )}
                  Simulate Successful Payment
                </Button>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label>Payment Method</Label>
                  <Select value={selectedMethod} onValueChange={setSelectedMethod}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select payment method" />
                    </SelectTrigger>
                    <SelectContent>
                      {paymentMethods.map((method) => (
                        <SelectItem key={method.id} value={method.id}>
                          <div className="flex items-center gap-2">
                            <method.icon className="h-4 w-4" />
                            <span>{method.name}</span>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {selectedMethod && (
                  <div className="rounded-lg border bg-muted/50 p-3">
                    {paymentMethods.find((m) => m.id === selectedMethod)?.description}
                  </div>
                )}

                <Button
                  className="w-full"
                  size="lg"
                  onClick={handleCreatePayment}
                  disabled={!selectedMethod || isProcessing}
                >
                  {isProcessing ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <CreditCard className="mr-2 h-4 w-4" />
                  )}
                  Create Payment
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

