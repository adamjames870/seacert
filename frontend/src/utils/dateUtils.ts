export const formatDate = (dateString: string | undefined | null): string => {
  if (!dateString) return 'N/A';
  
  const date = new Date(dateString);
  if (isNaN(date.getTime()) || date.getFullYear() <= 1) {
    return 'No Expiry'; // Or handle as needed for different contexts
  }

  const day = String(date.getDate()).padStart(2, '0');
  const month = date.toLocaleString('en-GB', { month: 'short' });
  const year = date.getFullYear();

  return `${day}-${month}-${year}`;
};

export const calculateExpiryDate = (issuedDate: string, validityMonths: number | null): string | null => {
  if (!issuedDate || validityMonths === null || validityMonths === undefined) {
    return null;
  }

  const date = new Date(issuedDate);
  if (isNaN(date.getTime())) {
    return null;
  }

  // Add months to the issued date
  const expiryDate = new Date(date);
  expiryDate.setMonth(expiryDate.getMonth() + validityMonths);
  
  // Usually, expiry is the day before (e.g., issued Jan 1, 1 year validity -> expires Dec 31)
  // But common practice can vary. Let's stick to simple month addition for now 
  // unless the requirement specifies otherwise. 
  // Let's subtract 1 day to be "the day before"
  expiryDate.setDate(expiryDate.getDate() - 1);

  return expiryDate.toISOString().split('T')[0];
};

export const calculateDaysInYear = (startDate: string, endDate: string, year: number): number => {
  const start = new Date(startDate);
  const end = new Date(endDate);
  
  // Create year boundaries in UTC to match 'YYYY-MM-DD' parsing
  const yearStart = new Date(`${year}-01-01`);
  const yearEnd = new Date(`${year}-12-31`);
  
  // Determine the overlap range
  const overlapStart = start > yearStart ? start : yearStart;
  const overlapEnd = end < yearEnd ? end : yearEnd;
  
  // If no overlap
  if (overlapStart > overlapEnd) {
    return 0;
  }
  
  // Calculate difference in days (inclusive)
  const diffTime = Math.abs(overlapEnd.getTime() - overlapStart.getTime());
  const diffDays = Math.round(diffTime / (1000 * 60 * 60 * 24)) + 1;
  
  return diffDays;
};
