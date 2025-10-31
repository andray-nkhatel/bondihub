export interface RentalAgreement {
  id: string;
  house_id: string;
  tenant_id: string;
  start_date: string;
  end_date: string;
  rent_amount: number;
  deposit: number;
  status: AgreementStatus;
  created_at: string;
  updated_at: string;
  house?: House;
  tenant?: User;
  payments?: Payment[];
}

export type AgreementStatus = 'active' | 'terminated' | 'expired';

export interface CreateRentalAgreementRequest {
  house_id: string;
  tenant_id: string;
  start_date: string;
  end_date: string;
  rent_amount: number;
  deposit: number;
}

export interface UpdateRentalAgreementRequest {
  status?: AgreementStatus;
}

export interface RentalAgreementFilters {
  status?: string;
  page?: number;
  limit?: number;
}

export interface RentalAgreementListResponse {
  agreements: RentalAgreement[];
  pagination: PaginationInfo;
}

// Import related types
import { House } from './house.model';
import { User } from './user.model';
import { Payment } from './payment.model';
import { PaginationInfo } from './house.model';
