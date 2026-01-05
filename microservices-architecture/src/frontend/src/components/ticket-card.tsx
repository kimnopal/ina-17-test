'use client';

import { useState } from 'react';
import type { Ticket } from '@/types';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Minus, Plus, ShoppingCart } from 'lucide-react';

interface TicketCardProps {
  ticket: Ticket;
  onBook: (ticketId: string, quantity: number) => void;
  isBooking: boolean;
  isAuthenticated: boolean;
}

export function TicketCard({ ticket, onBook, isBooking, isAuthenticated }: TicketCardProps) {
  const [quantity, setQuantity] = useState(1);

  const formattedPrice = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(ticket.price);

  const handleQuantityChange = (delta: number) => {
    const newQuantity = quantity + delta;
    if (newQuantity >= 1 && newQuantity <= ticket.quota) {
      setQuantity(newQuantity);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value) || 1;
    if (value >= 1 && value <= ticket.quota) {
      setQuantity(value);
    }
  };

  const isOutOfStock = ticket.quota <= 0;

  return (
    <Card className={isOutOfStock ? 'opacity-60' : ''}>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">{ticket.category}</CardTitle>
          <Badge variant={isOutOfStock ? 'destructive' : 'secondary'}>
            {isOutOfStock ? 'Sold Out' : `${ticket.quota} left`}
          </Badge>
        </div>
        <CardDescription className="text-xl font-semibold text-foreground">
          {formattedPrice}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {!isOutOfStock && (
          <>
            <div className="space-y-2">
              <Label htmlFor={`quantity-${ticket.id}`}>Quantity</Label>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleQuantityChange(-1)}
                  disabled={quantity <= 1 || isBooking}
                >
                  <Minus className="h-4 w-4" />
                </Button>
                <Input
                  id={`quantity-${ticket.id}`}
                  type="number"
                  min={1}
                  max={ticket.quota}
                  value={quantity}
                  onChange={handleInputChange}
                  className="w-20 text-center"
                  disabled={isBooking}
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleQuantityChange(1)}
                  disabled={quantity >= ticket.quota || isBooking}
                >
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Total</span>
              <span className="font-semibold">
                {new Intl.NumberFormat('id-ID', {
                  style: 'currency',
                  currency: 'IDR',
                  minimumFractionDigits: 0,
                }).format(ticket.price * quantity)}
              </span>
            </div>
            <Button
              className="w-full"
              onClick={() => onBook(ticket.id, quantity)}
              disabled={isBooking || !isAuthenticated}
            >
              <ShoppingCart className="mr-2 h-4 w-4" />
              {!isAuthenticated ? 'Login to Book' : isBooking ? 'Processing...' : 'Book Now'}
            </Button>
          </>
        )}
      </CardContent>
    </Card>
  );
}

