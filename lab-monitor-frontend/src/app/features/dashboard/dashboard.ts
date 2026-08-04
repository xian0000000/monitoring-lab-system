import { Component, DestroyRef, OnInit, computed, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';

import { DeviceService } from '../../core/services/device.service';
import { LabService } from '../../core/services/lab.service';
import { ScanService } from '../../core/services/scan.service';
import { Device, NewDevice } from '../../core/models/device.model';
import { Lab, NewLab } from '../../core/models/lab.model';

@Component({
  selector: 'app-dashboard',
  imports: [FormsModule, DatePipe],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css'
})
export class Dashboard implements OnInit {
  private deviceService = inject(DeviceService);
  private labService = inject(LabService);
  private scanService = inject(ScanService);
  private destroyRef = inject(DestroyRef);

  readonly devices = signal<Device[]>([]);
  readonly labs = signal<Lab[]>([]);
  readonly loading = signal(true);
  readonly error = signal<string | null>(null);
  readonly scanning = signal(false);
  readonly connected = this.scanService.connected;

  readonly onlineCount = computed(() => this.devices().filter((d) => d.status === 'online').length);
  readonly offlineCount = computed(() => this.devices().filter((d) => d.status !== 'online').length);
  readonly totalCount = computed(() => this.devices().length);

  readonly sortedDevices = computed(() =>
    [...this.devices()].sort((a, b) => a.ip_address.localeCompare(b.ip_address, undefined, { numeric: true }))
  );

  newDevice: NewDevice = { name: '', ip_address: '', lab_id: 0 };
  newLab: NewLab = { name: '', location: '', capacity: 0 };
  showDeviceForm = false;
  showLabForm = false;

  ngOnInit(): void {
    this.loadDevices();
    this.loadLabs();

    this.scanService.connect();
    this.scanService.results$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((result) => {
      this.devices.update((current) => {
        const idx = current.findIndex((d) => d.ip_address === result.ip);
        if (idx === -1) {
          const fresh: Device = {
            id: -Date.now(),
            name: result.ip,
            ip_address: result.ip,
            lab_id: 0,
            status: result.status,
            last_seen: result.scanned_at
          };
          return [...current, fresh];
        }
        const updated = [...current];
        updated[idx] = { ...updated[idx], status: result.status, last_seen: result.scanned_at };
        return updated;
      });
    });

    this.destroyRef.onDestroy(() => this.scanService.disconnect());
  }

  loadDevices(): void {
    this.loading.set(true);
    this.deviceService.getDevices().subscribe({
      next: (devices) => {
        this.devices.set(devices ?? []);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Gagal memuat daftar device. Pastikan backend jalan.');
        this.loading.set(false);
      }
    });
  }

  loadLabs(): void {
    this.labService.getLabs().subscribe({
      next: (labs) => this.labs.set(labs ?? []),
      error: () => {
        /* labs bersifat opsional untuk tampilan, cukup diamkan jika gagal */
      }
    });
  }

  labName(labId: number): string {
    return this.labs().find((l) => l.id === labId)?.name ?? '-';
  }

  triggerScan(): void {
    this.scanning.set(true);
    this.scanService.triggerScan().subscribe({
      next: () => setTimeout(() => this.scanning.set(false), 1500),
      error: () => this.scanning.set(false)
    });
  }

  submitDevice(): void {
    if (!this.newDevice.name || !this.newDevice.ip_address) {
      return;
    }
    this.deviceService.createDevice(this.newDevice).subscribe({
      next: (device) => {
        this.devices.update((current) => [...current, device]);
        this.newDevice = { name: '', ip_address: '', lab_id: 0 };
        this.showDeviceForm = false;
      },
      error: () => this.error.set('Gagal menambah device.')
    });
  }

  submitLab(): void {
    if (!this.newLab.name) {
      return;
    }
    this.labService.createLab(this.newLab).subscribe({
      next: (lab) => {
        this.labs.update((current) => [...current, lab]);
        this.newLab = { name: '', location: '', capacity: 0 };
        this.showLabForm = false;
      },
      error: () => this.error.set('Gagal menambah lab.')
    });
  }

  removeDevice(device: Device): void {
    if (device.id < 0) {
      return;
    }
    if (!confirm(`Hapus device "${device.name}"?`)) {
      return;
    }
    this.deviceService.deleteDevice(device.id).subscribe({
      next: () => this.devices.update((current) => current.filter((d) => d.id !== device.id)),
      error: () => this.error.set('Gagal menghapus device.')
    });
  }
}
