export interface Payment {
  id: string;
  agreement_id: string;
  amount: number;
  payment_date: string;
  method: PaymentMethod;
  reference_no: string;
  status: PaymentStatus;
  commission: number;
  created_at: string;
  updated_at: string;
  agreement?: RentalAgreement;
}

export type PaymentMethod = 'MTN' | 'Airtel' | 'Cash' | 'Bank';

export type PaymentStatus = 'pending' | 'completed' | 'failed' | 'refunded';

export interface CreatePaymentRequest {
  agreement_id: string;
  amount: number;
  method: PaymentMethod;
  reference_no?: string;
}

export interface PaymentResult {
  success: boolean;
  transaction_id: string;
  reference_no: string;
  status: string;
  message: string;
}

export interface PaymentFilters {
  status?: string;
  method?: string;
  page?: number;
  limit?: number;
}

export interface PaymentListResponse {
  payments: Payment[];
  pagination: PaginationInfo;
}

export interface PaymentStats {
  total_payments: number;
  total_amount: number;
  completed_payments: number;
  completed_amount: number;
  pending_payments: number;
  failed_payments: number;
  payments_by_method: PaymentMethodStats[];
}

export interface PaymentMethodStats {
  method: string;
  count: number;
  amount: number;
}

// Import related types
import { RentalAgreement } from './rental.model';
import { PaginationInfo } from './house.model';
