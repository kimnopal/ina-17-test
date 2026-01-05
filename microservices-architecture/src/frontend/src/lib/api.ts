import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  AuthUserResponse,
  EventsResponse,
  EventResponse,
  TicketsResponse,
  CreateBookingRequest,
  BookingResponse,
  BookingsResponse,
  CreatePaymentRequest,
  PaymentResponse,
} from "@/types";

// Service base URLs
const USER_SERVICE_URL =
  process.env.NEXT_PUBLIC_USER_SERVICE_URL || "http://localhost:3001";
const BOOKING_SERVICE_URL =
  process.env.NEXT_PUBLIC_BOOKING_SERVICE_URL || "http://localhost:3002";
const PAYMENT_SERVICE_URL =
  process.env.NEXT_PUBLIC_PAYMENT_SERVICE_URL || "http://localhost:3003";

// Helper to get auth token
const getAccessToken = (): string | null => {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("access_token");
};

// Generic fetch wrapper
async function fetchAPI<T>(url: string, options: RequestInit = {}): Promise<T> {
  const token = getAccessToken();

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  if (token) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Something went wrong");
  }

  return data;
}

// ==================== User Service API ====================

export const userAPI = {
  register: (data: RegisterRequest) =>
    fetchAPI<{ message: string; data: { id: string; username: string } }>(
      `${USER_SERVICE_URL}/api/v1/users/`,
      {
        method: "POST",
        body: JSON.stringify(data),
      }
    ),

  login: (data: LoginRequest) =>
    fetchAPI<LoginResponse>(`${USER_SERVICE_URL}/api/v1/login`, {
      method: "POST",
      body: JSON.stringify(data),
    }),

  logout: (refreshToken: string) =>
    fetchAPI<{ message: string }>(`${USER_SERVICE_URL}/api/v1/logout`, {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),

  getAuthenticatedUser: () =>
    fetchAPI<AuthUserResponse>(`${USER_SERVICE_URL}/api/v1/users/auth`),

  refreshToken: (refreshToken: string) =>
    fetchAPI<{ message: string; access_token: string; expires_in: number }>(
      `${USER_SERVICE_URL}/api/v1/refresh`,
      {
        method: "POST",
        body: JSON.stringify({ refresh_token: refreshToken }),
      }
    ),
};

// ==================== Booking Service API ====================

export const bookingAPI = {
  // Events
  getAllEvents: () =>
    fetchAPI<EventsResponse>(`${BOOKING_SERVICE_URL}/api/v1/events/`),

  getEventById: (id: string) =>
    fetchAPI<EventResponse>(`${BOOKING_SERVICE_URL}/api/v1/events/${id}`),

  getEventTickets: (eventId: string) =>
    fetchAPI<TicketsResponse>(
      `${BOOKING_SERVICE_URL}/api/v1/events/${eventId}/tickets`
    ),

  // Bookings
  createBooking: (data: CreateBookingRequest) =>
    fetchAPI<BookingResponse>(`${BOOKING_SERVICE_URL}/api/v1/bookings/`, {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getAllBookings: () =>
    fetchAPI<BookingsResponse>(`${BOOKING_SERVICE_URL}/api/v1/bookings/`),

  getBookingById: (id: string) =>
    fetchAPI<BookingResponse>(`${BOOKING_SERVICE_URL}/api/v1/bookings/${id}`),
};

// ==================== Payment Service API ====================

export const paymentAPI = {
  createPayment: (data: CreatePaymentRequest) =>
    fetchAPI<PaymentResponse>(`${PAYMENT_SERVICE_URL}/api/v1/payments/`, {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getPaymentById: (id: string) =>
    fetchAPI<PaymentResponse>(`${PAYMENT_SERVICE_URL}/api/v1/payments/${id}`),

  // Simulate payment gateway webhook (for testing)
  simulatePaymentSuccess: (paymentId: string) =>
    fetchAPI<{ message: string }>(
      `${PAYMENT_SERVICE_URL}/api/v1/payments/webhook/payment-gateway`,
      {
        method: "POST",
        body: JSON.stringify({
          payment_id: paymentId,
          status: "PAID",
        }),
      }
    ),
};
