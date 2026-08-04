export interface Device {
  id: number;
  name: string;
  ip_address: string;
  lab_id: number;
  status: 'online' | 'offline' | 'unknown' | string;
  last_seen: string | null;
}

export interface NewDevice {
  name: string;
  ip_address: string;
  lab_id: number;
}

export interface ScanResult {
  ip: string;
  status: 'online' | 'offline' | string;
  scanned_at: string;
}
