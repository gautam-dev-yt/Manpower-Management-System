import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { format, formatDistanceToNow, differenceInDays } from 'date-fns';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Date formatting utilities
export function formatDate(date: string | Date): string {
  return format(new Date(date), 'MMM dd, yyyy');
}

export function formatDateTime(date: string | Date): string {
  return format(new Date(date), 'MMM dd, yyyy HH:mm');
}

export function getRelativeTime(date: string | Date): string {
  return formatDistanceToNow(new Date(date), { addSuffix: true });
}

export function getDaysUntil(date: string | Date): number {
  return differenceInDays(new Date(date), new Date());
}

// Document expiry status helper
export function getExpiryStatus(
  expiryDate: string
): 'expired' | 'urgent' | 'warning' | 'valid' {
  const daysLeft = getDaysUntil(expiryDate);

  if (daysLeft < 0) return 'expired';
  if (daysLeft <= 7) return 'urgent';
  if (daysLeft <= 30) return 'warning';
  return 'valid';
}

export function getExpiryBadgeColor(
  status: 'expired' | 'urgent' | 'warning' | 'valid'
): string {
  switch (status) {
    case 'expired':
      return 'bg-red-100 text-red-800 border-red-200';
    case 'urgent':
      return 'bg-orange-100 text-orange-800 border-orange-200';
    case 'warning':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'valid':
      return 'bg-green-100 text-green-800 border-green-200';
  }
}

// File size formatter
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

// Phone number formatter
export function formatPhoneNumber(phone: string): string {
  // Format international phone numbers for display
  if (phone.startsWith('+')) {
    return phone;
  }
  return `+${phone}`;
}
