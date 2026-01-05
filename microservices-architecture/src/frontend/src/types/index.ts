// User types
export interface User {
  id: string;
  username: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  message: string;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface AuthUserResponse {
  message: string;
  data: User;
}

// Event types
export interface Event {
  id: string;
  name: string;
  description: string;
  event_date: string;
  created_at: string;
}

export interface EventsResponse {
  message: string;
  data: Event[];
}

export interface EventResponse {
  message: string;
  data: Event;
}

// Ticket types
export interface Ticket {
  id: string;
  event_id: string;
  category: string;
  price: number;
  quota: number;
}

export interface TicketsResponse {
  message: string;
  data: Ticket[];
}

// Booking types
export interface Booking {
  id: string;
  user_id: string;
  event_id: string;
  ticket_id: string;
  quantity: number;
  total_amount: number;
  status: 'PENDING' | 'PAID' | 'CONFIRMED' | 'CANCELLED';
  expired_at?: string;
  created_at: string;
}

export interface CreateBookingRequest {
  event_id: string;
  ticket_id: string;
  quantity: number;
}

export interface BookingResponse {
  message: string;
  data: Booking;
}

export interface BookingsResponse {
  message: string;
  data: Booking[];
}

// Payment types
export interface Payment {
  id: string;
  booking_id: string;
  user_id: string;
  amount: number;
  currency: string;
  payment_method: string;
  status: 'PENDING' | 'PAID' | 'FAILED' | 'EXPIRED';
  expired_at?: string;
  paid_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePaymentRequest {
  booking_id: string;
  amount: number;
  payment_method: string;
}

export interface PaymentResponse {
  message: string;
  data: Payment;
}

export interface PaymentsResponse {
  message: string;
  data: Payment[];
}

// API Error
export interface APIError {
  error: string;
}

